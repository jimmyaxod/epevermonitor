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
