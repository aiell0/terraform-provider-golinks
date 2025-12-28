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
				Config: testAccProviderConfig(t) + `
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
					resource.TestCheckResourceAttr("golinks_link.test", "variable_link", "false"),
					resource.TestCheckResourceAttr("golinks_link.test", "pinned", "false"),
					resource.TestCheckResourceAttr("golinks_link.test", "private", "false"),
					resource.TestCheckResourceAttr("golinks_link.test", "public", "false"),
					resource.TestCheckResourceAttr("golinks_link.test", "unlisted", "false"),
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
				Config: testAccProviderConfig(t) + `
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
					resource.TestCheckResourceAttr("golinks_link.test", "variable_link", "false"),
					resource.TestCheckResourceAttr("golinks_link.test", "pinned", "false"),
					resource.TestCheckResourceAttr("golinks_link.test", "private", "false"),
					resource.TestCheckResourceAttr("golinks_link.test", "public", "false"),
					resource.TestCheckResourceAttr("golinks_link.test", "unlisted", "false"),
					resource.TestCheckNoResourceAttr("golinks_link.test", "aliases"),
					resource.TestCheckNoResourceAttr("golinks_link.test", "geolinks"),
					resource.TestCheckNoResourceAttr("golinks_link.test", "multilinks"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestLinkResourceOptionalAttributes(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderConfig(t) + `
resource "golinks_link" "options" {
	name        = "testlink-options"
	url         = "https://google.com"
	description = "Link with optional settings"
	unlisted    = true
	public      = true
	private     = false
	tags        = ["testing", "tag2"]
}

resource "golinks_link" "private_only" {
	name        = "testlink-private"
	url         = "https://google.com"
	description = "Link with private access"
	unlisted    = true
	public      = false
	private     = true
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("golinks_link.options", "unlisted", "true"),
					resource.TestCheckResourceAttr("golinks_link.options", "public", "true"),
					resource.TestCheckResourceAttr("golinks_link.options", "private", "false"),
					resource.TestCheckResourceAttr("golinks_link.options", "format", "false"),
					resource.TestCheckResourceAttr("golinks_link.options", "hyphens", "false"),

					resource.TestCheckResourceAttr("golinks_link.options", "tags.#", "2"),
					resource.TestCheckResourceAttr("golinks_link.options", "tags.0", "testing"),
					resource.TestCheckResourceAttr("golinks_link.options", "tags.1", "tag2"),
					resource.TestCheckResourceAttr("golinks_link.private_only", "private", "true"),
					resource.TestCheckResourceAttr("golinks_link.private_only", "public", "false"),
					resource.TestCheckResourceAttr("golinks_link.private_only", "unlisted", "true"),
				),
			},
		},
	})
}
