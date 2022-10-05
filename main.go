package main

import (
	"iot-master-influxdb/internal"
	"log"
)

func main() {
	err := internal.OpenMQTT()
	if err != nil {
		log.Fatal(err)
	}

	internal.OpenInfluxdb()

	internal.OpenWeb()
}
