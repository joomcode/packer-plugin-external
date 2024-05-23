Type: `external-raw`

The `external-raw` data source allows an external program to act as a data source,
exposing arbitrary data for use elsewhere in the Packer configuration.

~> **Warning** This mechanism is provided as an "escape hatch" for exceptional
situations where a first-class Packer provider is not more appropriate.
Its capabilities are limited in comparison to a true data source, and
implementing a data source via an external program is likely to hurt the
portability of your Packer configuration by creating dependencies on
external programs and libraries that may not be available (or may need to
be used differently) on different operating systems.

## Example Usage
```hcl
data "external-raw" "example" {
  program = ["rev"]
  query = "hello"
}

locals {
  result = data.external-raw.example.result
}

source "null" "example" {
  communicator = "none"
}

build {
  sources = ["source.null.example"]

  provisioner "shell-local" {
    inline = [
      "echo ${trimspace(local.result)} == olleh",
    ]
  }
}
```

## External Program Protocol

The protocol is similar to the one used by the
[external](/packer/integrations/joomcode/external/latest/components/data-source/external) data source.
However, query and result are plaintext strings instead of JSON objects.
Refer to the [external](/packer/integrations/joomcode/external/latest/components/data-source/external) doc
for more details.

`external-raw` should be used over `external` in cases where the external program
does not support JSON input and/or output.

## Argument Reference

### Required

- `command` ([]string) - A list of strings, whose first element is the program
  to run and whose subsequent elements are optional command line arguments
  to the program. Packer does not execute the program through a shell, so
  it is not necessary to escape shell metacharacters nor add quotes around
  arguments containing spaces.

### Optional

- `working_dir` (string) - Working directory of the program.
  If not supplied, the program will run in the current directory.
- `query` (string) - An input to pass to the external program.
  If not supplied, the program will not receive anything on standard input.

### Output

- `result` (string) - The output returned from the external program.
