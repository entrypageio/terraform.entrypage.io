# Terraform Provider for entrypage.io

A minimal Terraform provider to create entrypage.io users via the public API.

## Usage

```hcl
terraform {
  required_providers {
    entrypage_io = {
      source  = "entrypage_io/entrypage_io"
      version = "0.1.1"
    }
  }
}

provider "entrypage_io" {
  # or set ENTRYPAGE_API_KEY in your environment
  api_key = var.entrypage_api_key
}

resource "entrypage_io_user" "example" {
  domain              = "acme.entrypage.io"
  preferred_username  = "jane.doe"
  name                = "Jane Doe"
}
```

Import (format `domain/sub`):

```sh
terraform import entrypage_io_user.example acme.entrypage.io/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee
```

## Development

Build and test locally:

```sh
go test ./...
```

The provider uses the `x-api-key` header (configure via `api_key` or `ENTRYPAGE_API_KEY`). `base_url` can be overridden for tests but defaults to `https://api.entrypage.io`.

Limitations: the API has no direct `GET /user/{sub}` endpoint. The `Read` operation therefore searches via `GET /v1/domain/{domain}/user?search=<sub>` and removes state when the user is not found.

See `DEV_NOTES.md` for a step-by-step local testing guide (building the binary, local plugin paths, and example usage).
