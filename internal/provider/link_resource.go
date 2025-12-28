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
	_ resource.ResourceWithModifyPlan  = &linkResource{}
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
				Computed:    true,
				Description: "The timestamp of the last update to the golink.",
			},
			"variable_link": schema.BoolAttribute{
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Indicates if the link is a variable link.",
			},
			"pinned": schema.BoolAttribute{
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Indicates if the link is pinned.",
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
				Description: "If true, the link is unlisted. Private links are always unlisted.",
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
				Attributes:  UserResourceSchemaAttributes,
			},
		},
	}
}

func (r *linkResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan linkResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config linkResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	planUpdated := false

	privateKnown := !plan.Private.IsNull() && !plan.Private.IsUnknown()
	privateTrue := privateKnown && plan.Private.ValueBool()

	if privateTrue {
		configUnlistedKnown := !config.Unlisted.IsNull() && !config.Unlisted.IsUnknown()
		if configUnlistedKnown && !config.Unlisted.ValueBool() {
			resp.Diagnostics.AddAttributeError(
				path.Root("unlisted"),
				"Private Links Must Be Unlisted",
				"When `private` is true, the GoLinks API always sets `unlisted` to true. Please remove the explicit `unlisted = false` configuration or set it to true.",
			)
			return
		}

		if plan.Unlisted.IsNull() || plan.Unlisted.IsUnknown() || !plan.Unlisted.ValueBool() {
			plan.Unlisted = types.BoolValue(true)
			planUpdated = true
		}
	}

	hyphensKnown := !plan.Hyphens.IsNull() && !plan.Hyphens.IsUnknown()
	hyphensTrue := hyphensKnown && plan.Hyphens.ValueBool()

	if hyphensTrue {
		resp.Diagnostics.AddAttributeError(
			path.Root("hyphens"),
			"Hyphens Not Supported",
			"The GoLinks API currently ignores the `hyphens` flag, so the provider cannot manage it. Remove this attribute from configuration.",
		)
		return
	}

	formatKnown := !plan.Format.IsNull() && !plan.Format.IsUnknown()
	formatTrue := formatKnown && plan.Format.ValueBool()
	if formatTrue {
		resp.Diagnostics.AddAttributeError(
			path.Root("format"),
			"Format Not Supported",
			"The GoLinks API currently ignores the `format` flag, so the provider cannot manage it. Remove this attribute from configuration.",
		)
		return
	}

	if planUpdated {
		diags = resp.Plan.Set(ctx, plan)
		resp.Diagnostics.Append(diags...)
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

	privateVal := !plan.Private.IsNull() && !plan.Private.IsUnknown() && plan.Private.ValueBool()
	publicVal := !plan.Public.IsNull() && !plan.Public.IsUnknown() && plan.Public.ValueBool()
	formatVal := !plan.Format.IsNull() && !plan.Format.IsUnknown() && plan.Format.ValueBool()
	hyphensVal := !plan.Hyphens.IsNull() && !plan.Hyphens.IsUnknown() && plan.Hyphens.ValueBool()
	unlistedVal := !plan.Unlisted.IsNull() && !plan.Unlisted.IsUnknown() && plan.Unlisted.ValueBool()

	if privateVal && !unlistedVal {
		unlistedVal = true
		plan.Unlisted = types.BoolValue(true)
	}

	link.Unlisted = BoolToInt(unlistedVal)
	link.Private = BoolToInt(privateVal)
	link.Public = BoolToInt(publicVal)
	link.Hyphens = BoolToInt(hyphensVal)
	link.Format = BoolToInt(formatVal)

	link.Tags = append(link.Tags, plan.Tags...)

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

	privateVal := !plan.Private.IsNull() && !plan.Private.IsUnknown() && plan.Private.ValueBool()
	publicVal := !plan.Public.IsNull() && !plan.Public.IsUnknown() && plan.Public.ValueBool()
	formatVal := !plan.Format.IsNull() && !plan.Format.IsUnknown() && plan.Format.ValueBool()
	hyphensVal := !plan.Hyphens.IsNull() && !plan.Hyphens.IsUnknown() && plan.Hyphens.ValueBool()
	unlistedVal := !plan.Unlisted.IsNull() && !plan.Unlisted.IsUnknown() && plan.Unlisted.ValueBool()

	if privateVal && !unlistedVal {
		unlistedVal = true
		plan.Unlisted = types.BoolValue(true)
	}

	link.Unlisted = BoolToInt(unlistedVal)
	link.Private = BoolToInt(privateVal)
	link.Public = BoolToInt(publicVal)
	link.Format = BoolToInt(formatVal)
	link.Hyphens = BoolToInt(hyphensVal)

	var tags []string
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

	// Log the link structure
	linkJSON, _ := json.MarshalIndent(link, "", "  ")
	tflog.Debug(ctx, "UpdateLinkRequest structure", map[string]interface{}{
		"link": string(linkJSON),
	})
	tflog.Debug(ctx, strconv.FormatInt(plan.Gid.ValueInt64(), 10))

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
