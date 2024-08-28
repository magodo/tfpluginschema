[![PkgGoDev](https://pkg.go.dev/badge/github.com/magodo/tfpluginschema)](https://pkg.go.dev/github.com/magodo/tfpluginschema)

# tfpluginschema

Terraform schema definition stands in the middle of the Terraform core schema and the Plugin SDK/FW schema.

## Why

The motivation for this is to add more information that is lost during the [conversion from plugin sdk to the terraform core schema](https://github.com/hashicorp/terraform-plugin-sdk/blob/6ffc92796f0716c07502e4d36aaafa5fd85e94cf/helper/schema/core_schema.go#L57). These information are fatal for developing tools that is oriented to the provider, rather than to the terraform core.

Specifically, we are:

1. Adding `Default` for the `Attribute`
1. Adding `Required`, `Optional`, `Computed` for the `BlockType` (SDK v2 only)
1. Adding `ExactlyOneOf`, `AtLeastOneOf`, `ConflictsWith` and `RequiredWith` for both `BlockType` and the `Attribute` (SDK v2 only)
1. Removing any other attributes
