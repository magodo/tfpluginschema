package fw

// Referencing: terraform-plugin-framework/internal/toproto6/schema_attribute.go@6246b0c15d868c2de240bd0ce7a49c861056fe3d

import (
	"context"
	"fmt"
	"sort"

	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	providerschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/magodo/tfpluginschema/internal/fw/fwschema"
	"github.com/magodo/tfpluginschema/schema"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

func ProviderSchemaAttribute(ctx context.Context, name string, path *tftypes.AttributePath, a providerschema.Attribute) (*schema.SchemaAttribute, error) {
	if !a.IsRequired() && !a.IsOptional() && !a.IsComputed() {
		return nil, path.NewErrorf("must have Required, Optional, or Computed set")
	}

	schemaAttribute := &schema.SchemaAttribute{
		Name:      name,
		Required:  a.IsRequired(),
		Optional:  a.IsOptional(),
		Computed:  a.IsComputed(),
		Sensitive: a.IsSensitive(),
	}
	tfType := a.GetType().TerraformType(ctx)
	b, err := tfType.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("marshalling tftype: %v", err)
	}
	typ, err := ctyjson.UnmarshalType(b)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling to cty type: %v", err)
	}
	schemaAttribute.Type = &typ

	nestedAttribute, ok := a.(providerschema.NestedAttribute)

	if !ok {
		return schemaAttribute, nil
	}

	object := &schema.SchemaObject{}
	nm := nestedAttribute.GetNestingMode()
	switch fwschema.NestingMode(nm) {
	case fwschema.NestingModeSingle:
		object.Nesting = schema.SchemaObjectNestingModeSingle
	case fwschema.NestingModeList:
		object.Nesting = schema.SchemaObjectNestingModeList
	case fwschema.NestingModeSet:
		object.Nesting = schema.SchemaObjectNestingModeSet
	case fwschema.NestingModeMap:
		object.Nesting = schema.SchemaObjectNestingModeMap
	default:
		return nil, path.NewErrorf("unrecognized nesting mode %v", nm)
	}

	for nestedName, nestedA := range nestedAttribute.GetNestedObject().GetAttributes() {
		nestedSchemaAttribute, err := ProviderSchemaAttribute(ctx, nestedName, path.WithAttributeName(nestedName), nestedA)

		if err != nil {
			return nil, err
		}

		object.Attributes = append(object.Attributes, nestedSchemaAttribute)
	}

	sort.Slice(object.Attributes, func(i, j int) bool {
		if object.Attributes[i] == nil {
			return true
		}

		if object.Attributes[j] == nil {
			return false
		}

		return object.Attributes[i].Name < object.Attributes[j].Name
	})

	schemaAttribute.NestedType = object
	schemaAttribute.Type = nil

	return schemaAttribute, nil
}

func ResourceSchemaAttribute(ctx context.Context, name string, path *tftypes.AttributePath, a resourceschema.Attribute) (*schema.SchemaAttribute, error) {
	if !a.IsRequired() && !a.IsOptional() && !a.IsComputed() {
		return nil, path.NewErrorf("must have Required, Optional, or Computed set")
	}

	schemaAttribute := &schema.SchemaAttribute{
		Name:      name,
		Required:  a.IsRequired(),
		Optional:  a.IsOptional(),
		Computed:  a.IsComputed(),
		Sensitive: a.IsSensitive(),
	}

	switch a := a.(type) {
	case resourceschema.BoolAttribute:
		if a.Default != nil {
			var resp defaults.BoolResponse
			a.Default.DefaultBool(ctx, defaults.BoolRequest{}, &resp)
			schemaAttribute.Default = resp.PlanValue
		}
	case resourceschema.Float32Attribute:
		if a.Default != nil {
			var resp defaults.Float32Response
			a.Default.DefaultFloat32(ctx, defaults.Float32Request{}, &resp)
			schemaAttribute.Default = resp.PlanValue
		}
	case resourceschema.Float64Attribute:
		if a.Default != nil {
			var resp defaults.Float64Response
			a.Default.DefaultFloat64(ctx, defaults.Float64Request{}, &resp)
			schemaAttribute.Default = resp.PlanValue
		}
	case resourceschema.Int32Attribute:
		if a.Default != nil {
			var resp defaults.Int32Response
			a.Default.DefaultInt32(ctx, defaults.Int32Request{}, &resp)
			schemaAttribute.Default = resp.PlanValue
		}
	case resourceschema.Int64Attribute:
		if a.Default != nil {
			var resp defaults.Int64Response
			a.Default.DefaultInt64(ctx, defaults.Int64Request{}, &resp)
			schemaAttribute.Default = resp.PlanValue
		}
	case resourceschema.NumberAttribute:
		if a.Default != nil {
			var resp defaults.NumberResponse
			a.Default.DefaultNumber(ctx, defaults.NumberRequest{}, &resp)
			schemaAttribute.Default = resp.PlanValue
		}
	case resourceschema.StringAttribute:
		if a.Default != nil {
			var resp defaults.StringResponse
			a.Default.DefaultString(ctx, defaults.StringRequest{}, &resp)
			schemaAttribute.Default = resp.PlanValue
		}
	case resourceschema.ListAttribute:
		if a.Default != nil {
			var resp defaults.ListResponse
			a.Default.DefaultList(ctx, defaults.ListRequest{}, &resp)
			schemaAttribute.Default = resp.PlanValue
		}
	case resourceschema.MapAttribute:
		if a.Default != nil {
			var resp defaults.MapResponse
			a.Default.DefaultMap(ctx, defaults.MapRequest{}, &resp)
			schemaAttribute.Default = resp.PlanValue
		}
	case resourceschema.SetAttribute:
		if a.Default != nil {
			var resp defaults.SetResponse
			a.Default.DefaultSet(ctx, defaults.SetRequest{}, &resp)
			schemaAttribute.Default = resp.PlanValue
		}
	case resourceschema.SingleNestedAttribute:
		if a.Default != nil {
			var resp defaults.ObjectResponse
			a.Default.DefaultObject(ctx, defaults.ObjectRequest{}, &resp)
			schemaAttribute.Default = resp.PlanValue
		}
	case resourceschema.ListNestedAttribute:
		if a.Default != nil {
			var resp defaults.ListResponse
			a.Default.DefaultList(ctx, defaults.ListRequest{}, &resp)
			schemaAttribute.Default = resp.PlanValue
		}
	case resourceschema.MapNestedAttribute:
		if a.Default != nil {
			var resp defaults.MapResponse
			a.Default.DefaultMap(ctx, defaults.MapRequest{}, &resp)
			schemaAttribute.Default = resp.PlanValue
		}
	case resourceschema.SetNestedAttribute:
		if a.Default != nil {
			var resp defaults.SetResponse
			a.Default.DefaultSet(ctx, defaults.SetRequest{}, &resp)
			schemaAttribute.Default = resp.PlanValue
		}
	case resourceschema.ObjectAttribute:
		if a.Default != nil {
			var resp defaults.ObjectResponse
			a.Default.DefaultObject(ctx, defaults.ObjectRequest{}, &resp)
			schemaAttribute.Default = resp.PlanValue
		}
	case resourceschema.DynamicAttribute:
		if a.Default != nil {
			var resp defaults.DynamicResponse
			a.Default.DefaultDynamic(ctx, defaults.DynamicRequest{}, &resp)
			schemaAttribute.Default = resp.PlanValue
		}
	default:
		return nil, path.NewErrorf("unhandled type for default value: %T", a)
	}

	tfType := a.GetType().TerraformType(ctx)
	b, err := tfType.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("marshalling tftype: %v", err)
	}
	typ, err := ctyjson.UnmarshalType(b)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling to cty type: %v", err)
	}
	schemaAttribute.Type = &typ

	nestedAttribute, ok := a.(resourceschema.NestedAttribute)

	if !ok {
		return schemaAttribute, nil
	}

	object := &schema.SchemaObject{}
	nm := nestedAttribute.GetNestingMode()
	switch fwschema.NestingMode(nm) {
	case fwschema.NestingModeSingle:
		object.Nesting = schema.SchemaObjectNestingModeSingle
	case fwschema.NestingModeList:
		object.Nesting = schema.SchemaObjectNestingModeList
	case fwschema.NestingModeSet:
		object.Nesting = schema.SchemaObjectNestingModeSet
	case fwschema.NestingModeMap:
		object.Nesting = schema.SchemaObjectNestingModeMap
	default:
		return nil, path.NewErrorf("unrecognized nesting mode %v", nm)
	}

	for nestedName, nestedA := range nestedAttribute.GetNestedObject().GetAttributes() {
		nestedSchemaAttribute, err := ResourceSchemaAttribute(ctx, nestedName, path.WithAttributeName(nestedName), nestedA)

		if err != nil {
			return nil, err
		}

		object.Attributes = append(object.Attributes, nestedSchemaAttribute)
	}

	sort.Slice(object.Attributes, func(i, j int) bool {
		if object.Attributes[i] == nil {
			return true
		}

		if object.Attributes[j] == nil {
			return false
		}

		return object.Attributes[i].Name < object.Attributes[j].Name
	})

	schemaAttribute.NestedType = object
	schemaAttribute.Type = nil

	return schemaAttribute, nil
}

func DatasourceSchemaAttribute(ctx context.Context, name string, path *tftypes.AttributePath, a datasourceschema.Attribute) (*schema.SchemaAttribute, error) {
	if !a.IsRequired() && !a.IsOptional() && !a.IsComputed() {
		return nil, path.NewErrorf("must have Required, Optional, or Computed set")
	}

	schemaAttribute := &schema.SchemaAttribute{
		Name:      name,
		Required:  a.IsRequired(),
		Optional:  a.IsOptional(),
		Computed:  a.IsComputed(),
		Sensitive: a.IsSensitive(),
	}
	tfType := a.GetType().TerraformType(ctx)
	b, err := tfType.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("marshalling tftype: %v", err)
	}
	typ, err := ctyjson.UnmarshalType(b)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling to cty type: %v", err)
	}
	schemaAttribute.Type = &typ

	nestedAttribute, ok := a.(datasourceschema.NestedAttribute)

	if !ok {
		return schemaAttribute, nil
	}

	object := &schema.SchemaObject{}
	nm := nestedAttribute.GetNestingMode()
	switch fwschema.NestingMode(nm) {
	case fwschema.NestingModeSingle:
		object.Nesting = schema.SchemaObjectNestingModeSingle
	case fwschema.NestingModeList:
		object.Nesting = schema.SchemaObjectNestingModeList
	case fwschema.NestingModeSet:
		object.Nesting = schema.SchemaObjectNestingModeSet
	case fwschema.NestingModeMap:
		object.Nesting = schema.SchemaObjectNestingModeMap
	default:
		return nil, path.NewErrorf("unrecognized nesting mode %v", nm)
	}

	for nestedName, nestedA := range nestedAttribute.GetNestedObject().GetAttributes() {
		nestedSchemaAttribute, err := DatasourceSchemaAttribute(ctx, nestedName, path.WithAttributeName(nestedName), nestedA)

		if err != nil {
			return nil, err
		}

		object.Attributes = append(object.Attributes, nestedSchemaAttribute)
	}

	sort.Slice(object.Attributes, func(i, j int) bool {
		if object.Attributes[i] == nil {
			return true
		}

		if object.Attributes[j] == nil {
			return false
		}

		return object.Attributes[i].Name < object.Attributes[j].Name
	})

	schemaAttribute.NestedType = object
	schemaAttribute.Type = nil

	return schemaAttribute, nil
}
