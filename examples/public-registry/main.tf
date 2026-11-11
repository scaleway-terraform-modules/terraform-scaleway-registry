provider "scaleway" {}

module "registry" {
  source = "../.."

  name        = "example-public-registry"
  description = "Public images of the example project"
  is_public   = true
}
