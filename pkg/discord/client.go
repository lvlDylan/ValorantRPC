package discord

import (
	"ValorantRPC/pkg/assets"
	"ValorantRPC/pkg/client"
	"ValorantRPC/pkg/errors"
	"ValorantRPC/pkg/i18n"
	"ValorantRPC/pkg/logging"

	discord "github.com/ffx64/discord-rpc-go/client"
)

func New(appID string, translator i18n.Translator, registry *assets.Registry, logger logging.Logger) *Client {
	return &Client{translator: translator, appID: appID, registry: registry, logger: logger}
}

func (c *Client) Connect() error {
	if c.connected {
		return nil
	}

	rpcClient := discord.NewClient(c.appID)
	if err := rpcClient.Login(); err != nil {
		return errors.New("discord.Connect", err)
	}

	c.rpc = rpcClient
	c.connected = true
	return nil
}

func (c *Client) Disconnect() {
	if c.connected && c.rpc != nil {
		c.rpc.SetActivity(discord.Activity{})
		c.rpc.Close()
		c.connected = false
		c.rpc = nil
	}
}

func (c *Client) Update(state *client.GameState) error {
	if !c.connected {
		err := c.Connect()
		if err != nil {
			return errors.New("discord.Update", err)
		}
	}

	activity := c.mapGameStateToActivity(state)
	if err := c.rpc.SetActivity(activity); err != nil {
		c.connected = false
		return errors.New("discord.Update", err)
	}

	return nil
}
