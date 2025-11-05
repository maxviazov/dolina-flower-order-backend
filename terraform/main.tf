terraform {
  required_version = ">= 1.0"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

provider "google" {
  project = var.gcp_project_id
  region  = var.gcp_region
}

# Enable required Google Cloud services
resource "google_project_service" "run" {
  service = "run.googleapis.com"
}

resource "google_project_service" "sqladmin" {
  service = "sqladmin.googleapis.com"
}

resource "google_project_service" "artifactregistry" {
  service = "artifactregistry.googleapis.com"
}

resource "google_project_service" "cloudbuild" {
  service = "cloudbuild.googleapis.com"
}

# Artifact Registry for Docker images
resource "google_artifact_registry_repository" "repo" {
  location      = var.gcp_region
  repository_id = "${var.service_name}-repo"
  format        = "DOCKER"
  description   = "Docker repository for ${var.service_name}"

  depends_on = [google_project_service.artifactregistry]
}

# Cloud SQL for PostgreSQL
resource "google_sql_database_instance" "db_instance" {
  name             = "${var.service_name}-db-instance"
  database_version = "POSTGRES_15"
  region           = var.gcp_region

  settings {
    tier = "db-g1-small" # Budget-friendly tier
    ip_configuration {
      ipv4_enabled    = true
      private_network = "projects/${var.gcp_project_id}/global/networks/default"
    }
  }

  deletion_protection = false # Set to true in production

  depends_on = [google_project_service.sqladmin]
}

resource "google_sql_database" "database" {
  instance = google_sql_database_instance.db_instance.name
  name     = var.db_name
}

resource "google_sql_user" "db_user" {
  instance = google_sql_database_instance.db_instance.name
  name     = var.db_user
  password = var.db_password
}

# Cloud Run service
resource "google_cloud_run_v2_service" "service" {
  name     = var.service_name
  location = var.gcp_region
  ingress  = "INGRESS_TRAFFIC_ALL"

  template {
    scaling {
      min_instance_count = 0 # Scale to zero for cost savings
      max_instance_count = 2
    }

    containers {
      image = "${google_artifact_registry_repository.repo.location}-docker.pkg.dev/${var.gcp_project_id}/${google_artifact_registry_repository.repo.repository_id}/${var.service_name}:latest"
      ports {
        container_port = 8080
      }

      env {
        name  = "SERVER_PORT"
        value = "8080"
      }
      env {
        name  = "DB_HOST"
        value = google_sql_database_instance.db_instance.private_ip_address
      }
      env {
        name  = "DB_PORT"
        value = "5432"
      }
      env {
        name  = "DB_USER"
        value = google_sql_user.db_user.name
      }
      env {
        name  = "DB_PASSWORD"
        value = google_sql_user.db_user.password
      }
      env {
        name  = "DB_NAME"
        value = google_sql_database.database.name
      }
      env {
        name  = "DB_SSL_MODE"
        value = "disable" # For private IP, SSL is not strictly needed but recommended for production
      }
      env {
        name  = "GIN_MODE"
        value = "release"
      }
    }

    vpc_access {
      connector = google_vpc_access_connector.connector.id
      egress    = "ALL_TRAFFIC"
    }
  }

  depends_on = [
    google_project_service.run,
    google_sql_database_instance.db_instance,
  ]
}

# Serverless VPC Access Connector
resource "google_vpc_access_connector" "connector" {
  name          = "${var.service_name}-vpc-connector"
  region        = var.gcp_region
  ip_cidr_range = "10.8.0.0/28"
  network       = "default"
}


# Allow public access to Cloud Run
resource "google_cloud_run_service_iam_member" "public_access" {
  location = google_cloud_run_v2_service.service.location
  project  = google_cloud_run_v2_service.service.project
  service  = google_cloud_run_v2_service.service.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}
