# Purge disabled

This example shows how to create a container registry namespace without the purge function.

With `purge_enabled = false`, the module creates neither the function nor its function namespace, cron trigger, IAM application, policy and API key. No API key secret is stored in the Terraform state, and old tags are never deleted.

## Usage

Configure the [scaleway provider](https://www.terraform.io/docs/providers/scaleway/index.html) with your own credentials, then:

```sh
terraform init
terraform apply
```
