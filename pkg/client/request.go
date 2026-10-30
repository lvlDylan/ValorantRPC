package client

import (
	"ValorantRPC/pkg/errors"
	"fmt"
	"net/http"
)

func (c *Client) doLocalRequest(method string, endpoint string) (*http.Response, error) {
	url := fmt.Sprintf("%s://127.0.0.1:%d%s", c.lockfile.Protocol, c.lockfile.Port, endpoint)
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth("riot", c.lockfile.Password)

	response, err := c.localClient.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		response.Body.Close()
		return nil, fmt.Errorf("%w: status code %d", errors.ErrRiotApiRequestFailed, response.StatusCode)
	}

	return response, nil
}

func (c *Client) doRemoteRequest(method string, endpoint string) (*http.Response, error) {
	url := fmt.Sprintf("https://glz-%s-1.%s.a.pvp.net%s", c.region, c.region, endpoint)
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Riot-ClientPlatform", "ew0KCSJwbGF0Zm9ybVR5cGUiOiAiUEMiLA0KCSJwbGF0Zm9ybU9TIjogIldpbmRvd3MiLA0KCSJwbGF0Zm9ybU9TVmVyc2lvbiI6ICIxMC4wLjE5MDQyLjEuMjU2LjY0Yml0IiwNCgkicGxhdGZvcm1DaGlwc2V0IjogIlVua25vd24iDQp9")
	request.Header.Set("X-Riot-ClientVersion", c.version)
	request.Header.Set("X-Riot-Entitlements-JWT", c.entitlementsToken)
	request.Header.Set("Authorization", "Bearer "+c.bearerToken)

	response, err := c.remoteClient.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		response.Body.Close()
		return nil, fmt.Errorf("%w: status code %d", errors.ErrRiotApiRequestFailed, response.StatusCode)
	}

	return response, nil
}
