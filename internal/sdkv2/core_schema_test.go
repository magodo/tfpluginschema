package sdkv2

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	sdkschema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/magodo/tfpluginschema/schema"
	"github.com/zclconf/go-cty/cty"
)

// A modified version based on: github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema/core_schema_test.go

var (
	typeComparer  = cmp.Comparer(cty.Type.Equals)
	valueComparer = cmp.Comparer(cty.Value.RawEquals)
	equateEmpty   = cmpopts.EquateEmpty()
)

func testSchema(block *schema.SchemaBlock) *schema.SchemaBlock {
	if block.Attributes == nil {
		block.Attributes = []*schema.SchemaAttribute{}
	}

	if block.BlockTypes == nil {
		block.BlockTypes = []*schema.SchemaNestedBlock{}
	}

	// Intentionally remove the logic that adding "id" implicitly.

	return block
}

func testResource(res *schema.Schema) *schema.Schema {
	if res.Block == nil {
		res.Block = testSchema(&schema.SchemaBlock{})
	}
	return res
}

func testProvider(p *schema.ProviderSchema) *schema.ProviderSchema {
	if p.Provider == nil {
		p.Provider = &schema.Schema{Block: testSchema(&schema.SchemaBlock{})}
	}
	if p.ResourceSchemas == nil {
		p.ResourceSchemas = make(map[string]*schema.Schema)
	}
	if p.DataSourceSchemas == nil {
		p.DataSourceSchemas = make(map[string]*schema.Schema)
	}
	return p
}

func TestFromSchemaMap(t *testing.T) {
	tests := map[string]struct {
		Schema map[string]*sdkschema.Schema
		Want   *schema.SchemaBlock
	}{
		"empty": {
			map[string]*sdkschema.Schema{},
			testSchema(&schema.SchemaBlock{}),
		},
		"primitives": {
			map[string]*sdkschema.Schema{
				"int": {
					Type:     sdkschema.TypeInt,
					Required: true,
				},
				"float": {
					Type:     sdkschema.TypeFloat,
					Optional: true,
				},
				"bool": {
					Type:     sdkschema.TypeBool,
					Computed: true,
				},
				"string": {
					Type:     sdkschema.TypeString,
					Optional: true,
					Computed: true,
				},
			},
			testSchema(&schema.SchemaBlock{
				Attributes: []*schema.SchemaAttribute{
					{
						Name:     "bool",
						Type:     ToPtr(cty.Bool),
						Computed: true,
						ForceNew: ToPtr(false),
					},
					{
						Name:     "float",
						Type:     ToPtr(cty.Number),
						Optional: true,
						ForceNew: ToPtr(false),
					},
					{
						Name:     "int",
						Type:     ToPtr(cty.Number),
						Required: true,
						ForceNew: ToPtr(false),
					},
					{
						Name:     "string",
						Type:     ToPtr(cty.String),
						Optional: true,
						Computed: true,
						ForceNew: ToPtr(false),
					},
				},
				BlockTypes: []*schema.SchemaNestedBlock{},
			}),
		},
		"simple collections": {
			map[string]*sdkschema.Schema{
				"list": {
					Type:     sdkschema.TypeList,
					Required: true,
					Elem: &sdkschema.Schema{
						Type: sdkschema.TypeInt,
					},
				},
				"set": {
					Type:     sdkschema.TypeSet,
					Optional: true,
					Elem: &sdkschema.Schema{
						Type: sdkschema.TypeString,
					},
				},
				"map": {
					Type:     sdkschema.TypeMap,
					Optional: true,
					Elem: &sdkschema.Schema{
						Type: sdkschema.TypeBool,
					},
				},
				"map_default_type": {
					Type:     sdkschema.TypeMap,
					Optional: true,
				},
			},
			testSchema(&schema.SchemaBlock{
				Attributes: []*schema.SchemaAttribute{
					{
						Name:     "list",
						Type:     ToPtr(cty.List(cty.Number)),
						Required: true,
						ForceNew: ToPtr(false),
					},
					{
						Name:     "map",
						Type:     ToPtr(cty.Map(cty.Bool)),
						Optional: true,
						ForceNew: ToPtr(false),
					},
					{
						Name:     "map_default_type",
						Type:     ToPtr(cty.Map(cty.String)),
						Optional: true,
						ForceNew: ToPtr(false),
					},
					{
						Name:     "set",
						Type:     ToPtr(cty.Set(cty.String)),
						Optional: true,
						ForceNew: ToPtr(false),
					},
				},
				BlockTypes: []*schema.SchemaNestedBlock{},
			}),
		},
		"incorrectly-specified collections": {
			map[string]*sdkschema.Schema{
				"list": {
					Type:     sdkschema.TypeList,
					Required: true,
					Elem:     sdkschema.TypeInt,
				},
				"set": {
					Type:     sdkschema.TypeSet,
					Optional: true,
					Elem:     sdkschema.TypeString,
				},
				"map": {
					Type:     sdkschema.TypeMap,
					Optional: true,
					Elem:     sdkschema.TypeBool,
				},
			},
			testSchema(&schema.SchemaBlock{
				Attributes: []*schema.SchemaAttribute{
					{
						Name:     "list",
						Type:     ToPtr(cty.List(cty.Number)),
						Required: true,
						ForceNew: ToPtr(false),
					},
					{
						Name:     "map",
						Type:     ToPtr(cty.Map(cty.Bool)),
						Optional: true,
						ForceNew: ToPtr(false),
					},
					{
						Name:     "set",
						Type:     ToPtr(cty.Set(cty.String)),
						Optional: true,
						ForceNew: ToPtr(false),
					},
				},
				BlockTypes: []*schema.SchemaNestedBlock{},
			}),
		},
		"sub-resource collections": {
			map[string]*sdkschema.Schema{
				"list": {
					Type:     sdkschema.TypeList,
					Required: true,
					Elem: &sdkschema.Resource{
						Schema: map[string]*sdkschema.Schema{},
					},
					MinItems: 1,
					MaxItems: 2,
				},
				"set": {
					Type:     sdkschema.TypeSet,
					Required: true,
					Elem: &sdkschema.Resource{
						Schema: map[string]*sdkschema.Schema{},
					},
				},
				"map": {
					Type:     sdkschema.TypeMap,
					Optional: true,
					Elem: &sdkschema.Resource{
						Schema: map[string]*sdkschema.Schema{},
					},
				},
			},
			testSchema(&schema.SchemaBlock{
				Attributes: []*schema.SchemaAttribute{
					{
						Name:     "map",
						Type:     ToPtr(cty.Map(cty.String)),
						Optional: true,
						ForceNew: ToPtr(false),
					},
				},
				BlockTypes: []*schema.SchemaNestedBlock{
					{
						TypeName: "list",
						Required: ToPtr(true),  // NEW
						ForceNew: ToPtr(false), // NEW
						Optional: ToPtr(false), // NEW
						Computed: ToPtr(false), // New
						Nesting:  schema.SchemaNestedBlockNestingModeList,
						Block:    &schema.SchemaBlock{},
						MinItems: 1,
						MaxItems: 2,
					},
					{
						TypeName: "set",
						Required: ToPtr(true),  // NEW
						ForceNew: ToPtr(false), // NEW
						Optional: ToPtr(false), // NEW
						Computed: ToPtr(false), // New
						Nesting:  schema.SchemaNestedBlockNestingModeSet,
						Block:    &schema.SchemaBlock{},
						MinItems: 1,
					},
				},
			}),
		},
		"sub-resource collections minitems+optional": {
			map[string]*sdkschema.Schema{
				"list": {
					Type:     sdkschema.TypeList,
					Optional: true,
					Elem: &sdkschema.Resource{
						Schema: map[string]*sdkschema.Schema{},
					},
					MinItems: 1,
					MaxItems: 1,
				},
				"set": {
					Type:     sdkschema.TypeSet,
					Optional: true,
					Elem: &sdkschema.Resource{
						Schema: map[string]*sdkschema.Schema{},
					},
					MinItems: 1,
					MaxItems: 1,
				},
			},
			testSchema(&schema.SchemaBlock{
				Attributes: []*schema.SchemaAttribute{},
				BlockTypes: []*schema.SchemaNestedBlock{
					{
						TypeName: "list",
						Optional: ToPtr(true),  // NEW
						Required: ToPtr(false), // NEW
						ForceNew: ToPtr(false), // NEW
						Computed: ToPtr(false), // New
						Nesting:  schema.SchemaNestedBlockNestingModeList,
						Block:    &schema.SchemaBlock{},
						MinItems: 0,
						MaxItems: 1,
					},
					{
						TypeName: "set",
						Optional: ToPtr(true),  // NEW
						Required: ToPtr(false), // NEW
						ForceNew: ToPtr(false), // NEW
						Computed: ToPtr(false), // New
						Nesting:  schema.SchemaNestedBlockNestingModeSet,
						Block:    &schema.SchemaBlock{},
						MinItems: 0,
						MaxItems: 1,
					},
				},
			}),
		},
		"sub-resource collections minitems+computed": {
			map[string]*sdkschema.Schema{
				"list": {
					Type:     sdkschema.TypeList,
					Computed: true,
					Elem: &sdkschema.Resource{
						Schema: map[string]*sdkschema.Schema{},
					},
					MinItems: 1,
					MaxItems: 1,
				},
				"set": {
					Type:     sdkschema.TypeSet,
					Computed: true,
					Elem: &sdkschema.Resource{
						Schema: map[string]*sdkschema.Schema{},
					},
					MinItems: 1,
					MaxItems: 1,
				},
			},
			testSchema(&schema.SchemaBlock{
				Attributes: []*schema.SchemaAttribute{
					{
						Name:     "list",
						Type:     ToPtr(cty.List(cty.EmptyObject)),
						Computed: true,
						ForceNew: ToPtr(false),
					},
					{
						Name:     "set",
						Type:     ToPtr(cty.Set(cty.EmptyObject)),
						Computed: true,
						ForceNew: ToPtr(false),
					},
				},
			}),
		},
		"nested attributes and blocks": {
			map[string]*sdkschema.Schema{
				"foo": {
					Type:     sdkschema.TypeList,
					Required: true,
					Elem: &sdkschema.Resource{
						Schema: map[string]*sdkschema.Schema{
							"bar": {
								Type:     sdkschema.TypeList,
								Required: true,
								Elem: &sdkschema.Schema{
									Type: sdkschema.TypeList,
									Elem: &sdkschema.Schema{
										Type: sdkschema.TypeString,
									},
								},
							},
							"baz": {
								Type:     sdkschema.TypeSet,
								Optional: true,
								Elem: &sdkschema.Resource{
									Schema: map[string]*sdkschema.Schema{},
								},
							},
						},
					},
				},
			},
			testSchema(&schema.SchemaBlock{
				Attributes: []*schema.SchemaAttribute{},
				BlockTypes: []*schema.SchemaNestedBlock{
					{
						TypeName: "foo",
						Required: ToPtr(true),  // NEW
						ForceNew: ToPtr(false), // NEW
						Optional: ToPtr(false), // NEW
						Computed: ToPtr(false), // NEW
						Nesting:  schema.SchemaNestedBlockNestingModeList,
						Block: &schema.SchemaBlock{
							Attributes: []*schema.SchemaAttribute{
								{
									Name:     "bar",
									Type:     ToPtr(cty.List(cty.List(cty.String))),
									Required: true,
									ForceNew: ToPtr(false),
								},
							},
							BlockTypes: []*schema.SchemaNestedBlock{
								{
									TypeName: "baz",
									Nesting:  schema.SchemaNestedBlockNestingModeSet,
									Block:    &schema.SchemaBlock{},
									Optional: ToPtr(true),  // NEW
									Required: ToPtr(false), // NEW
									ForceNew: ToPtr(false), // NEW
									Computed: ToPtr(false), // NEW
								},
							},
						},
						MinItems: 1,
					},
				},
			}),
		},
		"sensitive": {
			map[string]*sdkschema.Schema{
				"string": {
					Type:      sdkschema.TypeString,
					Optional:  true,
					Sensitive: true,
				},
			},
			testSchema(&schema.SchemaBlock{
				Attributes: []*schema.SchemaAttribute{
					{
						Name:      "string",
						Type:      ToPtr(cty.String),
						Optional:  true,
						Sensitive: true,
						ForceNew:  ToPtr(false),
					},
				},
				BlockTypes: []*schema.SchemaNestedBlock{},
			}),
		},
		"conditionally required on": {
			map[string]*sdkschema.Schema{
				"string": {
					Type:     sdkschema.TypeString,
					Required: true,
					DefaultFunc: func() (interface{}, error) {
						return nil, nil
					},
				},
			},
			testSchema(&schema.SchemaBlock{
				Attributes: []*schema.SchemaAttribute{
					{
						Name:     "string",
						Type:     ToPtr(cty.String),
						Required: true,
						ForceNew: ToPtr(false),
					},
				},
				BlockTypes: []*schema.SchemaNestedBlock{},
			}),
		},
		"conditionally required off": {
			map[string]*sdkschema.Schema{
				"string": {
					Type:     sdkschema.TypeString,
					Required: true,
					DefaultFunc: func() (interface{}, error) {
						return "boop", nil
					},
				},
			},
			testSchema(&schema.SchemaBlock{
				Attributes: []*schema.SchemaAttribute{
					{
						Name:     "string",
						Type:     ToPtr(cty.String),
						Optional: true,
						ForceNew: ToPtr(false),
					},
				},
				BlockTypes: []*schema.SchemaNestedBlock{},
			}),
		},
		"conditionally required error": {
			map[string]*sdkschema.Schema{
				"string": {
					Type:     sdkschema.TypeString,
					Required: true,
					DefaultFunc: func() (interface{}, error) {
						return nil, fmt.Errorf("placeholder error")
					},
				},
			},
			testSchema(&schema.SchemaBlock{
				Attributes: []*schema.SchemaAttribute{
					{
						Name:     "string",
						Type:     ToPtr(cty.String),
						Optional: true,
						ForceNew: ToPtr(false),
					},
				},
				BlockTypes: []*schema.SchemaNestedBlock{},
			}),
		},

		// NEW
		// Following test cases are newly added for tfpluginschema
		"default value": {
			map[string]*sdkschema.Schema{
				"int": {
					Type:     sdkschema.TypeInt,
					Optional: true,
					Default:  1,
				},
				"float": {
					Type:     sdkschema.TypeFloat,
					Optional: true,
					Default:  1.0,
				},
				"bool": {
					Type:     sdkschema.TypeBool,
					Optional: true,
					Default:  true,
				},
				"string": {
					Type:     sdkschema.TypeString,
					Optional: true,
					Default:  "foo",
				},
			},
			testSchema(&schema.SchemaBlock{
				Attributes: []*schema.SchemaAttribute{
					{
						Name:     "bool",
						Type:     ToPtr(cty.Bool),
						Optional: true,
						Default:  true,
						ForceNew: ToPtr(false),
					},
					{
						Name:     "float",
						Type:     ToPtr(cty.Number),
						Optional: true,
						Default:  1.0,
						ForceNew: ToPtr(false),
					},
					{
						Name:     "int",
						Type:     ToPtr(cty.Number),
						Optional: true,
						Default:  1,
						ForceNew: ToPtr(false),
					},
					{
						Name:     "string",
						Type:     ToPtr(cty.String),
						Optional: true,
						Default:  "foo",
						ForceNew: ToPtr(false),
					},
				},
				BlockTypes: []*schema.SchemaNestedBlock{},
			}),
		},
		"cross attribute constraints": {
			map[string]*sdkschema.Schema{
				"a1": {
					Type:          sdkschema.TypeInt,
					Optional:      true,
					ConflictsWith: []string{"a1", "a2"},
				},
				"a2": {
					Type:          sdkschema.TypeInt,
					Optional:      true,
					ConflictsWith: []string{"a1", "a2"},
				},
				"b1": {
					Type:         sdkschema.TypeInt,
					Optional:     true,
					RequiredWith: []string{"b1", "b2"},
				},
				"b2": {
					Type:         sdkschema.TypeInt,
					Optional:     true,
					RequiredWith: []string{"b1", "b2"},
				},
				"c1": {
					Type:         sdkschema.TypeInt,
					Optional:     true,
					ExactlyOneOf: []string{"c1", "c2"},
				},
				"c2": {
					Type:         sdkschema.TypeInt,
					Optional:     true,
					ExactlyOneOf: []string{"c1", "c2"},
				},
				"d1": {
					Type:         sdkschema.TypeInt,
					Optional:     true,
					AtLeastOneOf: []string{"d1", "d2"},
				},
				"d2": {
					Type:         sdkschema.TypeInt,
					Optional:     true,
					AtLeastOneOf: []string{"d1", "d2"},
				},
			},
			testSchema(&schema.SchemaBlock{
				Attributes: []*schema.SchemaAttribute{
					{
						Name:          "a1",
						Type:          ToPtr(cty.Number),
						Optional:      true,
						ConflictsWith: []string{"a1", "a2"},
						ForceNew:      ToPtr(false),
					},
					{
						Name:          "a2",
						Type:          ToPtr(cty.Number),
						Optional:      true,
						ConflictsWith: []string{"a1", "a2"},
						ForceNew:      ToPtr(false),
					},
					{
						Name:         "b1",
						Type:         ToPtr(cty.Number),
						Optional:     true,
						RequiredWith: []string{"b1", "b2"},
						ForceNew:     ToPtr(false),
					},
					{
						Name:         "b2",
						Type:         ToPtr(cty.Number),
						Optional:     true,
						RequiredWith: []string{"b1", "b2"},
						ForceNew:     ToPtr(false),
					},
					{
						Name:         "c1",
						Type:         ToPtr(cty.Number),
						Optional:     true,
						ExactlyOneOf: []string{"c1", "c2"},
						ForceNew:     ToPtr(false),
					},
					{
						Name:         "c2",
						Type:         ToPtr(cty.Number),
						Optional:     true,
						ExactlyOneOf: []string{"c1", "c2"},
						ForceNew:     ToPtr(false),
					},
					{
						Name:         "d1",
						Type:         ToPtr(cty.Number),
						Optional:     true,
						AtLeastOneOf: []string{"d1", "d2"},
						ForceNew:     ToPtr(false),
					},
					{
						Name:         "d2",
						Type:         ToPtr(cty.Number),
						Optional:     true,
						AtLeastOneOf: []string{"d1", "d2"},
						ForceNew:     ToPtr(false),
					},
				},
				BlockTypes: []*schema.SchemaNestedBlock{},
			}),
		},
		"cross block constraints": {
			map[string]*sdkschema.Schema{
				"a1": {
					Type:     sdkschema.TypeList,
					Optional: true,
					Elem: &sdkschema.Resource{
						Schema: map[string]*sdkschema.Schema{},
					},
					ConflictsWith: []string{"a1", "a2"},
				},
				"a2": {
					Type:     sdkschema.TypeList,
					Optional: true,
					Elem: &sdkschema.Resource{
						Schema: map[string]*sdkschema.Schema{},
					},
					ConflictsWith: []string{"a1", "a2"},
				},
				"b1": {
					Type:     sdkschema.TypeList,
					Optional: true,
					Elem: &sdkschema.Resource{
						Schema: map[string]*sdkschema.Schema{},
					},
					RequiredWith: []string{"b1", "b2"},
				},
				"b2": {
					Type:     sdkschema.TypeList,
					Optional: true,
					Elem: &sdkschema.Resource{
						Schema: map[string]*sdkschema.Schema{},
					},
					RequiredWith: []string{"b1", "b2"},
				},
				"c1": {
					Type:     sdkschema.TypeList,
					Optional: true,
					Elem: &sdkschema.Resource{
						Schema: map[string]*sdkschema.Schema{},
					},
					ExactlyOneOf: []string{"c1", "c2"},
				},
				"c2": {
					Type:     sdkschema.TypeList,
					Optional: true,
					Elem: &sdkschema.Resource{
						Schema: map[string]*sdkschema.Schema{},
					},
					ExactlyOneOf: []string{"c1", "c2"},
				},
				"d1": {
					Type:     sdkschema.TypeList,
					Optional: true,
					Elem: &sdkschema.Resource{
						Schema: map[string]*sdkschema.Schema{},
					},
					AtLeastOneOf: []string{"d1", "d2"},
				},
				"d2": {
					Type:     sdkschema.TypeList,
					Optional: true,
					Elem: &sdkschema.Resource{
						Schema: map[string]*sdkschema.Schema{},
					},
					AtLeastOneOf: []string{"d1", "d2"},
				},
			},
			testSchema(&schema.SchemaBlock{
				Attributes: []*schema.SchemaAttribute{},
				BlockTypes: []*schema.SchemaNestedBlock{
					{
						TypeName:      "a1",
						Nesting:       schema.SchemaNestedBlockNestingModeList,
						Block:         &schema.SchemaBlock{},
						Optional:      ToPtr(true),
						ConflictsWith: []string{"a1", "a2"},
						Required:      ToPtr(false),
						Computed:      ToPtr(false),
						ForceNew:      ToPtr(false),
					},
					{
						TypeName:      "a2",
						Nesting:       schema.SchemaNestedBlockNestingModeList,
						Block:         &schema.SchemaBlock{},
						Optional:      ToPtr(true),
						ConflictsWith: []string{"a1", "a2"},
						Required:      ToPtr(false),
						Computed:      ToPtr(false),
						ForceNew:      ToPtr(false),
					},
					{
						TypeName:     "b1",
						Nesting:      schema.SchemaNestedBlockNestingModeList,
						Block:        &schema.SchemaBlock{},
						Optional:     ToPtr(true),
						RequiredWith: []string{"b1", "b2"},
						Required:     ToPtr(false),
						Computed:     ToPtr(false),
						ForceNew:     ToPtr(false),
					},
					{
						TypeName:     "b2",
						Nesting:      schema.SchemaNestedBlockNestingModeList,
						Block:        &schema.SchemaBlock{},
						Optional:     ToPtr(true),
						RequiredWith: []string{"b1", "b2"},
						Required:     ToPtr(false),
						Computed:     ToPtr(false),
						ForceNew:     ToPtr(false),
					},
					{
						TypeName:     "c1",
						Nesting:      schema.SchemaNestedBlockNestingModeList,
						Block:        &schema.SchemaBlock{},
						Optional:     ToPtr(true),
						ExactlyOneOf: []string{"c1", "c2"},
						Required:     ToPtr(false),
						Computed:     ToPtr(false),
						ForceNew:     ToPtr(false),
					},
					{
						TypeName:     "c2",
						Nesting:      schema.SchemaNestedBlockNestingModeList,
						Block:        &schema.SchemaBlock{},
						Optional:     ToPtr(true),
						ExactlyOneOf: []string{"c1", "c2"},
						Required:     ToPtr(false),
						Computed:     ToPtr(false),
						ForceNew:     ToPtr(false),
					},
					{
						TypeName:     "d1",
						Nesting:      schema.SchemaNestedBlockNestingModeList,
						Block:        &schema.SchemaBlock{},
						Optional:     ToPtr(true),
						AtLeastOneOf: []string{"d1", "d2"},
						Required:     ToPtr(false),
						Computed:     ToPtr(false),
						ForceNew:     ToPtr(false),
					},
					{
						TypeName:     "d2",
						Nesting:      schema.SchemaNestedBlockNestingModeList,
						Block:        &schema.SchemaBlock{},
						Optional:     ToPtr(true),
						AtLeastOneOf: []string{"d1", "d2"},
						Required:     ToPtr(false),
						Computed:     ToPtr(false),
						ForceNew:     ToPtr(false),
					},
				},
			}),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := FromSchemaMap(test.Schema)
			if !cmp.Equal(got, test.Want, equateEmpty, typeComparer) {
				t.Error(cmp.Diff(got, test.Want, equateEmpty, typeComparer))
			}
		})
	}
}

func TestFromResource(t *testing.T) {
	tests := map[string]struct {
		Resource *sdkschema.Resource
		Want     *schema.Schema
	}{
		"empty": {
			&sdkschema.Resource{},
			testResource(&schema.Schema{}),
		},
		"primitives": {
			&sdkschema.Resource{
				SchemaVersion: 1,
				Schema: map[string]*sdkschema.Schema{
					"int": {
						Type:     sdkschema.TypeInt,
						Required: true,
					},
					"float": {
						Type:     sdkschema.TypeFloat,
						Optional: true,
					},
					"bool": {
						Type:     sdkschema.TypeBool,
						Computed: true,
					},
					"string": {
						Type:     sdkschema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
			testResource(&schema.Schema{
				Version: 1,
				Block: &schema.SchemaBlock{
					Attributes: []*schema.SchemaAttribute{
						{
							Name:     "bool",
							Type:     ToPtr(cty.Bool),
							Computed: true,
							ForceNew: ToPtr(false),
						},
						{
							Name:     "float",
							Type:     ToPtr(cty.Number),
							Optional: true,
							ForceNew: ToPtr(false),
						},
						{
							Name:     "int",
							Type:     ToPtr(cty.Number),
							Required: true,
							ForceNew: ToPtr(false),
						},
						{
							Name:     "string",
							Type:     ToPtr(cty.String),
							Optional: true,
							Computed: true,
							ForceNew: ToPtr(false),
						},
					},
					BlockTypes: []*schema.SchemaNestedBlock{},
				},
			}),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := FromResource(test.Resource)
			if !cmp.Equal(got, test.Want, equateEmpty, typeComparer) {
				t.Error(cmp.Diff(got, test.Want, equateEmpty, typeComparer))
			}
		})
	}
}

func TestFromProvider(t *testing.T) {
	tests := map[string]struct {
		Provider *sdkschema.Provider
		Want     *schema.ProviderSchema
	}{
		"empty": {
			&sdkschema.Provider{},
			testProvider(&schema.ProviderSchema{}),
		},
		"full": {
			&sdkschema.Provider{
				Schema: map[string]*sdkschema.Schema{
					"a": {
						Type:     sdkschema.TypeInt,
						Required: true,
					},
				},
				ResourcesMap: map[string]*sdkschema.Resource{
					"foo": {
						SchemaVersion: 1,
						Schema: map[string]*sdkschema.Schema{
							"b": {
								Type:     sdkschema.TypeInt,
								Required: true,
							},
						},
					},
				},
				DataSourcesMap: map[string]*sdkschema.Resource{
					"bar": {
						SchemaVersion: 1,
						Schema: map[string]*sdkschema.Schema{
							"c": {
								Type:     sdkschema.TypeInt,
								Required: true,
							},
						},
					},
				},
			},
			testProvider(&schema.ProviderSchema{
				Provider: &schema.Schema{
					Block: &schema.SchemaBlock{
						Attributes: []*schema.SchemaAttribute{
							{
								Name:     "a",
								Type:     ToPtr(cty.Number),
								Required: true,
								ForceNew: ToPtr(false),
							},
						},
						BlockTypes: []*schema.SchemaNestedBlock{},
					},
				},
				ResourceSchemas: map[string]*schema.Schema{
					"foo": testResource(&schema.Schema{
						Version: 1,
						Block: &schema.SchemaBlock{
							Attributes: []*schema.SchemaAttribute{
								{
									Name:     "b",
									Type:     ToPtr(cty.Number),
									Required: true,
									ForceNew: ToPtr(false),
								},
							},
							BlockTypes: []*schema.SchemaNestedBlock{},
						},
					}),
				},
				DataSourceSchemas: map[string]*schema.Schema{
					"bar": testResource(&schema.Schema{
						Version: 1,
						Block: &schema.SchemaBlock{
							Attributes: []*schema.SchemaAttribute{
								{
									Name:     "c",
									Type:     ToPtr(cty.Number),
									Required: true,
									ForceNew: ToPtr(false),
								},
							},
							BlockTypes: []*schema.SchemaNestedBlock{},
						},
					}),
				},
			}),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := FromProvider(test.Provider)
			if !cmp.Equal(got, test.Want, equateEmpty, typeComparer) {
				t.Error(cmp.Diff(got, test.Want, equateEmpty, typeComparer))
			}
		})
	}
}

func ToPtr[T any](v T) *T {
	return &v
}
