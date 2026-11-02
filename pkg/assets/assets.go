package assets

import (
	"ValorantRPC/pkg/errors"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func (r *Registry) fetchAgents(httpClient *http.Client) error {
	res, err := httpClient.Get("https://valorant-api.com/v1/agents?isPlayableCharacter=true")
	if err != nil {
		return fmt.Errorf("could not fetch agents: %w", err)
	}
	defer res.Body.Close()

	var agentsResponse ResponseAgents
	if err := json.NewDecoder(res.Body).Decode(&agentsResponse); err != nil {
		return fmt.Errorf("could not decode agents: %w", err)
	}

	for _, agent := range agentsResponse.Data {
		r.agents[strings.ToLower(agent.UUID)] = Agent{
			UUID:        agent.UUID,
			DisplayName: agent.DisplayName,
			Asset:       fmt.Sprintf("agent_%s", strings.ToLower(agent.DisplayName)),
		}
	}

	return nil
}

func (r *Registry) fetchMaps(httpClient *http.Client) error {
	res, err := httpClient.Get("https://valorant-api.com/v1/maps")
	if err != nil {
		return fmt.Errorf("could not fetch maps: %w", err)
	}
	defer res.Body.Close()

	var mapsResponse ResponseMaps
	if err := json.NewDecoder(res.Body).Decode(&mapsResponse); err != nil {
		return fmt.Errorf("could not decode maps: %w", err)
	}

	for _, m := range mapsResponse.Data {
		r.maps[strings.ToLower(m.MapUrl)] = Map{
			UUID:        m.UUID,
			DisplayName: m.DisplayName,
			Asset:       fmt.Sprintf("map_%s", strings.ToLower(m.DisplayName)),
		}
	}

	return nil
}

func NewRegistry() (*Registry, error) {
	r := &Registry{
		agents: make(map[string]Agent),
		maps:   make(map[string]Map),
	}

	httpClient := &http.Client{Timeout: 10 * time.Second}

	if err := r.fetchAgents(httpClient); err != nil {
		return nil, errors.New("assets.NewRegistry", err)
	}

	if err := r.fetchMaps(httpClient); err != nil {
		return nil, errors.New("assets.NewRegistry", err)
	}

	return r, nil
}

func (r *Registry) GetAgent(uuid string) Agent {
	if agent, ok := r.agents[strings.ToLower(uuid)]; ok {
		return agent
	}
	return Agent{
		UUID:        uuid,
		DisplayName: "Unknown Agent",
		Asset:       "agent_unknown",
	}
}

func (r *Registry) GetMap(mapPath string) Map {
	if m, ok := r.maps[strings.ToLower(mapPath)]; ok {
		return m
	}
	return Map{
		UUID:        "",
		DisplayName: "Unknown Map",
		Asset:       "map_unknown",
	}
}
