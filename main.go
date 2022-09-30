package main

import "log"

func main() {
	err := OpenMQTT()
	if err != nil {
		log.Fatal(err)
	}

	OpenInfluxdb()

	OpenWeb()
}
