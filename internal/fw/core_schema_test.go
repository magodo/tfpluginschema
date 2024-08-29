package fw_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	providerschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/magodo/tfpluginschema/internal/fw"
	"github.com/magodo/tfpluginschema/schema"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
)

var (
	typeComparer  = cmp.Comparer(cty.Type.Equals)
	valueComparer = cmp.Comparer(cty.Value.RawEquals)
	equateEmpty   = cmpopts.EquateEmpty()
)

var _ provider.Provider = &TestProvider{}

type TestProvider struct{}

// Configure implements provider.Provider.
func (t *TestProvider) Configure(context.Context, provider.ConfigureRequest, *provider.ConfigureResponse) {
	panic("unimplemented")
}

// DataSources implements provider.Provider.
func (t *TestProvider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource {
			return &TestDatasource{}
		},
	}
}

// Metadata implements provider.Provider.
func (t *TestProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "foo"
}

// Resources implements provider.Provider.
func (t *TestProvider) Resources(context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource {
			return &TestResource{}
		},
	}
}

// Schema implements provider.Provider.
func (t *TestProvider) Schema(ctx context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = providerschema.Schema{
		Attributes: map[string]providerschema.Attribute{
			"bool": providerschema.BoolAttribute{
				Required:  true,
				Optional:  true,
				Sensitive: true,
			},
			"number": providerschema.NumberAttribute{
				Required: true,
			},
			"string": providerschema.StringAttribute{
				Required: true,
			},
			"list": providerschema.ListAttribute{
				Required:    true,
				ElementType: basetypes.StringType{},
			},
			"set": providerschema.SetAttribute{
				Required:    true,
				ElementType: basetypes.StringType{},
			},
			"map": providerschema.MapAttribute{
				Required:    true,
				ElementType: basetypes.StringType{},
			},
			"single_nested": providerschema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]providerschema.Attribute{
					"string": providerschema.StringAttribute{
						Required: true,
					},
				},
			},
			"list_nested": providerschema.ListNestedAttribute{
				Required: true,
				NestedObject: providerschema.NestedAttributeObject{
					Attributes: map[string]providerschema.Attribute{
						"string": providerschema.StringAttribute{
							Required: true,
						},
					},
				},
			},
			"object": providerschema.ObjectAttribute{
				Required: true,
				AttributeTypes: map[string]attr.Type{
					"string": basetypes.StringType{},
				},
			},
		},
		Blocks: map[string]providerschema.Block{
			"single": providerschema.SingleNestedBlock{
				Attributes: map[string]providerschema.Attribute{
					"string": providerschema.StringAttribute{
						Required: true,
					},
				},
			},
			"list": providerschema.ListNestedBlock{
				NestedObject: providerschema.NestedBlockObject{
					Attributes: map[string]providerschema.Attribute{
						"string": providerschema.StringAttribute{
							Required: true,
						},
					},
				},
			},
			"set": providerschema.SetNestedBlock{
				NestedObject: providerschema.NestedBlockObject{
					Attributes: map[string]providerschema.Attribute{
						"string": providerschema.StringAttribute{
							Required: true,
						},
					},
				},
			},
		},
	}
}

var _ resource.Resource = &TestResource{}

type TestResource struct{}

// Create implements resource.Resource.
func (t *TestResource) Create(context.Context, resource.CreateRequest, *resource.CreateResponse) {
	panic("unimplemented")
}

// Delete implements resource.Resource.
func (t *TestResource) Delete(context.Context, resource.DeleteRequest, *resource.DeleteResponse) {
	panic("unimplemented")
}

// Metadata implements resource.Resource.
func (t *TestResource) Metadata(ctx context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "foo_resource"
}

// Read implements resource.Resource.
func (t *TestResource) Read(context.Context, resource.ReadRequest, *resource.ReadResponse) {
	panic("unimplemented")
}

// Schema implements resource.Resource.
func (t *TestResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceschema.Schema{
		Attributes: map[string]resourceschema.Attribute{
			"bool": resourceschema.BoolAttribute{
				Required:  true,
				Optional:  true,
				Sensitive: true,
				Computed:  true,
				Default:   booldefault.StaticBool(true),
			},
			"number": resourceschema.NumberAttribute{
				Required: true,
			},
			"string": resourceschema.StringAttribute{
				Required: true,
				Default:  stringdefault.StaticString("foo"),
			},
			"list": resourceschema.ListAttribute{
				Required:    true,
				ElementType: basetypes.StringType{},
				Default:     listdefault.StaticValue(basetypes.NewListValueMust(basetypes.StringType{}, []attr.Value{basetypes.NewStringValue("a")})),
			},
			"set": resourceschema.SetAttribute{
				Required:    true,
				ElementType: basetypes.StringType{},
				Default:     setdefault.StaticValue(basetypes.NewSetValueMust(basetypes.StringType{}, []attr.Value{basetypes.NewStringValue("a")})),
			},
			"map": resourceschema.MapAttribute{
				Required:    true,
				ElementType: basetypes.StringType{},
				Default:     mapdefault.StaticValue(basetypes.NewMapValueMust(basetypes.StringType{}, map[string]attr.Value{"a": basetypes.NewStringValue("a")})),
			},
			"single_nested": resourceschema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]resourceschema.Attribute{
					"string": resourceschema.StringAttribute{
						Required: true,
					},
				},
				Default: objectdefault.StaticValue(basetypes.NewObjectValueMust(map[string]attr.Type{"a": basetypes.StringType{}}, map[string]attr.Value{"a": basetypes.NewStringValue("a")})),
			},
			"list_nested": resourceschema.ListNestedAttribute{
				Required: true,
				NestedObject: resourceschema.NestedAttributeObject{
					Attributes: map[string]resourceschema.Attribute{
						"string": resourceschema.StringAttribute{
							Required: true,
						},
					},
				},
				Default: listdefault.StaticValue(basetypes.NewListValueMust(basetypes.ObjectType{AttrTypes: map[string]attr.Type{"a": basetypes.StringType{}}}, []attr.Value{basetypes.NewObjectValueMust(map[string]attr.Type{"a": basetypes.StringType{}}, map[string]attr.Value{"a": basetypes.NewStringValue("a")})})),
			},
			"object": resourceschema.ObjectAttribute{
				Required: true,
				AttributeTypes: map[string]attr.Type{
					"string": basetypes.StringType{},
				},
				Default: objectdefault.StaticValue(basetypes.NewObjectValueMust(map[string]attr.Type{"a": basetypes.StringType{}}, map[string]attr.Value{"a": basetypes.NewStringValue("a")})),
			},
		},
		Blocks: map[string]resourceschema.Block{
			"single": resourceschema.SingleNestedBlock{
				Attributes: map[string]resourceschema.Attribute{
					"string": resourceschema.StringAttribute{
						Required: true,
					},
				},
			},
			"list": resourceschema.ListNestedBlock{
				NestedObject: resourceschema.NestedBlockObject{
					Attributes: map[string]resourceschema.Attribute{
						"string": resourceschema.StringAttribute{
							Required: true,
						},
					},
				},
			},
			"set": resourceschema.SetNestedBlock{
				NestedObject: resourceschema.NestedBlockObject{
					Attributes: map[string]resourceschema.Attribute{
						"string": resourceschema.StringAttribute{
							Required: true,
						},
					},
				},
			},
		},
	}
}

// Update implements resource.Resource.
func (t *TestResource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {
	panic("unimplemented")
}

var _ datasource.DataSource = &TestDatasource{}

type TestDatasource struct{}

// Metadata implements datasource.DataSource.
func (t *TestDatasource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "foo_resource"
}

// Read implements datasource.DataSource.
func (t *TestDatasource) Read(context.Context, datasource.ReadRequest, *datasource.ReadResponse) {
	panic("unimplemented")
}

// Schema implements datasource.DataSource.
func (t *TestDatasource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceschema.Schema{
		Attributes: map[string]datasourceschema.Attribute{
			"bool": datasourceschema.BoolAttribute{
				Required:  true,
				Optional:  true,
				Sensitive: true,
				Computed:  true,
			},
			"number": datasourceschema.NumberAttribute{
				Required: true,
			},
			"string": datasourceschema.StringAttribute{
				Required: true,
			},
			"list": datasourceschema.ListAttribute{
				Required:    true,
				ElementType: basetypes.StringType{},
			},
			"set": datasourceschema.SetAttribute{
				Required:    true,
				ElementType: basetypes.StringType{},
			},
			"map": datasourceschema.MapAttribute{
				Required:    true,
				ElementType: basetypes.StringType{},
			},
			"single_nested": datasourceschema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]datasourceschema.Attribute{
					"string": datasourceschema.StringAttribute{
						Required: true,
					},
				},
			},
			"list_nested": datasourceschema.ListNestedAttribute{
				Required: true,
				NestedObject: datasourceschema.NestedAttributeObject{
					Attributes: map[string]datasourceschema.Attribute{
						"string": datasourceschema.StringAttribute{
							Required: true,
						},
					},
				},
			},
			"object": datasourceschema.ObjectAttribute{
				Required: true,
				AttributeTypes: map[string]attr.Type{
					"string": basetypes.StringType{},
				},
			},
		},
		Blocks: map[string]datasourceschema.Block{
			"single": datasourceschema.SingleNestedBlock{
				Attributes: map[string]datasourceschema.Attribute{
					"string": datasourceschema.StringAttribute{
						Required: true,
					},
				},
			},
			"list": datasourceschema.ListNestedBlock{
				NestedObject: datasourceschema.NestedBlockObject{
					Attributes: map[string]datasourceschema.Attribute{
						"string": datasourceschema.StringAttribute{
							Required: true,
						},
					},
				},
			},
			"set": datasourceschema.SetNestedBlock{
				NestedObject: datasourceschema.NestedBlockObject{
					Attributes: map[string]datasourceschema.Attribute{
						"string": datasourceschema.StringAttribute{
							Required: true,
						},
					},
				},
			},
		},
	}
}

func TestFromProvider(t *testing.T) {
	got, err := fw.FromProvider(&TestProvider{})
	require.NoError(t, err)

	want := &schema.ProviderSchema{
		Provider: &schema.Schema{
			Block: &schema.SchemaBlock{
				Attributes: []*schema.SchemaAttribute{
					{
						Name:      "bool",
						Type:      &cty.Bool,
						Required:  true,
						Optional:  true,
						Sensitive: true,
					},
					{
						Name:     "list",
						Required: true,
						Type:     ToPtr(cty.List(cty.String)),
					},
					{
						Name:     "list_nested",
						Required: true,
						NestedType: &schema.SchemaObject{
							Attributes: []*schema.SchemaAttribute{
								{
									Name:     "string",
									Type:     &cty.String,
									Required: true,
								},
							},
							Nesting: schema.SchemaObjectNestingModeList,
						},
					},
					{
						Name:     "map",
						Required: true,
						Type:     ToPtr(cty.Map(cty.String)),
					},
					{
						Name:     "number",
						Required: true,
						Type:     &cty.Number,
					},
					{
						Name:     "object",
						Required: true,
						Type:     ToPtr(cty.Object(map[string]cty.Type{"string": cty.String})),
					},
					{
						Name:     "set",
						Required: true,
						Type:     ToPtr(cty.Set(cty.String)),
					},
					{
						Name:     "single_nested",
						Required: true,
						NestedType: &schema.SchemaObject{
							Attributes: []*schema.SchemaAttribute{
								{
									Name:     "string",
									Type:     &cty.String,
									Required: true,
								},
							},
							Nesting: schema.SchemaObjectNestingModeSingle,
						},
					},
					{
						Name:     "string",
						Required: true,
						Type:     &cty.String,
					},
				},
				BlockTypes: []*schema.SchemaNestedBlock{
					{
						TypeName: "list",
						Nesting:  schema.SchemaNestedBlockNestingModeList,
						Block: &schema.SchemaBlock{
							Attributes: []*schema.SchemaAttribute{
								{
									Name:     "string",
									Type:     &cty.String,
									Required: true,
								},
							},
						},
					},
					{
						TypeName: "set",
						Nesting:  schema.SchemaNestedBlockNestingModeSet,
						Block: &schema.SchemaBlock{
							Attributes: []*schema.SchemaAttribute{
								{
									Name:     "string",
									Type:     &cty.String,
									Required: true,
								},
							},
						},
					},
					{
						TypeName: "single",
						Nesting:  schema.SchemaNestedBlockNestingModeSingle,
						Block: &schema.SchemaBlock{
							Attributes: []*schema.SchemaAttribute{
								{
									Name:     "string",
									Type:     &cty.String,
									Required: true,
								},
							},
						},
					},
				},
			},
		},
		ResourceSchemas: map[string]*schema.Schema{
			"foo_resource": {
				Block: &schema.SchemaBlock{
					Attributes: []*schema.SchemaAttribute{
						{
							Name:      "bool",
							Type:      &cty.Bool,
							Required:  true,
							Optional:  true,
							Sensitive: true,
							Computed:  true,
							Default:   true,
						},
						{
							Name:     "list",
							Required: true,
							Type:     ToPtr(cty.List(cty.String)),
							Default:  []interface{}{"a"},
						},
						{
							Name:     "list_nested",
							Required: true,
							NestedType: &schema.SchemaObject{
								Attributes: []*schema.SchemaAttribute{
									{
										Name:     "string",
										Type:     &cty.String,
										Required: true,
									},
								},
								Nesting: schema.SchemaObjectNestingModeList,
							},
							Default: []interface{}{map[string]interface{}{"a": "a"}},
						},
						{
							Name:     "map",
							Required: true,
							Type:     ToPtr(cty.Map(cty.String)),
							Default:  map[string]interface{}{"a": "a"},
						},
						{
							Name:     "number",
							Required: true,
							Type:     &cty.Number,
						},
						{
							Name:     "object",
							Required: true,
							Type:     ToPtr(cty.Object(map[string]cty.Type{"string": cty.String})),
							Default:  map[string]interface{}{"a": "a"},
						},
						{
							Name:     "set",
							Required: true,
							Type:     ToPtr(cty.Set(cty.String)),
							Default:  []interface{}{"a"},
						},
						{
							Name:     "single_nested",
							Required: true,
							NestedType: &schema.SchemaObject{
								Attributes: []*schema.SchemaAttribute{
									{
										Name:     "string",
										Type:     &cty.String,
										Required: true,
									},
								},
								Nesting: schema.SchemaObjectNestingModeSingle,
							},
							Default: map[string]interface{}{"a": "a"},
						},
						{
							Name:     "string",
							Required: true,
							Type:     &cty.String,
							Default:  "foo",
						},
					},
					BlockTypes: []*schema.SchemaNestedBlock{
						{
							TypeName: "list",
							Nesting:  schema.SchemaNestedBlockNestingModeList,
							Block: &schema.SchemaBlock{
								Attributes: []*schema.SchemaAttribute{
									{
										Name:     "string",
										Type:     &cty.String,
										Required: true,
									},
								},
							},
						},
						{
							TypeName: "set",
							Nesting:  schema.SchemaNestedBlockNestingModeSet,
							Block: &schema.SchemaBlock{
								Attributes: []*schema.SchemaAttribute{
									{
										Name:     "string",
										Type:     &cty.String,
										Required: true,
									},
								},
							},
						},
						{
							TypeName: "single",
							Nesting:  schema.SchemaNestedBlockNestingModeSingle,
							Block: &schema.SchemaBlock{
								Attributes: []*schema.SchemaAttribute{
									{
										Name:     "string",
										Type:     &cty.String,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
		},
		DataSourceSchemas: map[string]*schema.Schema{
			"foo_resource": {
				Block: &schema.SchemaBlock{
					Attributes: []*schema.SchemaAttribute{
						{
							Name:      "bool",
							Type:      &cty.Bool,
							Required:  true,
							Optional:  true,
							Sensitive: true,
							Computed:  true,
						},
						{
							Name:     "list",
							Required: true,
							Type:     ToPtr(cty.List(cty.String)),
						},
						{
							Name:     "list_nested",
							Required: true,
							NestedType: &schema.SchemaObject{
								Attributes: []*schema.SchemaAttribute{
									{
										Name:     "string",
										Type:     &cty.String,
										Required: true,
									},
								},
								Nesting: schema.SchemaObjectNestingModeList,
							},
						},
						{
							Name:     "map",
							Required: true,
							Type:     ToPtr(cty.Map(cty.String)),
						},
						{
							Name:     "number",
							Required: true,
							Type:     &cty.Number,
						},
						{
							Name:     "object",
							Required: true,
							Type:     ToPtr(cty.Object(map[string]cty.Type{"string": cty.String})),
						},
						{
							Name:     "set",
							Required: true,
							Type:     ToPtr(cty.Set(cty.String)),
						},
						{
							Name:     "single_nested",
							Required: true,
							NestedType: &schema.SchemaObject{
								Attributes: []*schema.SchemaAttribute{
									{
										Name:     "string",
										Type:     &cty.String,
										Required: true,
									},
								},
								Nesting: schema.SchemaObjectNestingModeSingle,
							},
						},
						{
							Name:     "string",
							Required: true,
							Type:     &cty.String,
						},
					},
					BlockTypes: []*schema.SchemaNestedBlock{
						{
							TypeName: "list",
							Nesting:  schema.SchemaNestedBlockNestingModeList,
							Block: &schema.SchemaBlock{
								Attributes: []*schema.SchemaAttribute{
									{
										Name:     "string",
										Type:     &cty.String,
										Required: true,
									},
								},
							},
						},
						{
							TypeName: "set",
							Nesting:  schema.SchemaNestedBlockNestingModeSet,
							Block: &schema.SchemaBlock{
								Attributes: []*schema.SchemaAttribute{
									{
										Name:     "string",
										Type:     &cty.String,
										Required: true,
									},
								},
							},
						},
						{
							TypeName: "single",
							Nesting:  schema.SchemaNestedBlockNestingModeSingle,
							Block: &schema.SchemaBlock{
								Attributes: []*schema.SchemaAttribute{
									{
										Name:     "string",
										Type:     &cty.String,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if !cmp.Equal(got, want, equateEmpty, typeComparer) {
		t.Error(cmp.Diff(got, want, equateEmpty, typeComparer))
	}
}

func ToPtr[T any](v T) *T {
	return &v
}
