pod "first" {
  constraint {
    "${meta.first}" = "true"
  }
}

pod "second" {
  constraint {
    "${meta.second}" = "true"
  }
  resource "port" "8080" {

  }
}