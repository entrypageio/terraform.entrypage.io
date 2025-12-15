package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/entrypage/terraform-provider-entrypage/internal/provider"
)

// Replace with proper versioning when releasing.
const version = "0.1.0"

func main() {
	err := providerserver.Serve(context.Background(), provider.New(version), providerserver.ServeOpts{
		Address: "registry.terraform.io/entrypage_io/entrypage_io",
	})
	if err != nil {
		log.Fatal(err)
	}
}
