package client

import (
	"ValorantRPC/pkg/lockfile"
	"net/http"
	"time"
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
	Version   string `json:"version"`
	ProductID string `json:"productId"`
	ExitCode  *int   `json:"exitCode"`
}

type externalSessionResponse map[string]session

type chatSessionResponse struct {
	Loaded    bool   `json:"loaded"`
	ProductID string `json:"product_id"`
	PUUID     string `json:"puuid"`
}

type regionResponse struct {
	Region string `json:"region"`
}
type GameState struct {
	MatchMode    string // "competitive", "unranked", "deathmatch"....
	MapPath      string // /Game/Maps/Ascent/Ascent
	IsCustomGame bool   // TRUE or FALSE
	AllyScore    int
	EnemyScore   int
	MatchID      string
	AgentID      string
	StartTime    time.Time
}

type PresenceResponse struct {
	Presence []struct {
		Puuid   string `json:"puuid"`
		Private string `json:"private"`
		Product string `json:"product"`
	} `json:"presences"`
}
type Private struct {
	PartyOwnerMatchScoreAllyTeam  int    `json:"partyOwnerMatchScoreAllyTeam"`
	PartyOwnerMatchScoreEnemyTeam int    `json:"partyOwnerMatchScoreEnemyTeam"`
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
	MatchID          string   `json:"MatchID"`
	State            string   `json:"State"`
	MapUrl           string   `json:"MapID"`
	ModeID           string   `json:"ModeID"`
	ProvisioningFlow string   `json:"ProvisioningFlow"`
	Players          []Player `json:"Players"`
}
