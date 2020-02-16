package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() {
	flag.StringVar(&token, "token", "", "The API token for deconz")
	flag.StringVar(&deconzHost, "deconz_host", "localhost", "The host address of the deconz instance")
	flag.IntVar(&deconzPort, "deconz_port", 0, "The port on which deconz is available")
	flag.IntVar(&port, "port", 2112, "The port on which this application is started")
	flag.BoolVar(&verbose, "verbose", false, "Verbose logging")
}

func evalFlags() error {
	if token == "" {
		return fmt.Errorf("deconz token is required")
	} else if deconzHost == "" {
		return fmt.Errorf("deconz host is required")
	} else if deconzPort == 0 {
		return fmt.Errorf("deconz port is required")
	}

	return nil
}

func main() {
	flag.Parse()
	if err := evalFlags(); err != nil {
		log.Fatalf("%s", err)
	}
	recordMetrics()
	serve()
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

var (
	deconzHost = ""
	deconzPort = 0
	port       = 0
	token      = ""
	verbose    = false
	labels     = []string{"name", "uid", "manufacturer", "model", "type"}
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
	}, labels)
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
	errorCtr = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "deconz",
		Subsystem: "sensor",
		Name:      "errors",
		Help:      "Failures to retrieve data from API",
	})
)

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
						btryMetric.With(labels).Set(float64(sensor.Config.Battery))
						tmpMetric.With(labels).Set(float64(sensor.State.Temperature) / 100.0)
					case "ZHAHumidity":
						humidMetric.With(labels).Set(float64(sensor.State.Humidity) / 100.0)
					case "ZHAPressure":
						pressureMetric.With(labels).Set(float64(sensor.State.Pressure))
					}
				}
			}

			time.Sleep(5 * time.Second)
		}
	}()
}

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
