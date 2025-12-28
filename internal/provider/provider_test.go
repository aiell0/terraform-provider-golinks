// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func testAccProviderConfig(t *testing.T) string {
	t.Helper()

	token := os.Getenv("GOLINKS_TOKEN")
	if token == "" {
		t.Skip("set GOLINKS_TOKEN to run acceptance tests")
	}

	return fmt.Sprintf(`
provider "golinks" {
  token = %q
}
`, token)
}

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"golinks": providerserver.NewProtocol6WithError(New("test")()),
	}
)
