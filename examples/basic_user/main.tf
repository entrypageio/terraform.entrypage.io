terraform {
  required_providers {
    entrypage = {
      source  = "entrypage/entrypage"
      version = "0.1.1"
    }
  }
}

provider "entrypage" {
  api_key = var.entrypage_api_key
}

resource "entrypage_user" "example" {
  domain             = var.domain
  preferred_username = var.preferred_username
  name               = var.name
}

output "user_id" {
  value       = entrypage_user.example.id
  description = "User subject identifier (sub)"
}
