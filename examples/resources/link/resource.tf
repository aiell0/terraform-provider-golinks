resource "golinks_link" "this" {
  name        = "tftest"
  url         = "https://google.com"
  description = "test golink"
  unlisted    = false
  public      = false
  private     = false
  format      = false
  hyphens     = false
  tags        = ["testing"]
}
