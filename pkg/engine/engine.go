package engine

import (
	"ValorantRPC/pkg/client"
	"ValorantRPC/pkg/discord"
	"ValorantRPC/pkg/logging"
	"time"
)

func New(c *client.Client, d *discord.Client, logger logging.Logger, interval time.Duration) *Engine {
	return &Engine{
		client:         c,
		discord:        d,
		logger:         logger,
		interval:       interval,
		currentMatchID: "",
	}
}

func (e *Engine) Start() {
	e.logger.Info("Starting application. Press CTRL+C to stop.")
	e.discord.Connect()

	ticker := time.NewTicker(e.interval)
	defer ticker.Stop()

	e.tick()

	for range ticker.C {
		e.tick()
	}
}

func (e *Engine) tick() {
	state, err := e.client.GetGameState()
	if err != nil {
		e.logger.Error("Error getting game state", "error", err)
		return
	}

	if state.SessionLoopState == "INGAME" {
		if state.MatchID != "" && state.MatchID != e.currentMatchID {
			e.currentMatchID = state.MatchID
			e.matchStartTime = time.Now()
			e.logger.Info("New match started", "match_id", state.MatchID)
		}
		state.StartTime = e.matchStartTime
	} else {
		e.currentMatchID = ""
		e.matchStartTime = time.Time{}
		state.StartTime = e.matchStartTime
	}

	if err := e.discord.Update(state); err != nil {
		e.logger.Error("Error updating discord", "error", err)
	}

	e.logger.Debug("RPC updated", "state", state.SessionLoopState, "mode", state.MatchMode, "map", state.MapPath)
}
