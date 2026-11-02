package client

import (
	"ValorantRPC/pkg/errors"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) Auth() error {
	response, err := c.doLocalRequest(http.MethodGet, "/entitlements/v1/token")
	if err != nil {
		return errors.New("client.Auth", err)
	}
	defer response.Body.Close()

	var authResponse entitlementsResponse
	if err := json.NewDecoder(response.Body).Decode(&authResponse); err != nil {
		cause := fmt.Errorf("%w: %v", errors.ErrRiotApiMalformed, err)
		return errors.New("client.Auth", cause)
	}

	c.bearerToken = authResponse.AccessToken
	c.entitlementsToken = authResponse.Token
	c.puuid = authResponse.Subject

	return nil
}

func (c *Client) FetchVersion() error {
	response, err := c.doLocalRequest(http.MethodGet, "/product-session/v1/external-sessions")
	if err != nil {
		return errors.New("client.FetchVersion", err)
	}
	defer response.Body.Close()

	var sessions externalSessionResponse
	if err := json.NewDecoder(response.Body).Decode(&sessions); err != nil {
		cause := fmt.Errorf("%w: %v", errors.ErrRiotApiMalformed, err)
		return errors.New("client.FetchVersion", cause)
	}

	for _, session := range sessions {
		if session.Version != "" {
			c.version = session.Version
			break
		}
	}

	return nil
}

func (c *Client) FetchRegion() error {
	response, err := c.doLocalRequest(http.MethodGet, "/riotclient/region-locale")
	if err != nil {
		return errors.New("client.FetchRegion", err)
	}
	defer response.Body.Close()

	var region regionResponse
	if err := json.NewDecoder(response.Body).Decode(&region); err != nil {
		cause := fmt.Errorf("%w: %v", errors.ErrRiotApiMalformed, err)
		return errors.New("client.FetchRegion", cause)
	}

	c.region = region.Region

	return nil
}
