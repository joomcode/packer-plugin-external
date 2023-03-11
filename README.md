# Packer external data source plugin

A plugin for [Packer](https://www.packer.io/) which provides access to external commands. Compatible with Packer >= 1.7.0

## Usage
```hcl
packer {
  required_plugins {
    external = {
      version = ">= 0.0.2"
      source  = "github.com/joomcode/external"
    }
  }
}

data "external" "example" {
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
  external_result = data.external.example.result["foo"] # "val1"
  raw_result      = data.external-raw.example.result # "olleh\n"
}
```

See docs for more detailed information.

## Running Acceptance Tests

```bash
make testacc
```

This will run the acceptance tests for all plugins in this set.

# Requirements

-	[Go](https://go.dev/doc/install) >= 1.17
