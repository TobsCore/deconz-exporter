package deconz

import (
	"encoding/json"
	"fmt"
)

// RestSensor wraps the response of ZigBee
type RestSensor struct {
	Config       SensorConfig    `json:"config"`
	Endpoint     int             `json:"ep"`
	Etag         string          `json:"etag"`
	Manufacturer string          `json:"manufacturername"`
	ModelID      string          `json:"modelid"`
	Name         string          `json:"name"`
	Mode         int             `json:"mode"`
	State        json.RawMessage `json:"state"`
	Swversion    string          `json:"swversion"`
	Type         string          `json:"type"`
	UID          string          `json:"uniqueid"`
}

// SensorConfig represents a configuration of elements that are changed
type SensorConfig struct {
	Battery   int  `json:"battery"`
	Offset    int  `json:"offset"`
	TurnedOn  bool `json:"on"`
	Reachable bool `json:"reachable"`
}

// Sensor is the actual object, that is used by this implementation
type Sensor struct {
	Name         string
	Type         string
	State        interface{}
	Battery      int
	TurnedOn     bool
	Reachable    bool
	Manufacturer string
	UID          string
	ModelID      string
}

// Create will create an actual sensor object for the given deconz sensor
func (d *RestSensor) Create() *Sensor {
	new := Sensor{
		Name:         d.Name,
		Type:         d.Type,
		Battery:      d.Config.Battery,
		TurnedOn:     d.Config.TurnedOn,
		Reachable:    d.Config.Reachable,
		Manufacturer: d.Manufacturer,
		UID:          d.UID,
		ModelID:      d.ModelID,
	}

	state := ParseState(d.Type, d.State)
	new.State = state
	return &new
}

// ParseState is responsible for parsing the given raw json blob for a specific
// sensor type
func ParseState(sensorType string, raw json.RawMessage) interface{} {
	switch sensorType {
	case "ZHATemperature":
		var s ZHATemperature
		_ = json.Unmarshal(raw, &s)
		return &s
	case "ZHAPressure":
		var s ZHAPressure
		_ = json.Unmarshal(raw, &s)
		return &s
	case "ZHAHumidity":
		var s ZHAHumidity
		_ = json.Unmarshal(raw, &s)
		return &s
	case "Daylight":
		var s Daylight
		_ = json.Unmarshal(raw, &s)
		return &s
	default:
		panic("Unknown sensor type " + sensorType)
	}

}

// String is used to print the sensor
func (d RestSensor) String() string {
	return fmt.Sprintf("%s (%s)", d.Name, d.Type)
}

// State is for embedding into event states
type State struct {
	Lastupdated string
}

// ZHATemperature contains temperature information
type ZHATemperature struct {
	State
	Temperature int
}

// ZHAPressure contains pressure information
type ZHAPressure struct {
	State
	Pressure int
}

// ZHAHumidity contains humidity information
type ZHAHumidity struct {
	State
	Humidity int
}

// Daylight contains information about the daylight sensor
type Daylight struct {
	State
	Daylight bool
	Dark     bool
	Status   int
}
