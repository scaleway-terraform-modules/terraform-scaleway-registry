# Create a serverless namespace for maintenance functions
resource "scaleway_function_namespace" "maintenance" {
  name        = format("%s-maintenance", var.name)
  description = format("Namespace for maintenance functions for %s registry", var.name)
  project_id  = scaleway_registry_namespace.this.project_id
  region      = scaleway_registry_namespace.this.region
}

# Create an IAM application for the function
resource "scaleway_iam_application" "registry_purge" {
  name        = format("%s-registry-purge", var.name)
  description = format("Application for the registry %s purge function", var.name)
}

# Create an IAM policy to allow the function to access the registry
resource "scaleway_iam_policy" "registry_purge" {
  name           = format("%s-registry-purge-policy", var.name)
  description    = format("Policy for registry %s purge function", var.name)
  application_id = scaleway_iam_application.registry_purge.id

  # Rule to allow the function to read and delete images in the registry namespace
  rule {
    project_ids          = [scaleway_registry_namespace.this.project_id]
    permission_set_names = ["ContainerRegistryFullAccess"]
  }
}

# Create an API key for the function
resource "scaleway_iam_api_key" "registry_purge" {
  application_id     = scaleway_iam_application.registry_purge.id
  description        = format("API key for registry %s purge function", var.name)
  default_project_id = scaleway_registry_namespace.this.project_id
}

# Archive the function code
resource "archive_file" "registry_purge" {
  type        = "zip"
  source_dir  = "${path.module}/functions/purge"
  output_path = "${path.module}/functions/purge.zip"
  excludes    = [".git", ".gitignore", "README.md", "vendor"]
}

# Deploy the function
resource "scaleway_function" "registry_purge" {
  name         = format("%s-purge", var.name)
  description  = format("Purges old images from the container registry namespace %s", var.name)
  deploy       = true
  namespace_id = scaleway_function_namespace.maintenance.id
  handler      = "Handle"
  http_option  = "redirected"
  privacy      = "private"
  runtime      = "go123"
  timeout      = var.purge_timeout

  memory_limit = 128
  max_scale    = 1
  min_scale    = 0

  zip_file = archive_file.registry_purge.output_path
  zip_hash = archive_file.registry_purge.output_sha256

  environment_variables = {
    REGISTRY_NAMESPACE    = var.name
    RETENTION_DAYS        = var.purge_retention_days
    PRESERVE_TAG_PATTERNS = join(",", var.purge_preserve_tag_patterns)
    DRY_RUN               = var.purge_dry_run
    LOG_LEVEL             = var.purge_log_level
  }

  secret_environment_variables = {
    REGISTRY_ACCESS_KEY = scaleway_iam_api_key.registry_purge.access_key
    REGISTRY_SECRET_KEY = scaleway_iam_api_key.registry_purge.secret_key
    REGISTRY_PROJECT_ID = scaleway_registry_namespace.this.project_id
    REGISTRY_REGION     = scaleway_registry_namespace.this.region
  }
}

# Create a cron trigger for the function only if cleanup_schedule is not empty
resource "scaleway_function_cron" "registry_purge" {
  count       = var.purge_schedule != "" ? 1 : 0
  function_id = scaleway_function.registry_purge.id
  schedule    = var.purge_schedule
  args        = "{}"
}
