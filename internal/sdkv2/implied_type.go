package sdkv2

// A modified version based on: github.com/hashicorp/terraform-plugin-sdk/v2/internal/configs/configschema/implied_type.go

import (
	"github.com/magodo/tfpluginschema/schema"
	"github.com/zclconf/go-cty/cty"
)

func ImpliedType(b *schema.Block) cty.Type {
	if b == nil {
		return cty.EmptyObject
	}

	atys := make(map[string]cty.Type)

	for name, attrS := range b.Attributes {
		atys[name] = attrS.Type
	}

	for name, blockS := range b.NestedBlocks {
		if _, exists := atys[name]; exists {
			panic("invalid schema, blocks and attributes cannot have the same name")
		}

		childType := ImpliedType(blockS.Block)

		switch blockS.NestingMode {
		case schema.NestingSingle, schema.NestingGroup:
			atys[name] = childType
		case schema.NestingList:
			if childType.HasDynamicTypes() {
				atys[name] = cty.DynamicPseudoType
			} else {
				atys[name] = cty.List(childType)
			}
		case schema.NestingSet:
			if childType.HasDynamicTypes() {
				panic("can't use cty.DynamicPseudoType inside a block type with NestingSet")
			}
			atys[name] = cty.Set(childType)
		case schema.NestingMap:
			if childType.HasDynamicTypes() {
				atys[name] = cty.DynamicPseudoType
			} else {
				atys[name] = cty.Map(childType)
			}
		default:
			panic("invalid nesting type")
		}
	}

	return cty.Object(atys)
}
