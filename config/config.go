package config

import (
	"github.com/iot-master-contribe/influxdb/influx"
	"github.com/zgwit/iot-master/v3/model"
	"github.com/zgwit/iot-master/v3/pkg/log"
	"github.com/zgwit/iot-master/v3/pkg/mqtt"
	"github.com/zgwit/iot-master/v3/pkg/web"
	"gopkg.in/yaml.v2"
	"os"
)

type Configure struct {
	Web      web.Options    `json:"web"`
	Mqtt     mqtt.Options   `json:"mqtt"`
	Log      log.Options    `json:"log"`
	Influxdb influx.Options `json:"influxdb"`
	Apps     []model.App    `json:"apps"`
}

var Config = Configure{
	Web:  web.Default(),
	Mqtt: mqtt.Default(),
	Log:  log.Default(),
	Apps: []model.App{
		{
			Id:      "influx",
			Name:    "InfluxDB",
			Address: "http://localhost:40003",
			Entries: []model.AppEntry{
				{Name: "历史查询", Path: ""},
				{Name: "连接配置", Path: "config"},
			},
		},
	},
}

func init() {
	Config.Web.Addr = ":40003"
}

// Load 加载
func Load(filename string) error {
	// 如果没有文件，则使用默认信息创建
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return Store(filename)
	} else {
		y, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
			return err
		}
		defer y.Close()

		d := yaml.NewDecoder(y)
		err = d.Decode(&Config)
		if err != nil {
			log.Fatal(err)
			return err
		}

		return nil
	}
}

func Store(filename string) error {
	y, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755) //os.Create(filename)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer y.Close()

	e := yaml.NewEncoder(y)
	defer e.Close()

	err = e.Encode(&Config)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
