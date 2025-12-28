package client

type CreateLinkRequest struct {
	URL         string    `json:"url"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
	Unlisted    int32     `json:"unlisted,omitempty"`
	Private     int32     `json:"private,omitempty"`
	Public      int32     `json:"public,omitempty"`
	Format      int32     `json:"format,omitempty"`
	Hyphens     int32     `json:"hyphens,omitempty"`
	Aliases     []string  `json:"aliases,omitempty"`
	Geolinks    []Geolink `json:"geolinks,omitempty"`
}

type UpdateLinkRequest struct {
	Gid         int64     `json:"gid"`
	URL         string    `json:"url"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
	Unlisted    int32     `json:"unlisted,omitempty"`
	Private     int32     `json:"private,omitempty"`
	Public      int32     `json:"public,omitempty"`
	Format      int32     `json:"format,omitempty"`
	Hyphens     int32     `json:"hyphens,omitempty"`
	Aliases     []string  `json:"aliases,omitempty"`
	Geolinks    []Geolink `json:"geolinks,omitempty"`
}

type Geolink struct {
	Location string `json:"location"`
	URL      string `json:"url"`
}

type AuthStruct struct {
	Token string `json:"token"`
}

type AuthResponse struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

type GolinksResponse struct {
	Metadata MetadataResponse `json:"metadata"`
	Results  []GolinkResponse `json:"results"`
}

type MetadataResponse struct {
	Limit        int64         `json:"limit"`
	Offset       int64         `json:"offset"`
	Count        int64         `json:"count"`
	TotalResults int64         `json:"total_results"`
	Links        LinksResponse `json:"links"`
}

type LinksResponse struct {
	Prev string `json:"prev"`
	Next string `json:"next"`
}

type GolinkResponse struct {
	Gid          int64                `json:"gid"`
	Cid          int64                `json:"cid"`
	User         UserResponse         `json:"user"`
	URL          string               `json:"url"`
	Name         string               `json:"name"`
	Description  string               `json:"description"`
	Tags         []TagResponse        `json:"tags"`
	Unlisted     int32                `json:"unlisted"`
	VariableLink int32                `json:"variable_link"`
	Pinned       int32                `json:"pinned"`
	Format       int32                `json:"format"`
	Hyphens      int32                `json:"hyphens"`
	RedirectHits RedirectHitsResponse `json:"redirect_hits"`
	CreatedAt    int64                `json:"created_at"`
	UpdatedAt    int64                `json:"updated_at"`
}

type UserResponse struct {
	Uid          int64  `json:"uid"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	UserImageURL string `json:"user_image_url"`
}

type TagResponse struct {
	Tid  int64  `json:"tid"`
	Name string `json:"name"`
}

type RedirectHitsResponse struct {
	Daily   int64 `json:"daily"`
	Weekly  int64 `json:"weekly"`
	Monthly int64 `json:"monthly"`
	Alltime int64 `json:"alltime"`
}
