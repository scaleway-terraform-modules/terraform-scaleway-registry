locals {
  # Number of instances of each purge resource: none when the purge is disabled
  purge_count = var.purge_enabled ? 1 : 0
}

# Create a serverless namespace for maintenance functions
resource "scaleway_function_namespace" "maintenance" {
  count       = local.purge_count
  name        = format("%s-maintenance", var.name)
  description = format("Namespace for maintenance functions for %s registry", var.name)
  project_id  = scaleway_registry_namespace.this.project_id
  region      = scaleway_registry_namespace.this.region
}

# Create an IAM application for the function
resource "scaleway_iam_application" "registry_purge" {
  count       = local.purge_count
  name        = format("%s-registry-purge", var.name)
  description = format("Application for the registry %s purge function", var.name)
}

# Create an IAM policy to allow the function to access the registry
resource "scaleway_iam_policy" "registry_purge" {
  count          = local.purge_count
  name           = format("%s-registry-purge-policy", var.name)
  description    = format("Policy for registry %s purge function", var.name)
  application_id = scaleway_iam_application.registry_purge[0].id

  # Rule to allow the function to read and delete images in the registry namespace
  rule {
    project_ids          = [scaleway_registry_namespace.this.project_id]
    permission_set_names = ["ContainerRegistryFullAccess"]
  }
}

# Create an API key for the function
resource "scaleway_iam_api_key" "registry_purge" {
  count              = local.purge_count
  application_id     = scaleway_iam_application.registry_purge[0].id
  description        = format("API key for registry %s purge function", var.name)
  default_project_id = scaleway_registry_namespace.this.project_id
}

# Archive the function code
resource "archive_file" "registry_purge" {
  count       = local.purge_count
  type        = "zip"
  source_dir  = "${path.module}/functions/purge"
  output_path = "${path.module}/functions/purge-${var.name}.zip" # One archive per registry: module calls sourced from the same path share path.module
  excludes    = [".git", ".gitignore", "README.md", "vendor"]
}

# Deploy the function
resource "scaleway_function" "registry_purge" {
  count        = local.purge_count
  name         = format("%s-purge", var.name)
  description  = format("Purges old images from the container registry namespace %s", var.name)
  deploy       = true
  namespace_id = scaleway_function_namespace.maintenance[0].id
  handler      = "Handle"
  http_option  = "redirected"
  privacy      = "private"
  runtime      = "go126" # Must provide the Go version required by functions/purge/go.mod: the CI builds with it.
  timeout      = var.purge_timeout

  memory_limit = 128
  max_scale    = 1
  min_scale    = 0

  zip_file = archive_file.registry_purge[0].output_path
  zip_hash = archive_file.registry_purge[0].output_sha256

  environment_variables = {
    REGISTRY_NAMESPACE    = var.name
    RETENTION_DAYS        = var.purge_retention_days
    PRESERVE_TAG_PATTERNS = join("|", var.purge_preserve_tag_patterns)
    DRY_RUN               = var.purge_dry_run
    LOG_LEVEL             = var.purge_log_level
  }

  secret_environment_variables = {
    REGISTRY_ACCESS_KEY = scaleway_iam_api_key.registry_purge[0].access_key
    REGISTRY_SECRET_KEY = scaleway_iam_api_key.registry_purge[0].secret_key
    REGISTRY_PROJECT_ID = scaleway_registry_namespace.this.project_id
    REGISTRY_REGION     = scaleway_registry_namespace.this.region
  }
}

# Create a cron trigger for the function only if the purge is enabled and purge_schedule is not empty
resource "scaleway_function_cron" "registry_purge" {
  count       = var.purge_enabled && var.purge_schedule != "" ? 1 : 0
  function_id = scaleway_function.registry_purge[0].id
  schedule    = var.purge_schedule
  args        = "{}"
}

# Keep the purge resources created before they became optional, instead of replacing them
moved {
  from = scaleway_function_namespace.maintenance
  to   = scaleway_function_namespace.maintenance[0]
}

moved {
  from = scaleway_iam_application.registry_purge
  to   = scaleway_iam_application.registry_purge[0]
}

moved {
  from = scaleway_iam_policy.registry_purge
  to   = scaleway_iam_policy.registry_purge[0]
}

moved {
  from = scaleway_iam_api_key.registry_purge
  to   = scaleway_iam_api_key.registry_purge[0]
}

moved {
  from = archive_file.registry_purge
  to   = archive_file.registry_purge[0]
}

moved {
  from = scaleway_function.registry_purge
  to   = scaleway_function.registry_purge[0]
}
