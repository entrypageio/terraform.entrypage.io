package provider

import (
	"context"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	providerschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const defaultBaseURL = "https://api.entrypage.io"

// Ensure the provider implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &entrypageProvider{}
)

// entrypageProvider is the provider implementation.
type entrypageProvider struct {
	version string
}

// New returns a new provider instance.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &entrypageProvider{
			version: version,
		}
	}
}

// entrypageProviderModel maps provider schema data to a Go type.
type entrypageProviderModel struct {
	APIKey  types.String `tfsdk:"api_key"`
	BaseURL types.String `tfsdk:"base_url"`
}

// Metadata returns the provider type name.
func (p *entrypageProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "entrypage"
	resp.Version = p.version
}

// Schema defines the provider-level configuration schema.
func (p *entrypageProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = providerschema.Schema{
		Attributes: map[string]providerschema.Attribute{
			"api_key": providerschema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				MarkdownDescription: "API key for entrypage.io. " +
					"Can also be set via the `ENTRYPAGE_API_KEY` environment variable.",
			},
			"base_url": providerschema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Base URL for the Entrypage API. Defaults to `https://api.entrypage.io`.",
			},
		},
	}
}

// Configure prepares a configured API client for data sources and resources.
func (p *entrypageProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config entrypageProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiKey := strings.TrimSpace(os.Getenv("ENTRYPAGE_API_KEY"))
	if !config.APIKey.IsNull() && !config.APIKey.IsUnknown() {
		apiKey = strings.TrimSpace(config.APIKey.ValueString())
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing API key",
			"Please provide an API key via the provider `api_key` attribute or the ENTRYPAGE_API_KEY environment variable.",
		)
		return
	}

	baseURL := defaultBaseURL
	if !config.BaseURL.IsNull() && !config.BaseURL.IsUnknown() && config.BaseURL.ValueString() != "" {
		baseURL = strings.TrimRight(config.BaseURL.ValueString(), "/")
	}

	client := newClient(apiKey, baseURL)

	resp.DataSourceData = client
	resp.ResourceData = client
}

// Resources defines provider resources.
func (p *entrypageProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewUserResource,
	}
}

// DataSources defines provider data sources.
func (p *entrypageProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}
