package fw

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/magodo/tfpluginschema/schema"
)

func FromProvider(p provider.Provider) (*schema.ProviderSchema, error) {
	ctx := context.Background()

	var providerMetadataResp provider.MetadataResponse
	p.Metadata(ctx, provider.MetadataRequest{}, &providerMetadataResp)

	var providerSchemaResp provider.SchemaResponse
	p.Schema(ctx, provider.SchemaRequest{}, &providerSchemaResp)
	if providerSchemaResp.Diagnostics.HasError() {
		return nil, fmt.Errorf("getting provider schema: %#v", providerSchemaResp.Diagnostics)
	}

	var resources []resource.Resource
	for _, rf := range p.Resources(ctx) {
		resources = append(resources, rf())
	}

	var datasources []datasource.DataSource
	for _, df := range p.DataSources(ctx) {
		datasources = append(datasources, df())
	}

	providerSchema, err := ProviderSchema(ctx, providerSchemaResp.Schema)
	if err != nil {
		return nil, fmt.Errorf("converting provider schema: %v", err)
	}

	ret := &schema.ProviderSchema{
		Provider:          providerSchema,
		ResourceSchemas:   map[string]*schema.Schema{},
		DataSourceSchemas: map[string]*schema.Schema{},
	}

	for _, res := range resources {
		var metadataResp resource.MetadataResponse
		res.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: providerMetadataResp.TypeName}, &metadataResp)

		var schemaResp resource.SchemaResponse
		res.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
		if schemaResp.Diagnostics.HasError() {
			return nil, fmt.Errorf("getting resource schema: %#v", schemaResp.Diagnostics)
		}
		sch, err := ResourceSchema(ctx, schemaResp.Schema)
		if err != nil {
			return nil, fmt.Errorf("converting resource schema (%s): %v", metadataResp.TypeName, err)
		}
		ret.ResourceSchemas[metadataResp.TypeName] = sch
	}
	for _, ds := range datasources {
		var metadataResp datasource.MetadataResponse
		ds.Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: providerMetadataResp.TypeName}, &metadataResp)

		var schemaResp datasource.SchemaResponse
		ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)
		if schemaResp.Diagnostics.HasError() {
			return nil, fmt.Errorf("getting datasource schema: %#v", schemaResp.Diagnostics)
		}
		sch, err := DatasourceSchema(ctx, schemaResp.Schema)
		if err != nil {
			return nil, fmt.Errorf("converting datasource schema (%s): %v", metadataResp.TypeName, err)
		}
		ret.DataSourceSchemas[metadataResp.TypeName] = sch
	}
	return ret, nil
}
