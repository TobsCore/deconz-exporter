package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gorilla/websocket"
)

func init() {
	token = os.Getenv("DECONZ_TOKEN")
	deconzHost = os.Getenv("DECONZ_HOST")
	parsedDeconzPort, err := strconv.Atoi(os.Getenv("DECONZ_PORT"))
	if err != nil {
		log.Fatalf("%s must be integer", "DECONZ_PORT")
	}
	deconzPort = parsedDeconzPort
	portEnv := os.Getenv("DECONZ_APP_PORT")
	if portEnv != "" {
		// If the environment variable is set, use this value.
		// If it is not set, the default value for the port is used.
		parsedPort, err := strconv.Atoi(portEnv)
		if err != nil {
			log.Fatalf("%s must be integer", "DECONZ_APP_PORT")
		}
		port = parsedPort
	}
	flag.BoolVar(&verbose, "verbose", false, "Verbose logging")
}

func evalVars() error {
	if token == "" {
		return fmt.Errorf("DECONZ_TOKEN is required")
	} else if deconzHost == "" {
		return fmt.Errorf("DECONZ_HOST is required")
	} else if deconzPort == 0 {
		return fmt.Errorf("DECONZ_PORT is required")
	} else if port == 0 {
		return fmt.Errorf("DECONZ_APP_PORT is required")
	}

	return nil
}

func main() {
	flag.Parse()
	if err := evalVars(); err != nil {
		log.Fatalf("%s", err)
	}
	//recordMetrics()
	openWS()
	//serve()
}

func serve() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
            <head><title>Deconz Exporter</title></head>
            <body>
            <h1>Deconz Exporter</h1>
            <p><a href="/metrics">Metrics</a></p>
            </body>
            </html>`))
	})
	instance := fmt.Sprintf(":%d", port)
	fmt.Println("Starting Server on", instance)
	log.Fatal(http.ListenAndServe(instance, nil))
}

const (
	// DefaultPort describes the port the application will use if no dedicated port
	// is defined as an environment variable
	DefaultPort int = 8080
)

var (
	deconzHost = ""
	deconzPort = 0
	port       = DefaultPort
	token      = ""
	verbose    = false // The default value is set via the `flag` package
	labels     = []string{"name", "uid", "manufacturer", "model", "type"}
	labelsArbi = []string{"name", "manufacturer", "model"}
	tmpMetric  = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "deconz",
		Subsystem: "sensor",
		Name:      "temperature",
		Help:      "Temperature of sensor in Celsius",
	}, labels)
	btryMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "deconz",
		Subsystem: "sensor",
		Name:      "battery",
		Help:      "Battery level of sensor in percent",
	}, labelsArbi)
	humidMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "deconz",
		Subsystem: "sensor",
		Name:      "humidity",
		Help:      "Humidity of sensor in percent",
	}, labels)
	pressureMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "deconz",
		Subsystem: "sensor",
		Name:      "pressure",
		Help:      "Air pressure in hectopascal (hPa)",
	}, labels)
	lastUpdMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "deconz",
		Subsystem: "sensor",
		Name:      "sinceUpdate",
		Help:      "The time since the last update that was received from this sensor",
	}, labelsArbi)
	errorCtr = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "deconz",
		Subsystem: "sensor",
		Name:      "errors",
		Help:      "Failures to retrieve data from API",
	})
)

func openWS() {
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

// recordMetrics starts a runner, that will collect the metrics and place them
// in the corresponding struct, so can fetch these metrics. For this to work,
// the sensor data are fetched every 5 seconds.
func recordMetrics() {
	url := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", deconzHost, deconzPort),
		Path:   fmt.Sprintf("api/%s/sensors", token),
	}
	go func() {
		for {
			sensors, err := pollSensors(url)
			if err != nil {
				errorCtr.Inc()
				if verbose {
					fmt.Println("failed to retrieve data", err)
				}
			} else {
				if verbose {
					fmt.Printf("%+v", sensors)
				}
				for _, sensor := range sensors {
					labels := prometheus.Labels{"name": sensor.Name, "uid": sensor.UID, "manufacturer": sensor.Manufacturer, "model": sensor.ModelID, "type": sensor.Type}

					switch sensor.Type {
					case "ZHATemperature":
						tmpMetric.With(labels).Set(float64(sensor.State.Temperature) / 100.0)
					case "ZHAHumidity":
						humidMetric.With(labels).Set(float64(sensor.State.Humidity) / 100.0)
					case "ZHAPressure":
						pressureMetric.With(labels).Set(float64(sensor.State.Pressure))
					}

					collectArbitraryData(sensor, labels)
				}
			}

			time.Sleep(5 * time.Second)
		}
	}()
}

// collectArbitraryData will collect additional information about a sensor and
// expose it for prometheus. Such information includes:
// battery: The sensors battery state
// last update: The time since the last update has been retrieved
// The labels are used by prometheus as metadata, to identify the sensor.
func collectArbitraryData(sensor Sensor, labels prometheus.Labels) {
	delete(labels, "type")
	delete(labels, "uid")
	btryMetric.With(labels).Set(float64(sensor.Config.Battery))

	// Parse Date
	format := "2006-01-02T15:04:05"
	lastUpdate, err := time.Parse(format, sensor.State.Lastupdated)
	if err != nil {
		log.Fatalf("Failed to parse date %s, %s", sensor.State.Lastupdated, err)
	}

	timeDiff := time.Now().Sub(lastUpdate)
	lastUpdMetric.With(labels).Set(float64(timeDiff.Seconds()))
}

// pollSensors tries to retrieve the sensor data from the given URL. It will
// return an error, if the data can not be fetched, otherwise a map of the
// sensor data is returned, where the key is an ID, given by deconz.
// The sensor contains only information about one metric, such as temperature
// or air pressure. If one sensor has the capability to provide multiple 
// data points, each of these data points are returned as one entry in the map.
func pollSensors(url url.URL) (map[string]Sensor, error) {
	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var s map[string]Sensor
	if err := json.Unmarshal(body, &s); err != nil {
		return nil, err
	}

	return s, err
}
