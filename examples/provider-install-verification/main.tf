terraform {
  required_providers {
    golinks = {
      source = "hashicorp.com/edu/golinks"
    }
  }
  required_version = ">= 1.8.0"
}

provider "golinks" {
  token = "ff26c97d4e4409b6fc13b1780795e05815fdb3421e5a986001fcf746e2f58010"
}

# data "golinks_links" "all" {}

resource "golinks_link" "this" {
  name        = "tftest"
  url         = "https://google.com"
  description = "This is a test link"
}

# output "golinks" {
#   value = data.golinks_links.all
# }
