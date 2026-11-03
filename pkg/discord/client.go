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

func (c *Client) ClearActivity() {
	if c.connected && c.rpc != nil {
		c.rpc.SetActivity(discord.Activity{})
	}
}

func isSameActivity(a, b discord.Activity) bool {
	if a.Type != b.Type || a.Details != b.Details || a.State != b.State {
		return false
	}

	if (a.Timestamps == nil) != (b.Timestamps == nil) {
		return false
	}
	if a.Timestamps != nil && b.Timestamps != nil {
		if a.Timestamps.Start != b.Timestamps.Start {
			return false
		}
	}

	if (a.Assets == nil) != (b.Assets == nil) {
		return false
	}

	if a.Assets != nil && b.Assets != nil {
		if a.Assets.LargeImage != b.Assets.LargeImage ||
			a.Assets.LargeText != b.Assets.LargeText ||
			a.Assets.SmallImage != b.Assets.SmallImage ||
			a.Assets.SmallText != b.Assets.SmallText {
			return false
		}
	}

	return true
}

func (c *Client) Update(state *client.GameState) error {
	if !c.connected {
		err := c.Connect()
		if err != nil {
			return errors.New("discord.Update", err)
		}
	}

	activity := c.mapGameStateToActivity(state)
	if isSameActivity(c.lastActivity, activity) {
		return nil
	}

	if err := c.rpc.SetActivity(activity); err != nil {
		c.connected = false
		return errors.New("discord.Update", err)
	}

	c.lastActivity = activity
	c.logger.Debug("RPC updated", "details", activity.Details, "state", activity.State)
	return nil
}
