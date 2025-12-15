# Dev notes: local testing the provider

This is a quick checklist to run the provider locally (after the troubleshooting we did).

## Build the provider
```sh
go build -o /tmp/terraform-provider-entrypage_io_v0.1.1
```

## Make Terraform find the provider
Use the local plugin directory layout (macOS arm64 path shown; adjust for your OS/arch):
```sh
mkdir -p /tmp/registry.terraform.io/entrypage_io/entrypage_io/0.1.1/darwin_arm64
cp /tmp/terraform-provider-entrypage_io_v0.1.1 /tmp/registry.terraform.io/entrypage_io/entrypage_io/0.1.1/darwin_arm64/terraform-provider-entrypage_io_v0.1.1
chmod +x /tmp/registry.terraform.io/entrypage_io/entrypage_io/0.1.1/darwin_arm64/terraform-provider-entrypage_io_v0.1.1
```

Optional: same for OpenTofu naming (if needed):
```sh
mkdir -p /tmp/registry.opentofu.org/entrypage_io/entrypage_io/0.1.1/darwin_arm64
cp /tmp/terraform-provider-entrypage_io_v0.1.1 /tmp/registry.opentofu.org/entrypage_io/entrypage_io/0.1.1/darwin_arm64/terraform-provider-entrypage_io_v0.1.1
chmod +x /tmp/registry.opentofu.org/entrypage_io/entrypage_io/0.1.1/darwin_arm64/terraform-provider-entrypage_io_v0.1.1
```

## Init and apply example
From `examples/basic_user`:
```sh
rm -rf .terraform .terraform.lock.hcl terraform.tfstate terraform.tfstate.backup
terraform init -plugin-dir=/tmp
terraform apply -auto-approve
```

Fill `terraform.tfvars` (or env) with:
```hcl
entrypage_api_key  = "<your-api-key>"
domain             = "your.domain.entrypage.io"
preferred_username = "jane.doe"
name               = "Jane Doe"
```

## Notes
- The provider accepts both 200 and 201 responses on create; 409 conflicts will try to adopt an existing user by preferred_username.
- We keep `version = "0.1.1"` aligned between the build output name and `required_providers` in `examples/basic_user/main.tf`.
- Delete test users in the sandbox via `DELETE /v1/domain/{domain}/user/{sub}` if you need a clean slate.
