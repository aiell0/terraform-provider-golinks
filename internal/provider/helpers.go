package provider

import (
	"strconv"
	"time"

	"terraform-provider-golinks/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func IntToBool(i int32) bool {
	return i == 1
}

func BoolToInt(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

func UserToObject(user client.UserResponse) types.Object {
	obj, _ := types.ObjectValue(UserAttrTypes, map[string]attr.Value{
		"uid":            types.Int64Value(user.Uid),
		"first_name":     types.StringValue(user.FirstName),
		"last_name":      types.StringValue(user.LastName),
		"username":       types.StringValue(user.Username),
		"email":          types.StringValue(user.Email),
		"user_image_url": types.StringValue(user.UserImageURL),
	})
	return obj
}

func MapLinkResponseToModel(resp *client.GolinkResponse, model *linkResourceModel, setLastUpdated bool) {
	model.ID = types.StringValue(strconv.FormatInt(resp.Gid, 10))
	model.Gid = types.Int64Value(resp.Gid)
	model.Cid = types.Int64Value(resp.Cid)
	model.URL = types.StringValue(resp.URL)
	model.Name = types.StringValue(resp.Name)
	model.Description = types.StringValue(resp.Description)
	model.Unlisted = types.BoolValue(IntToBool(resp.Unlisted))
	model.VariableLink = types.BoolValue(IntToBool(resp.VariableLink))
	model.Pinned = types.BoolValue(IntToBool(resp.Pinned))
	model.Format = types.BoolValue(IntToBool(resp.Format))
	model.Hyphens = types.BoolValue(IntToBool(resp.Hyphens))
	model.CreatedAt = types.Int64Value(resp.CreatedAt)
	model.UpdatedAt = types.Int64Value(resp.UpdatedAt)
	model.User = UserToObject(resp.User)

	if setLastUpdated {
		model.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	}

	var tags []string
	for _, tag := range resp.Tags {
		tags = append(tags, tag.Name)
	}
	model.Tags = tags
}
