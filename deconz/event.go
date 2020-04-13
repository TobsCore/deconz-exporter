package deconz

import "encoding/json"

// Event represents events that are triggered by Deconz web sockets
type Event struct {
	MessageType  string          `json:"t"`
	EventType    string          `json:"e"`
	ID           int             `json:"id,string"`
	UniqueID     string          `json:"uniqueid"`
	ResourceType string          `json:"r"`
	GroupID      string          `json:"gid"`
	SceneID      string          `json:"scid"`
	Name         string          `json:"name"`
	State        json.RawMessage `json:"state"`
	Config       SensorConfig    `json:"config"`
	Sensor       struct {
		Config struct {
			Battery   int  `json:"battery"`
			On        bool `json:"on"`
			Reachable bool `json:"reachable"`
		} `json:"config"`
		Ep               int    `json:"ep"`
		Etag             string `json:"etag"`
		ID               int    `json:"id,string"`
		Manufacturername string `json:"manufacturername"`
		Mode             int    `json:"mode"`
		Modelid          string `json:"modelid"`
		Name             string `json:"name"`
		State            struct {
			Buttonevent interface{} `json:"buttonevent"`
			Lastupdated string      `json:"lastupdated"`
		} `json:"state"`
		Type     string `json:"type"`
		Uniqueid string `json:"uniqueid"`
	} `json:"sensor"`
}
