variable "entrypage_api_key" {
  description = "API key for entrypage.io (or set ENTRYPAGE_API_KEY)."
  type        = string
  sensitive   = true
}

variable "domain" {
  description = "Entrypage domain (e.g. acme.entrypage.io)."
  type        = string
}

variable "preferred_username" {
  description = "Preferred username for the user."
  type        = string
}

variable "name" {
  description = "Display name for the user."
  type        = string
}
