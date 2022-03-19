package main

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/goburrow/modbus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// statusBatteryTemp
type statusBatteryTempType int

const (
	NormalTemp statusBatteryTempType = iota
	OverTemp
	LowTemp
)

func (me statusBatteryTempType) String() string {
	return [...]string{"NormalTemp", "OverTemp", "LowTemp"}[me]
}

// statusBatteryVolt
type statusBatteryVoltType int

const (
	NormalVolt statusBatteryVoltType = iota
	OverVolt
	UnderVolt
	LowVoltDisconnect
	FaultVolt
)

func (me statusBatteryVoltType) String() string {
	return [...]string{"NormalVolt", "OverVolt", "UnderVolt", "LowVoltDisconnect", "FaultVolt"}[me]
}

// statusChargingStatus
type statusChargingStatusType int

const (
	NoCharging statusChargingStatusType = iota
	FaultCharging
	PromoteCharging
	EqualibriumCharging
)

func (me statusChargingStatusType) String() string {
	return [...]string{"NoCharging", "Fault", "PromoteCharging", "EqualibriumCharging"}[me]
}

// statusChargingInputVoltStatus
type statusChargingInputVoltStatusType int

const (
	NormalInputVolt statusChargingInputVoltStatusType = iota
	NoPowerInputVolt
	HigherInputVolt
	ErrorInputVolt
)

func (me statusChargingInputVoltStatusType) String() string {
	return [...]string{"NormalInputVolt", "NoPowerInputVolt", "HigherInputVolt", "ErrorInputVolt"}[me]
}

// Epever
type Epever struct {
	device  string
	handler *modbus.RTUClientHandler
	client  modbus.Client

	ratedInputVoltage float64
	ratedInputCurrent float64
	ratedInputPower   float64

	ratedBatteryVoltage float64
	ratedBatteryCurrent float64
	ratedBatteryPower   float64

	chargeVoltage float64
	chargeCurrent float64
	chargePower   float64

	batteryVoltage float64
	batteryCurrent float64
	batteryPower   float64

	loadVoltage float64
	loadCurrent float64
	loadPower   float64

	tempBattery  float64
	tempInside   float64
	tempHeatsink float64

	batteryPercent    float64
	tempRemoteBattery float64
	tempBattery2      float64

	statusBattery                   uint16
	statusBatteryWrongID            bool
	statusBatteryResistanceAbnormal bool
	statusBatteryTemp               statusBatteryTempType
	statusBatteryVolt               statusBatteryVoltType

	statusCharging                         uint16
	statusChargingRunning                  bool
	statusChargingStatus                   statusChargingStatusType
	statusChargingLoadOpenCircuit          bool
	statusChargingLoadMosfetShort          bool
	statusChargingLoadShort                bool
	statusChargingLoadOverCurrent          bool
	statusChargingInputOverCurrent         bool
	statusChargingAntiReverseMosfetShort   bool
	statusChargingOrAntiReverseMosfetShort bool
	statusChargingMosfetShort              bool
	statusChargingInputVoltStatus          statusChargingInputVoltStatusType

	statusDischarging uint16

	histBatteryVoltageTodayMax float64
	histBatteryVoltageTodayMin float64
	histConsumedToday          float64
	histConsumedMonth          float64
	histConsumedYear           float64
	histConsumed               float64
	histGeneratedToday         float64
	histGeneratedMonth         float64
	histGeneratedYear          float64
	histGenerated              float64

	batteryNetVoltage float64
	batteryNetCurrent float64

	// TODO
	batteryConfigBatteryType uint16
	batteryConfigCapacity    uint16

	batteryConfigTempCoef                          float64
	batteryConfigOverVoltDisconnect                float64
	batteryConfigChargingLimitVoltage              float64
	batteryConfigOverVoltageReconnect              float64
	batteryConfigEqualizeChargingVoltage           float64
	batteryConfigBoostChargingVoltage              float64
	batteryConfigFloatChargingVoltage              float64
	batteryConfigBoostReconnectChargingVoltage     float64
	batteryConfigLowVoltageReconnectVoltage        float64
	batteryConfigUnderVoltageWarningRecoverVoltage float64
	batteryConfigUnderVoltageWarningVoltage        float64
	batteryConfigLowVoltageDisconnectVoltage       float64
	batteryConfigDischargingLimitVoltage           float64

	chargeEqualizationDuration uint16
	chargeBoostDuration        uint16
	chargeEqualizePeriodDays   uint16

	RTCsec   uint16
	RTCmin   uint16
	RTChour  uint16
	RTCday   uint16
	RTCmonth uint16
	RTCyear  uint16
}

// Show this epever monitor as a string
func (e *Epever) String() string {
	rInput := fmt.Sprintf("Rated input %.2fV %.2fA %.2fW", e.ratedInputVoltage, e.ratedInputCurrent, e.ratedInputPower)
	rBattery := fmt.Sprintf("Rated battery %.2fV %.2fA %.2fW", e.ratedBatteryVoltage, e.ratedBatteryCurrent, e.ratedBatteryPower)

	batteryStatus := fmt.Sprintf("%s,%s", e.statusBatteryTemp, e.statusBatteryVolt)
	if e.statusBatteryWrongID {
		batteryStatus = fmt.Sprintf("%s%s", batteryStatus, " WrongID")
	}
	if e.statusBatteryResistanceAbnormal {
		batteryStatus = fmt.Sprintf("%s%s", batteryStatus, " ResAbnormal")
	}

	chargingStatus := fmt.Sprintf("%s,%s", e.statusChargingStatus, e.statusChargingInputVoltStatus)
	if e.statusChargingRunning {
		chargingStatus = fmt.Sprintf("%s%s", chargingStatus, " Running")
	}
	if e.statusChargingLoadOpenCircuit {
		chargingStatus = fmt.Sprintf("%s%s", chargingStatus, " LoadOpenCircuit")
	}
	if e.statusChargingLoadMosfetShort {
		chargingStatus = fmt.Sprintf("%s%s", chargingStatus, " LoadMosfetShort")
	}
	if e.statusChargingLoadShort {
		chargingStatus = fmt.Sprintf("%s%s", chargingStatus, " LoadShort")
	}
	if e.statusChargingLoadOverCurrent {
		chargingStatus = fmt.Sprintf("%s%s", chargingStatus, " LoadOverCurrent")
	}
	if e.statusChargingInputOverCurrent {
		chargingStatus = fmt.Sprintf("%s%s", chargingStatus, " InputOverCurrent")
	}
	if e.statusChargingAntiReverseMosfetShort {
		chargingStatus = fmt.Sprintf("%s%s", chargingStatus, " AntiReverseMosfetShort")
	}
	if e.statusChargingOrAntiReverseMosfetShort {
		chargingStatus = fmt.Sprintf("%s%s", chargingStatus, " ChargingOrAntiReverseMosfetShort")
	}
	if e.statusChargingMosfetShort {
		chargingStatus = fmt.Sprintf("%s%s", chargingStatus, " ChargingMosfetShort")
	}

	chargeData := fmt.Sprintf("Charge %.2fV %.2fA %.2fW [%b] %s", e.chargeVoltage, e.chargeCurrent, e.chargePower, e.statusCharging, chargingStatus)
	batteryData := fmt.Sprintf("Battery %.2fV %.2fA %.2fW (%.2f percent) [%b] %s\nBattery Net %.2fV %.2fA dayVoltRange %.2fV - %.2fV",
		e.batteryVoltage,
		e.batteryCurrent,
		e.batteryPower,
		e.batteryPercent,
		e.statusBattery,
		batteryStatus,
		e.batteryNetVoltage,
		e.batteryNetCurrent,
		e.histBatteryVoltageTodayMin,
		e.histBatteryVoltageTodayMax,
	)
	loadData := fmt.Sprintf("Load %.2fV %.2fA %.2fW [%b]", e.loadVoltage, e.loadCurrent, e.loadPower, e.statusDischarging)

	tempData := fmt.Sprintf("Temp battery:%.2fc inside:%.2fc heatsink:%.2fc battery2:%.2fc", e.tempBattery, e.tempInside, e.tempHeatsink, e.tempBattery2)

	hGenerated := fmt.Sprintf("Generated %.2fkwh Day %.2fkwh Mon %.2fkwh Year %.2fkwh Total", e.histGeneratedToday, e.histGeneratedMonth, e.histGeneratedYear, e.histGenerated)
	hConsumed := fmt.Sprintf("Consumed %.2fkwh Day %.2fkwh Mon %.2fkwh Year %.2fkwh Total", e.histConsumedToday, e.histConsumedMonth, e.histConsumedYear, e.histConsumed)

	batConfig := fmt.Sprintf(" - OverVolt(Disconnect %.2f Reconnect %.2f)\n"+
		" - LowVoltage(Disconnect %.2f Reconnect %.2f)\n"+
		" - UnderVolt(Warning %.2f Recover %.2f)\n"+
		" - Charge(boost %.2f float %.2f equalize %.2f)\n"+
		" - ChargingLimit %.2f BoostReconnect %.2f DischargingLimit %.2f\n"+
		" - BatteryConfig %d (USR/SEAL/GEL/FLOOD) Capacity %dAh",
		e.batteryConfigOverVoltDisconnect,
		e.batteryConfigOverVoltageReconnect,
		e.batteryConfigLowVoltageDisconnectVoltage,
		e.batteryConfigLowVoltageReconnectVoltage,
		e.batteryConfigUnderVoltageWarningVoltage,
		e.batteryConfigUnderVoltageWarningRecoverVoltage,
		e.batteryConfigBoostChargingVoltage,
		e.batteryConfigFloatChargingVoltage,
		e.batteryConfigEqualizeChargingVoltage,
		e.batteryConfigChargingLimitVoltage,
		e.batteryConfigBoostReconnectChargingVoltage,
		e.batteryConfigDischargingLimitVoltage,

		e.batteryConfigBatteryType,
		e.batteryConfigCapacity,
	)

	chargeConfig := fmt.Sprintf("ChargeConfig Equalization %d mins boost %d mins. Equalization period %d days",
		e.chargeEqualizationDuration,
		e.chargeBoostDuration,
		e.chargeEqualizePeriodDays,
	)

	dateTime := fmt.Sprintf("RTC %4d-%2d-%2d %2d:%2d:%2d",
		e.RTCyear,
		e.RTCmonth,
		e.RTCday,
		e.RTChour,
		e.RTCmin,
		e.RTCsec)

	return fmt.Sprintf("EPEVER %s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n", dateTime, rInput, rBattery, chargeData, batteryData, loadData, tempData, hConsumed, hGenerated, batConfig, chargeConfig)
}

// Prometheus metrics
var (
	ratedInputVoltage = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_rated_input_voltage",
		Help: "Rated input voltage"})
	ratedInputCurrent = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_rated_input_current",
		Help: "Rated input current"})
	ratedInputPower = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_rated_input_power",
		Help: "Rated input power"})

	pvVoltage = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_pv_voltage",
		Help: "PV array voltage"})
	pvCurrent = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_pv_current",
		Help: "PV array current"})
	pvPower = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_pv_power",
		Help: "PV array power"})
	loadVoltage = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_load_voltage",
		Help: "Load voltage"})
	loadCurrent = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_load_current",
		Help: "Load current"})
	loadPower = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_load_power",
		Help: "Load power"})
	batVoltage = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_bat_voltage",
		Help: "Battery array voltage"})
	batCurrent = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_bat_current",
		Help: "Battery array current"})
	batPower = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_bat_power",
		Help: "Battery array power"})

	tempBattery = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_temp_battery",
		Help: "Temperature battery"})
	tempInside = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_temp_inside",
		Help: "Temperature inside"})
	tempHeatsink = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_temp_heatsink",
		Help: "Temperature heatsink"})
	tempRemoteBattery = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_temp_remote_battery",
		Help: "Temperature remote battery"})

	batteryPercent = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_battery_percent",
		Help: "Battery percent"})

	consumedToday = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_consumed_today",
		Help: "Consumed today"})
	consumedMonth = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_consumed_month",
		Help: "Consumed month"})
	consumedYear = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_consumed_year",
		Help: "Consumed year"})
	consumedTotal = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_consumed_total",
		Help: "Consumed total"})

	generatedToday = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_generated_today",
		Help: "Generated today"})
	generatedMonth = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_generated_month",
		Help: "Generated month"})
	generatedYear = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_generated_year",
		Help: "Generated year"})
	generatedTotal = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_generated_total",
		Help: "Generated total"})

	batteryNetCurrent = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_battery_net_current",
		Help: "Battery net current"})

	batteryNetVoltage = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_battery_net_voltage",
		Help: "Battery net voltage"})

	solarConfigNum = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_config_num",
		Help: "Number of panels"})
	solarConfigTotalPower = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_config_total_power",
		Help: "Total max power"})
	solarConfigBatteryNum = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_config_battery_num",
		Help: "Number of batteries"})

	statBatteryWrongID = promauto.NewGauge(prometheus.GaugeOpts{Name: "status_battery_wrong_id",
		Help: "Status Battery Wrong ID"})
	statBatteryResistanceAbnormal = promauto.NewGauge(prometheus.GaugeOpts{Name: "status_battery_resistance_abnormal",
		Help: "Status Battery Resistance Abnormal",
	})
	statBatteryTemp = promauto.NewGauge(prometheus.GaugeOpts{Name: "status_battery_temp",
		Help: "Status Battery Temp"})
	statBatteryVolt = promauto.NewGauge(prometheus.GaugeOpts{Name: "status_battery_volt",
		Help: "Status Battery Temp"})
	statChargingRunning = promauto.NewGauge(prometheus.GaugeOpts{Name: "status_charging_running",
		Help: "Status Charging Running"})
	statChargingLoadOpenCircuit = promauto.NewGauge(prometheus.GaugeOpts{Name: "status_charging_load_open_circuit",
		Help: "Status Charging Load Open Circuit"})
	statChargingLoadMosfetShort = promauto.NewGauge(prometheus.GaugeOpts{Name: "status_charging_load_mosfet_short",
		Help: "Status Charging Load Mosfet Short"})
	statChargingLoadShort = promauto.NewGauge(prometheus.GaugeOpts{Name: "status_charging_load_short",
		Help: "Status Charging Load Short"})
	statChargingLoadOverCurrent = promauto.NewGauge(prometheus.GaugeOpts{Name: "status_charging_load_over_current",
		Help: "Status Charging Load Over Current"})
	statChargingInputOverCurrent = promauto.NewGauge(prometheus.GaugeOpts{Name: "status_charging_input_over_current",
		Help: "Status Charging Input Over Current"})
	statChargingAntiReverseMosfetShort = promauto.NewGauge(prometheus.GaugeOpts{Name: "status_charging_anti_reverse_mosfet_short",
		Help: "Status Charging Anti Reverse Mosfet Short"})
	statChargingOrAntiReverseMosfetShort = promauto.NewGauge(prometheus.GaugeOpts{Name: "status_charging_or_anti_reverse_mosfet_short",
		Help: "Status Charging Or Anit Reverse Mosfet Short"})
	statChargingMosfetShort = promauto.NewGauge(prometheus.GaugeOpts{Name: "status_charging_mosfet_short",
		Help: "Status Charging Mosfet Short"})
	statChargingStatus = promauto.NewGauge(prometheus.GaugeOpts{Name: "status_charging_status",
		Help: "Status Charging Status"})
	statChargingInputVoltStatus = promauto.NewGauge(prometheus.GaugeOpts{Name: "status_charging_input_volt_status",
		Help: "Status Charging Input Volt Status"})

	configEqualizationDuration = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_config_equalization_duration",
		Help: "Config Equalization Duration"})
	configBoostDuration = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_config_boost_duration",
		Help: "Config Boost Duration"})
	configEqualizationPeriod = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_config_equalization_period",
		Help: "Config Equalization Period"})

	batConfigOverVoltDisconnect = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_battery_config_over_voltage_disconnect",
		Help: "Config Over Voltage Disconnect"})
	batConfigChargingLimitVoltage = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_battery_config_charging_limit_voltage",
		Help: "Config Charging Limit Voltage"})
	batConfigOverVoltageReconnect = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_battery_config_over_voltage_reconnect",
		Help: "Config Over Voltage Reconnect"})
	batConfigEqualizeChargingVoltage = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_battery_config_equalize_charging_voltage",
		Help: "Config Equalize Charging Voltage"})
	batConfigBoostChargingVoltage = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_battery_config_boost_charging_voltage",
		Help: "Config Boost Charging Voltage"})
	batConfigFloatChargingVoltage = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_battery_config_float_charging_voltage",
		Help: "Config Float Charging Voltage"})
	batConfigBoostReconnectChargingVoltage = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_battery_config_boost_reconnect_charging_voltage",
		Help: "Config Boost Reconnect Charging Voltage"})
	batConfigLowVoltageReconnectVoltage = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_battery_config_low_voltage_reconnect_voltage",
		Help: "Config Low Voltage Reconnect Voltage"})
	batConfigUnderVoltageWarningRecoverVoltage = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_battery_config_under_voltage_warning_reconnect_voltage",
		Help: "Config Under Voltage Warning Reconnect Voltage"})
	batConfigUnderVoltageWarningVoltage = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_battery_config_under_voltage_warning_voltage",
		Help: "Config Under Voltage Warning Voltage"})
	batConfigLowVoltageDisconnectVoltage = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_battery_config_low_voltage_disconnect_voltage",
		Help: "Config Low Voltage Disconnect Voltage"})
	batConfigDischargingLimitVoltage = promauto.NewGauge(prometheus.GaugeOpts{Name: "solar_battery_config_discharging_limit_voltage",
		Help: "Config Discharging Limit Voltage"})
)

// PushMetrics to prometheus
func (e *Epever) PushMetrics() {
	ratedInputVoltage.Set(e.ratedInputVoltage)
	ratedInputCurrent.Set(e.ratedInputCurrent)
	ratedInputPower.Set(e.ratedBatteryPower)
	pvVoltage.Set(e.chargeVoltage)
	pvCurrent.Set(e.chargeCurrent)
	pvPower.Set(e.chargePower)
	loadVoltage.Set(e.loadVoltage)
	loadCurrent.Set(e.loadCurrent)
	loadPower.Set(e.loadPower)
	batVoltage.Set(e.batteryVoltage)
	batCurrent.Set(e.batteryCurrent)
	batPower.Set(e.batteryPower)
	tempBattery.Set(e.tempBattery)
	tempInside.Set(e.tempInside)
	tempHeatsink.Set(e.tempHeatsink)
	tempRemoteBattery.Set(e.tempRemoteBattery)
	batteryPercent.Set(e.batteryPercent / 100)

	consumedToday.Set(e.histConsumedToday)
	consumedMonth.Set(e.histConsumedMonth)
	consumedYear.Set(e.histConsumedYear)
	consumedTotal.Set(e.histConsumed)

	generatedToday.Set(e.histGeneratedToday)
	generatedMonth.Set(e.histGeneratedMonth)
	generatedYear.Set(e.histGeneratedYear)
	generatedTotal.Set(e.histGenerated)

	batteryNetCurrent.Set(e.batteryNetCurrent)
	batteryNetVoltage.Set(e.batteryNetVoltage)

	if e.statusBatteryResistanceAbnormal {
		statBatteryResistanceAbnormal.Set(1)
	} else {
		statBatteryResistanceAbnormal.Set(0)
	}
	if e.statusBatteryWrongID {
		statBatteryWrongID.Set(1)
	} else {
		statBatteryWrongID.Set(0)
	}

	statBatteryTemp.Set(float64(e.statusBatteryTemp))
	statBatteryVolt.Set(float64(e.statusBatteryVolt))

	statChargingStatus.Set(float64(e.statusChargingStatus))
	statChargingInputVoltStatus.Set(float64(e.statusChargingInputVoltStatus))

	if e.statusChargingRunning {
		statChargingRunning.Set(1)
	} else {
		statChargingRunning.Set(0)
	}
	if e.statusChargingLoadOpenCircuit {
		statChargingLoadOpenCircuit.Set(1)
	} else {
		statChargingLoadOpenCircuit.Set(0)
	}
	if e.statusChargingLoadMosfetShort {
		statChargingLoadMosfetShort.Set(1)
	} else {
		statChargingLoadMosfetShort.Set(0)
	}
	if e.statusChargingLoadShort {
		statChargingLoadShort.Set(1)
	} else {
		statChargingLoadShort.Set(0)
	}
	if e.statusChargingLoadOverCurrent {
		statChargingLoadOverCurrent.Set(1)
	} else {
		statChargingLoadOverCurrent.Set(0)
	}
	if e.statusChargingInputOverCurrent {
		statChargingInputOverCurrent.Set(1)
	} else {
		statChargingInputOverCurrent.Set(0)
	}
	if e.statusChargingAntiReverseMosfetShort {
		statChargingAntiReverseMosfetShort.Set(1)
	} else {
		statChargingAntiReverseMosfetShort.Set(0)
	}
	if e.statusChargingOrAntiReverseMosfetShort {
		statChargingOrAntiReverseMosfetShort.Set(1)
	} else {
		statChargingOrAntiReverseMosfetShort.Set(0)
	}
	if e.statusChargingMosfetShort {
		statChargingMosfetShort.Set(1)
	} else {
		statChargingMosfetShort.Set(0)
	}

	configEqualizationDuration.Set(float64(e.chargeEqualizationDuration))
	configEqualizationPeriod.Set(float64(e.chargeEqualizePeriodDays))
	configBoostDuration.Set(float64(e.chargeBoostDuration))

	batConfigOverVoltDisconnect.Set(e.batteryConfigOverVoltDisconnect)
	batConfigChargingLimitVoltage.Set(e.batteryConfigChargingLimitVoltage)
	batConfigOverVoltageReconnect.Set(e.batteryConfigOverVoltageReconnect)
	batConfigEqualizeChargingVoltage.Set(e.batteryConfigEqualizeChargingVoltage)
	batConfigBoostChargingVoltage.Set(e.batteryConfigBoostChargingVoltage)
	batConfigFloatChargingVoltage.Set(e.batteryConfigFloatChargingVoltage)
	batConfigBoostReconnectChargingVoltage.Set(e.batteryConfigBoostReconnectChargingVoltage)
	batConfigLowVoltageReconnectVoltage.Set(e.batteryConfigLowVoltageReconnectVoltage)
	batConfigUnderVoltageWarningRecoverVoltage.Set(e.batteryConfigUnderVoltageWarningRecoverVoltage)
	batConfigUnderVoltageWarningVoltage.Set(e.batteryConfigUnderVoltageWarningVoltage)
	batConfigLowVoltageDisconnectVoltage.Set(e.batteryConfigLowVoltageDisconnectVoltage)
	batConfigDischargingLimitVoltage.Set(e.batteryConfigDischargingLimitVoltage)
}

// Create a new Epever using the given device eg "/dev/ttyXRUSB0"
func NewEpever(device string) *Epever {
	return &Epever{device: device}
}

// Connect
func (e *Epever) Connect() {
	if e.handler != nil {
		fmt.Printf("Closing existing connection.\n")
		e.handler.Close()
	}

	e.handler = modbus.NewRTUClientHandler(e.device)
	e.handler.BaudRate = 115200
	e.handler.DataBits = 8
	e.handler.Parity = "N"
	e.handler.StopBits = 1
	e.handler.SlaveId = 1
	e.handler.Timeout = 10 * time.Second

	for {
		err := e.handler.Connect()
		if err == nil {
			e.client = modbus.NewClient(e.handler)
			fmt.Printf("Connected to epever on %s\n", e.device)
			return
		}
		fmt.Printf("Error connecting using %s. Waiting... %v\n", e.device, err)
		time.Sleep(10 * time.Second)
	}

}

// Read some input registers and reconnect/retry if needed.
func (e *Epever) readWithRetry(address uint16, quantity uint16) (results []byte) {
	for {
		if e.client == nil {
			e.Connect()
		}
		data, err := e.client.ReadInputRegisters(address, quantity)
		if err == nil {
			return data
		} else {
			fmt.Printf("Error readWithRetry: %x - %v\n", address, err)
			e.Connect()
		}
	}
}

// Read some holding registers and reconnect/retry if needed.
func (e *Epever) readHoldingWithRetry(address uint16, quantity uint16) (results []byte) {
	for {
		if e.client == nil {
			e.Connect()
		}
		data, err := e.client.ReadHoldingRegisters(address, quantity)
		if err == nil {
			return data
		} else {
			fmt.Printf("Error readHoldingWithRetry: %x - %v\n", address, err)
			e.Connect()
		}
	}
}

// Refresh gets latest stats
func (e *Epever) Refresh() error {
	// Grab some stats...

	// client is ready for reading stuff...
	ratedInput := e.readWithRetry(REGRatedInputVoltage, 4)
	e.ratedInputVoltage = float64(binary.BigEndian.Uint16(ratedInput)) / 100
	e.ratedInputCurrent = float64(binary.BigEndian.Uint16(ratedInput[2:])) / 100
	e.ratedInputPower = float64(uint32(binary.BigEndian.Uint16(ratedInput[4:]))|
		(uint32(binary.BigEndian.Uint16(ratedInput[6:]))<<16)) / 100

	ratedBattery := e.readWithRetry(REGRatedBatteryVoltage, 4)
	e.ratedBatteryVoltage = float64(binary.BigEndian.Uint16(ratedBattery)) / 100
	e.ratedBatteryCurrent = float64(binary.BigEndian.Uint16(ratedBattery[2:])) / 100
	e.ratedBatteryPower = float64(uint32(binary.BigEndian.Uint16(ratedBattery[4:]))|
		(uint32(binary.BigEndian.Uint16(ratedBattery[6:]))<<16)) / 100

	//
	chargeData := e.readWithRetry(REGChargeVoltage, 4)
	e.chargeVoltage = float64(binary.BigEndian.Uint16(chargeData)) / 100
	e.chargeCurrent = float64(binary.BigEndian.Uint16(chargeData[2:])) / 100
	e.chargePower = float64(uint32(binary.BigEndian.Uint16(chargeData[4:]))|
		(uint32(binary.BigEndian.Uint16(chargeData[6:]))<<16)) / 100

	batteryData := e.readWithRetry(REGBatteryVoltage, 4)
	e.batteryVoltage = float64(binary.BigEndian.Uint16(batteryData)) / 100
	e.batteryCurrent = float64(binary.BigEndian.Uint16(batteryData[2:])) / 100
	e.batteryPower = float64(uint32(binary.BigEndian.Uint16(batteryData[4:]))|
		(uint32(binary.BigEndian.Uint16(batteryData[6:]))<<16)) / 100

	loadData := e.readWithRetry(REGLoadVoltage, 4)
	e.loadVoltage = float64(binary.BigEndian.Uint16(loadData)) / 100
	e.loadCurrent = float64(binary.BigEndian.Uint16(loadData[2:])) / 100
	e.loadPower = float64(uint32(binary.BigEndian.Uint16(loadData[4:]))|
		(uint32(binary.BigEndian.Uint16(loadData[6:]))<<16)) / 100

	tempData := e.readWithRetry(REGTempBattery, 3)
	e.tempBattery = float64(binary.BigEndian.Uint16(tempData)) / 100
	e.tempInside = float64(binary.BigEndian.Uint16(tempData[2:])) / 100
	e.tempHeatsink = float64(binary.BigEndian.Uint16(tempData[4:])) / 100

	e.batteryPercent = float64(binary.BigEndian.Uint16(e.readWithRetry(REGBatteryPercent, 1)))
	e.tempRemoteBattery = float64(binary.BigEndian.Uint16(e.readWithRetry(REGTempRemoteBattery, 1))) / 100

	e.tempBattery2 = float64(binary.BigEndian.Uint16(e.readWithRetry(REGTempBattery2, 1))) / 100

	statuses := e.readWithRetry(REGBatteryStatus, 3)
	e.statusBattery = binary.BigEndian.Uint16(statuses)

	e.statusBatteryWrongID = ((e.statusBattery >> 15) & 1) == 1
	e.statusBatteryResistanceAbnormal = ((e.statusBattery >> 8) & 1) == 1
	e.statusBatteryTemp = statusBatteryTempType((e.statusBattery >> 4) & 0b1111)
	e.statusBatteryVolt = statusBatteryVoltType(e.statusBattery & 0b1111)

	e.statusCharging = binary.BigEndian.Uint16(statuses[2:])

	e.statusChargingRunning = (e.statusCharging & 1) == 1

	e.statusChargingLoadOpenCircuit = ((e.statusCharging >> 5) & 1) == 1

	e.statusChargingLoadMosfetShort = ((e.statusCharging >> 7) & 1) == 1
	e.statusChargingLoadShort = ((e.statusCharging >> 8) & 1) == 1
	e.statusChargingLoadOverCurrent = ((e.statusCharging >> 9) & 1) == 1
	e.statusChargingInputOverCurrent = ((e.statusCharging >> 10) & 1) == 1
	e.statusChargingAntiReverseMosfetShort = ((e.statusCharging >> 11) & 1) == 1
	e.statusChargingOrAntiReverseMosfetShort = ((e.statusCharging >> 12) & 1) == 1
	e.statusChargingMosfetShort = ((e.statusCharging >> 13) & 1) == 1

	e.statusChargingInputVoltStatus = statusChargingInputVoltStatusType((e.statusCharging >> 14) & 0b11)
	e.statusChargingStatus = statusChargingStatusType((e.statusCharging >> 2) & 0b11)

	e.statusDischarging = binary.BigEndian.Uint16(statuses[4:])

	historicalData := e.readWithRetry(REGBatteryVoltageTodayMax, 18)

	e.histBatteryVoltageTodayMax = float64(binary.BigEndian.Uint16(historicalData)) / 100
	e.histBatteryVoltageTodayMin = float64(binary.BigEndian.Uint16(historicalData[2:])) / 100

	e.histConsumedToday = float64(uint32(binary.BigEndian.Uint16(historicalData[4:]))|
		(uint32(binary.BigEndian.Uint16(historicalData[6:]))<<16)) / 100
	e.histConsumedMonth = float64(uint32(binary.BigEndian.Uint16(historicalData[8:]))|
		(uint32(binary.BigEndian.Uint16(historicalData[10:]))<<16)) / 100
	e.histConsumedYear = float64(uint32(binary.BigEndian.Uint16(historicalData[12:]))|
		(uint32(binary.BigEndian.Uint16(historicalData[14:]))<<16)) / 100
	e.histConsumed = float64(uint32(binary.BigEndian.Uint16(historicalData[16:]))|
		(uint32(binary.BigEndian.Uint16(historicalData[18:]))<<16)) / 100
	e.histGeneratedToday = float64(uint32(binary.BigEndian.Uint16(historicalData[20:]))|
		(uint32(binary.BigEndian.Uint16(historicalData[22:]))<<16)) / 100
	e.histGeneratedMonth = float64(uint32(binary.BigEndian.Uint16(historicalData[24:]))|
		(uint32(binary.BigEndian.Uint16(historicalData[26:]))<<16)) / 100
	e.histGeneratedYear = float64(uint32(binary.BigEndian.Uint16(historicalData[28:]))|
		(uint32(binary.BigEndian.Uint16(historicalData[30:]))<<16)) / 100
	e.histGenerated = float64(uint32(binary.BigEndian.Uint16(historicalData[32:]))|
		(uint32(binary.BigEndian.Uint16(historicalData[34:]))<<16)) / 100

	batteryNetData := e.readWithRetry(REGBatteryNetVoltage, 3)

	// TODO: Check these, they need to be signed
	e.batteryNetVoltage = float64(binary.BigEndian.Uint16(batteryNetData)) / 100

	// Net Current
	loNCurrent := (0xffff & int32(binary.BigEndian.Uint16(batteryNetData[2:]))) // 0000-ffff
	hiNCurrent := (int32(binary.BigEndian.Uint16(batteryNetData[4:])) << 16)
	netCurrentVal := hiNCurrent | loNCurrent
	e.batteryNetCurrent = float64(netCurrentVal) / 100

	batteryConfigData := e.readHoldingWithRetry(REGBatteryType, 15)

	e.batteryConfigBatteryType = binary.BigEndian.Uint16(batteryConfigData)  // 9000
	e.batteryConfigCapacity = binary.BigEndian.Uint16(batteryConfigData[2:]) // 9001

	e.batteryConfigTempCoef = float64(binary.BigEndian.Uint16(batteryConfigData[4:])) / 100                           // 9002
	e.batteryConfigOverVoltDisconnect = float64(binary.BigEndian.Uint16(batteryConfigData[6:])) / 100                 // 9003
	e.batteryConfigChargingLimitVoltage = float64(binary.BigEndian.Uint16(batteryConfigData[8:])) / 100               // 9004
	e.batteryConfigOverVoltageReconnect = float64(binary.BigEndian.Uint16(batteryConfigData[10:])) / 100              // 9005
	e.batteryConfigEqualizeChargingVoltage = float64(binary.BigEndian.Uint16(batteryConfigData[12:])) / 100           // 9006
	e.batteryConfigBoostChargingVoltage = float64(binary.BigEndian.Uint16(batteryConfigData[14:])) / 100              // 9007
	e.batteryConfigFloatChargingVoltage = float64(binary.BigEndian.Uint16(batteryConfigData[16:])) / 100              // 9008
	e.batteryConfigBoostReconnectChargingVoltage = float64(binary.BigEndian.Uint16(batteryConfigData[18:])) / 100     // 9009
	e.batteryConfigLowVoltageReconnectVoltage = float64(binary.BigEndian.Uint16(batteryConfigData[20:])) / 100        // 900a
	e.batteryConfigUnderVoltageWarningRecoverVoltage = float64(binary.BigEndian.Uint16(batteryConfigData[22:])) / 100 // 900b
	e.batteryConfigUnderVoltageWarningVoltage = float64(binary.BigEndian.Uint16(batteryConfigData[24:])) / 100        // 900c
	e.batteryConfigLowVoltageDisconnectVoltage = float64(binary.BigEndian.Uint16(batteryConfigData[26:])) / 100       // 900d
	e.batteryConfigDischargingLimitVoltage = float64(binary.BigEndian.Uint16(batteryConfigData[28:])) / 100           // 900e

	e.chargeEqualizationDuration = binary.BigEndian.Uint16(e.readHoldingWithRetry(REGBatteryEqualizeDuration, 1))
	e.chargeBoostDuration = binary.BigEndian.Uint16(e.readHoldingWithRetry(REGBatteryBoostDuration, 1))
	e.chargeEqualizePeriodDays = binary.BigEndian.Uint16(e.readHoldingWithRetry(REGBatteryEqualizePeriodDays, 1))

	rtcData := e.readHoldingWithRetry(REGRTCSecMin, 3)

	e.RTCsec = uint16(rtcData[1])
	e.RTCmin = uint16(rtcData[0])

	e.RTChour = uint16(rtcData[3])
	e.RTCday = uint16(rtcData[2])

	e.RTCmonth = uint16(rtcData[5])
	e.RTCyear = uint16(rtcData[4])

	/*
	   const REGBatteryRatedVoltage = 0x9067
	   const REGBatteryDischarge = 0x906d
	   const REGBatteryChargeDepth = 0x906e
	   const REGBatteryChargingMode = 0x9070

	   const REGBatteryEqualizePeriodDays = 0x9016

	   // RTC
	   const REGRTCSecMin = 0x9013
	   const REGRTCHourDay = 0x9014
	   const REGRTCMonthYear = 0x9015
	*/

	return nil
}
