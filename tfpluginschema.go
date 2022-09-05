package tfpluginschema

import (
	sdkschema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/magodo/tfpluginschema/internal/sdkv2"
	"github.com/magodo/tfpluginschema/schema"
)

// FromSDKv2SchemasMap converts the schema map from the schema defined in the plugin sdk v2 to the schema defined in tfpluginschema.
func FromSDKv2SchemaMap(m map[string]*sdkschema.Schema) *schema.Block {
	return sdkv2.FromSchemaMap(m)
}
