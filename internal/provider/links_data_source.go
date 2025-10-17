package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/panoplytechnology/golinks-client-go"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &linksDataSource{}
	_ datasource.DataSourceWithConfigure = &linksDataSource{}
)

// LinksDataSource is a helper function to simplify the provider implementation.
func NewLinksDataSource() datasource.DataSource {
	return &linksDataSource{}
}

// linksDataSource is the data source implementation.
type linksDataSource struct {
	client *golinks.Client
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
	Tags         []tagModel   `tfsdk:"tags"`
	Unlisted     types.Int64  `tfsdk:"unlisted"`
	VariableLink types.Int64  `tfsdk:"variable_link"`
	Pinned       types.Int64  `tfsdk:"pinned"`
	RedirectHits types.Object `tfsdk:"redirect_hits"`
	CreatedAt    types.Int64  `tfsdk:"created_at"`
	UpdatedAt    types.Int64  `tfsdk:"updated_at"`
}

// userModel maps user data.
type userModel struct {
	Uid          types.Int64  `tfsdk:"uid"`
	FirstName    types.String `tfsdk:"first_name"`
	LastName     types.String `tfsdk:"last_name"`
	Username     types.String `tfsdk:"username"`
	Email        types.String `tfsdk:"email"`
	UserImageURL types.String `tfsdk:"user_image_url"`
}

// tagModel maps tag data.
type tagModel struct {
	Tid  types.Int64  `tfsdk:"tid"`
	Name types.String `tfsdk:"name"`
}

// redirectHitsModel maps redirect hits data.
type redirectHitsModel struct {
	Daily   types.Int64 `tfsdk:"daily"`
	Weekly  types.Int64 `tfsdk:"weekly"`
	Monthly types.Int64 `tfsdk:"monthly"`
	Alltime types.Int64 `tfsdk:"alltime"`
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
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"uid": schema.Int64Attribute{
									Computed: true,
								},
								"first_name": schema.StringAttribute{
									Computed: true,
								},
								"last_name": schema.StringAttribute{
									Computed: true,
								},
								"username": schema.StringAttribute{
									Computed: true,
								},
								"email": schema.StringAttribute{
									Computed: true,
								},
								"user_image_url": schema.StringAttribute{
									Computed: true,
								},
							},
						},
						"url": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"tags": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"tid": schema.Int64Attribute{
										Computed: true,
									},
									"name": schema.StringAttribute{
										Computed: true,
									},
								},
							},
						},
						"unlisted": schema.Int64Attribute{
							Computed: true,
						},
						"variable_link": schema.Int64Attribute{
							Computed: true,
						},
						"pinned": schema.Int64Attribute{
							Computed: true,
						},
						"redirect_hits": schema.SingleNestedAttribute{
							Computed: true,
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
							Computed: true,
						},
						"updated_at": schema.Int64Attribute{
							Computed: true,
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

	golinksResp, err := d.client.GetGolinks()
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
		// Map user
		userObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
			"uid":            types.Int64Type,
			"first_name":     types.StringType,
			"last_name":      types.StringType,
			"username":       types.StringType,
			"email":          types.StringType,
			"user_image_url": types.StringType,
		}, userModel{
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

		// Map redirect hits
		redirectHitsObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
			"daily":   types.Int64Type,
			"weekly":  types.Int64Type,
			"monthly": types.Int64Type,
			"alltime": types.Int64Type,
		}, redirectHitsModel{
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
			Unlisted:     types.Int64Value(golink.Unlisted),
			VariableLink: types.Int64Value(golink.VariableLink),
			Pinned:       types.Int64Value(golink.Pinned),
			RedirectHits: redirectHitsObj,
			CreatedAt:    types.Int64Value(golink.CreatedAt),
			UpdatedAt:    types.Int64Value(golink.UpdatedAt),
		}

		// Map tags
		for _, tag := range golink.Tags {
			golinkState.Tags = append(golinkState.Tags, tagModel{
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

	client, ok := req.ProviderData.(*golinks.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *golinks.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}
