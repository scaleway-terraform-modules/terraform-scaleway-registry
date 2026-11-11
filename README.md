# Terraform / Scaleway

## Purpose

This repository is used to manage a container registry on scaleway using terraform.

## Usage

- Setup the [scaleway provider](https://www.terraform.io/docs/providers/scaleway/index.html) in your tf file.
- Include this module in your tf file. Refer to [documentation](https://www.terraform.io/docs/modules/sources.html#generic-git-repository).

```hcl
module "my_registry" {
  source  = "scaleway-terraform-modules/registry/scaleway"
  version = "~> 1.0"

}
```

See the [examples](examples) folder for runnable configurations: a [simple](examples/simple-registry) private registry, a [public](examples/public-registry) one, a purge that [actually deletes tags](examples/purge-enabled), and a registry [without the purge function](examples/purge-disabled).

## Purge function

Along with the namespace, the module deploys a Scaleway Function that deletes old tags from it. It runs on the `purge_schedule` cron schedule (every day at midnight by default; set it to `""` to only run it by hand), and logs what it does in Cockpit.

**The function runs in dry-run mode by default:** it only logs the tags it would delete. Set `purge_dry_run = false` to actually delete them.

The function calls Scaleway's API with an API key the module creates for it: an IAM application with the `ContainerRegistryFullAccess` permission set on the namespace's **whole project**, and an API key that never expires. The key's secret is stored in the Terraform state, so keep the state in a protected backend. If you don't need the purge, set `purge_enabled = false`: the module then creates neither the function nor its IAM resources, and stores no key.

For each tag of each image in the namespace, the function:

1. Dates the tag with its last update date. If the tag has no update date, it logs a warning and uses its creation date instead. If it has no creation date either, it logs a second warning and keeps the tag: its age is unknown.
2. Keeps the tag if it is more recent than `purge_retention_days` days.
3. Keeps the tag if its name matches one of the `purge_preserve_tag_patterns` regular expressions (by default `latest`, `latest-*` and [semantic versions](https://semver.org/) such as `1.2.3`).
4. Deletes the tag otherwise.

At the end of each run, it logs how many tags were deleted, skipped, preserved, or failed to be deleted. It also returns this summary as JSON, with status `200`, or `500` when a tag couldn't be deleted or an image's tags couldn't be listed. When the namespace or its images can't be listed, it stops and returns `500` with the error.

<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement_terraform) | >= 1.11 |
| <a name="requirement_archive"></a> [archive](#requirement_archive) | >= 2.4.0 |
| <a name="requirement_scaleway"></a> [scaleway](#requirement_scaleway) | >= 2.79.0 |

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
| <a name="input_project_id"></a> [project_id](#input_project_id) | ID of the project the namespace is associated with. Resource will be created in the project set at the provider level if null. | `string` | `null` | no |
| <a name="input_purge_dry_run"></a> [purge_dry_run](#input_purge_dry_run) | Prevent deletion of tags & images. | `bool` | `true` | no |
| <a name="input_purge_enabled"></a> [purge_enabled](#input_purge_enabled) | Whether to deploy the purge function, with its function namespace, cron trigger, and the IAM application, policy and API key it uses. When disabled, none of them is created, so no API key secret is stored in the Terraform state. | `bool` | `true` | no |
| <a name="input_purge_log_level"></a> [purge_log_level](#input_purge_log_level) | Log level of the purge function. Refer to the Go [slog package documentation](https://pkg.go.dev/log/slog#Level) for accepted values. | `number` | `0` | no |
| <a name="input_purge_preserve_tag_patterns"></a> [purge_preserve_tag_patterns](#input_purge_preserve_tag_patterns) | List of regex patterns for tags to preserve. By default, preserves semantic versioning tags and `latest` tags. | `list(string)` | ```[ "^latest((-.+)?)$", "^(0|[1-9]\\d*)\\.(0|[1-9]\\d*)\\.(0|[1-9]\\d*)(?:-((?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\\.(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\\+([0-9a-zA-Z-]+(?:\\.[0-9a-zA-Z-]+)*))?$" ]``` | no |
| <a name="input_purge_retention_days"></a> [purge_retention_days](#input_purge_retention_days) | Number of days to keep a tag after its last update (or its creation, if it has no update date) before deleting it. Tags without either date are never deleted. | `number` | `30` | no |
| <a name="input_purge_schedule"></a> [purge_schedule](#input_purge_schedule) | Cron schedule of the purge function. Set to an empty string to only run it by hand. | `string` | `"0 0 * * *"` | no |
| <a name="input_purge_timeout"></a> [purge_timeout](#input_purge_timeout) | Timeout of the purge function, in seconds. | `number` | `300` | no |
| <a name="input_region"></a> [region](#input_region) | Region in which the namespace should be created. Resource will be created in the region set at the provider level if null. | `string` | `null` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_maintenance_namespace_id"></a> [maintenance_namespace_id](#output_maintenance_namespace_id) | ID of the maintenance function namespace. Null when `purge_enabled` is false. |
| <a name="output_namespace_endpoint"></a> [namespace_endpoint](#output_namespace_endpoint) | Address of the namespace, as used by `docker login`, `docker push` and `docker pull` (`rg.{region}.scw.cloud/{name}`). |
| <a name="output_namespace_id"></a> [namespace_id](#output_namespace_id) | ID of the namespace. |
| <a name="output_namespace_name"></a> [namespace_name](#output_namespace_name) | Name of the namespace. |
| <a name="output_namespace_project_id"></a> [namespace_project_id](#output_namespace_project_id) | ID of the project containing the namespace. |
| <a name="output_namespace_region"></a> [namespace_region](#output_namespace_region) | Region of the namespace. |
| <a name="output_purge_application_id"></a> [purge_application_id](#output_purge_application_id) | ID of the IAM application used by the purge function. Null when `purge_enabled` is false. |
| <a name="output_purge_cron_id"></a> [purge_cron_id](#output_purge_cron_id) | ID of the cron trigger of the purge function. Null when `purge_enabled` is false or `purge_schedule` is empty. |
| <a name="output_purge_function_endpoint"></a> [purge_function_endpoint](#output_purge_function_endpoint) | Domain name of the purge function (without `https://`). Null when `purge_enabled` is false. |
| <a name="output_purge_function_id"></a> [purge_function_id](#output_purge_function_id) | ID of the purge function. Null when `purge_enabled` is false. |
<!-- END_TF_DOCS -->

## Authors

Module is maintained with help from [the community](https://github.com/scaleway-terraform-modules/terraform-scaleway-registry/graphs/contributors).

## License

Mozilla Public License 2.0 Licensed. See [LICENSE](https://github.com/scaleway-terraform-modules/terraform-scaleway-registry/tree/master/LICENSE) for full details.
