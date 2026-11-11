output "namespace_endpoint" {
  description = "Address of the namespace, as used by `docker login`, `docker push` and `docker pull`."
  value       = module.registry.namespace_endpoint
}

output "namespace_id" {
  description = "ID of the namespace."
  value       = module.registry.namespace_id
}

output "purge_function_id" {
  description = "ID of the purge function."
  value       = module.registry.purge_function_id
}
