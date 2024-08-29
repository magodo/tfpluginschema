package sdkv2

// A modified version based on: github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema/core_schema.go

import (
	"fmt"
	"sort"

	sdkschema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/magodo/tfpluginschema/schema"
	"github.com/zclconf/go-cty/cty"
)

func FromSchemaMap(m map[string]*sdkschema.Schema) *schema.SchemaBlock {
	if len(m) == 0 {
		return &schema.SchemaBlock{}
	}

	ret := &schema.SchemaBlock{
		Attributes: []*schema.SchemaAttribute{},
		BlockTypes: []*schema.SchemaNestedBlock{},
	}

	for name, ps := range m {
		if ps.Elem == nil {
			ret.Attributes = append(ret.Attributes, fromProviderSchemaAttribute(name, ps))
			continue
		}
		if ps.Type == sdkschema.TypeMap {
			if _, isResource := ps.Elem.(*sdkschema.Resource); isResource {
				sch := *ps
				sch.Elem = &sdkschema.Schema{
					Type: sdkschema.TypeString,
				}
				ret.Attributes = append(ret.Attributes, fromProviderSchemaAttribute(name, &sch))
				continue
			}
		}
		switch ps.ConfigMode {
		case sdkschema.SchemaConfigModeAttr:
			ret.Attributes = append(ret.Attributes, fromProviderSchemaAttribute(name, ps))
		case sdkschema.SchemaConfigModeBlock:
			ret.BlockTypes = append(ret.BlockTypes, fromProviderSchemaBlock(name, ps))
		default: // SchemaConfigModeAuto, or any other invalid value
			if ps.Computed && !ps.Optional {
				// Computed-only schemas are always handled as attributes,
				// because they never appear in configuration.
				ret.Attributes = append(ret.Attributes, fromProviderSchemaAttribute(name, ps))
				continue
			}
			switch ps.Elem.(type) {
			case *sdkschema.Schema, sdkschema.ValueType:
				ret.Attributes = append(ret.Attributes, fromProviderSchemaAttribute(name, ps))
			case *sdkschema.Resource:
				ret.BlockTypes = append(ret.BlockTypes, fromProviderSchemaBlock(name, ps))
			default:
				// Should never happen for a valid schema
				panic(fmt.Errorf("invalid Schema.Elem %#v; need *schema.Schema or *schema.Resource", ps.Elem))
			}
		}
	}

	sort.Slice(ret.Attributes, func(i, j int) bool {
		return ret.Attributes[i].Name < ret.Attributes[j].Name
	})

	sort.Slice(ret.BlockTypes, func(i, j int) bool {
		return ret.BlockTypes[i].TypeName < ret.BlockTypes[j].TypeName
	})

	return ret
}

func fromProviderSchemaAttribute(name string, ps *sdkschema.Schema) *schema.SchemaAttribute {
	reqd := ps.Required
	opt := ps.Optional
	if reqd && ps.DefaultFunc != nil {
		v, err := ps.DefaultFunc()
		if err != nil || (err == nil && v != nil) {
			reqd = false
			opt = true
		}
	}
	typ := fromProviderSchemaType(ps)

	return &schema.SchemaAttribute{
		Name:     name,
		Type:     &typ,
		Optional: opt,
		Required: reqd,
		Computed: ps.Computed,
		ForceNew: &ps.ForceNew,

		Default:   ps.Default,
		Sensitive: ps.Sensitive,

		ConflictsWith: ps.ConflictsWith,
		ExactlyOneOf:  ps.ExactlyOneOf,
		AtLeastOneOf:  ps.AtLeastOneOf,
		RequiredWith:  ps.RequiredWith,
	}
}

func fromProviderSchemaBlock(name string, ps *sdkschema.Schema) *schema.SchemaNestedBlock {
	ret := &schema.SchemaNestedBlock{
		TypeName: name,
		Required: &ps.Required,
		Optional: &ps.Optional,
		Computed: &ps.Computed,
		ForceNew: &ps.ForceNew,

		ConflictsWith: ps.ConflictsWith,
		ExactlyOneOf:  ps.ExactlyOneOf,
		AtLeastOneOf:  ps.AtLeastOneOf,
		RequiredWith:  ps.RequiredWith,
	}

	if nested := FromResource(ps.Elem.(*sdkschema.Resource)); nested != nil {
		ret.Block = nested.Block
	}

	switch ps.Type {
	case sdkschema.TypeList:
		ret.Nesting = schema.SchemaNestedBlockNestingModeList
	case sdkschema.TypeSet:
		ret.Nesting = schema.SchemaNestedBlockNestingModeSet
	case sdkschema.TypeMap:
		ret.Nesting = schema.SchemaNestedBlockNestingModeMap
	default:
		// Should never happen for a valid schema
		panic(fmt.Errorf("invalid s.Type %s for s.Elem being resource", ps.Type))
	}

	ret.MinItems = ps.MinItems
	ret.MaxItems = ps.MaxItems

	if ps.Required && ps.MinItems == 0 {
		// configschema doesn't have a "required" representation for nested
		// blocks, but we can fake it by requiring at least one item.
		ret.MinItems = 1
	}
	if ps.Optional && ps.MinItems > 0 {
		// Historically helper/schema would ignore MinItems if Optional were
		// set, so we must mimic this behavior here to ensure that providers
		// relying on that undocumented behavior can continue to operate as
		// they did before.
		ret.MinItems = 0
	}
	if ps.Computed && !ps.Optional {
		// MinItems/MaxItems are meaningless for computed nested blocks, since
		// they are never set by the user anyway. This ensures that we'll never
		// generate weird errors about them.
		ret.MinItems = 0
		ret.MaxItems = 0
	}

	return ret
}

func fromProviderSchemaType(ps *sdkschema.Schema) cty.Type {
	switch ps.Type {
	case sdkschema.TypeString:
		return cty.String
	case sdkschema.TypeBool:
		return cty.Bool
	case sdkschema.TypeInt, sdkschema.TypeFloat:
		return cty.Number
	case sdkschema.TypeList, sdkschema.TypeSet, sdkschema.TypeMap:
		var elemType cty.Type
		switch set := ps.Elem.(type) {
		case *sdkschema.Schema:
			elemType = fromProviderSchemaType(set)
		case sdkschema.ValueType:
			elemType = fromProviderSchemaType(&sdkschema.Schema{Type: set})
		case *sdkschema.Resource:
			elemType = ImpliedType(FromResource(set).Block)
		default:
			if set != nil {
				panic(fmt.Errorf("invalid Schema.Elem %#v; need *schema.Schema or *schema.Resource", ps.Elem))
			}
			elemType = cty.String
		}
		switch ps.Type {
		case sdkschema.TypeList:
			return cty.List(elemType)
		case sdkschema.TypeSet:
			return cty.Set(elemType)
		case sdkschema.TypeMap:
			return cty.Map(elemType)
		default:
			panic("invalid collection type")
		}
	default:
		panic(fmt.Errorf("invalid Schema.Type %s", ps.Type))
	}
}

func FromResource(res *sdkschema.Resource) *schema.Schema {
	ret := &schema.Schema{
		Version: int64(res.SchemaVersion),
		Block:   FromSchemaMap(res.Schema),
	}
	return ret
}

func FromProvider(p *sdkschema.Provider) *schema.ProviderSchema {
	ret := &schema.ProviderSchema{
		Provider: &schema.Schema{
			Block: FromSchemaMap(p.Schema),
		},
		ResourceSchemas:   map[string]*schema.Schema{},
		DataSourceSchemas: map[string]*schema.Schema{},
	}

	for name, res := range p.ResourcesMap {
		ret.ResourceSchemas[name] = FromResource(res)
	}
	for name, res := range p.DataSourcesMap {
		ret.DataSourceSchemas[name] = FromResource(res)
	}
	return ret
}
