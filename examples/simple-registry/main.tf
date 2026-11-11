provider "scaleway" {}

module "registry" {
  source = "../.."

  name = "example-simple-registry"
}
