package client

import (
	"ValorantRPC/pkg/errors"
	"ValorantRPC/pkg/lockfile"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"time"
)

func New() *Client {
	remoteHttpClient := http.DefaultClient
	localHttpClient := &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	return &Client{
		localClient:  localHttpClient,
		remoteClient: remoteHttpClient,
	}
}

func (c *Client) IsValorantRunning() bool {
	lf, err := lockfile.Read()
	if err != nil {
		c.Reset()
		return false
	}

	c.lockfile = lf

	response, err := c.doLocalRequest(http.MethodGet, "/chat/v4/presences")
	if err != nil {
		c.Reset()
		return false
	}
	defer response.Body.Close()

	var presences PresenceResponse
	if err := json.NewDecoder(response.Body).Decode(&presences); err != nil {
		c.Reset()
		return false
	}

	for _, presence := range presences.Presence {
		err := c.EnsureConnected()
		if err != nil {
			return false
		}
		if presence.Puuid == c.puuid && presence.Product == "valorant" {
			return true
		}
	}

	c.Reset()
	return false
}

func (c *Client) IsConnected() bool {
	return c.bearerToken != "" && c.lockfile != nil
}

func (c *Client) EnsureConnected() error {
	if c.IsConnected() {
		return nil
	}

	lf, err := lockfile.Read()
	if err != nil {
		return errors.New("client.EnsureConnected", err)
	}

	c.lockfile = lf

	if err := c.Auth(); err != nil {
		c.Reset()
		return errors.New("client.EnsureConnected", err)
	}

	if err := c.FetchRegion(); err != nil {
		c.Reset()
		return errors.New("client.EnsureConnected", err)
	}

	if err := c.FetchVersion(); err != nil {
		c.Reset()
		return errors.New("client.EnsureConnected", err)
	}

	return nil
}

func (c *Client) Reset() {
	c.lockfile = nil
	c.bearerToken = ""
	c.entitlementsToken = ""
	c.puuid = ""
	c.region = ""
	c.version = ""
}
