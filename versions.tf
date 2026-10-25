terraform {
  required_version = ">= 1.11"
  required_providers {
    scaleway = {
      source  = "scaleway/scaleway"
      version = ">= 2.79.0"
    }
    archive = {
      source  = "hashicorp/archive"
      version = ">= 2.0.0"
    }
  }
}
