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

func testSchema(block *schema.Block) *schema.Block {
	if block.Attributes == nil {
		block.Attributes = make(map[string]*schema.Attribute)
	}

	if block.NestedBlocks == nil {
		block.NestedBlocks = make(map[string]*schema.NestedBlock)
	}

	// Intentionally remove the logic that adding "id" implicitly.

	return block
}

func testResource(res *schema.Resource) *schema.Resource {
	if res.Block == nil {
		res.Block = testSchema(&schema.Block{})
	}
	return res
}

func testProvider(p *schema.ProviderSchema) *schema.ProviderSchema {
	if p.Provider == nil {
		p.Provider = &schema.Schema{Block: testSchema(&schema.Block{})}
	}
	if p.ResourceSchemas == nil {
		p.ResourceSchemas = make(map[string]*schema.Resource)
	}
	if p.DataSourceSchemas == nil {
		p.DataSourceSchemas = make(map[string]*schema.Resource)
	}
	return p
}

func TestFromSchemaMap(t *testing.T) {
	tests := map[string]struct {
		Schema map[string]*sdkschema.Schema
		Want   *schema.Block
	}{
		"empty": {
			map[string]*sdkschema.Schema{},
			testSchema(&schema.Block{}),
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
			testSchema(&schema.Block{
				Attributes: map[string]*schema.Attribute{
					"int": {
						Type:     cty.Number,
						Required: true,
					},
					"float": {
						Type:     cty.Number,
						Optional: true,
					},
					"bool": {
						Type:     cty.Bool,
						Computed: true,
					},
					"string": {
						Type:     cty.String,
						Optional: true,
						Computed: true,
					},
				},
				NestedBlocks: map[string]*schema.NestedBlock{},
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
			testSchema(&schema.Block{
				Attributes: map[string]*schema.Attribute{
					"list": {
						Type:     cty.List(cty.Number),
						Required: true,
					},
					"set": {
						Type:     cty.Set(cty.String),
						Optional: true,
					},
					"map": {
						Type:     cty.Map(cty.Bool),
						Optional: true,
					},
					"map_default_type": {
						Type:     cty.Map(cty.String),
						Optional: true,
					},
				},
				NestedBlocks: map[string]*schema.NestedBlock{},
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
			testSchema(&schema.Block{
				Attributes: map[string]*schema.Attribute{
					"list": {
						Type:     cty.List(cty.Number),
						Required: true,
					},
					"set": {
						Type:     cty.Set(cty.String),
						Optional: true,
					},
					"map": {
						Type:     cty.Map(cty.Bool),
						Optional: true,
					},
				},
				NestedBlocks: map[string]*schema.NestedBlock{},
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
			testSchema(&schema.Block{
				Attributes: map[string]*schema.Attribute{
					"map": {
						Type:     cty.Map(cty.String),
						Optional: true,
					},
				},
				NestedBlocks: map[string]*schema.NestedBlock{
					"list": {
						Required:    true, // NEW
						NestingMode: schema.NestingList,
						Block:       &schema.Block{},
						MinItems:    1,
						MaxItems:    2,
					},
					"set": {
						Required:    true, // NEW
						NestingMode: schema.NestingSet,
						Block:       &schema.Block{},
						MinItems:    1,
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
			testSchema(&schema.Block{
				Attributes: map[string]*schema.Attribute{},
				NestedBlocks: map[string]*schema.NestedBlock{
					"list": {
						Optional:    true, // NEW
						NestingMode: schema.NestingList,
						Block:       &schema.Block{},
						MinItems:    0,
						MaxItems:    1,
					},
					"set": {
						Optional:    true, // NEW
						NestingMode: schema.NestingSet,
						Block:       &schema.Block{},
						MinItems:    0,
						MaxItems:    1,
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
			testSchema(&schema.Block{
				Attributes: map[string]*schema.Attribute{
					"list": {
						Type:     cty.List(cty.EmptyObject),
						Computed: true,
					},
					"set": {
						Type:     cty.Set(cty.EmptyObject),
						Computed: true,
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
			testSchema(&schema.Block{
				Attributes: map[string]*schema.Attribute{},
				NestedBlocks: map[string]*schema.NestedBlock{
					"foo": {
						Required:    true, // NEW
						NestingMode: schema.NestingList,
						Block: &schema.Block{
							Attributes: map[string]*schema.Attribute{
								"bar": {
									Type:     cty.List(cty.List(cty.String)),
									Required: true,
								},
							},
							NestedBlocks: map[string]*schema.NestedBlock{
								"baz": {
									NestingMode: schema.NestingSet,
									Block:       &schema.Block{},
									Optional:    true, // NEW
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
			testSchema(&schema.Block{
				Attributes: map[string]*schema.Attribute{
					"string": {
						Type:      cty.String,
						Optional:  true,
						Sensitive: true,
					},
				},
				NestedBlocks: map[string]*schema.NestedBlock{},
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
			testSchema(&schema.Block{
				Attributes: map[string]*schema.Attribute{
					"string": {
						Type:     cty.String,
						Required: true,
					},
				},
				NestedBlocks: map[string]*schema.NestedBlock{},
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
			testSchema(&schema.Block{
				Attributes: map[string]*schema.Attribute{
					"string": {
						Type:     cty.String,
						Optional: true,
					},
				},
				NestedBlocks: map[string]*schema.NestedBlock{},
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
			testSchema(&schema.Block{
				Attributes: map[string]*schema.Attribute{
					"string": {
						Type:     cty.String,
						Optional: true,
					},
				},
				NestedBlocks: map[string]*schema.NestedBlock{},
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
			testSchema(&schema.Block{
				Attributes: map[string]*schema.Attribute{
					"int": {
						Type:     cty.Number,
						Optional: true,
						Default:  1,
					},
					"float": {
						Type:     cty.Number,
						Optional: true,
						Default:  1.0,
					},
					"bool": {
						Type:     cty.Bool,
						Optional: true,
						Default:  true,
					},
					"string": {
						Type:     cty.String,
						Optional: true,
						Default:  "foo",
					},
				},
				NestedBlocks: map[string]*schema.NestedBlock{},
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
			testSchema(&schema.Block{
				Attributes: map[string]*schema.Attribute{
					"a1": {
						Type:          cty.Number,
						Optional:      true,
						ConflictsWith: []string{"a1", "a2"},
					},
					"a2": {
						Type:          cty.Number,
						Optional:      true,
						ConflictsWith: []string{"a1", "a2"},
					},
					"b1": {
						Type:         cty.Number,
						Optional:     true,
						RequiredWith: []string{"b1", "b2"},
					},
					"b2": {
						Type:         cty.Number,
						Optional:     true,
						RequiredWith: []string{"b1", "b2"},
					},
					"c1": {
						Type:         cty.Number,
						Optional:     true,
						ExactlyOneOf: []string{"c1", "c2"},
					},
					"c2": {
						Type:         cty.Number,
						Optional:     true,
						ExactlyOneOf: []string{"c1", "c2"},
					},
					"d1": {
						Type:         cty.Number,
						Optional:     true,
						AtLeastOneOf: []string{"d1", "d2"},
					},
					"d2": {
						Type:         cty.Number,
						Optional:     true,
						AtLeastOneOf: []string{"d1", "d2"},
					},
				},
				NestedBlocks: map[string]*schema.NestedBlock{},
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
			testSchema(&schema.Block{
				Attributes: map[string]*schema.Attribute{},
				NestedBlocks: map[string]*schema.NestedBlock{
					"a1": {
						NestingMode:   schema.NestingList,
						Block:         &schema.Block{},
						Optional:      true,
						ConflictsWith: []string{"a1", "a2"},
					},
					"a2": {
						NestingMode:   schema.NestingList,
						Block:         &schema.Block{},
						Optional:      true,
						ConflictsWith: []string{"a1", "a2"},
					},
					"b1": {
						NestingMode:  schema.NestingList,
						Block:        &schema.Block{},
						Optional:     true,
						RequiredWith: []string{"b1", "b2"},
					},
					"b2": {
						NestingMode:  schema.NestingList,
						Block:        &schema.Block{},
						Optional:     true,
						RequiredWith: []string{"b1", "b2"},
					},
					"c1": {
						NestingMode:  schema.NestingList,
						Block:        &schema.Block{},
						Optional:     true,
						ExactlyOneOf: []string{"c1", "c2"},
					},
					"c2": {
						NestingMode:  schema.NestingList,
						Block:        &schema.Block{},
						Optional:     true,
						ExactlyOneOf: []string{"c1", "c2"},
					},
					"d1": {
						NestingMode:  schema.NestingList,
						Block:        &schema.Block{},
						Optional:     true,
						AtLeastOneOf: []string{"d1", "d2"},
					},
					"d2": {
						NestingMode:  schema.NestingList,
						Block:        &schema.Block{},
						Optional:     true,
						AtLeastOneOf: []string{"d1", "d2"},
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
		Want     *schema.Resource
	}{
		"empty": {
			&sdkschema.Resource{},
			testResource(&schema.Resource{}),
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
			testResource(&schema.Resource{
				SchemaVersion: 1,
				Block: &schema.Block{
					Attributes: map[string]*schema.Attribute{
						"int": {
							Type:     cty.Number,
							Required: true,
						},
						"float": {
							Type:     cty.Number,
							Optional: true,
						},
						"bool": {
							Type:     cty.Bool,
							Computed: true,
						},
						"string": {
							Type:     cty.String,
							Optional: true,
							Computed: true,
						},
					},
					NestedBlocks: map[string]*schema.NestedBlock{},
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
					Block: &schema.Block{
						Attributes: map[string]*schema.Attribute{
							"a": {
								Type:     cty.Number,
								Required: true,
							},
						},
						NestedBlocks: map[string]*schema.NestedBlock{},
					},
				},
				ResourceSchemas: map[string]*schema.Resource{
					"foo": testResource(&schema.Resource{
						SchemaVersion: 1,
						Block: &schema.Block{
							Attributes: map[string]*schema.Attribute{
								"b": {
									Type:     cty.Number,
									Required: true,
								},
							},
							NestedBlocks: map[string]*schema.NestedBlock{},
						},
					}),
				},
				DataSourceSchemas: map[string]*schema.Resource{
					"bar": testResource(&schema.Resource{
						SchemaVersion: 1,
						Block: &schema.Block{
							Attributes: map[string]*schema.Attribute{
								"c": {
									Type:     cty.Number,
									Required: true,
								},
							},
							NestedBlocks: map[string]*schema.NestedBlock{},
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
