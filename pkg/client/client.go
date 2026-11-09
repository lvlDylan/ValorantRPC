package client

import (
	"ValorantRPC/pkg/errors"
	"ValorantRPC/pkg/lockfile"
	"crypto/tls"
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
	_, err := lockfile.Read()
	if err != nil {
		return false
	}
	return true
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
