package discord

import (
	"ValorantRPC/pkg/client"
	"fmt"
	"strings"

	discord "github.com/ffx64/discord-rpc-go/client"
)

func (c *Client) mapGameStateToActivity(state *client.GameState) discord.Activity {
	if state == nil {
		return discord.Activity{}
	}

	if state.MatchMode == "MENUS" || state.MatchID == "" {
		return discord.Activity{
			Type:    discord.Playing,
			Details: c.translator.T("rpc.menus.details"),
			State:   c.translator.T("rpc.menus.state"),
			Assets: &discord.Assets{
				LargeImage: "valorant_logo",
				LargeText:  "VALORANT",
			},
			Timestamps: &discord.Timestamps{
				Start: state.StartTime.Unix(),
			},
		}
	}

	modeTranslationKey := fmt.Sprintf("queue.%s", strings.ToLower(state.MatchMode))
	modeTranslated := c.translator.T(modeTranslationKey)
	if modeTranslated == modeTranslationKey {
		modeTranslated = state.MatchMode
	}

	gameMap := c.registry.GetMap(state.MapPath)
	agent := c.registry.GetAgent(state.AgentID)

	detailsText := c.translator.T("rpc.ingame.details", gameMap.DisplayName, modeTranslated)
	stateText := c.translator.T("rpc.ingame.state", agent.DisplayName, state.AllyScore, state.EnemyScore)

	act := discord.Activity{
		Type:    discord.Playing,
		Details: detailsText,
		State:   stateText,
		Assets: &discord.Assets{
			LargeImage: gameMap.Asset,
			LargeText:  gameMap.DisplayName,
			SmallImage: agent.Asset,
			SmallText:  agent.DisplayName,
		},
		Timestamps: &discord.Timestamps{
			Start: state.StartTime.Unix(),
		},
	}

	return act
}
