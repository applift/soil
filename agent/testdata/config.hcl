id = "agent-1"

resource "port.default" {
  allocator = "range"
  config "default" {
    "minor" = "10000"
    "major" = "10100"
  }
}


meta {
  "consul" = "true"
  "consul-client" = "true"
  "field" = "all,consul"
}
