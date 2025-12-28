// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) SignIn(ctx context.Context) (*AuthResponse, error) {
	if c.Auth.Token == "" {
		return nil, fmt.Errorf("token is required")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Auth.Token))

	_, err = c.doRequest(req)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: c.Auth.Token}, nil
}
