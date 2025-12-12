package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/entrypage/terraform-provider-entrypage/internal/provider"
)

// Replace with proper versioning when releasing.
const version = "0.1.1"

func main() {
	err := providerserver.Serve(context.Background(), provider.New(version), providerserver.ServeOpts{
		Address: "registry.terraform.io/entrypage/entrypage",
	})
	if err != nil {
		log.Fatal(err)
	}
}
