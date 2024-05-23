# For full specification on the configuration of this file visit:
# https://github.com/hashicorp/integration-template#metadata-configuration
integration {
  name = "External"
  description = "The External plugin is able to communicate with external commands."
  identifier = "packer/joomcode/external"
  component {
    type = "data-source"
    name = "External Raw"
    slug = "raw"
  }
  component {
    type = "data-source"
    name = "External"
    slug = "external"
  }
}
