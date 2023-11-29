The External plugin is able to communicate with external commands.

### Installation

To install this plugin, copy and paste this code into your Packer configuration, then run [`packer init`](https://www.packer.io/docs/commands/init).

```hcl
packer {
  required_plugins {
    veertu-anka = {
      version = "> 0.0.2"
      source  = "github.com/joomcode/external"
    }
  }
}
```

Alternatively, you can use `packer plugins install` to manage installation of this plugin.

```sh
$ packer plugins install github.com/joomcode/external
```

### Components

#### Data Sources
- [external](/packer/integrations/joomcode/external/latest/components/data-source/external) - Communicate with external commands
  using JSON protocol.
- [external-raw](/packer/integrations/joomcode/external/latest/components/data-source/raw) - Communicate with external commands
  using plaintext protocol.

