// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"terraform-provider-golinks/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &golinksProvider{}
)

// golinksProviderModel maps provider schema data to a Go type.
type golinksProviderModel struct {
	Token types.String `tfsdk:"token"`
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &golinksProvider{
			version: version,
		}
	}
}

// golinksProvider is the provider implementation.
type golinksProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *golinksProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "golinks"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *golinksProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Description: "API Token for authenticating with the GoLinks API.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

// Configure prepares a golinks API client for data sources and resources.
func (p *golinksProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	tflog.Info(ctx, "Configuring GoLinks client")

	// Retrieve provider data from configuration var config golinksProviderModel
	var config golinksProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown GoLinks API Host",
			"The provider cannot create the GoLinks API client as there is an unknown configuration value for the GoLinks API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the GOLINKS_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	token := os.Getenv("GOLINKS_TOKEN")

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing GoLinks API Host",
			"The provider cannot create the GoLinks API client as there is a missing or empty value for the GoLinks API host. "+
				"Set the token value in the configuration or use the GOLINKS_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "golinks_token", token)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "golinks_token")
	tflog.Debug(ctx, "Creating GoLinks client")

	// Create a new GoLinks client using the configuration values
	client, err := client.NewClient(ctx, &token)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create GoLinks API Client",
			"An unexpected error occurred when creating the GoLinks API client. "+
				"If TEST the error is not clear, please contact the provider developers.\n\n"+
				"GoLinks Client Error: "+err.Error(),
		)
		return
	}

	// Make the GoLinks client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured GoLinks client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *golinksProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		LinksDataSource,
		LinkDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *golinksProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewLinkResource,
	}
}
