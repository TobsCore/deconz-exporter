package main

// Sensor wraps the response of ZigBee
type Sensor struct {
	Config struct {
		Battery   int  `json:"battery"`
		Offset    int  `json:"offset"`
		TurnedOn  bool `json:"on"`
		Reachable bool `json:"reachable"`
		Heatsetpoint int `json:"heatsetpoint"`
		Mode      string `json:"mode"`
	} `json:"config"`
	Endpoint     int    `json:"ep"`
	Etag         string `json:"etag"`
	Manufacturer string `json:"manufacturername"`
	ModelID      string `json:"modelid"`
	Name         string `json:"name"`
	Mode         int    `json:"mode"`
	State        struct {
		Lastupdated string `json:"lastupdated"`
		Temperature int    `json:"temperature"`
		Humidity    int    `json:"humidity"`
		Pressure    int    `json:"pressure"`
		Open		bool   `json:open`
		On			bool   `json:on`
		Valve		int    `json:valve`
		Power		int	   `json:power`
	} `json:"state"`
	Swversion string `json:"swversion"`
	Type      string `json:"type"`
	UID       string `json:"uniqueid"`
}
