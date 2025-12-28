package provider

import (
	"context"
	"fmt"

	"terraform-provider-golinks/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &linksDataSource{}
	_ datasource.DataSourceWithConfigure = &linksDataSource{}
)

// LinksDataSource is a helper function to simplify the provider implementation.
func LinksDataSource() datasource.DataSource {
	return &linksDataSource{}
}

// linksDataSource is the data source implementation.
type linksDataSource struct {
	client *client.Client
}

// Metadata returns the data source type name.
func (d *linksDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_links"
}

// golinksDataSourceModel maps the data source schema data.
type golinksDataSourceModel struct {
	Metadata types.Object  `tfsdk:"metadata"`
	Results  []golinkModel `tfsdk:"results"`
}

// metadataModel maps metadata schema data.
type metadataModel struct {
	Limit        types.Int64  `tfsdk:"limit"`
	Offset       types.Int64  `tfsdk:"offset"`
	Count        types.Int64  `tfsdk:"count"`
	TotalResults types.Int64  `tfsdk:"total_results"`
	Links        types.Object `tfsdk:"links"`
}

// linksModel maps pagination links data.
type linksModel struct {
	Prev types.String `tfsdk:"prev"`
	Next types.String `tfsdk:"next"`
}

// golinkModel maps golink schema data.
type golinkModel struct {
	Gid          types.Int64  `tfsdk:"gid"`
	Cid          types.Int64  `tfsdk:"cid"`
	User         types.Object `tfsdk:"user"`
	URL          types.String `tfsdk:"url"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Tags         []TagModel   `tfsdk:"tags"`
	Unlisted     types.Int32  `tfsdk:"unlisted"`
	VariableLink types.Int32  `tfsdk:"variable_link"`
	Pinned       types.Int32  `tfsdk:"pinned"`
	RedirectHits types.Object `tfsdk:"redirect_hits"`
	CreatedAt    types.Int64  `tfsdk:"created_at"`
	UpdatedAt    types.Int64  `tfsdk:"updated_at"`
}

// Schema defines the schema for the data source.
func (d *linksDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"metadata": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"limit": schema.Int64Attribute{
						Computed: true,
					},
					"offset": schema.Int64Attribute{
						Computed: true,
					},
					"count": schema.Int64Attribute{
						Computed: true,
					},
					"total_results": schema.Int64Attribute{
						Computed: true,
					},
					"links": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"prev": schema.StringAttribute{
								Computed: true,
							},
							"next": schema.StringAttribute{
								Computed: true,
							},
						},
					},
				},
			},
			"results": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"gid": schema.Int64Attribute{
							Computed: true,
						},
						"cid": schema.Int64Attribute{
							Computed: true,
						},
						"user": schema.SingleNestedAttribute{
							Computed:    true,
							Attributes:  UserDataSourceSchemaAttributes,
							Description: "The user who created the GoLink.",
						},
						"url": schema.StringAttribute{
							Computed:    true,
							Description: "The destination URL the GoLink points to.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the GoLink.",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "The description of the GoLink.",
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
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *linksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state golinksDataSourceModel

	golinksResp, err := d.client.GetGolinks(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read GoLinks",
			err.Error(),
		)
		return
	}

	// Map metadata
	linksObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
		"prev": types.StringType,
		"next": types.StringType,
	}, linksModel{
		Prev: types.StringValue(golinksResp.Metadata.Links.Prev),
		Next: types.StringValue(golinksResp.Metadata.Links.Next),
	})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	metadataObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
		"limit":         types.Int64Type,
		"offset":        types.Int64Type,
		"count":         types.Int64Type,
		"total_results": types.Int64Type,
		"links": types.ObjectType{AttrTypes: map[string]attr.Type{
			"prev": types.StringType,
			"next": types.StringType,
		}},
	}, metadataModel{
		Limit:        types.Int64Value(golinksResp.Metadata.Limit),
		Offset:       types.Int64Value(golinksResp.Metadata.Offset),
		Count:        types.Int64Value(golinksResp.Metadata.Count),
		TotalResults: types.Int64Value(golinksResp.Metadata.TotalResults),
		Links:        linksObj,
	})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Metadata = metadataObj

	// Map response body to model
	for _, golink := range golinksResp.Results {
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

		golinkState := golinkModel{
			Gid:          types.Int64Value(golink.Gid),
			Cid:          types.Int64Value(golink.Cid),
			User:         userObj,
			URL:          types.StringValue(golink.URL),
			Name:         types.StringValue(golink.Name),
			Description:  types.StringValue(golink.Description),
			Unlisted:     types.Int32Value(golink.Unlisted),
			VariableLink: types.Int32Value(golink.VariableLink),
			Pinned:       types.Int32Value(golink.Pinned),
			RedirectHits: redirectHitsObj,
			CreatedAt:    types.Int64Value(golink.CreatedAt),
			UpdatedAt:    types.Int64Value(golink.UpdatedAt),
		}

		// Map tags
		for _, tag := range golink.Tags {
			golinkState.Tags = append(golinkState.Tags, TagModel{
				Tid:  types.Int64Value(tag.Tid),
				Name: types.StringValue(tag.Name),
			})
		}

		state.Results = append(state.Results, golinkState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *linksDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
