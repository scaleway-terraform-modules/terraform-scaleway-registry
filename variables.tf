variable "description" {
  type        = string
  description = "Description of the namespace."
  default     = null
}

variable "is_public" {
  type        = bool
  description = "Whether the images stored in the namespace should be downloadable publicly (docker pull)."
  default     = false
  nullable    = false
}

variable "name" {
  type        = string
  description = "Unique name of the namespace."
}

variable "project_id" {
  type        = string
  description = "ID of the project the namespace is associated with. Ressource will be created in the project set at the provider level if null."
  default     = null
}

variable "purge_dry_run" {
  type        = bool
  description = "Prevent deletion of tags & images."
  default     = true
  nullable    = false
}

variable "purge_enabled" {
  type        = bool
  description = "Whether to deploy the purge function, with its function namespace, cron trigger, and the IAM application, policy and API key it uses. When disabled, none of them is created, so no API key secret is stored in the Terraform state."
  default     = true
  nullable    = false
}

variable "purge_log_level" {
  type        = number
  description = "Log Level of the purge function. Refer to the Go [slog package documentation](https://pkg.go.dev/log/slog#Level) for accepted values"
  default     = 0
  nullable    = false
}

variable "purge_preserve_tag_patterns" {
  type        = list(string)
  description = "List of regex patterns for tags to preserve. By default, preserves semantic versioning tags and `latest` tags."
  default = [
    "^latest((-.+)?)$",
    "^(0|[1-9]\\d*)\\.(0|[1-9]\\d*)\\.(0|[1-9]\\d*)(?:-((?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\\.(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\\+([0-9a-zA-Z-]+(?:\\.[0-9a-zA-Z-]+)*))?$"
  ]
  nullable = false
}

variable "purge_retention_days" {
  type        = number
  description = "Number of days to keep a tag after its last update (or its creation, if it has no update date) before deleting it. Tags without either date are never deleted."
  default     = 30
  nullable    = false
}

variable "purge_schedule" {
  type        = string
  description = "Cron schedule of the purge function. Set to an empty string to only run it by hand."
  default     = "0 0 * * *"
  nullable    = false
}

variable "purge_timeout" {
  type        = number
  description = "Timeout of the purge function, in seconds."
  default     = 300
  nullable    = false
}

variable "region" {
  type        = string
  description = "Region in which the namespace should be created. Ressource will be created in the region set at the provider level if null."
  default     = null
}
