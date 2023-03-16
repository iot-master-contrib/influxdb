package config

import (
	"github.com/iot-master-contribe/influxdb/influx"
	"github.com/zgwit/iot-master/v3/model"
	"github.com/zgwit/iot-master/v3/pkg/log"
	"github.com/zgwit/iot-master/v3/pkg/mqtt"
	"gopkg.in/yaml.v2"
	"os"
	"runtime"
)

var Config = Configure{
	Web: Web{
		Addr: ":60001",
	},
	MQTT: mqtt.Default(),
	Log:  log.Default(),
	Apps: []model.App{
		{
			Id:      "influx",
			Name:    "InfluxDB",
			Address: "http://localhost:60001",
			Entries: []model.AppEntry{
				{Name: "配置", Path: "config"},
			},
		},
		{
			Id:      "history",
			Hidden:  true,
			Address: "http://localhost:60001",
		},
	},
}

type Configure struct {
	Web      Web            `json:"web"`
	MQTT     mqtt.Options   `json:"mqtt"`
	Log      log.Options    `json:"log"`
	Influxdb influx.Options `json:"influxdb"`
	Apps     []model.App    `json:"apps"`
}

type Web struct {
	Addr string `yaml:"addr"`
}

// Load 加载
func Load(filename string) error {
	// 如果没有文件，则使用默认信息创建
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		//Config.MQTT.Url = "unix://[" + url.PathEscape(filepath.Join(os.TempDir(), "iot-master.sock")) + "]"
		//Config.MQTT.Url = "unix://" + url.PathEscape(filepath.Join(os.TempDir(), "iot-master.sock"))
		if runtime.GOOS == "windows" {
			Config.MQTT.Url = ":1843"
		}
		return Store(filename)
		//return nil
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
