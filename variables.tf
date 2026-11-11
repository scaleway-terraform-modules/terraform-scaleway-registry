variable "name" {
  description = "Unique name of the namespace."
  type        = string
}

variable "description" {
  description = "Description of the namespace."
  type        = string
  default     = null
}

variable "is_public" {
  description = "Whether the images stored in the namespace should be downloadable publicly (docker pull)."
  type        = bool
  default     = false
  nullable    = false
}

variable "region" {
  description = "Region in which the namespace should be created. Ressource will be created in the region set at the provider level if null."
  type        = string
  default     = null
}

variable "project_id" {
  description = "ID of the project the namespace is associated with. Ressource will be created in the project set at the provider level if null."
  type        = string
  default     = null
}

# Purge function variables
variable "purge_enabled" {
  description = "Whether to deploy the purge function, with its function namespace, cron trigger, and the IAM application, policy and API key it uses. When disabled, none of them is created, so no API key secret is stored in the Terraform state."
  type        = bool
  default     = true
  nullable    = false
}

variable "purge_retention_days" {
  description = "Number of days to keep a tag after its last update (or its creation, if it has no update date) before deleting it. Tags without either date are never deleted."
  type        = number
  default     = 30
  nullable    = false
}

variable "purge_preserve_tag_patterns" {
  description = "List of regex patterns for tags to preserve. By default, preserves semantic versioning tags and `latest` tags."
  type        = list(string)
  default = [
    "^latest((-.+)?)$",
    "^(0|[1-9]\\d*)\\.(0|[1-9]\\d*)\\.(0|[1-9]\\d*)(?:-((?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\\.(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\\+([0-9a-zA-Z-]+(?:\\.[0-9a-zA-Z-]+)*))?$"
  ]
  nullable = false
}

variable "purge_schedule" {
  description = "Cron schedule of the purge function. Set to an empty string to only run it by hand."
  type        = string
  default     = "0 0 * * *"
  nullable    = false
}

variable "purge_timeout" {
  description = "Timeout of the purge function, in seconds."
  type        = number
  default     = 300
  nullable    = false
}

variable "purge_dry_run" {
  description = "Prevent deletion of tags & images."
  type        = bool
  default     = true
  nullable    = false
}

variable "purge_log_level" {
  description = "Log Level of the purge function. Refer to the Go [slog package documentation](https://pkg.go.dev/log/slog#Level) for accepted values"
  type        = number
  default     = 0
  nullable    = false
}
