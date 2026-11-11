provider "scaleway" {}

module "registry" {
  source = "../.."

  name = "example-purge-enabled"

  # Actually delete tags, instead of only logging them
  purge_dry_run = false
  # Delete tags not updated for 2 weeks
  purge_retention_days = 14
  # Every Sunday at 03:00
  purge_schedule = "0 3 * * 0"
  # Never delete `latest`, `main`, nor version tags such as `1.2.3` or `v1.2.3`
  purge_preserve_tag_patterns = [
    "^latest$",
    "^main$",
    "^v?\\d+\\.\\d+\\.\\d+$",
  ]
}
