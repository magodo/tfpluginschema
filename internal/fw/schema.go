package fw

// Referencing: terraform-plugin-framework/internal/toproto6/schema.go@6246b0c15d868c2de240bd0ce7a49c861056fe3d

import (
	"context"
	"sort"

	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	providerschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/magodo/tfpluginschema/schema"
)

func ProviderSchema(ctx context.Context, s providerschema.Schema) (*schema.Schema, error) {
	result := &schema.Schema{
		Version: s.GetVersion(),
	}

	var attrs []*schema.SchemaAttribute
	var blocks []*schema.SchemaNestedBlock

	for name, attr := range s.GetAttributes() {
		a, err := ProviderSchemaAttribute(ctx, name, tftypes.NewAttributePath().WithAttributeName(name), attr)

		if err != nil {
			return nil, err
		}

		attrs = append(attrs, a)
	}

	for name, block := range s.GetBlocks() {
		proto6, err := ProviderBlock(ctx, name, tftypes.NewAttributePath().WithAttributeName(name), block)

		if err != nil {
			return nil, err
		}

		blocks = append(blocks, proto6)
	}

	sort.Slice(attrs, func(i, j int) bool {
		if attrs[i] == nil {
			return true
		}

		if attrs[j] == nil {
			return false
		}

		return attrs[i].Name < attrs[j].Name
	})

	sort.Slice(blocks, func(i, j int) bool {
		if blocks[i] == nil {
			return true
		}

		if blocks[j] == nil {
			return false
		}

		return blocks[i].TypeName < blocks[j].TypeName
	})

	result.Block = &schema.SchemaBlock{
		// core doesn't do anything with version, as far as I can tell,
		// so let's not set it.
		Attributes: attrs,
		BlockTypes: blocks,
	}

	return result, nil
}

func ResourceSchema(ctx context.Context, s resourceschema.Schema) (*schema.Schema, error) {
	result := &schema.Schema{
		Version: s.GetVersion(),
	}

	var attrs []*schema.SchemaAttribute
	var blocks []*schema.SchemaNestedBlock

	for name, attr := range s.GetAttributes() {
		a, err := ResourceSchemaAttribute(ctx, name, tftypes.NewAttributePath().WithAttributeName(name), attr)

		if err != nil {
			return nil, err
		}

		attrs = append(attrs, a)
	}

	for name, block := range s.GetBlocks() {
		proto6, err := ResourceBlock(ctx, name, tftypes.NewAttributePath().WithAttributeName(name), block)

		if err != nil {
			return nil, err
		}

		blocks = append(blocks, proto6)
	}

	sort.Slice(attrs, func(i, j int) bool {
		if attrs[i] == nil {
			return true
		}

		if attrs[j] == nil {
			return false
		}

		return attrs[i].Name < attrs[j].Name
	})

	sort.Slice(blocks, func(i, j int) bool {
		if blocks[i] == nil {
			return true
		}

		if blocks[j] == nil {
			return false
		}

		return blocks[i].TypeName < blocks[j].TypeName
	})

	result.Block = &schema.SchemaBlock{
		// core doesn't do anything with version, as far as I can tell,
		// so let's not set it.
		Attributes: attrs,
		BlockTypes: blocks,
	}

	return result, nil
}

func DatasourceSchema(ctx context.Context, s datasourceschema.Schema) (*schema.Schema, error) {
	result := &schema.Schema{
		Version: s.GetVersion(),
	}

	var attrs []*schema.SchemaAttribute
	var blocks []*schema.SchemaNestedBlock

	for name, attr := range s.GetAttributes() {
		a, err := DatasourceSchemaAttribute(ctx, name, tftypes.NewAttributePath().WithAttributeName(name), attr)

		if err != nil {
			return nil, err
		}

		attrs = append(attrs, a)
	}

	for name, block := range s.GetBlocks() {
		proto6, err := DatasourceBlock(ctx, name, tftypes.NewAttributePath().WithAttributeName(name), block)

		if err != nil {
			return nil, err
		}

		blocks = append(blocks, proto6)
	}

	sort.Slice(attrs, func(i, j int) bool {
		if attrs[i] == nil {
			return true
		}

		if attrs[j] == nil {
			return false
		}

		return attrs[i].Name < attrs[j].Name
	})

	sort.Slice(blocks, func(i, j int) bool {
		if blocks[i] == nil {
			return true
		}

		if blocks[j] == nil {
			return false
		}

		return blocks[i].TypeName < blocks[j].TypeName
	})

	result.Block = &schema.SchemaBlock{
		// core doesn't do anything with version, as far as I can tell,
		// so let's not set it.
		Attributes: attrs,
		BlockTypes: blocks,
	}

	return result, nil
}
