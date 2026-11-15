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

  validation {
    condition     = !var.purge_enabled || can(regex("^[a-z0-9-]+$", var.name))
    error_message = "'name' must only contain lowercase letters, digits and dashes when 'purge_enabled' is true: it is used in the names of the purge function and its namespace."
  }

  validation {
    condition     = !var.purge_enabled || length(var.name) <= 28
    error_message = "'name' must be at most 28 characters long when 'purge_enabled' is true: the purge function is named '<name>-purge', and function names are limited to 34 characters."
  }
}

variable "project_id" {
  type        = string
  description = "ID of the project the namespace is associated with. Resource will be created in the project set at the provider level if null."
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
  description = "Log level of the purge function. Refer to the Go [slog package documentation](https://pkg.go.dev/log/slog#Level) for accepted values."
  default     = 0
  nullable    = false

  validation {
    condition     = contains([-4, 0, 4, 8], var.purge_log_level)
    error_message = "'purge_log_level' must be one of `-4` (Debug), `0` (Info), `4` (Warn) or `8` (Error)."
  }
}

variable "purge_preserve_tag_patterns" {
  type        = list(string)
  description = "List of regex patterns ([RE2 syntax](https://github.com/google/re2/wiki/Syntax)) for tags to preserve. Must contain at least one non-empty pattern. By default, preserves semantic versioning tags and `latest` tags."
  default = [
    "^latest((-.+)?)$",
    "^(0|[1-9]\\d*)\\.(0|[1-9]\\d*)\\.(0|[1-9]\\d*)(?:-((?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\\.(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\\+([0-9a-zA-Z-]+(?:\\.[0-9a-zA-Z-]+)*))?$"
  ]
  nullable = false

  validation {
    condition     = length(var.purge_preserve_tag_patterns) > 0
    error_message = "'purge_preserve_tag_patterns' must contain at least one pattern: an empty list preserves every tag. Set 'purge_dry_run' to true or 'purge_enabled' to false to keep every tag."
  }

  validation {
    condition     = alltrue([for pattern in var.purge_preserve_tag_patterns : pattern != ""])
    error_message = "'purge_preserve_tag_patterns' must not contain empty patterns: an empty pattern matches, and preserves, every tag."
  }

  validation {
    condition     = alltrue([for pattern in var.purge_preserve_tag_patterns : can(regexall(pattern, ""))])
    error_message = "'purge_preserve_tag_patterns' values must be valid regular expressions (RE2 syntax)."
  }
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
  description = "Region in which the namespace should be created. Resource will be created in the region set at the provider level if null."
  default     = null
}
