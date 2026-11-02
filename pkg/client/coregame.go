package client

import (
	"ValorantRPC/pkg/errors"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetGameMatch() (*GamePlayerResponse, error) {
	response, err := c.doRemoteRequest(http.MethodGet, "/core-game/v1/players/"+c.puuid)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil, nil
		}
		return nil, errors.New("coregame.GetGameMatchID", errors.ErrRiotApiRequestFailed)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("coregame.GetGameMatchID", errors.ErrRiotApiRequestFailed)
	}

	var gamePlayerResponse GamePlayerResponse
	if err := json.NewDecoder(response.Body).Decode(&gamePlayerResponse); err != nil {
		cause := fmt.Errorf("%w: %v", errors.ErrRiotApiRequestFailed, err)
		return nil, errors.New("coregame.GetGameMatchID", cause)
	}

	return &gamePlayerResponse, nil
}

func (c *Client) GetGameMatchDetails(matchID string) (*GameDetailsResponse, error) {
	response, err := c.doRemoteRequest(http.MethodGet, "/core-game/v1/matches/"+matchID)
	if err != nil {
		return nil, errors.New("coregame.GetGameMatchDetails", errors.ErrRiotApiRequestFailed)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("coregame.GetGameMatchDetails", errors.ErrRiotApiRequestFailed)
	}

	var gameMatchDetails GameDetailsResponse
	if err := json.NewDecoder(response.Body).Decode(&gameMatchDetails); err != nil {
		cause := fmt.Errorf("%w: %v", errors.ErrRiotApiRequestFailed, err)
		return nil, errors.New("coregame.GetGameMatchDetails", cause)
	}

	return &gameMatchDetails, nil
}
