// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"terraform-provider-golinks/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &linkDataSource{}
	_ datasource.DataSourceWithConfigure = &linkDataSource{}
)

// LinksDataSource is a helper function to simplify the provider implementation.
func LinkDataSource() datasource.DataSource {
	return &linkDataSource{}
}

// linkDataSource is the data source implementation.
type linkDataSource struct {
	client *client.Client
}

type linkDataSourceModel struct {
	Name         types.String `tfsdk:"name"`
	Gid          types.Int64  `tfsdk:"gid"`
	Cid          types.Int64  `tfsdk:"cid"`
	User         types.Object `tfsdk:"user"`
	URL          types.String `tfsdk:"url"`
	Description  types.String `tfsdk:"description"`
	Tags         []TagModel   `tfsdk:"tags"`
	Unlisted     types.Int64  `tfsdk:"unlisted"`
	VariableLink types.Int64  `tfsdk:"variable_link"`
	Pinned       types.Int64  `tfsdk:"pinned"`
	RedirectHits types.Object `tfsdk:"redirect_hits"`
	CreatedAt    types.Int64  `tfsdk:"created_at"`
	UpdatedAt    types.Int64  `tfsdk:"updated_at"`
}

// Metadata returns the data source type name.
func (d *linkDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_link"
}

// Schema defines the schema for the data source.
func (d *linkDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a single GoLink by name.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the GoLink to retrieve.",
			},
			"gid": schema.Int64Attribute{
				Computed:    true,
				Description: "The GoLink ID returned by the API.",
			},
			"cid": schema.Int64Attribute{
				Computed:    true,
				Description: "The company ID for the GoLink.",
			},
			"url": schema.StringAttribute{
				Computed:    true,
				Description: "The destination URL the GoLink points to.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "The description of the GoLink.",
			},
			"user": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The user who created the GoLink.",
				Attributes:  UserDataSourceSchemaAttributes,
			},
			"tags": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The tags associated with the GoLink.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tid": schema.Int64Attribute{
							Computed:    true,
							Description: "The unique ID of the tag.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the tag.",
						},
					},
				},
			},
			"unlisted": schema.Int64Attribute{
				Computed:    true,
				Description: "Indicates if the GoLink is unlisted.",
			},
			"variable_link": schema.Int64Attribute{
				Computed:    true,
				Description: "Indicates if the GoLink is a variable link.",
			},
			"pinned": schema.Int64Attribute{
				Computed:    true,
				Description: "Indicates if the GoLink is pinned.",
			},
			"redirect_hits": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The redirect hits statistics for the GoLink.",
				Attributes: map[string]schema.Attribute{
					"daily": schema.Int64Attribute{
						Computed: true,
					},
					"weekly": schema.Int64Attribute{
						Computed: true,
					},
					"monthly": schema.Int64Attribute{
						Computed: true,
					},
					"alltime": schema.Int64Attribute{
						Computed: true,
					},
				},
			},
			"created_at": schema.Int64Attribute{
				Computed:    true,
				Description: "The timestamp when the GoLink was created.",
			},
			"updated_at": schema.Int64Attribute{
				Computed:    true,
				Description: "The timestamp when the GoLink was last updated.",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *linkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state linkDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Name.IsNull() || state.Name.IsUnknown() || state.Name.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing GoLink Name",
			"The data source requires the `name` attribute to identify which GoLink to retrieve.",
		)
		return
	}

	golinksResp, err := d.client.GetGolinksByName(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read GoLink",
			err.Error(),
		)
		return
	}

	golink := golinksResp

	userObj, diags := types.ObjectValueFrom(ctx, UserAttrTypes, UserModel{
		Uid:          types.Int64Value(golink.User.Uid),
		FirstName:    types.StringValue(golink.User.FirstName),
		LastName:     types.StringValue(golink.User.LastName),
		Username:     types.StringValue(golink.User.Username),
		Email:        types.StringValue(golink.User.Email),
		UserImageURL: types.StringValue(golink.User.UserImageURL),
	})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	redirectHitsObj, diags := types.ObjectValueFrom(ctx, RedirectHitsAttrTypes, RedirectHitsModel{
		Daily:   types.Int64Value(golink.RedirectHits.Daily),
		Weekly:  types.Int64Value(golink.RedirectHits.Weekly),
		Monthly: types.Int64Value(golink.RedirectHits.Monthly),
		Alltime: types.Int64Value(golink.RedirectHits.Alltime),
	})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tags := make([]TagModel, 0, len(golink.Tags))
	for _, tag := range golink.Tags {
		tags = append(tags, TagModel{
			Tid:  types.Int64Value(tag.Tid),
			Name: types.StringValue(tag.Name),
		})
	}

	state.Gid = types.Int64Value(golink.Gid)
	state.Cid = types.Int64Value(golink.Cid)
	state.URL = types.StringValue(golink.URL)
	state.Name = types.StringValue(golink.Name)
	state.Description = types.StringValue(golink.Description)
	state.Unlisted = types.Int64Value(int64(golink.Unlisted))
	state.VariableLink = types.Int64Value(int64(golink.VariableLink))
	state.Pinned = types.Int64Value(int64(golink.Pinned))
	state.RedirectHits = redirectHitsObj
	state.CreatedAt = types.Int64Value(golink.CreatedAt)
	state.UpdatedAt = types.Int64Value(golink.UpdatedAt)
	state.User = userObj
	state.Tags = tags

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Configure adds the provider configured client to the data source.
func (d *linkDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *golinks.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}
