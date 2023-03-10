data "external" "test" {
  program = ["jq", "-r", "{\"my_key1\": \"my_\\(.key1)\"}"]
  query = {
    key1 = "val1"
  }
}

locals {
  result = data.external.test.result
}

source "null" "test" {
  communicator = "none"
}

build {
  sources = ["source.null.test"]

  provisioner "shell-local" {
    inline = [
      "echo 'result[my_key1]: ${local.result["my_key1"]}'",
    ]
  }
}
