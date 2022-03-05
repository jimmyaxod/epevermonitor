package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Some configuration consts we'll expose as metrics
const SOLAR_CONFIG_PANEL_NUM = 4
const SOLAR_CONFIG_MAX_POWER = 400
const SOLAR_CONFIG_BATTERY_NUM = 2

// main
func main() {

	ep := NewEpever()

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":2112", nil)

	fmt.Printf("Listening on port 2112...\n")

	// periodically update the prometheus regs
	ticker := time.NewTicker(1 * time.Minute)

	err := ep.Refresh()
	fmt.Printf("Epever %v %v\n", err, ep)

	for {
		select {
		case <-ticker.C:
			err := ep.Refresh()
			fmt.Printf("Epever %v %v\n", err, ep)
			ep.PushMetrics()
			// Do some statics
			solarConfigNum.Set(SOLAR_CONFIG_PANEL_NUM)
			solarConfigTotalPower.Set(SOLAR_CONFIG_MAX_POWER)
			solarConfigBatteryNum.Set(SOLAR_CONFIG_BATTERY_NUM)
		}
	}

}
