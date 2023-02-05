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
