package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) SignUp(auth AuthStruct) (*AuthResponse, error) {
	if auth.Token == "" {
		return nil, fmt.Errorf("define token")
	}
	rb, err := json.Marshal(auth)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/signup", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	ar := AuthResponse{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}

	return &ar, nil
}

func (c *Client) SignIn() (*AuthResponse, error) {
	if c.Auth.Token == "" {
		return nil, fmt.Errorf("define token")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Auth.Token))

	_, err = c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	ar := AuthResponse{
		Token: c.Auth.Token,
	}

	return &ar, nil
}

func (c *Client) GetUserTokenSignIn(auth AuthStruct) (*AuthResponse, error) {
	if c.Auth.Token == "" {
		return nil, fmt.Errorf("define username and password")
	}
	rb, err := json.Marshal(auth)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/signin", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, errors.New("Unable to login")
	}

	ar := AuthResponse{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}

	return &ar, nil
}

func (c *Client) SignOut(authToken *string) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/signout", c.HostURL), strings.NewReader(string("")))
	if err != nil {
		return err
	}

	body, err := c.doRequest(req, authToken)
	if err != nil {
		return err
	}

	if string(body) != "Signed out user" {
		return errors.New(string(body))
	}

	return nil
}
