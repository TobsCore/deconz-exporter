package deconz

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

// Deconz is a container for the instance's configuration and the application's
// state, which sensors have which values, etc.
type Deconz struct {
	Host    string
	Port    int
	Token   string
	Sensors map[int]Sensor
	Config  Config
}

// Init is used for an initial configuration. It fetches the configuration and
// the initial sensor data. After init is run, updates are received via web
// sockets
func Init(host string, port int, token string) (*Deconz, error) {
	s := Deconz{Host: host, Port: port, Token: token}
	err := s.loadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed load config, %s", err)
	}

	err = s.initSensors()
	if err != nil {
		return nil, fmt.Errorf("failed init sensors, %s", err)
	}

	return &s, nil
}

func (d *Deconz) initSensors() error {
	url := d.apiURL("sensors")
	resp, err := http.Get(url.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var resSensors map[int]RestSensor
	if err := json.Unmarshal(body, &resSensors); err != nil {
		return err
	}

	d.Sensors = map[int]Sensor{}
	for id, sensor := range resSensors {
		d.Sensors[id] = *sensor.Create()
	}

	return nil
}

func (d *Deconz) loadConfig() error {
	url := d.apiURL("config")
	resp, err := http.Get(url.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &d.Config); err != nil {
		return err
	}

	return nil
}

// ListenToEvents will open a websocket connection and update its internal
// state, when updates are retrieved.
func (d *Deconz) ListenToEvents() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "192.168.0.222:8443", Path: ""}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
			// Initialize Battery in order to detect, if it's overwritten
			e := Event{Config: SensorConfig{Battery: -1}}
			if err := json.Unmarshal(message, &e); err != nil {
				log.Println("parse error:", err)
				return
			}

			d.handleEvent(e)
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func (d *Deconz) handleEvent(e Event) {
	sensor := d.Sensors[e.ID]
	switch e.EventType {
	case "added":
		// TODO: Implement this
		break
	case "changed":
		if e.State != nil {
			state := ParseState(sensor.Type, e.State)
			sensor.State = state
		}
		// The valid battery level values are [0..100], therefore having the
		// default -1 value means, that the config is not actually updated.
		if e.Config.Battery != -1 {
			sensor.Battery = e.Config.Battery
			sensor.Reachable = e.Config.Reachable
			sensor.TurnedOn = e.Config.TurnedOn
			log.Printf("update config value of %s (event %+v)", sensor.Name, e)
		}
		Collect(sensor)
		log.Printf("update for %s", sensor.Name)
		break
	case "deleted":
		// TODO: Implement this
		break
	case "scene-called":
		// TODO: Implement this
		break
	}
}

// apiURL builds the URL for the given endpoint.
func (d *Deconz) apiURL(endpoint string) url.URL {
	return url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", d.Host, d.Port),
		Path:   fmt.Sprintf("api/%s/%s", d.Token, endpoint),
	}
}
