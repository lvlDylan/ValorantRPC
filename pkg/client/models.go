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
