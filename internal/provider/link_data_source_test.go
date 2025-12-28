package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestLinkDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderConfig(t) + `
resource "golinks_link" "test" {
  name        = "testlink-datasource"
  url         = "https://golinks.io"
  description = "Link fetched via data source"
}

data "golinks_link" "test" {
  name = golinks_link.test.name
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.golinks_link.test", "name", "testlink-datasource"),
					resource.TestCheckResourceAttr("data.golinks_link.test", "url", "https://golinks.io"),
					resource.TestCheckResourceAttr("data.golinks_link.test", "description", "Link fetched via data source"),
					resource.TestCheckResourceAttrSet("data.golinks_link.test", "gid"),
					resource.TestCheckResourceAttrSet("data.golinks_link.test", "cid"),
					resource.TestCheckResourceAttrSet("data.golinks_link.test", "redirect_hits.daily"),
				),
			},
		},
	})
}
