package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tobscore/deconz-exporter/deconz"
)

type appConf struct {
	host  string
	port  int
	token string
}

const (
	// DefaultPort describes the port the application will use if no dedicated
	// port is defined as an environment variable
	DefaultPort int = 8080
)

var (
	port    = DefaultPort
	verbose = false // The default value is set via the `flag` package
	d       = deconz.Deconz{}
	c       = appConf{}
)

func init() {
	c.token = os.Getenv("DECONZ_TOKEN")
	c.host = os.Getenv("DECONZ_HOST")
	parsedDeconzPort, err := strconv.Atoi(os.Getenv("DECONZ_PORT"))
	if err != nil {
		log.Fatalf("%s must be integer", "DECONZ_PORT")
	}
	c.port = parsedDeconzPort
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
	switch {
	case c.token == "":
		return fmt.Errorf("DECONZ_TOKEN is required")
	case c.host == "":
		return fmt.Errorf("DECONZ_HOST is required")
	case c.port == 0:
		return fmt.Errorf("DECONZ_PORT is required")
	case port == 0:
		return fmt.Errorf("DECONZ_APP_PORT is required")
	}

	return nil
}

func main() {
	flag.Parse()
	if err := evalVars(); err != nil {
		log.Fatalf("%s", err)
	}
	d, err := deconz.Init(c.host, c.port, c.token)
	if err != nil {
		log.Fatalf("failed init, %s", err)
	}
	log.Printf("Connected to %s at %s:%d", d.Config.Name, d.Host, d.Port)
	log.Printf("(%d sensors)", len(d.Sensors))
	for i, sensor := range d.Sensors {
		if verbose {
			log.Println(i, sensor)
		}
		deconz.Collect(sensor)
	}
	go serve()
	d.ListenToEvents()
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
	log.Println("Starting Server on", instance)
	log.Fatal(http.ListenAndServe(instance, nil))
}
