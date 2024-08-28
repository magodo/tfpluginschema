package fw

// Referencing: terraform-plugin-framework/internal/toproto6/block.go@6246b0c15d868c2de240bd0ce7a49c861056fe3d

import (
	"context"
	"sort"

	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	providerschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/magodo/tfpluginschema/internal/fw/fwschema"
	"github.com/magodo/tfpluginschema/schema"
)

func ProviderBlock(ctx context.Context, name string, path *tftypes.AttributePath, b providerschema.Block) (*schema.SchemaNestedBlock, error) {
	schemaNestedBlock := &schema.SchemaNestedBlock{
		Block:    &schema.SchemaBlock{},
		TypeName: name,
	}

	nm := b.GetNestingMode()
	switch fwschema.BlockNestingMode(nm) {
	case fwschema.BlockNestingModeList:
		schemaNestedBlock.Nesting = schema.SchemaNestedBlockNestingModeList
	case fwschema.BlockNestingModeSet:
		schemaNestedBlock.Nesting = schema.SchemaNestedBlockNestingModeSet
	case fwschema.BlockNestingModeSingle:
		schemaNestedBlock.Nesting = schema.SchemaNestedBlockNestingModeSingle
	default:
		return nil, path.NewErrorf("unrecognized nesting mode %v", nm)
	}

	nestedBlockObject := b.GetNestedObject()

	for attrName, attr := range nestedBlockObject.GetAttributes() {
		attrPath := path.WithAttributeName(attrName)
		attrProto6, err := ProviderSchemaAttribute(ctx, attrName, attrPath, attr)

		if err != nil {
			return nil, err
		}

		schemaNestedBlock.Block.Attributes = append(schemaNestedBlock.Block.Attributes, attrProto6)
	}

	for blockName, block := range nestedBlockObject.GetBlocks() {
		blockPath := path.WithAttributeName(blockName)
		blockProto6, err := ProviderBlock(ctx, blockName, blockPath, block)

		if err != nil {
			return nil, err
		}

		schemaNestedBlock.Block.BlockTypes = append(schemaNestedBlock.Block.BlockTypes, blockProto6)
	}

	sort.Slice(schemaNestedBlock.Block.Attributes, func(i, j int) bool {
		if schemaNestedBlock.Block.Attributes[i] == nil {
			return true
		}

		if schemaNestedBlock.Block.Attributes[j] == nil {
			return false
		}

		return schemaNestedBlock.Block.Attributes[i].Name < schemaNestedBlock.Block.Attributes[j].Name
	})

	sort.Slice(schemaNestedBlock.Block.BlockTypes, func(i, j int) bool {
		if schemaNestedBlock.Block.BlockTypes[i] == nil {
			return true
		}

		if schemaNestedBlock.Block.BlockTypes[j] == nil {
			return false
		}

		return schemaNestedBlock.Block.BlockTypes[i].TypeName < schemaNestedBlock.Block.BlockTypes[j].TypeName
	})

	return schemaNestedBlock, nil
}

func ResourceBlock(ctx context.Context, name string, path *tftypes.AttributePath, b resourceschema.Block) (*schema.SchemaNestedBlock, error) {
	schemaNestedBlock := &schema.SchemaNestedBlock{
		Block:    &schema.SchemaBlock{},
		TypeName: name,
	}

	nm := b.GetNestingMode()
	switch fwschema.BlockNestingMode(nm) {
	case fwschema.BlockNestingModeList:
		schemaNestedBlock.Nesting = schema.SchemaNestedBlockNestingModeList
	case fwschema.BlockNestingModeSet:
		schemaNestedBlock.Nesting = schema.SchemaNestedBlockNestingModeSet
	case fwschema.BlockNestingModeSingle:
		schemaNestedBlock.Nesting = schema.SchemaNestedBlockNestingModeSingle
	default:
		return nil, path.NewErrorf("unrecognized nesting mode %v", nm)
	}

	nestedBlockObject := b.GetNestedObject()

	for attrName, attr := range nestedBlockObject.GetAttributes() {
		attrPath := path.WithAttributeName(attrName)
		attrProto6, err := ResourceSchemaAttribute(ctx, attrName, attrPath, attr)

		if err != nil {
			return nil, err
		}

		schemaNestedBlock.Block.Attributes = append(schemaNestedBlock.Block.Attributes, attrProto6)
	}

	for blockName, block := range nestedBlockObject.GetBlocks() {
		blockPath := path.WithAttributeName(blockName)
		blockProto6, err := ResourceBlock(ctx, blockName, blockPath, block)

		if err != nil {
			return nil, err
		}

		schemaNestedBlock.Block.BlockTypes = append(schemaNestedBlock.Block.BlockTypes, blockProto6)
	}

	sort.Slice(schemaNestedBlock.Block.Attributes, func(i, j int) bool {
		if schemaNestedBlock.Block.Attributes[i] == nil {
			return true
		}

		if schemaNestedBlock.Block.Attributes[j] == nil {
			return false
		}

		return schemaNestedBlock.Block.Attributes[i].Name < schemaNestedBlock.Block.Attributes[j].Name
	})

	sort.Slice(schemaNestedBlock.Block.BlockTypes, func(i, j int) bool {
		if schemaNestedBlock.Block.BlockTypes[i] == nil {
			return true
		}

		if schemaNestedBlock.Block.BlockTypes[j] == nil {
			return false
		}

		return schemaNestedBlock.Block.BlockTypes[i].TypeName < schemaNestedBlock.Block.BlockTypes[j].TypeName
	})

	return schemaNestedBlock, nil
}

func DatasourceBlock(ctx context.Context, name string, path *tftypes.AttributePath, b datasourceschema.Block) (*schema.SchemaNestedBlock, error) {
	schemaNestedBlock := &schema.SchemaNestedBlock{
		Block:    &schema.SchemaBlock{},
		TypeName: name,
	}

	nm := b.GetNestingMode()
	switch fwschema.BlockNestingMode(nm) {
	case fwschema.BlockNestingModeList:
		schemaNestedBlock.Nesting = schema.SchemaNestedBlockNestingModeList
	case fwschema.BlockNestingModeSet:
		schemaNestedBlock.Nesting = schema.SchemaNestedBlockNestingModeSet
	case fwschema.BlockNestingModeSingle:
		schemaNestedBlock.Nesting = schema.SchemaNestedBlockNestingModeSingle
	default:
		return nil, path.NewErrorf("unrecognized nesting mode %v", nm)
	}

	nestedBlockObject := b.GetNestedObject()

	for attrName, attr := range nestedBlockObject.GetAttributes() {
		attrPath := path.WithAttributeName(attrName)
		attrProto6, err := DatasourceSchemaAttribute(ctx, attrName, attrPath, attr)

		if err != nil {
			return nil, err
		}

		schemaNestedBlock.Block.Attributes = append(schemaNestedBlock.Block.Attributes, attrProto6)
	}

	for blockName, block := range nestedBlockObject.GetBlocks() {
		blockPath := path.WithAttributeName(blockName)
		blockProto6, err := DatasourceBlock(ctx, blockName, blockPath, block)

		if err != nil {
			return nil, err
		}

		schemaNestedBlock.Block.BlockTypes = append(schemaNestedBlock.Block.BlockTypes, blockProto6)
	}

	sort.Slice(schemaNestedBlock.Block.Attributes, func(i, j int) bool {
		if schemaNestedBlock.Block.Attributes[i] == nil {
			return true
		}

		if schemaNestedBlock.Block.Attributes[j] == nil {
			return false
		}

		return schemaNestedBlock.Block.Attributes[i].Name < schemaNestedBlock.Block.Attributes[j].Name
	})

	sort.Slice(schemaNestedBlock.Block.BlockTypes, func(i, j int) bool {
		if schemaNestedBlock.Block.BlockTypes[i] == nil {
			return true
		}

		if schemaNestedBlock.Block.BlockTypes[j] == nil {
			return false
		}

		return schemaNestedBlock.Block.BlockTypes[i].TypeName < schemaNestedBlock.Block.BlockTypes[j].TypeName
	})

	return schemaNestedBlock, nil
}
