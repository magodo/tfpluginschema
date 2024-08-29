package sdkv2

// A modified version based on: github.com/hashicorp/terraform-plugin-sdk/v2/internal/configs/configschema/implied_type.go

import (
	"github.com/magodo/tfpluginschema/schema"
	"github.com/zclconf/go-cty/cty"
)

func ImpliedType(b *schema.SchemaBlock) cty.Type {
	if b == nil {
		return cty.EmptyObject
	}

	atys := make(map[string]cty.Type)

	for _, attrS := range b.Attributes {
		atys[attrS.Name] = *attrS.Type
	}

	for _, blockS := range b.BlockTypes {
		if _, exists := atys[blockS.TypeName]; exists {
			panic("invalid schema, blocks and attributes cannot have the same name")
		}

		childType := ImpliedType(blockS.Block)

		switch blockS.Nesting {
		case schema.SchemaNestedBlockNestingModeSingle, schema.SchemaNestedBlockNestingModeGroup:
			atys[blockS.TypeName] = childType
		case schema.SchemaNestedBlockNestingModeList:
			if childType.HasDynamicTypes() {
				atys[blockS.TypeName] = cty.DynamicPseudoType
			} else {
				atys[blockS.TypeName] = cty.List(childType)
			}
		case schema.SchemaNestedBlockNestingModeSet:
			if childType.HasDynamicTypes() {
				panic("can't use cty.DynamicPseudoType inside a block type with NestingSet")
			}
			atys[blockS.TypeName] = cty.Set(childType)
		case schema.SchemaNestedBlockNestingModeMap:
			if childType.HasDynamicTypes() {
				atys[blockS.TypeName] = cty.DynamicPseudoType
			} else {
				atys[blockS.TypeName] = cty.Map(childType)
			}
		default:
			panic("invalid nesting type")
		}
	}

	return cty.Object(atys)
}
