package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type UserModel struct {
	Uid          types.Int64  `tfsdk:"uid"`
	FirstName    types.String `tfsdk:"first_name"`
	LastName     types.String `tfsdk:"last_name"`
	Username     types.String `tfsdk:"username"`
	Email        types.String `tfsdk:"email"`
	UserImageURL types.String `tfsdk:"user_image_url"`
}

type TagModel struct {
	Tid  types.Int64  `tfsdk:"tid"`
	Name types.String `tfsdk:"name"`
}

type RedirectHitsModel struct {
	Daily   types.Int64 `tfsdk:"daily"`
	Weekly  types.Int64 `tfsdk:"weekly"`
	Monthly types.Int64 `tfsdk:"monthly"`
	Alltime types.Int64 `tfsdk:"alltime"`
}

type GeolinkModel struct {
	Location types.String `tfsdk:"location"`
	URL      types.String `tfsdk:"url"`
}

var UserAttrTypes = map[string]attr.Type{
	"uid":            types.Int64Type,
	"first_name":     types.StringType,
	"last_name":      types.StringType,
	"username":       types.StringType,
	"email":          types.StringType,
	"user_image_url": types.StringType,
}

var TagAttrTypes = map[string]attr.Type{
	"tid":  types.Int64Type,
	"name": types.StringType,
}

var RedirectHitsAttrTypes = map[string]attr.Type{
	"daily":   types.Int64Type,
	"weekly":  types.Int64Type,
	"monthly": types.Int64Type,
	"alltime": types.Int64Type,
}

var GeolinkAttrTypes = map[string]attr.Type{
	"location": types.StringType,
	"url":      types.StringType,
}
