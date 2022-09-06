package tfpluginschema

import (
	sdkschema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/magodo/tfpluginschema/internal/sdkv2"
	"github.com/magodo/tfpluginschema/schema"
)

// FromSDKv2Provider converts the provider from the schema defined in the plugin sdk v2 to the schema defined in tfpluginschema.
func FromSDKv2Provider(p *sdkschema.Provider) *schema.ProviderSchema {
	return sdkv2.FromProvider(p)
}

// FromSDKv2Resource converts the resource from the schema defined in the plugin sdk v2 to the schema defined in tfpluginschema.
func FromSDKv2Resource(res *sdkschema.Resource) *schema.Resource {
	return sdkv2.FromResource(res)
}

// FromSDKv2SchemasMap converts the schema map from the schema defined in the plugin sdk v2 to the schema defined in tfpluginschema.
func FromSDKv2SchemaMap(m map[string]*sdkschema.Schema) *schema.Block {
	return sdkv2.FromSchemaMap(m)
}
