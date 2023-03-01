package main

import (
	"iot-master-influxdb/internal"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	app, _ := filepath.Abs(os.Args[0])
	ext := filepath.Ext(os.Args[0])
	cfg := strings.TrimSuffix(app, ext) + ".yaml" //替换后缀名.exe为.yaml
	err := internal.Load(cfg)

	if err != nil {
		log.Fatal(err)
	}

	err = internal.OpenMQTT()
	if err != nil {
		log.Fatal(err)
	}

	internal.OpenInfluxdb()

	internal.OpenWeb()
}
