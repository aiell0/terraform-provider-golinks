// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	HostURL                = "https://api.golinks.io"
	contentTypeFormEncoded = "application/x-www-form-urlencoded"
)

type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Auth       AuthStruct
	Token      string
}

func NewClient(ctx context.Context, token *string) (*Client, error) {
	if token == nil {
		return nil, fmt.Errorf("token is required")
	}

	c := Client{
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
		HostURL:    HostURL,
		Token:      *token,
		Auth:       AuthStruct{Token: *token},
	}

	ar, err := c.SignIn(ctx)
	if err != nil {
		return nil, err
	}

	c.Token = ar.Token

	return &c, nil
}

func (c *Client) GetGolinks(ctx context.Context) (*GolinksResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/golinks", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	var resp GolinksResponse
	if err := c.doRequestJSON(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GetGolinksByName(ctx context.Context, name string) (*GolinkResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/golinks", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Set("name", name)
	req.URL.RawQuery = query.Encode()

	var resp GolinkResponse
	if err := c.doRequestJSON(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func buildCreateLinkFormData(link CreateLinkRequest) url.Values {
	formData := url.Values{}
	formData.Set("name", link.Name)
	formData.Set("url", link.URL)
	if link.Description != "" {
		formData.Set("description", link.Description)
	}
	formData.Set("unlisted", strconv.Itoa(int(link.Unlisted)))
	formData.Set("public", strconv.Itoa(int(link.Public)))
	formData.Set("private", strconv.Itoa(int(link.Private)))
	formData.Set("format", strconv.Itoa(int(link.Format)))
	if link.Format == 1 {
		formData.Set("hyphens", strconv.Itoa(int(link.Hyphens)))
	}
	for _, alias := range link.Aliases {
		formData.Add("aliases", alias)
	}
	for _, tag := range link.Tags {
		formData.Add("tags[]", tag)
	}
	for i, geo := range link.Geolinks {
		formData.Set(fmt.Sprintf("geolinks[%d][location]", i), geo.Location)
		formData.Set(fmt.Sprintf("geolinks[%d][url]", i), geo.URL)
	}
	return formData
}

func buildUpdateLinkFormData(link UpdateLinkRequest) url.Values {
	formData := url.Values{}
	formData.Set("gid", strconv.FormatInt(link.Gid, 10))
	formData.Set("name", link.Name)
	formData.Set("url", link.URL)
	formData.Set("description", link.Description)
	formData.Set("unlisted", strconv.Itoa(int(link.Unlisted)))
	formData.Set("public", strconv.Itoa(int(link.Public)))
	formData.Set("private", strconv.Itoa(int(link.Private)))
	formData.Set("format", strconv.Itoa(int(link.Format)))
	if link.Format == 1 {
		formData.Set("hyphens", strconv.Itoa(int(link.Hyphens)))
	}
	for _, alias := range link.Aliases {
		formData.Add("aliases", alias)
	}
	for _, tag := range link.Tags {
		formData.Add("tags[]", tag)
	}
	for i, geo := range link.Geolinks {
		formData.Set(fmt.Sprintf("geolinks[%d][location]", i), geo.Location)
		formData.Set(fmt.Sprintf("geolinks[%d][url]", i), geo.URL)
	}
	return formData
}

func (c *Client) CreateLink(ctx context.Context, link CreateLinkRequest) (*GolinkResponse, error) {
	formData := buildCreateLinkFormData(link)

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/golinks", c.HostURL), strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentTypeFormEncoded)

	var resp GolinkResponse
	if err := c.doRequestJSON(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) UpdateLink(ctx context.Context, link UpdateLinkRequest) (*GolinkResponse, error) {
	formData := buildUpdateLinkFormData(link)

	req, err := http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf("%s/golinks", c.HostURL), strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentTypeFormEncoded)

	var resp GolinkResponse
	if err := c.doRequestJSON(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) DeleteLink(ctx context.Context, gid int64) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("%s/golinks?gid=%d", c.HostURL, gid), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", contentTypeFormEncoded)

	var resp GolinkResponse
	return c.doRequestJSON(req, &resp)
}

func (c *Client) GetLink(ctx context.Context, gid string) (*GolinkResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/golinks/%s", c.HostURL, gid), nil)
	if err != nil {
		return nil, err
	}

	var resp GolinkResponse
	if err := c.doRequestJSON(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))

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

func (c *Client) doRequestJSON(req *http.Request, v interface{}) error {
	body, err := c.doRequest(req)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}
