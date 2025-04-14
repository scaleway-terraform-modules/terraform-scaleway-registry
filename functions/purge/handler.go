package handler

import (
	"context"
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

// Default semver pattern for tag preservation
const DefaultTagPattern string = `^latest((-.+)?)$|^(0|[1-9]\\d*)\\.(0|[1-9]\\d*)\\.(0|[1-9]\\d*)(?:-((?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\\.(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\\+([0-9a-zA-Z-]+(?:\\.[0-9a-zA-Z-]+)*))?$`

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

// RegistryClient handles operations with the Scaleway Registry API
type RegistryClient struct {
	client    *registry.API
	config    *Config
	namespace *registry.Namespace
}

var config Config
var scw_client *scw.Client

var (
	ErrMissingAccessKey         = &ConfigError{Field: "REGISTRY_ACCESS_KEY", Message: "Scaleway access key is required"}
	ErrMissingProjectID         = &ConfigError{Field: "REGISTRY_PROJECT_ID", Message: "Scaleway project ID is required"}
	ErrMissingRegion            = &ConfigError{Field: "REGISTRY_REGION", Message: "Region is required"}
	ErrMissingRegistryNamespace = &ConfigError{Field: "REGISTRY_NAMESPACE", Message: "Registry namespace is required"}
	ErrMissingSecretKey         = &ConfigError{Field: "REGISTRY_SECRET_KEY", Message: "Scaleway secret key is required"}
)

func init() {
	slog.Info("===== Runing 'init' =====")
	// Set log level based on environment variable
	var logLevel slog.Level
	slogLevel, err := strconv.Atoi(os.Getenv("LOG_LEVEL"))
	if err != nil {
		slog.Warn("Failed to parse LOG_LEVEL. Enforcing level to Info")
		slogLevel = 0
	}
	switch slogLevel {
	case -4.:
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
		slog.Error("Scaleway access key is required")
		panic(ErrMissingAccessKey)
	}
	if config.ProjectID == "" {
		slog.Error("Scaleway project ID is required")
		panic(ErrMissingProjectID)
	}
	if config.Region == "" {
		slog.Error("Region is required")
		panic(ErrMissingRegion)
	}
	if config.RegistryNamespace == "" {
		slog.Error("Registry namespace is required")
		panic(ErrMissingRegistryNamespace)
	}
	if config.SecretKey == "" {
		slog.Error("Scaleway secret key is required")
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
		fmt.Errorf("Failed to parse PRESERVE_TAG_PATTERNS: %w. Enforcing SemVer 2.0.0 tags", err)
		config.PreserveTagPatterns = regexp.MustCompile(DefaultTagPattern)
	}

	slog.Info(fmt.Sprintf("Dry Run mode: %t.", config.DryRun))
	slog.Info(fmt.Sprintf("Retention Period: %d.", config.RetentionDays))
	slog.Info(fmt.Sprintf("Tags to preserve: %s", config.PreserveTagPatterns.String()))

	scw_client, err = scw.NewClient(
		scw.WithAuth(config.AccessKey, config.SecretKey),
		scw.WithDefaultRegion(scw.Region(config.Region)),
		scw.WithDefaultProjectID(config.ProjectID),
	)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to create Scaleway client: %w.", err))
		panic(err)
	}
}

func Handle(w http.ResponseWriter, r *http.Request) {
	slog.Info("===== Runing 'Handle' =====")
	// Get Retention date
	retentionDate := time.Now().AddDate(0, 0, -config.RetentionDays)
	slog.Info(fmt.Sprintf("Purging images older than: %s", retentionDate.Format(time.RFC3339)))

	// Create registry API client
	registryAPI := registry.NewAPI(scw_client)
	// Create registry client
	registryClient := &RegistryClient{
		client: registryAPI,
		config: &config,
	}

	// Get namespace
	namespace, err := registryClient.getNamespace(context.Background())
	if err != nil {
		slog.Error(fmt.Sprintf("failed to get registry namespace: %w.", err))
		panic(err)
	}
	registryClient.namespace = namespace

	// Get all images
	images, err := registryClient.ListImages(context.Background())
	if err != nil {
		slog.Error(fmt.Sprintf("failed to get images from namespace: %w.", err))
		panic(err)
	}

	// Track statistics
	var errorImages int
	// var deletedImages, preservedImages, skippedImages int
	var deletedTags, errorTags, preservedTags, skippedTags int

	// Process each image
	for _, image := range images {
		// List tags for the image
		tags, err := registryClient.ListTags(context.Background(), image)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to list tags for image %s: %v", image.ID, err))
			errorImages++
			continue
		}

		// Process each tag
		for _, tag := range tags {
			// Check if tag is older than retention date
			if tag.UpdatedAt.Before(retentionDate) {
				slog.Debug(fmt.Sprintf("Tag %s:%s updated at %s is older than retention date %s",
					image.Name, tag.Name, tag.UpdatedAt.Format(time.RFC3339), retentionDate.Format(time.RFC3339)))

				// Check if tag is protected
				if config.PreserveTagPatterns.MatchString(tag.Name) {
					slog.Debug(fmt.Sprintf("Tag %s match protection pattern. Preserving", tag.Name))
					preservedTags++
				} else {
					slog.Debug(fmt.Sprintf("Tag %s doesn't match protection pattern.", tag.Name))
					// Delete tag if not dry run
					if err := registryClient.DeleteTag(context.Background(), tag); err != nil {
						slog.Error(fmt.Sprintf("Failed to delete tag %s (%s): %v", tag.Name, tag.ID, err))
						errorTags++
						continue
					}
					deletedTags++
				}
			} else {
				slog.Debug(fmt.Sprintf("Tag %s updated at %s is newer than retention date %s. Skipping",
					tag.Name, tag.UpdatedAt.Format(time.RFC3339), retentionDate.Format(time.RFC3339)))
				skippedTags++
			}
		}
	}

	// Log summary
	slog.Info(fmt.Sprintf("Purge summary: %d tags deleted, %d tags skipped, %d tags preserved, %d errors",
		deletedTags, skippedTags, preservedTags, errorTags))
	slog.Info(fmt.Sprintf("Purge summary: %d images errors",
		errorImages))
	// slog.Info(fmt.Sprintf("Purge summary: %d images deleted, %d images skipped, %d images preserved, %d errors",
	// 	deletedImages, skippedImages, preservedImages, errorImages))
}

// DeleteTag deletes a tag
func (r *RegistryClient) DeleteTag(ctx context.Context, tag *registry.Tag) error {
	if r.config.DryRun {
		slog.Debug(fmt.Sprintf("[DRY RUN] Would delete tag: %s (%s)", tag.Name, tag.ID))
		return nil
	}

	slog.Info(fmt.Sprintf("Deleting tag: %s (%s)", tag.Name, tag.ID))
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
	slog.Info(fmt.Sprintf("Getting namespace: %s.", r.config.RegistryNamespace))

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

// ListImages lists all images in the registry namespace
func (r *RegistryClient) ListImages(ctx context.Context) ([]*registry.Image, error) {
	slog.Info(fmt.Sprintf("Listing images in namespace: %s", r.namespace.Name))

	var allImages []*registry.Image
	var page int32 = 1
	var pageSize uint32 = 100

	for {
		res, err := r.client.ListImages(&registry.ListImagesRequest{
			Region:      scw.Region(r.config.Region),
			NamespaceID: &r.namespace.ID,
			PageSize:    &pageSize, // uint32
			Page:        &page,     // int32
		}, scw.WithContext(ctx))
		if err != nil {
			return nil, fmt.Errorf("failed to list images: %w", err)
		}

		allImages = append(allImages, res.Images...)

		// Check if we've reached the last page
		if len(res.Images) < int(pageSize) {
			break
		}

		// Move to the next page
		page++
	}

	slog.Info(fmt.Sprintf("Found %d images in namespace: %s", len(allImages), r.namespace.Name))
	return allImages, nil
}

// ListTags lists all tags for an image
func (r *RegistryClient) ListTags(ctx context.Context, image *registry.Image) ([]*registry.Tag, error) {
	slog.Info(fmt.Sprintf("Listing tags for image: %s", image.Name))

	var allTags []*registry.Tag
	var page int32 = 1
	var pageSize uint32 = 100

	for {
		res, err := r.client.ListTags(&registry.ListTagsRequest{
			Region:   scw.Region(r.config.Region),
			ImageID:  image.ID,
			PageSize: &pageSize,
			Page:     &page,
		}, scw.WithContext(ctx))
		if err != nil {
			return nil, fmt.Errorf("failed to list tags for image %s: %w", image.Name, err)
		}

		allTags = append(allTags, res.Tags...)

		// Check if we've reached the last page
		if len(res.Tags) < int(pageSize) {
			break
		}

		// Move to the next page
		page++
	}

	slog.Info(fmt.Sprintf("Found %d tags for image: %s", len(allTags), image.Name))
	return allTags, nil
}
