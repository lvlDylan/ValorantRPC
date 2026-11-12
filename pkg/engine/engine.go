package engine

import (
	"ValorantRPC/pkg/client"
	"ValorantRPC/pkg/discord"
	"ValorantRPC/pkg/logging"
	"context"
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

func (e *Engine) Start(ctx context.Context, triggerScan chan struct{}, statusChan chan bool) {
	e.logger.Info("Engine started successfully.")

	ticker := time.NewTicker(e.interval)
	defer ticker.Stop()

	e.tick(statusChan)

	for {
		select {
		case <-ctx.Done():
			e.logger.Info("Shutting down engine...")
			e.discord.Disconnect()
			return
		case <-ticker.C:
			e.tick(statusChan)
		case <-triggerScan:
			e.logger.Info("Triggering scan...")
			e.tick(statusChan)
		}
	}
}

func (e *Engine) tick(statusChan chan bool) {

	isRunning := e.client.IsValorantRunning()
	if isRunning {
		statusChan <- true
		if err := e.discord.Connect(); err == nil {
			state, err := e.client.GetGameState()
			if err != nil {
				e.logger.Error("Error getting game state", "error", err)
				return
			}

			if state.MatchID != "" {
				if state.MatchID != e.currentMatchID {
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
		}
	} else {
		e.discord.ClearActivity()
		statusChan <- false
	}
}
