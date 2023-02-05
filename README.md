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
| <a name="requirement_scaleway"></a> [scaleway](#requirement_scaleway) | >= 2.0.0 |

## Resources

| Name | Type |
|------|------|
| [scaleway_registry_namespace.this](https://registry.terraform.io/providers/scaleway/scaleway/latest/docs/resources/registry_namespace) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_name"></a> [name](#input_name) | Unique name of the namespace. | `string` | n/a | yes |
| <a name="input_description"></a> [description](#input_description) | Description of the namespace. | `string` | `null` | no |
| <a name="input_is_public"></a> [is_public](#input_is_public) | Whether the images stored in the namespace should be downloadable publicly (docker pull). | `bool` | `false` | no |
| <a name="input_project_id"></a> [project_id](#input_project_id) | ID of the project the namespace is associated with. Ressource will be created in the project set at the provider level if null. | `string` | `null` | no |
| <a name="input_region"></a> [region](#input_region) | Region in which the namespace should be created. Ressource will be created in the region set at the provider level if null. | `string` | `null` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_endpoint"></a> [endpoint](#output_endpoint) | URL of the namespace. |
| <a name="output_id"></a> [id](#output_id) | ID of the namespace. |
<!-- END_TF_DOCS -->

## Authors

Module is maintained with help from [the community](https://github.com/scaleway-terraform-modules/terraform-scaleway-registry/graphs/contributors).

## License

Mozilla Public License 2.0 Licensed. See [LICENSE](https://github.com/scaleway-terraform-modules/terraform-scaleway-registry/tree/master/LICENSE) for full details.
