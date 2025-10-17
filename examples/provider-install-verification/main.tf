terraform {
  required_providers {
    golinks = {
      source = "hashicorp.com/edu/golinks"
    }
  }
}

provider "golinks" {
  token = "ff26c97d4e4409b6fc13b1780795e05815fdb3421e5a986001fcf746e2f58010"
}

data "golinks_links" "all" {}

output "golinks" {
  value = data.golinks_links.all
}
