# Public registry

This example shows how to create a public container registry namespace: anyone can pull its images, without logging in. Pushing images still requires a Scaleway API key.

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

To pull an image, no login is needed:

```sh
docker pull "$(terraform output -raw namespace_endpoint)/my-image:1.0.0"
```
