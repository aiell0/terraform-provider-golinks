package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestLinksDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccProviderConfig(t) + `data "golinks_links" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of links returned
					resource.TestCheckResourceAttr("data.golinks_links.test", "results.#", "50"),
				),
			},
		},
	})
}
