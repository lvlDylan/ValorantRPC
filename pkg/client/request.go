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
