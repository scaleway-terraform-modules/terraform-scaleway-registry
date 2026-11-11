provider "scaleway" {}

module "registry" {
  source = "../.."

  name = "example-purge-disabled"

  # No purge function, IAM application or API key
  purge_enabled = false
}
