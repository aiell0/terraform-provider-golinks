package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"terraform-provider-golinks/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &linkResource{}
	_ resource.ResourceWithConfigure   = &linkResource{}
	_ resource.ResourceWithImportState = &linkResource{}
)

// linksResource is the resource implementation.
type linkResource struct {
	client *client.Client
}

// golinkResourceModel maps the resource schema data.
type linkResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Gid          types.Int64  `tfsdk:"gid"`
	Cid          types.Int64  `tfsdk:"cid"`
	LastUpdated  types.String `tfsdk:"last_updated"`
	User         types.Object `tfsdk:"user"`
	URL          types.String `tfsdk:"url"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Tags         []string     `tfsdk:"tags"`
	Unlisted     types.Bool   `tfsdk:"unlisted"`
	Private      types.Bool   `tfsdk:"private"`
	Public       types.Bool   `tfsdk:"public"`
	VariableLink types.Bool   `tfsdk:"variable_link"`
	Pinned       types.Bool   `tfsdk:"pinned"`
	Format       types.Bool   `tfsdk:"format"`
	Hyphens      types.Bool   `tfsdk:"hyphens"`
	Aliases      types.List   `tfsdk:"aliases"`
	Geolinks     types.List   `tfsdk:"geolinks"`
	CreatedAt    types.Int64  `tfsdk:"created_at"`
	UpdatedAt    types.Int64  `tfsdk:"updated_at"`
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
			"variable_link": schema.BoolAttribute{
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"pinned": schema.BoolAttribute{
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"cid": schema.Int64Attribute{
				Computed:    true,
				Description: "The Company ID.",
			},
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
			"unlisted": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "If true, the link is unlisted. If false (default), shared with everyone in your organization.",
				Default:     booldefault.StaticBool(false),
			},
			"private": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "If true, the link is private. Links cannot change to or from private after creation.",
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"public": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "If true, the link can be accessed by people outside of your organization.",
				Default:     booldefault.StaticBool(false),
			},
			"format": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "If the value is true, invalid characters (e.g. punctuation) will be removed from the created go link name.",
				Default:     booldefault.StaticBool(false),
			},
			"hyphens": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "If the value is true, spaces will be replaced with hyphens in the go link name. If false, spaces will be removed. Requires format set to true.",
				Default:     booldefault.StaticBool(false),
			},
			"aliases": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Create multiple names for the same link with aliases.",
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
			"user": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The user who created the golink.",
				Attributes: map[string]schema.Attribute{
					"uid": schema.Int64Attribute{
						Computed:    true,
						Description: "The user ID.",
					},
					"first_name": schema.StringAttribute{
						Computed:    true,
						Description: "The user's first name.",
					},
					"last_name": schema.StringAttribute{
						Computed:    true,
						Description: "The user's last name.",
					},
					"username": schema.StringAttribute{
						Computed:    true,
						Description: "The user's username.",
					},
					"email": schema.StringAttribute{
						Computed:    true,
						Description: "The user's email address.",
					},
					"user_image_url": schema.StringAttribute{
						Computed:    true,
						Description: "URL to the user's profile image.",
					},
				},
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
	var link client.CreateLinkRequest
	link.URL = plan.URL.ValueString()
	link.Name = plan.Name.ValueString()
	link.Description = plan.Description.ValueString()
	link.Unlisted = BoolToInt(plan.Unlisted.ValueBool())
	link.Private = BoolToInt(plan.Private.ValueBool())
	link.Public = BoolToInt(plan.Public.ValueBool())
	link.Hyphens = BoolToInt(plan.Hyphens.ValueBool())
	link.Format = BoolToInt(plan.Format.ValueBool())

	var tags []string
	// for _, t := range plan.Tags {
	// 	tags = append(tags, t)
	// }
	// link.Tags = tags
	tags = append(tags, plan.Tags...)
	link.Tags = tags

	var aliases []string
	if !plan.Aliases.IsNull() && !plan.Aliases.IsUnknown() {
		diags := plan.Aliases.ElementsAs(ctx, &aliases, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	link.Aliases = aliases

	var geolinks []client.Geolink
	if !plan.Geolinks.IsNull() && !plan.Geolinks.IsUnknown() {
		var geolinkModels []GeolinkModel
		diags := plan.Geolinks.ElementsAs(ctx, &geolinkModels, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, gl := range geolinkModels {
			geolinks = append(geolinks, client.Geolink{
				Location: gl.Location.ValueString(),
				URL:      gl.URL.ValueString(),
			})
		}
	}
	link.Geolinks = geolinks

	// Create new link
	linkresponse, err := r.client.CreateLink(ctx, link)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating link",
			"Could not create link, unexpected error: "+err.Error(),
		)
		return
	}

	MapLinkResponseToModel(linkresponse, &plan, true)

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

	linkresponse, err := r.client.GetLink(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving link",
			"Could not get link, unexpected error: "+err.Error(),
		)
		return
	}

	MapLinkResponseToModel(linkresponse, &state, false)

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

	var link client.UpdateLinkRequest
	link.Gid = state.Gid.ValueInt64()
	link.URL = plan.URL.ValueString()
	link.Name = plan.Name.ValueString()
	link.Description = plan.Description.ValueString()

	var tags []string
	tags = append(tags, plan.Tags...)
	link.Tags = tags

	// Log the link structure
	linkJSON, _ := json.MarshalIndent(link, "", "  ")
	tflog.Debug(ctx, "UpdateLinkRequest structure", map[string]interface{}{
		"link": string(linkJSON),
	})
	tflog.Debug(ctx, strconv.FormatInt(plan.Gid.ValueInt64(), 10))

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
	_, err := r.client.UpdateLink(ctx, link)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Golink",
			"Could not update link, unexpected error: "+err.Error(),
		)
		return
	}

	linkresponse, err := r.client.GetLink(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving link",
			"Could not get link, unexpected error: "+err.Error(),
		)
		return
	}

	MapLinkResponseToModel(linkresponse, &plan, true)

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
	err := r.client.DeleteLink(ctx, state.Gid.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Golink",
			"Could not delete link, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *linkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Configure adds the provider configured client to the resource.
func (r *linkResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}
