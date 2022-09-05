package sdkv2

// A modified version based on: github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema/core_schema.go

import (
	"fmt"

	sdkschema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/magodo/tfpluginschema/schema"
	"github.com/zclconf/go-cty/cty"
)

func FromProviderSchemaMap(providerschemas map[string]*sdkschema.Schema) *schema.Block {
	if len(providerschemas) == 0 {
		return &schema.Block{}
	}

	ret := &schema.Block{
		Attributes:   map[string]*schema.Attribute{},
		NestedBlocks: map[string]*schema.NestedBlock{},
	}

	for name, ps := range providerschemas {
		if ps.Elem == nil {
			ret.Attributes[name] = fromProviderSchemaAttribute(ps)
			continue
		}
		if ps.Type == sdkschema.TypeMap {
			if _, isResource := ps.Elem.(*sdkschema.Resource); isResource {
				sch := *ps
				sch.Elem = &sdkschema.Schema{
					Type: sdkschema.TypeString,
				}
				ret.Attributes[name] = fromProviderSchemaAttribute(&sch)
				continue
			}
		}
		switch ps.ConfigMode {
		case sdkschema.SchemaConfigModeAttr:
			ret.Attributes[name] = fromProviderSchemaAttribute(ps)
		case sdkschema.SchemaConfigModeBlock:
			ret.NestedBlocks[name] = fromProviderSchemaBlock(ps)
		default: // SchemaConfigModeAuto, or any other invalid value
			if ps.Computed && !ps.Optional {
				// Computed-only schemas are always handled as attributes,
				// because they never appear in configuration.
				ret.Attributes[name] = fromProviderSchemaAttribute(ps)
				continue
			}
			switch ps.Elem.(type) {
			case *sdkschema.Schema, sdkschema.ValueType:
				ret.Attributes[name] = fromProviderSchemaAttribute(ps)
			case *sdkschema.Resource:
				ret.NestedBlocks[name] = fromProviderSchemaBlock(ps)
			default:
				// Should never happen for a valid schema
				panic(fmt.Errorf("invalid Schema.Elem %#v; need *schema.Schema or *schema.Resource", ps.Elem))
			}
		}
	}

	return ret
}

func fromProviderSchemaAttribute(ps *sdkschema.Schema) *schema.Attribute {
	reqd := ps.Required
	opt := ps.Optional
	if reqd && ps.DefaultFunc != nil {
		v, err := ps.DefaultFunc()
		if err != nil || (err == nil && v != nil) {
			reqd = false
			opt = true
		}
	}

	return &schema.Attribute{
		Type:     fromProviderSchemaType(ps),
		Optional: opt,
		Required: reqd,
		Computed: ps.Computed,
		ForceNew: ps.ForceNew,

		Default:   ps.Default,
		Sensitive: ps.Sensitive,

		ConflictsWith: ps.ConflictsWith,
		ExactlyOneOf:  ps.ExactlyOneOf,
		AtLeastOneOf:  ps.AtLeastOneOf,
		RequiredWith:  ps.RequiredWith,
	}
}

func fromProviderSchemaBlock(ps *sdkschema.Schema) *schema.NestedBlock {
	ret := &schema.NestedBlock{
		Required: ps.Required,
		Optional: ps.Optional,
		Computed: ps.Computed,
		ForceNew: ps.ForceNew,

		ConflictsWith: ps.ConflictsWith,
		ExactlyOneOf:  ps.ExactlyOneOf,
		AtLeastOneOf:  ps.AtLeastOneOf,
		RequiredWith:  ps.RequiredWith,
	}

	if nested := fromProviderResource(ps.Elem.(*sdkschema.Resource)); nested != nil {
		ret.Block = nested
	}

	switch ps.Type {
	case sdkschema.TypeList:
		ret.NestingMode = schema.NestingList
	case sdkschema.TypeSet:
		ret.NestingMode = schema.NestingSet
	case sdkschema.TypeMap:
		ret.NestingMode = schema.NestingMap
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
			elemType = ImpliedType(fromProviderResource(set))
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

func fromProviderResource(pr *sdkschema.Resource) *schema.Block {
	return FromProviderSchemaMap(pr.Schema)
}
