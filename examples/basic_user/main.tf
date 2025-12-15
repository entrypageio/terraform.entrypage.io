terraform {
  required_providers {
    entrypage_io = {
      source  = "entrypage_io/entrypage_io"
      version = "0.1.0"
    }
  }
}

provider "entrypage_io" {
  api_key = var.entrypage_api_key
}

resource "entrypage_io_user" "example" {
  domain             = var.domain
  preferred_username = var.preferred_username
  name               = var.name
}

output "user_id" {
  value       = entrypage_io_user.example.id
  description = "User subject identifier (sub)"
}
