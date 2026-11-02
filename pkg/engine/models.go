package engine

import (
	"ValorantRPC/pkg/client"
	"ValorantRPC/pkg/discord"
	"ValorantRPC/pkg/logging"
	"time"
)

type Engine struct {
	client         *client.Client
	discord        *discord.Client
	logger         logging.Logger
	interval       time.Duration
	matchStartTime time.Time
	currentMatchID string
}
