package client

import (
	"ValorantRPC/pkg/errors"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"unicode"
)

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}

	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func formatQueueId(provisioningFlow string, queueID string, isCustomGame bool) string {
	if isCustomGame || strings.EqualFold(provisioningFlow, "CustomGame") {
		return "Custom Game"
	}

	if queueID != "" {
		switch strings.ToLower(queueID) {
		case "competitive":
			return "Competitive"
		case "unrated":
			return "Unrated"
		case "swiftplay":
			return "Swiftplay"
		case "spikerush":
			return "Spike Rush"
		case "deathmatch":
			return "Deathmatch"
		case "hurm":
			return "Team Deathmatch"
		case "ggteam":
			return "Escalation"
		case "onefa":
			return "Replication"
		default:
			return capitalize(queueID)
		}
	}

	return "In match"
}

func (c *Client) getPlayerState() (*Private, error) {
	response, err := c.doLocalRequest(http.MethodGet, "/chat/v4/presences")
	if err != nil {
		return nil, errors.New("player_state.GetPlayerState", err)
	}
	defer response.Body.Close()

	var presences PresenceResponse
	if err := json.NewDecoder(response.Body).Decode(&presences); err != nil {
		return nil, errors.New("player_state.GetPlayerState", err)
	}

	for _, presence := range presences.Presence {
		if presence.Puuid == c.puuid && presence.Private != "" {
			decodedBytes, err := base64.StdEncoding.DecodeString(presence.Private)
			if err != nil {
				return nil, errors.New("player_state.GetPlayerState", err)
			}
			var private Private
			if err := json.Unmarshal(decodedBytes, &private); err != nil {
				return nil, errors.New("player_state.GetPlayerState", err)
			}

			return &private, nil
		}
	}

	return nil, errors.New("player_state.GetPlayerState", errors.ErrPlayerNotFound)
}

func findAgentIdByPlayerUUID(playerUUID string, players []Player) string {
	for _, player := range players {
		if player.Subject == playerUUID {
			return player.CharacterId
		}
	}

	return "agent_unknown"
}

func (c *Client) GetGameState() (*GameState, error) {

	match, err := c.GetGameMatch()
	if err != nil {
		return nil, errors.New("player_state.GetGameState", err)
	}

	if match == nil || match.MatchID == "" {
		return &GameState{
			MatchMode: "MENUS",
		}, nil
	}

	matchDetails, err := c.GetGameMatchDetails(match.MatchID)
	if err != nil {
		return nil, errors.New("player_state.GetGameState", err)
	}

	state := &GameState{
		MatchMode:    matchDetails.ModeID,
		MapPath:      matchDetails.MapUrl,
		IsCustomGame: matchDetails.ProvisioningFlow == "CustomGame",
		MatchID:      matchDetails.MatchID,
		StartTime:    time.Now(),
		AgentID:      findAgentIdByPlayerUUID(c.puuid, matchDetails.Players),
	}

	if presence, err := c.getPlayerState(); err == nil && presence != nil {
		state.AllyScore = presence.PartyOwnerMatchScoreAllyTeam
		state.EnemyScore = presence.PartyOwnerMatchScoreEnemyTeam
	}

	return state, nil
}
