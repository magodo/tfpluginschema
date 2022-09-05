package tfpluginschema

import (
	sdkschema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/magodo/tfpluginschema/internal/sdkv2"
	"github.com/magodo/tfpluginschema/schema"
)

func FromSDKv2ProviderSchemaMap(providerschemas map[string]*sdkschema.Schema) *schema.Block {
	return sdkv2.FromProviderSchemaMap(providerschemas)
}
