# Terraform / Scaleway

## Purpose

This repository is used to manage a container registry on scaleway using terraform.

## Usage

- Setup the [scaleway provider](https://www.terraform.io/docs/providers/scaleway/index.html) in your tf file.
- Include this module in your tf file. Refer to [documentation](https://www.terraform.io/docs/modules/sources.html#generic-git-repository).

```hcl
module "my_registry" {
  source  = "scaleway-terraform-modules/registry/scaleway"
  version = "0.0.1"

}
```

<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement_terraform) | >= 0.13 |
| <a name="requirement_archive"></a> [archive](#requirement_archive) | >= 2.0.0 |
| <a name="requirement_scaleway"></a> [scaleway](#requirement_scaleway) | >= 2.0.0 |

## Resources

| Name | Type |
|------|------|
| [archive_file.registry_purge](https://registry.terraform.io/providers/hashicorp/archive/latest/docs/resources/file) | resource |
| [scaleway_function.registry_purge](https://registry.terraform.io/providers/scaleway/scaleway/latest/docs/resources/function) | resource |
| [scaleway_function_cron.registry_purge](https://registry.terraform.io/providers/scaleway/scaleway/latest/docs/resources/function_cron) | resource |
| [scaleway_function_namespace.maintenance](https://registry.terraform.io/providers/scaleway/scaleway/latest/docs/resources/function_namespace) | resource |
| [scaleway_iam_api_key.registry_purge](https://registry.terraform.io/providers/scaleway/scaleway/latest/docs/resources/iam_api_key) | resource |
| [scaleway_iam_application.registry_purge](https://registry.terraform.io/providers/scaleway/scaleway/latest/docs/resources/iam_application) | resource |
| [scaleway_iam_policy.registry_purge](https://registry.terraform.io/providers/scaleway/scaleway/latest/docs/resources/iam_policy) | resource |
| [scaleway_registry_namespace.this](https://registry.terraform.io/providers/scaleway/scaleway/latest/docs/resources/registry_namespace) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_name"></a> [name](#input_name) | Unique name of the namespace. | `string` | n/a | yes |
| <a name="input_description"></a> [description](#input_description) | Description of the namespace. | `string` | `null` | no |
| <a name="input_is_public"></a> [is_public](#input_is_public) | Whether the images stored in the namespace should be downloadable publicly (docker pull). | `bool` | `false` | no |
| <a name="input_project_id"></a> [project_id](#input_project_id) | ID of the project the namespace is associated with. Ressource will be created in the project set at the provider level if null. | `string` | `null` | no |
| <a name="input_purge_dry_run"></a> [purge_dry_run](#input_purge_dry_run) | Prevent deletion of tags & images. | `bool` | `true` | no |
| <a name="input_purge_log_level"></a> [purge_log_level](#input_purge_log_level) | Log Level of the purge function. Refer to the Go [slog package documentation](https://pkg.go.dev/log/slog#Level) for accepted values | `number` | `0` | no |
| <a name="input_purge_preserve_tag_patterns"></a> [purge_preserve_tag_patterns](#input_purge_preserve_tag_patterns) | List of regex patterns for tags to preserve. By default, preserves semantic versioning tags and `latest` tags. | `list(string)` | ```[ "^latest((-.+)?)$", "^(0|[1-9]\\d*)\\.(0|[1-9]\\d*)\\.(0|[1-9]\\d*)(?:-((?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\\.(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\\+([0-9a-zA-Z-]+(?:\\.[0-9a-zA-Z-]+)*))?$" ]``` | no |
| <a name="input_purge_retention_days"></a> [purge_retention_days](#input_purge_retention_days) | Number of days to keep images before deletion. | `number` | `30` | no |
| <a name="input_purge_schedule"></a> [purge_schedule](#input_purge_schedule) | Cron schedule for the cleanup function. Set to empty string to disable scheduled cleanup. | `string` | `"0 0 * * *"` | no |
| <a name="input_purge_timeout"></a> [purge_timeout](#input_purge_timeout) | Timeout for the cleanup function in seconds. | `number` | `300` | no |
| <a name="input_region"></a> [region](#input_region) | Region in which the namespace should be created. Ressource will be created in the region set at the provider level if null. | `string` | `null` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_endpoint"></a> [endpoint](#output_endpoint) | URL of the namespace. |
| <a name="output_id"></a> [id](#output_id) | ID of the namespace. |
| <a name="output_maintenance_namespace_id"></a> [maintenance_namespace_id](#output_maintenance_namespace_id) | ID of the maintenance function namespace. |
| <a name="output_purge_application_id"></a> [purge_application_id](#output_purge_application_id) | ID of the IAM application for the cleanup function. |
| <a name="output_purge_cron_id"></a> [purge_cron_id](#output_purge_cron_id) | ID of the cron trigger for the cleanup function. |
| <a name="output_purge_function_endpoint"></a> [purge_function_endpoint](#output_purge_function_endpoint) | Endpoint URL of the cleanup function. |
| <a name="output_purge_function_id"></a> [purge_function_id](#output_purge_function_id) | ID of the cleanup function. |
<!-- END_TF_DOCS -->

## Authors

Module is maintained with help from [the community](https://github.com/scaleway-terraform-modules/terraform-scaleway-registry/graphs/contributors).

## License

Mozilla Public License 2.0 Licensed. See [LICENSE](https://github.com/scaleway-terraform-modules/terraform-scaleway-registry/tree/master/LICENSE) for full details.
