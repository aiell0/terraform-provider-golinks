package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/panoplytechnology/golinks-client-go"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &linkResource{}
	_ resource.ResourceWithConfigure = &linkResource{}
)

// linksResource is the resource implementation.
type linkResource struct {
	client *golinks.Client
}

// golinkResourceModel maps the resource schema data.
type linkResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Gid         types.Int64  `tfsdk:"gid"`
	Cid         types.Int64  `tfsdk:"cid"`
	LastUpdated types.String `tfsdk:"last_updated"`
	// User         UserModel         `tfsdk:"user"`
	URL          types.String `tfsdk:"url"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Tags         []TagModel   `tfsdk:"tags"`
	Unlisted     types.Int64  `tfsdk:"unlisted"`
	Private      types.Int64  `tfsdk:"private"`
	Public       types.Int64  `tfsdk:"public"`
	VariableLink types.Int64  `tfsdk:"variable_link"`
	Pinned       types.Int64  `tfsdk:"pinned"`
	// RedirectHits RedirectHitsModel `tfsdk:"redirect_hits"`
	Aliases    types.List     `tfsdk:"aliases"`
	Multilinks types.List     `tfsdk:"multilinks"`
	Geolinks   []GeolinkModel `tfsdk:"geolinks"`
	CreatedAt  types.Int64    `tfsdk:"created_at"`
	UpdatedAt  types.Int64    `tfsdk:"updated_at"`
}

type UserModel struct {
	Uid          types.Int64  `tfsdk:"uid"`
	firstName    types.String `tfsdk:"first_name"`
	lastName     types.String `tfsdk:"last_name"`
	Username     types.String `tfsdk:"username"`
	Email        types.String `tfsdk:"email"`
	UserImageUrl types.String `tfsdk:"user_image_url"`
}

type RedirectHitsModel struct {
	Daily   types.Int64 `tfsdk:"daily"`
	Weekly  types.Int64 `tfsdk:"weekly"`
	Monthly types.Int64 `tfsdk:"monthly"`
	AllTime types.Int64 `tfsdk:"all_time"`
}

type TagModel struct {
	Tid  types.Int64  `tfsdk:"tid"`
	Name types.String `tfsdk:"name"`
}

// geolinkModel maps geolink data.
type GeolinkModel struct {
	Location string
	URL      string
}

// NewLinksResource is a helper function to simplify the provider implementation.
func NewLinkResource() resource.Resource {
	return &linkResource{}
}

// Metadata returns the resource type name.
func (r *linkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_link"
}

// Schema defines the schema for the resource.
func (r *linkResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a GoLink.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"gid": schema.Int64Attribute{
				Computed:    true,
				Description: "The GoLink ID returned by the API.",
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"cid": schema.Int64Attribute{
				Computed:    true,
				Description: "The Company ID.",
			},
			// "user": schema.ListNestedAttribute{
			// 	Computed:    true,
			// 	Description: "User information.",
			// 	NestedObject: schema.NestedAttributeObject{
			// 		Attributes: map[string]schema.Attribute{
			// 			"uid": schema.Int64Attribute{
			// 				Computed:    true,
			// 				Description: "The unique user id.",
			// 			},
			// 			"first_name": schema.StringAttribute{
			// 				Computed:    true,
			// 				Description: "The first name of the user who owns the link (if found)",
			// 			},
			// 			"last_name": schema.StringAttribute{
			// 				Computed:    true,
			// 				Description: "The last name of the user who owns the link (if found)",
			// 			},
			// 			"username": schema.StringAttribute{
			// 				Computed:    true,
			// 				Description: "The username of the user who owns the link (if found)",
			// 			},
			// 			"email": schema.StringAttribute{
			// 				Computed:    true,
			// 				Description: "The email of the user who owns the link.",
			// 			},
			// 			"user_image_url": schema.StringAttribute{
			// 				Computed:    true,
			// 				Description: "The gravatar of the user who owns the link.",
			// 			},
			// 		},
			// },
			// },
			"url": schema.StringAttribute{
				Required:    true,
				Description: "The destination URL.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The link name.",
			},
			"description": schema.StringAttribute{
				Required:    true,
				Description: "Brief description of the link.",
			},
			"tags": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Organize your golinks and find the right ones quickly with tags.",
			},
			"unlisted": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "If 1, the link is unlisted. If 0 (default), shared with everyone in your organization.",
				Default:     int64default.StaticInt64(0),
			},
			"private": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "If 1, the link is private. Links cannot change to or from private after creation.",
				Default:     int64default.StaticInt64(0),
			},
			"public": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "If 1, the link can be accessed by people outside of your organization.",
				Default:     int64default.StaticInt64(0),
			},
			"variable_link": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Denotes if the link is a variable link.",
				Default:     int64default.StaticInt64(0),
			},
			"pinned": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Denotes if the link is pinned to the top of your GoLinks feed.",
				Default:     int64default.StaticInt64(0),
			},
			// "redirect_hits": schema.ListNestedAttribute{
			// 	Computed:    true,
			// 	Description: "The number of redirects for the golink.",
			// 	NestedObject: schema.NestedAttributeObject{
			// 		Attributes: map[string]schema.Attribute{
			// 			"daily": schema.Int64Attribute{
			// 				Computed:    true,
			// 				Description: "The number of daily redirects for the golink.",
			// 			},
			// 			"weekly": schema.Int64Attribute{
			// 				Computed:    true,
			// 				Description: "The number of weekly redirects for the golink.",
			// 			},
			// 			"monthly": schema.Int64Attribute{
			// 				Computed:    true,
			// 				Description: "The number of monthly redirects for the golink.",
			// 			},
			// 			"alltime": schema.Int64Attribute{
			// 				Computed:    true,
			// 				Description: "The number of lifetime redirects for the golink.",
			// 			},
			// 		},
			// 	},
			// },
			"aliases": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Create multiple names for the same link with aliases.",
			},
			"multilinks": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "The list of target links if the link is a multi link.",
			},
			"geolinks": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Create different destinations for a link depending on current location.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"location": schema.StringAttribute{
							Required:    true,
							Description: "Two-character ISO country code or 'US-XX' for US states.",
						},
						"url": schema.StringAttribute{
							Required:    true,
							Description: "The destination URL for this location.",
						},
					},
				},
			},
			"created_at": schema.Int64Attribute{
				Computed:    true,
				Description: "Unix timestamp when the golink was created.",
			},
			"updated_at": schema.Int64Attribute{
				Computed:    true,
				Description: "Unix timestamp when the golink was last updated.",
			},
		},
	}
}

// Create a new resource.
func (r *linkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan linkResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var link golinks.CreateLink
	// link.Uid = plan.User.Uid.ValueInt64()
	link.URL = plan.URL.ValueString()
	link.Name = plan.Name.ValueString()
	link.Description = plan.Description.ValueString()
	link.Unlisted = plan.Unlisted.ValueInt64()
	link.Private = plan.Private.ValueInt64()
	link.Public = plan.Public.ValueInt64()
	link.Format = 0
	link.Hyphens = 0

	// Generate API request body from plan
	var tags []golinks.Tag
	for _, t := range plan.Tags {
		tags = append(tags, golinks.Tag{
			Tid:  t.Tid.ValueInt64(),
			Name: t.Name.String(),
		})
	}

	var aliases []string
	if !plan.Aliases.IsNull() && !plan.Aliases.IsUnknown() {
		diags := plan.Aliases.ElementsAs(ctx, &aliases, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var geolinks []golinks.Geolink
	for _, gl := range plan.Geolinks {
		geolinks = append(geolinks, golinks.Geolink{
			Location: gl.Location,
			URL:      gl.URL,
		})
	}

	// Create new link
	linkresponse, err := r.client.CreateLink(link)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating link",
			"Could not create link, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(strconv.FormatInt(linkresponse.Gid, 10))
	plan.Gid = types.Int64Value(linkresponse.Gid)
	plan.Cid = types.Int64Value(linkresponse.Cid)
	plan.URL = types.StringValue(linkresponse.URL)
	plan.Name = types.StringValue(linkresponse.Name)
	plan.Description = types.StringValue(linkresponse.Description)
	plan.Unlisted = types.Int64Value(linkresponse.Unlisted)
	plan.VariableLink = types.Int64Value(linkresponse.VariableLink)
	plan.Pinned = types.Int64Value(linkresponse.Pinned)
	plan.CreatedAt = types.Int64Value(linkresponse.CreatedAt)
	plan.UpdatedAt = types.Int64Value(linkresponse.UpdatedAt)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	// plan.User = UserModel{
	// 	Uid:          types.Int64Value(linkresponse.User.Uid),
	// 	firstName:    types.StringValue(linkresponse.User.FirstName),
	// 	lastName:     types.StringValue(linkresponse.User.LastName),
	// 	Username:     types.StringValue(linkresponse.User.Username),
	// 	Email:        types.StringValue(linkresponse.User.Email),
	// 	UserImageUrl: types.StringValue(linkresponse.User.UserImageURL),
	// }
	//
	for index, tag := range linkresponse.Tags {
		plan.Tags[index] = TagModel{
			Tid:  types.Int64Value(tag.Tid),
			Name: types.StringValue(tag.Name),
		}
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *linkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from plan
	var state linkResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	linkresponse, err := r.client.GetLink(state.Gid.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving link",
			"Could not ge link, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	state.ID = types.StringValue(strconv.FormatInt(linkresponse.Gid, 10))
	state.Gid = types.Int64Value(linkresponse.Gid)
	state.Cid = types.Int64Value(linkresponse.Cid)
	state.URL = types.StringValue(linkresponse.URL)
	state.Name = types.StringValue(linkresponse.Name)
	state.Description = types.StringValue(linkresponse.Description)
	state.Unlisted = types.Int64Value(linkresponse.Unlisted)
	state.VariableLink = types.Int64Value(linkresponse.VariableLink)
	state.Pinned = types.Int64Value(linkresponse.Pinned)
	state.CreatedAt = types.Int64Value(linkresponse.CreatedAt)
	state.UpdatedAt = types.Int64Value(linkresponse.UpdatedAt)
	// plan.User = UserModel{
	// 	Uid:          types.Int64Value(linkresponse.User.Uid),
	// 	firstName:    types.StringValue(linkresponse.User.FirstName),
	// 	lastName:     types.StringValue(linkresponse.User.LastName),
	// 	Username:     types.StringValue(linkresponse.User.Username),
	// 	Email:        types.StringValue(linkresponse.User.Email),
	// 	UserImageUrl: types.StringValue(linkresponse.User.UserImageURL),
	// }
	//
	for index, tag := range linkresponse.Tags {
		state.Tags[index] = TagModel{
			Tid:  types.Int64Value(tag.Tid),
			Name: types.StringValue(tag.Name),
		}
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *linkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan linkResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state linkResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var link golinks.CreateLink
	// link.Uid = plan.User.Uid.ValueInt64()
	link.Gid = state.Gid.ValueInt64()
	link.URL = plan.URL.ValueString()
	link.Name = plan.Name.ValueString()
	link.Description = plan.Description.ValueString()

	// Log the link structure
	linkJSON, _ := json.MarshalIndent(link, "", "  ")
	tflog.Debug(ctx, "CreateLink structure", map[string]interface{}{
		"link": string(linkJSON),
	})
	tflog.Debug(ctx, strconv.FormatInt(plan.Gid.ValueInt64(), 10))

	// Generate API request body from plan
	// var tags []golinks.Tag
	// for _, t := range plan.Tags {
	// 	tags = append(tags, golinks.Tag{
	// 		Tid:  t.Tid.ValueInt64(),
	// 		Name: t.Name.String(),
	// 	})
	// }
	//
	// var aliases []string
	// if !plan.Aliases.IsNull() && !plan.Aliases.IsUnknown() {
	// 	diags := plan.Aliases.ElementsAs(ctx, &aliases, false)
	// 	resp.Diagnostics.Append(diags...)
	// 	if resp.Diagnostics.HasError() {
	// 		return
	// 	}
	// }
	//
	// var geolinks []golinks.Geolink
	// for _, gl := range plan.Geolinks {
	// 	geolinks = append(geolinks, golinks.Geolink{
	// 		Location: gl.Location,
	// 		URL:      gl.URL,
	// 	})
	// }

	// Update link
	_, err := r.client.UpdateLink(link)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Golink",
			"Could not update link, unexpected error: "+err.Error(),
		)
		return
	}

	linkresponse, err := r.client.GetLink(state.Gid.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving link",
			"Could not get link, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(strconv.FormatInt(linkresponse.Gid, 10))
	plan.Gid = types.Int64Value(linkresponse.Gid)
	plan.Cid = types.Int64Value(linkresponse.Cid)
	plan.URL = types.StringValue(linkresponse.URL)
	plan.Name = types.StringValue(linkresponse.Name)
	plan.Description = types.StringValue(linkresponse.Description)
	plan.Unlisted = types.Int64Value(linkresponse.Unlisted)
	plan.VariableLink = types.Int64Value(linkresponse.VariableLink)
	plan.Pinned = types.Int64Value(linkresponse.Pinned)
	plan.CreatedAt = types.Int64Value(linkresponse.CreatedAt)
	plan.UpdatedAt = types.Int64Value(linkresponse.UpdatedAt)

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *linkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state linkResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	err := r.client.DeleteLink(state.Gid.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Golink",
			"Could not delete link, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *linkResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}
