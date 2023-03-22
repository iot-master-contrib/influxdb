package main

import (
	"encoding/json"
	"github.com/iot-master-contribe/influxdb/config"
	"github.com/iot-master-contribe/influxdb/influx"
	"github.com/iot-master-contribe/influxdb/internal"
	"github.com/zgwit/iot-master/v3/pkg/log"
	"github.com/zgwit/iot-master/v3/pkg/mqtt"
	"os"
	"path/filepath"
	"strings"
)

// @title 历史数据库接口文档
// @version 1.0 版本
// @description API文档
// @BasePath /app/history/api/
// @query.collection.format multi
func main() {
	app, _ := filepath.Abs(os.Args[0])
	ext := filepath.Ext(os.Args[0])
	cfg := strings.TrimSuffix(app, ext) + ".yaml" //替换后缀名.exe为.yaml
	err := config.Load(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = log.Open(config.Config.Log)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	err = mqtt.Open(config.Config.MQTT)
	if err != nil {
		log.Fatal(err)
	}

	//注册应用
	for _, v := range config.Config.Apps {
		payload, _ := json.Marshal(v)
		_ = mqtt.Publish("master/register", payload, false, 0)
	}

	internal.SubscribeProperty(mqtt.Client)

	influx.Open(config.Config.Influxdb)

	internal.OpenWeb()
}
