output "id" {
  description = "ID of the namespace."
  value       = scaleway_registry_namespace.this.id
}

output "endpoint" {
  description = "URL of the namespace."
  value       = scaleway_registry_namespace.this.endpoint
}
