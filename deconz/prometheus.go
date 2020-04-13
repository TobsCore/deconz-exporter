package deconz

import (
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
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

// Collect is used to collect data points for sensors
func Collect(sensor Sensor) {
	labels := prometheus.Labels{"name": sensor.Name, "uid": sensor.UID, "manufacturer": sensor.Manufacturer, "model": sensor.ModelID, "type": sensor.Type}
	slimLabels := prometheus.Labels{"name": sensor.Name, "manufacturer": sensor.Manufacturer, "model": sensor.ModelID}

	switch state := sensor.State.(type) {
	case *ZHATemperature:
		tmpMetric.With(labels).Set(float64(state.Temperature) / 100.0)
		lastUpdMetric.With(slimLabels).Set(float64(timeDiff(state.Lastupdated).Seconds()))
	case *ZHAHumidity:
		humidMetric.With(labels).Set(float64(state.Humidity) / 100.0)
		lastUpdMetric.With(slimLabels).Set(float64(timeDiff(state.Lastupdated).Seconds()))
	case *ZHAPressure:
		pressureMetric.With(labels).Set(float64(state.Pressure))
		lastUpdMetric.With(slimLabels).Set(float64(timeDiff(state.Lastupdated).Seconds()))
	case *Daylight:
		lastUpdMetric.With(slimLabels).Set(float64(timeDiff(state.Lastupdated).Seconds()))
	default:
		log.Printf("unrecognized state %T", state)
	}

	collectBatteryData(sensor, slimLabels)

}

// collectBatteryData collects information about the battery state of the given sensor.
func collectBatteryData(sensor Sensor, labels prometheus.Labels) {
	btryMetric.With(labels).Set(float64(sensor.Battery))

}

func timeDiff(lastUpdate string) time.Duration {
	// Parse Date
	format := "2006-01-02T15:04:05"
	lastUpdateParsed, err := time.Parse(format, lastUpdate)
	if err != nil {
		log.Fatalf("Failed to parse date %s, %s", lastUpdate, err)
	}

	return time.Now().Sub(lastUpdateParsed)
}
