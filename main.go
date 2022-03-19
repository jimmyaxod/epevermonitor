package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Some configuration consts we'll expose as metrics alongside the epever metrics
const SOLAR_CONFIG_PANEL_NUM = 6
const SOLAR_CONFIG_MAX_POWER = 600
const SOLAR_CONFIG_BATTERY_NUM = 4

const PROMETHEUS_PORT = 2112
const UPDATE_PERIOD = 1 * time.Minute

// main
func main() {

	ep := NewEpever("/dev/ttyXRUSB0")

	// Setup prometheus
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(fmt.Sprintf(":%d", PROMETHEUS_PORT), nil)

	fmt.Printf("Listening on port %d...\n", PROMETHEUS_PORT)

	// periodically update the prometheus regs
	ticker := time.NewTicker(UPDATE_PERIOD)

	err := ep.Refresh()
	fmt.Printf("Epever %v %v\n", err, ep)

	for {
		select {
		case <-ticker.C:
			err := ep.Refresh()
			fmt.Printf("Epever %v %v\n", err, ep)
			ep.PushMetrics()
			// Push some statics metrics as well
			solarConfigNum.Set(SOLAR_CONFIG_PANEL_NUM)
			solarConfigTotalPower.Set(SOLAR_CONFIG_MAX_POWER)
			solarConfigBatteryNum.Set(SOLAR_CONFIG_BATTERY_NUM)
		}
	}

}
