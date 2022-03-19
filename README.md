# Epever Solar Monitor

This can be used to read Epever solar charge controllers, and relay the metrics to prometheus.
It has been tested with Epever Tracer 4210 but should work with others.

## Physical setup

Connection is with a RJ45 -> USB cable. eg 'CC-USB-RS485-150U'

## Linux driver

To get the driver working on linux you'll need to build a kernel module and insmod it.
See 'xr_usb_serial_common-1a'

## Operation

You should now be able to run this, and see various metrics and statistics from the charge controller.
They will also be exposed on an endpoint for prometheus. You can then setup grafana etc

## Sample output

From commandline:

```
EPEVER RTC   22- 3-19 17:26:21
Rated input 100.00V 40.00A 1040.00W
Rated battery 24.00V 40.00A 1040.00W
Charge 3.64V 0.00A 0.00W [1] NoCharging,NormalInputVolt Running
Battery 25.34V 0.00A 0.00W (56.00 percent) [0] NormalTemp,NormalVolt
Battery Net 25.34V -0.28A dayVoltRange 24.19V - 29.85V
Load 25.34V 0.31A 7.85W [1]
Temp battery:17.23c inside:19.68c heatsink:19.68c battery2:24.00c
Consumed 0.12kwh Day 0.99kwh Mon 0.99kwh Year 1.33kwh Total
Generated 1.33kwh Day 10.77kwh Mon 10.77kwh Year 17.12kwh Total
 - OverVolt(Disconnect 32.00 Reconnect 30.00)
 - LowVoltage(Disconnect 22.20 Reconnect 25.20)
 - UnderVolt(Warning 24.00 Recover 24.40)
 - Charge(boost 28.80 float 27.60 equalize 29.20)
 - ChargingLimit 30.00 BoostReconnect 26.40 DischargingLimit 21.20
 - BatteryConfig 0 (USR/SEAL/GEL/FLOOD) Capacity 200Ah
ChargeConfig Equalization 0 mins boost 120 mins. Equalization period 30 days
```

When hooked up to grafana:

![Grafana graphs](https://github.com/jimmyaxod/epevermonitor/raw/master/screenshot.png "Grafana graphs")