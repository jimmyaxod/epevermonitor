package main

// ====
// 30xx
const REGRatedInputVoltage = 0x3000
const REGRatedInputCurrent = 0x3001
const REGRatedInputPowerL = 0x3002
const REGRatedInputPowerH = 0x3003

const REGRatedBatteryVoltage = 0x3004
const REGRatedBatteryCurrent = 0x3005
const REGRatedBatteryPowerL = 0x3006
const REGRatedBatteryPowerH = 0x3007

//const REGChargingMode = 0x3008

// 3009
// 300a
// 300b
// 300c
//const REGRatedLoadVoltage = 0x300d
//const REGRatedLoadCurrent = 0x300e
//const REGRatedLoadPowerL = 0x300f
//const REGRatedLoadPowerH = 0x3010

// ====
// 31xx
const REGChargeVoltage = 0x3100
const REGChargeCurrent = 0x3101
const REGChargePowerL = 0x3102
const REGChargePowerH = 0x3103

const REGBatteryVoltage = 0x3104
const REGBatteryCurrent = 0x3105
const REGBatteryPowerL = 0x3106
const REGBatteryPowerH = 0x3107

// 3108
// 3109
// 310a
// 310b
const REGLoadVoltage = 0x310c
const REGLoadCurrent = 0x310d
const REGLoadPowerL = 0x310e
const REGLoadPowerH = 0x310f

const REGTempBattery = 0x3110
const REGTempInside = 0x3111
const REGTempHeatsink = 0x3112

// 3113
// 3114
// 3115
// 3116
// 3117
// 3118
// 3119
const REGBatteryPercent = 0x311a
const REGTempRemoteBattery = 0x311b

// 311c
const REGTempBattery2 = 0x311d

//const REGTempAmbient = 0x311e

// 311f

// ====
// 32xx
const REGBatteryStatus = 0x3200
const REGChargingStatus = 0x3201
const REGDischargingStatus = 0x3202

/*
D15
-
D14: 00
H
normal, 01
H
low
,
02H High, 03H
no access
Input volt error.
D13
-
D12
:
output power
:00
-
light
load,01
-
moderate,02
-
rated,03
-
overlo
ad
D11:
short circuit
D10:
unable to discharge
D9:
unable to stop discharging
D8:
output voltage abnormal
D7:
input overpressure
D6: high
voltage side short circuit
D5: boost overpressure
D4:
output overpressure
D1: 0 Normal, 1 Fault.
D0: 1 Running, 0 Standby
*/

// ====
// 33xx
const REGBatteryVoltageTodayMax = 0x3302
const REGBatteryVoltageTodayMin = 0x3303
const REGConsumedTodayL = 0x3304
const REGConsumedTodayH = 0x3305
const REGConsumedMonthL = 0x3306
const REGConsumedMonthH = 0x3307
const REGConsumedYearL = 0x3308
const REGConsumedYearH = 0x3309
const REGConsumedL = 0x330a
const REGConsumedH = 0x330b
const REGGeneratedTodayL = 0x330c
const REGGeneratedTodayH = 0x330d
const REGGeneratedMonthL = 0x330e
const REGGeneratedMonthH = 0x330f
const REGGeneratedYearL = 0x3310
const REGGeneratedYearH = 0x3311
const REGGeneratedL = 0x3312
const REGGeneratedH = 0x3313

const REGBatteryNetVoltage = 0x331a
const REGBatteryNetCurrentL = 0x331b
const REGBatteryNetCurrentH = 0x331c

// Holding registers
const REGBatteryType = 0x9000
const REGBatteryCapacity = 0x9001 // AH
const REGBatteryTempCoef = 0x9002
const REGBatteryOverVoltageDisconnect = 0x9003
const REGBatteryChargingLimitVoltage = 0x9004
const REGBatteryOverVoltageReconnect = 0x9005
const REGBatteryEqualizeChargingVoltage = 0x9006
const REGBatteryBoostChargingVoltage = 0x9007
const REGBatteryFloatChargingVoltage = 0x9008
const REGBatteryBoostReconnectChargingVoltage = 0x9009
const REGBatteryLowVoltageReconnectVoltage = 0x900a
const REGBatteryUnderVoltageWarningRecoverVoltage = 0x900b
const REGBatteryUnderVoltageWarningVoltage = 0x900c
const REGBatteryLowVoltageDisconnectVoltage = 0x900d
const REGBatteryDischargingLimitVoltage = 0x900e

const REGBatteryRatedVoltage = 0x9067
const REGBatteryEqualizeDuration = 0x906b
const REGBatteryBoostDuration = 0x906c
const REGBatteryDischarge = 0x906d
const REGBatteryChargeDepth = 0x906e
const REGBatteryChargingMode = 0x9070

const REGBatteryEqualizePeriodDays = 0x9016

// RTC
const REGRTCSecMin = 0x9013
const REGRTCHourDay = 0x9014
const REGRTCMonthYear = 0x9015
