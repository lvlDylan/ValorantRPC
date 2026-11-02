package client

import (
	"ValorantRPC/pkg/lockfile"
	"net/http"
)

type Client struct {
	localClient  *http.Client
	remoteClient *http.Client
	lockfile     *lockfile.Lockfile

	bearerToken       string
	entitlementsToken string
	puuid             string
	region            string
	version           string
}

type entitlementsResponse struct {
	AccessToken string `json:"accessToken"`
	Token       string `json:"token"`
	Subject     string `json:"subject"`
}

type session struct {
	Version string `json:"version"`
}

type externalSessionResponse map[string]session

type regionResponse struct {
	Region string `json:"region"`
}
type GameState struct {
	SessionLoopState string // MENUS, PREGAME OR INGAME
	MatchMode        string // "competitive", "unranked", "deathmatch"....
	MapPath          string // /Game/Maps/Ascent/Ascent
	IsCustomGame     bool   // TRUE or FALSE
	AllyScore        int
	EnemyScore       int
	MatchID          string
	AgentID          string
}

type PresenceResponse struct {
	Presence []struct {
		Puuid   string `json:"puuid"`
		Private string `json:"private"`
	} `json:"presences"`
}
type Private struct {
	IsCustomGame                  bool   `json:"isCustomGame"`
	SessionLoopState              string `json:"sessionLoopState"`
	PartyOwnerMatchScoreAllyTeam  int    `json:"partyOwnerMatchScoreAllyTeam"`
	PartyOwnerMatchScoreEnemyTeam int    `json:"partyOwnerMatchScoreEnemyTeam"`
	MatchMap                      string `json:"matchMap"`
	QueueID                       string `json:"queueId"`
	ProvisioningFlow              string `json:"provisioningFlow"`
}

type GamePlayerResponse struct {
	Subject string `json:"Subject"`
	MatchID string `json:"MatchID"`
}

type Player struct {
	Subject     string `json:"Subject"`
	TeamId      string `json:"TeamID"`
	CharacterId string `json:"CharacterID"`
}

type GameDetailsResponse struct {
	MatchID string   `json:"MatchID"`
	State   string   `json:"State"`
	MapUrl  string   `json:"MapID"`
	ModeID  string   `json:"ModeID"`
	Type    string   `json:"ProvisioningFlow"`
	Players []Player `json:"Players"`
}
