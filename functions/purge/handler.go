package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/scaleway/scaleway-sdk-go/api/registry/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// DefaultTagPattern preserves `latest`, `latest-*` and SemVer 2.0.0 tags, as the
// default of the purge_preserve_tag_patterns variable. It is a raw string: use
// single backslashes.
const DefaultTagPattern string = `^latest((-.+)?)$|^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`

// Config holds the configuration for the purge function
type Config struct {
	// AccessKey is the Scaleway API access key
	AccessKey string
	// SecretKey is the Scaleway API secret key
	SecretKey string
	// Region is the Scaleway region where the registry is located
	Region string
	// ProjectID is the Scaleway project ID
	ProjectID string
	// RegistryNamespace is the name of the registry namespace
	RegistryNamespace string
	// RetentionDays is the number of days to keep images before deletion
	RetentionDays int
	// DryRun if true, will not actually delete images
	DryRun bool
	// PreserveTagPatterns is a regex pattern for tags to preserve
	PreserveTagPatterns *regexp.Regexp
}

// ConfigError represents a configuration error
type ConfigError struct {
	Field   string
	Message string
}

// Error returns the message of the configuration error, prefixed by the faulty variable
func (e *ConfigError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// Summary holds the result of a purge run, returned as the HTTP response body
type Summary struct {
	// DryRun is true if tags were not actually deleted
	DryRun bool `json:"dry_run"`
	// RetentionDate is the date before which tags are deleted (RFC 3339)
	RetentionDate string `json:"retention_date"`
	// DeletedTags is the number of tags deleted, or that would be deleted in dry-run mode
	DeletedTags int `json:"deleted_tags"`
	// SkippedTags is the number of tags kept because they are recent or have no date
	SkippedTags int `json:"skipped_tags"`
	// PreservedTags is the number of tags kept because they match a preserved pattern
	PreservedTags int `json:"preserved_tags"`
	// ErrorTags is the number of tags that failed to be deleted
	ErrorTags int `json:"error_tags"`
	// ErrorImages is the number of images whose tags failed to be listed
	ErrorImages int `json:"error_images"`
	// Error is the error that stopped the run, if any
	Error string `json:"error,omitempty"`
}

// RegistryClient handles operations with the Scaleway Registry API
type RegistryClient struct {
	client    *registry.API
	config    *Config
	namespace *registry.Namespace
}

// listPageSize is the number of images or tags fetched per API call
const listPageSize = 100

var config Config
var scwClient *scw.Client

var (
	ErrMissingAccessKey         = &ConfigError{Field: "REGISTRY_ACCESS_KEY", Message: "Scaleway access key is required"}
	ErrMissingProjectID         = &ConfigError{Field: "REGISTRY_PROJECT_ID", Message: "Scaleway project ID is required"}
	ErrMissingRegion            = &ConfigError{Field: "REGISTRY_REGION", Message: "Region is required"}
	ErrMissingRegistryNamespace = &ConfigError{Field: "REGISTRY_NAMESPACE", Message: "Registry namespace is required"}
	ErrMissingSecretKey         = &ConfigError{Field: "REGISTRY_SECRET_KEY", Message: "Scaleway secret key is required"}
)

func init() {
	slog.Info("===== Running 'init' =====")
	// Set log level based on environment variable
	var logLevel slog.Level
	slogLevel, err := strconv.Atoi(os.Getenv("LOG_LEVEL"))
	if err != nil {
		slog.Warn("Failed to parse LOG_LEVEL. Enforcing level to Info")
		slogLevel = 0
	}
	switch slogLevel {
	case -4:
		logLevel = slog.LevelDebug
	case 0:
		logLevel = slog.LevelInfo
	case 4:
		logLevel = slog.LevelWarn
	case 8:
		logLevel = slog.LevelError
	default:
		// invalid log level, handle error
		slog.Warn("Invalid log level. Enforcing level to Info")
		logLevel = slog.LevelInfo
	}
	slog.SetLogLoggerLevel(logLevel)

	config.AccessKey = os.Getenv("REGISTRY_ACCESS_KEY")
	config.ProjectID = os.Getenv("REGISTRY_PROJECT_ID")
	config.Region = os.Getenv("REGISTRY_REGION")
	config.RegistryNamespace = os.Getenv("REGISTRY_NAMESPACE")
	config.SecretKey = os.Getenv("REGISTRY_SECRET_KEY")

	// Validate required configuration
	if config.AccessKey == "" {
		panic(ErrMissingAccessKey)
	}
	if config.ProjectID == "" {
		panic(ErrMissingProjectID)
	}
	if config.Region == "" {
		panic(ErrMissingRegion)
	}
	if config.RegistryNamespace == "" {
		panic(ErrMissingRegistryNamespace)
	}
	if config.SecretKey == "" {
		panic(ErrMissingSecretKey)
	}

	// Validate optional configuration
	config.DryRun, err = strconv.ParseBool(os.Getenv("DRY_RUN"))
	if err != nil {
		slog.Warn("Failed to parse DRY_RUN. Enforcing Dry Run mode")
		config.DryRun = true
	}
	config.RetentionDays, err = strconv.Atoi(os.Getenv("RETENTION_DAYS"))
	if err != nil || config.RetentionDays < 1 {
		slog.Warn("Failed to parse RETENTION_DAYS. Enforcing Dry Run mode with a 30 days retention period")
		config.DryRun = true
		config.RetentionDays = 30
	}

	config.PreserveTagPatterns, err = regexp.Compile(os.Getenv("PRESERVE_TAG_PATTERNS"))
	if err != nil {
		slog.Error("Failed to parse PRESERVE_TAG_PATTERNS. Enforcing Dry Run mode, and preserving latest and SemVer 2.0.0 tags", "error", err)
		config.DryRun = true
		config.PreserveTagPatterns = regexp.MustCompile(DefaultTagPattern)
	}

	slog.Info("Dry Run mode", "dry_run", config.DryRun)
	slog.Info("Retention period", "retention_days", config.RetentionDays)
	slog.Info("Tags to preserve", "pattern", config.PreserveTagPatterns.String())

	scwClient, err = scw.NewClient(
		scw.WithAuth(config.AccessKey, config.SecretKey),
		scw.WithDefaultRegion(scw.Region(config.Region)),
		scw.WithDefaultProjectID(config.ProjectID),
	)
	if err != nil {
		panic(fmt.Errorf("failed to create Scaleway client: %w", err))
	}
}

func Handle(w http.ResponseWriter, r *http.Request) {
	slog.Info("===== Running 'Handle' =====")
	// Get Retention date
	retentionDate := time.Now().AddDate(0, 0, -config.RetentionDays)
	slog.Info("Purging images older than the retention date", "retention_date", retentionDate.Format(time.RFC3339))

	// Cancel the API calls when the request ends, e.g. when the function times out
	ctx := r.Context()
	summary := Summary{
		DryRun:        config.DryRun,
		RetentionDate: retentionDate.Format(time.RFC3339),
	}

	// Create registry API client
	registryAPI := registry.NewAPI(scwClient)
	// Create registry client
	registryClient := &RegistryClient{
		client: registryAPI,
		config: &config,
	}

	// Get namespace
	namespace, err := registryClient.getNamespace(ctx)
	if err != nil {
		slog.Error("Failed to get registry namespace", "error", err)
		summary.Error = fmt.Sprintf("failed to get registry namespace: %v", err)
		writeSummary(w, http.StatusInternalServerError, summary)
		return
	}
	registryClient.namespace = namespace

	// Get all images
	images, err := registryClient.listImages(ctx)
	if err != nil {
		slog.Error("Failed to get images from namespace", "error", err)
		summary.Error = fmt.Sprintf("failed to get images from namespace: %v", err)
		writeSummary(w, http.StatusInternalServerError, summary)
		return
	}

	// Track statistics
	var errorImages int
	// var deletedImages, preservedImages, skippedImages int
	var deletedTags, errorTags, preservedTags, skippedTags int

	// Process each image
	for _, image := range images {
		// List tags for the image
		tags, err := registryClient.listTags(ctx, image)
		if err != nil {
			slog.Error("Failed to list tags", "image", image.Name, "image_id", image.ID, "error", err)
			errorImages++
			continue
		}

		// Process each tag
		for _, tag := range tags {
			// Date the tag with its update date, or its creation date if it has none
			tagDate := tag.UpdatedAt
			if tagDate == nil {
				slog.Warn("Tag has no update date. Using its creation date", "image", image.Name, "tag", tag.Name, "tag_id", tag.ID)
				tagDate = tag.CreatedAt
			}
			// Skip tags without any date: their age is unknown, so never delete them
			if tagDate == nil {
				slog.Warn("Tag has no creation date either. Skipping", "image", image.Name, "tag", tag.Name, "tag_id", tag.ID)
				skippedTags++
				continue
			}

			// Check if tag is older than retention date
			if tagDate.Before(retentionDate) {
				slog.Debug("Tag is older than the retention date",
					"image", image.Name, "tag", tag.Name, "updated_at", tagDate.Format(time.RFC3339), "retention_date", retentionDate.Format(time.RFC3339))

				// Check if tag is protected
				if config.PreserveTagPatterns.MatchString(tag.Name) {
					slog.Debug("Tag matches a preserved pattern. Preserving", "image", image.Name, "tag", tag.Name)
					preservedTags++
				} else {
					slog.Debug("Tag doesn't match any preserved pattern", "image", image.Name, "tag", tag.Name)
					// Delete tag if not dry run
					if err := registryClient.deleteTag(ctx, tag); err != nil {
						slog.Error("Failed to delete tag", "image", image.Name, "tag", tag.Name, "tag_id", tag.ID, "error", err)
						errorTags++
						continue
					}
					deletedTags++
				}
			} else {
				slog.Debug("Tag is newer than the retention date. Skipping",
					"image", image.Name, "tag", tag.Name, "updated_at", tagDate.Format(time.RFC3339), "retention_date", retentionDate.Format(time.RFC3339))
				skippedTags++
			}
		}
	}

	// Log summary
	slog.Info("Purge summary",
		"deleted_tags", deletedTags, "skipped_tags", skippedTags, "preserved_tags", preservedTags,
		"error_tags", errorTags, "error_images", errorImages)
	// slog.Info(fmt.Sprintf("Purge summary: %d images deleted, %d images skipped, %d images preserved, %d errors",
	// 	deletedImages, skippedImages, preservedImages, errorImages))

	// Return the summary, with an error status if any tag or image failed
	summary.DeletedTags = deletedTags
	summary.SkippedTags = skippedTags
	summary.PreservedTags = preservedTags
	summary.ErrorTags = errorTags
	summary.ErrorImages = errorImages
	status := http.StatusOK
	if errorTags > 0 || errorImages > 0 {
		status = http.StatusInternalServerError
	}
	writeSummary(w, status, summary)
}

// writeSummary writes the summary of a purge run as a JSON response
func writeSummary(w http.ResponseWriter, status int, summary Summary) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(summary); err != nil {
		slog.Error("Failed to write the response", "error", err)
	}
}

// deleteTag deletes a tag
func (r *RegistryClient) deleteTag(ctx context.Context, tag *registry.Tag) error {
	if r.config.DryRun {
		slog.Info("[DRY RUN] Would delete tag", "tag", tag.Name, "tag_id", tag.ID)
		return nil
	}

	slog.Info("Deleting tag", "tag", tag.Name, "tag_id", tag.ID)
	_, err := r.client.DeleteTag(&registry.DeleteTagRequest{
		Region: scw.Region(r.config.Region),
		TagID:  tag.ID,
	}, scw.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to delete tag %s (%s): %w", tag.Name, tag.ID, err)
	}

	return nil
}

// getNamespace gets the registry namespace
func (r *RegistryClient) getNamespace(ctx context.Context) (*registry.Namespace, error) {
	slog.Info("Getting namespace", "namespace", r.config.RegistryNamespace)

	// List namespaces to find the one with the matching name
	res, err := r.client.ListNamespaces(&registry.ListNamespacesRequest{
		Region:    scw.Region(r.config.Region),
		ProjectID: &r.config.ProjectID,
		Name:      &r.config.RegistryNamespace,
	}, scw.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	// Check if namespace exists
	if len(res.Namespaces) == 0 {
		return nil, fmt.Errorf("namespace not found: %s", r.config.RegistryNamespace)
	}

	// Return the first matching namespace
	return res.Namespaces[0], nil
}

// listImages lists all images in the registry namespace
func (r *RegistryClient) listImages(ctx context.Context) ([]*registry.Image, error) {
	slog.Info("Listing images", "namespace", r.namespace.Name)

	res, err := r.client.ListImages(&registry.ListImagesRequest{
		Region:      scw.Region(r.config.Region),
		NamespaceID: &r.namespace.ID,
		PageSize:    scw.Uint32Ptr(listPageSize),
	}, scw.WithContext(ctx), scw.WithAllPages())
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}

	slog.Info("Found images", "namespace", r.namespace.Name, "count", len(res.Images))
	return res.Images, nil
}

// listTags lists all tags for an image
func (r *RegistryClient) listTags(ctx context.Context, image *registry.Image) ([]*registry.Tag, error) {
	slog.Info("Listing tags", "image", image.Name)

	res, err := r.client.ListTags(&registry.ListTagsRequest{
		Region:   scw.Region(r.config.Region),
		ImageID:  image.ID,
		PageSize: scw.Uint32Ptr(listPageSize),
	}, scw.WithContext(ctx), scw.WithAllPages())
	if err != nil {
		return nil, fmt.Errorf("failed to list tags for image %s: %w", image.Name, err)
	}

	slog.Info("Found tags", "image", image.Name, "count", len(res.Tags))
	return res.Tags, nil
}
