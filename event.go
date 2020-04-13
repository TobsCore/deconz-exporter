package main

// Event represents events that are triggered by Deconz web sockets
type Event struct {
	MessageType  string `json:"t"`
	EventType    string `json:"e"`
	ID           int    `json:"id,string"`
	UniqueID     string `json:"uniqueid"`
	ResourceType string `json:"r"`
	GroupID      string `json:"gid"`
	SceneID      string `json:"scid"`
	Name         string `json:"name"`
	State        struct {
		Bri int       `json:"bri"`
		On  bool      `json:"on"`
		X   int       `json:"x"`
		Y   int       `json:"y"`
		Xy  []float64 `json:"xy"`
	} `json:"state"`
	Sensor struct {
		Config struct {
			Battery   int  `json:"battery"`
			On        bool `json:"on"`
			Reachable bool `json:"reachable"`
		} `json:"config"`
		Ep               int    `json:"ep"`
		Etag             string `json:"etag"`
		ID               string `json:"id"`
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
