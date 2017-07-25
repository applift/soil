id = "localhost-1"

resource "port.default" {
  allocator = "range"
  config {
    "minor" = "10000"
    "major" = "10100"
  }
}

resource "port.overlay" {
  allocator = "range"
  config {
    "minor" = "10000"
    "major" = "10100"
  }
}

meta {
  "consul" = "true"
  "consul-client" = "true"
  "field" = "all,consul"
}

exec = "ExecStart=/usr/bin/sleep inf"


// private pod
pod "one-1" {
  runtime = true
  count = 1
  constraint {
    "${meta.consul}" = "true"
  }
  unit "one-1-0.service" {
    create = "start"
    source = <<EOF
    [Unit]
    Description=%p

    [Service]
    ExecStart=/usr/bin/sleep inf

    [Install]
    WantedBy=default.target
  EOF
  }
}

