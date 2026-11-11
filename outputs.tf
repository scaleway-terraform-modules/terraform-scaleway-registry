output "id" {
  description = "ID of the namespace."
  value       = scaleway_registry_namespace.this.id
}

output "endpoint" {
  description = "URL of the namespace."
  value       = scaleway_registry_namespace.this.endpoint
}

output "maintenance_namespace_id" {
  description = "ID of the maintenance function namespace. Null when `purge_enabled` is false."
  value       = one(scaleway_function_namespace.maintenance[*].id)
}

# Cleanup function outputs
output "purge_application_id" {
  description = "ID of the IAM application for the cleanup function. Null when `purge_enabled` is false."
  value       = one(scaleway_iam_application.registry_purge[*].id)
}

output "purge_cron_id" {
  description = "ID of the cron trigger for the cleanup function. Null when `purge_enabled` is false or `purge_schedule` is empty."
  value       = one(scaleway_function_cron.registry_purge[*].id)
}

output "purge_function_endpoint" {
  description = "Endpoint URL of the cleanup function. Null when `purge_enabled` is false."
  value       = one(scaleway_function.registry_purge[*].domain_name)
}

output "purge_function_id" {
  description = "ID of the cleanup function. Null when `purge_enabled` is false."
  value       = one(scaleway_function.registry_purge[*].id)
}
