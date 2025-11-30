package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const HostURL string = "https://api.golinks.io"

type CreateLink struct {
	Gid         int64     `json:"gid,omitempty"`
	Cid         int64     `json:"cid,omitempty"`
	Uid         int64     `json:"uid,omitempty"`
	URL         string    `json:"url"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Tags        []Tag     `json:"tags,omitempty"`
	Unlisted    int64     `json:"unlisted,omitempty"`
	Private     int64     `json:"private,omitempty"`
	Public      int64     `json:"public,omitempty"`
	Format      int64     `json:"format,omitempty"`
	Hyphens     int64     `json:"hyphens,omitempty"`
	Aliases     []string  `json:"aliases,omitempty"`
	Geolinks    []Geolink `json:"geolinks,omitempty"`
	CreatedAt   int64     `json:"created_at,omitempty"`
	UpdatedAt   int64     `json:"updated_at,omitempty"`
}

type Geolink struct {
	Location string `json:"location"`
	URL      string `json:"url"`
}

type Tag struct {
	Tid  int64  `json:"tid,omitempty"`
	Name string `json:"name"`
}

type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Auth       AuthStruct
	Token      string
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
	Unlisted     int64                `json:"unlisted"`
	VariableLink int64                `json:"variable_link"`
	Pinned       int64                `json:"pinned"`
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

func NewClient(token *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
		HostURL:    HostURL,
	}

	if token != nil {
		c.Token = *token
	}

	c.Auth = AuthStruct{
		Token: *token,
	}

	ar, err := c.SignIn()
	if err != nil {
		return nil, err
	}

	c.Token = ar.Token

	return &c, nil
}

func (c *Client) GetGolinks() (*GolinksResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/golinks", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return nil, err
	}

	golinksResponse := GolinksResponse{}
	err = json.Unmarshal(body, &golinksResponse)
	if err != nil {
		return nil, err
	}

	return &golinksResponse, nil
}

func (c *Client) CreateLink(link CreateLink) (*GolinkResponse, error) {
	formData := url.Values{}
	formData.Set("name", link.Name)
	formData.Set("url", link.URL)
	if link.Description != "" {
		formData.Set("description", link.Description)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/golinks", c.HostURL), strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return nil, err
	}

	golinkResponse := GolinkResponse{}
	err = json.Unmarshal(body, &golinkResponse)
	if err != nil {
		return nil, err
	}

	return &golinkResponse, nil

}

func (c *Client) UpdateLink(link CreateLink) (*GolinkResponse, error) {
	formData := url.Values{}
	formData.Set("name", link.Name)
	formData.Set("url", link.URL)
	formData.Set("gid", strconv.FormatInt(link.Gid, 10))
	formData.Set("description", link.Description)

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/golinks", c.HostURL), strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return nil, err
	}

	golinkResponse := GolinkResponse{}
	err = json.Unmarshal(body, &golinkResponse)
	if err != nil {
		return nil, err
	}

	return &golinkResponse, nil
}

func (c *Client) DeleteLink(gid int64) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/golinks?gid=%d", c.HostURL, gid), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return err
	}

	golinkResponse := GolinkResponse{}
	err = json.Unmarshal(body, &golinkResponse)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetLink(gid string) (*GolinkResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/golinks/%s", c.HostURL, gid), nil)

	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return nil, err
	}

	golinkResponse := GolinkResponse{}
	err = json.Unmarshal(body, &golinkResponse)
	if err != nil {
		return nil, err
	}

	return &golinkResponse, nil

}

func (c *Client) doRequest(req *http.Request, authToken *string) ([]byte, error) {
	token := c.Token

	if authToken != nil {
		token = *authToken
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
