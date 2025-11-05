variable "gcp_project_id" {
  description = "The GCP project ID to deploy to."
  type        = string
}

variable "gcp_region" {
  description = "The GCP region to deploy to."
  type        = string
  default     = "us-central1"
}

variable "service_name" {
  description = "The name of the service."
  type        = string
  default     = "dolina-flower-order-backend"
}

variable "db_name" {
  description = "The name of the database."
  type        = string
  default     = "dolina_flowers"
}

variable "db_user" {
  description = "The username for the database."
  type        = string
  default     = "postgres"
}

variable "db_password" {
  description = "The password for the database user."
  type        = string
  sensitive   = true
}
