# Simple registry

This example shows how to create a private container registry namespace with the default settings.

The module also deploys the purge function, in dry-run mode: every day at midnight, it logs the tags older than 30 days that it would delete, without deleting them. See the [purge-enabled](../purge-enabled) example to actually delete tags, and the [purge-disabled](../purge-disabled) example to go without the function.

## Usage

Configure the [scaleway provider](https://www.terraform.io/docs/providers/scaleway/index.html) with your own credentials, then:

```sh
terraform init
terraform apply
```

To push an image, log in with a Scaleway API secret key, then tag and push it under the namespace's address:

```sh
REGISTRY=$(terraform output -raw namespace_endpoint)
docker login "${REGISTRY}" -u nologin --password-stdin <<< "${SCW_SECRET_KEY}"
docker tag my-image:1.0.0 "${REGISTRY}/my-image:1.0.0"
docker push "${REGISTRY}/my-image:1.0.0"
```
