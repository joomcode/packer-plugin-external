data "external-raw" "test" {
  program = ["rev"]
  query   = "foobar"
}

locals {
  result = data.external-raw.test.result
}

source "null" "test" {
  communicator = "none"
}

build {
  sources = ["source.null.test"]

  provisioner "shell-local" {
    inline = [
      "echo 'result: ${local.result}'"
    ]
  }
}
