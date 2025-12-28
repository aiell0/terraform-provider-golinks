terraform {
  required_providers {
    golinks = {
      source = "hashicorp.com/edu/golinks"
    }
  }
  required_version = ">= 1.8.0"
}

provider "golinks" {
  token = var.golinks_token
}

resource "golinks_link" "this" {
  name        = "tftest"
  url         = "https://google.com"
  description = "test golink"
  unlisted    = false
  public      = false
  private     = false
  format      = false
  hyphens     = false
  tags        = ["testing", "tag2"]
}

output "golinks" {
  value = golinks_link.this.name
}

data "golinks_links" "all" {}
data "golinks_link" "specific" {
  name = golinks_link.this.name
}
