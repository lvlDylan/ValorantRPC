package discord

import (
	"ValorantRPC/pkg/assets"
	"ValorantRPC/pkg/i18n"
	"ValorantRPC/pkg/logging"

	discord "github.com/ffx64/discord-rpc-go/client"
)

type Client struct {
	rpc        *discord.Client
	logger     *logging.Logger
	registry   *assets.Registry
	translator i18n.Translator

	appID     string
	connected bool
}
