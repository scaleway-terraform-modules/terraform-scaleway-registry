output "maintenance_namespace_id" {
  description = "ID of the maintenance function namespace. Null when `purge_enabled` is false."
  value       = one(scaleway_function_namespace.maintenance[*].id)
}

output "namespace_endpoint" {
  description = "Address of the namespace, as used by `docker login`, `docker push` and `docker pull` (`rg.{region}.scw.cloud/{name}`)."
  value       = scaleway_registry_namespace.this.endpoint
}

output "namespace_id" {
  description = "ID of the namespace."
  value       = scaleway_registry_namespace.this.id
}

output "namespace_name" {
  description = "Name of the namespace."
  value       = scaleway_registry_namespace.this.name
}

output "namespace_project_id" {
  description = "ID of the project containing the namespace."
  value       = scaleway_registry_namespace.this.project_id
}

output "namespace_region" {
  description = "Region of the namespace."
  value       = scaleway_registry_namespace.this.region
}

# Purge function outputs
output "purge_application_id" {
  description = "ID of the IAM application used by the purge function. Null when `purge_enabled` is false."
  value       = one(scaleway_iam_application.registry_purge[*].id)
}

output "purge_cron_id" {
  description = "ID of the cron trigger of the purge function. Null when `purge_enabled` is false or `purge_schedule` is empty."
  value       = one(scaleway_function_cron.registry_purge[*].id)
}

output "purge_function_endpoint" {
  description = "Domain name of the purge function (without `https://`). Null when `purge_enabled` is false."
  value       = one(scaleway_function.registry_purge[*].domain_name)
}

output "purge_function_id" {
  description = "ID of the purge function. Null when `purge_enabled` is false."
  value       = one(scaleway_function.registry_purge[*].id)
}
