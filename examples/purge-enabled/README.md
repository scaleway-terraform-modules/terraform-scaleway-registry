# Purge enabled

This example shows how to make the purge function actually delete old tags. **It deletes tags:** try it on a namespace whose images you can rebuild.

The function:

- runs every Sunday at 03:00, instead of every day at midnight;
- deletes the tags that haven't been updated for 14 days, instead of 30;
- never deletes `latest`, `main` and version tags such as `1.2.3` or `v1.2.3`. These patterns replace the default ones, which preserve `latest`, `latest-*` and semantic versions.

What it deletes is logged in Cockpit.

## Usage

Configure the [scaleway provider](https://www.terraform.io/docs/providers/scaleway/index.html) with your own credentials, then:

```sh
terraform init
terraform apply
```
