package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/svetob/terraform-provider-httpdownload/internal/client"
)

var (
	_ provider.Provider = &httpDownloadProvider{}
)

func New() provider.Provider {
	return &httpDownloadProvider{}
}

type httpDownloadProvider struct{}

func (p *httpDownloadProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "httpdownload"
}

func (p *httpDownloadProvider) Schema(_ context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Http download provider",
	}
}

func (p *httpDownloadProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	c := client.NewHttpDownloadClient()

	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *httpDownloadProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewHttpDownloadFileResource,
	}
}

func (p *httpDownloadProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
