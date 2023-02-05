resource "scaleway_registry_namespace" "this" {
  description = var.description
  is_public   = var.is_public
  name        = var.name
  project_id  = var.project_id
  region      = var.region
}
