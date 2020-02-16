package main

// Sensor wraps the response of ZigBee
type Sensor struct {
	Config struct {
		Battery   int  `json:"battery"`
		Offset    int  `json:"offset"`
		TurnedOn  bool `json:"on"`
		Reachable bool `json:"reachable"`
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
	} `json:"state"`
	Swversion string `json:"swversion"`
	Type      string `json:"type"`
	UID       string `json:"uniqueid"`
}
