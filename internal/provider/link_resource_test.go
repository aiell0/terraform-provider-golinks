package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestLinkResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "golinks_link" "test" {
	name = "testlink"
	url = "https://google.com"
	description = "This is a test link"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify link item
					resource.TestCheckResourceAttr("golinks_link.test", "description", "This is a test link"),
					resource.TestCheckResourceAttr("golinks_link.test", "name", "testlink"),
					resource.TestCheckResourceAttr("golinks_link.test", "url", "https://google.com"),
					resource.TestCheckResourceAttr("golinks_link.test", "variable_link", "0"),
					resource.TestCheckResourceAttr("golinks_link.test", "pinned", "0"),
					resource.TestCheckResourceAttr("golinks_link.test", "private", "0"),
					resource.TestCheckResourceAttr("golinks_link.test", "public", "0"),
					resource.TestCheckNoResourceAttr("golinks_link.test", "aliases"),
					resource.TestCheckNoResourceAttr("golinks_link.test", "geolinks"),
					resource.TestCheckNoResourceAttr("golinks_link.test", "multilinks"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "golinks_link.test",
				ImportState:       true,
				ImportStateVerify: true,
				// The last_updated attribute does not exist in the Golinks
				// API, therefore there is no value for it during import.
				ImportStateVerifyIgnore: []string{"last_updated", "private", "public"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "golinks_link" "test" {
	name = "testlink2"
	url = "https://google.com"
	description = "This is a new test link"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify first order item updated
					resource.TestCheckResourceAttr("golinks_link.test", "description", "This is a new test link"),
					resource.TestCheckResourceAttr("golinks_link.test", "name", "testlink2"),
					resource.TestCheckResourceAttr("golinks_link.test", "url", "https://google.com"),
					resource.TestCheckResourceAttr("golinks_link.test", "variable_link", "0"),
					resource.TestCheckResourceAttr("golinks_link.test", "pinned", "0"),
					resource.TestCheckResourceAttr("golinks_link.test", "private", "0"),
					resource.TestCheckResourceAttr("golinks_link.test", "public", "0"),
					resource.TestCheckNoResourceAttr("golinks_link.test", "aliases"),
					resource.TestCheckNoResourceAttr("golinks_link.test", "geolinks"),
					resource.TestCheckNoResourceAttr("golinks_link.test", "multilinks"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
