packer {
  required_plugins {
    external = {
      version = ">=v0.0.1"
      source  = "github.com/joomcode/external"
    }
  }
}

data "external-json" "example" {
  program = ["jq", "{ \"foo\": .key1 }"]
  query = {
    key1 = "val1"
  }
}

data "external-raw" "example" {
  program = ["rev"]
  query   = "hello"
}

locals {
  json_result = data.external-json.example.result
  raw_result  = data.external-raw.example.result
}

source "null" "example" {
  communicator = "none"
}

build {
  sources = ["source.null.example"]

  provisioner "shell-local" {
    inline = [
      "echo ${local.json_result["foo"]} == val1",
      "echo ${trimspace(local.raw_result)} == olleh",
    ]
  }
}
