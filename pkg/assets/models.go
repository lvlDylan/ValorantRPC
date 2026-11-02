package assets

type RemoteMap struct {
	UUID        string `json:"uuid"`
	DisplayName string `json:"displayName"`
	MapUrl      string `json:"mapUrl"`
}

type RemoteAgent struct {
	UUID        string `json:"uuid"`
	DisplayName string `json:"displayName"`
}

type ResponseAgents struct {
	Data []RemoteAgent `json:"data"`
}

type ResponseMaps struct {
	Data []RemoteMap `json:"data"`
}
type Agent struct {
	UUID        string
	DisplayName string
	Asset       string
}

type Map struct {
	UUID        string
	DisplayName string
	Asset       string
}
type Registry struct {
	agents map[string]Agent
	maps   map[string]Map
}
