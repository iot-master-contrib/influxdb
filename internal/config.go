package internal

import (
	"github.com/zgwit/iot-master/v3/model"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

var config = Config{
	Web: Web{
		Addr: ":60001",
	},
	MQTT: MQTT{
		Url:      "iot-master.sock", //开发时，改为:1843 方便调试
		ClientId: "iot-master-influxdb",
		Username: "",
		Password: "",
	},
	Apps: []model.App{
		{
			Id:   "$influx",
			Name: "Influxdb",
		},
		{
			Id: "$history",
		},
	},
}

type Config struct {
	Web      Web         `yaml:"web"`
	MQTT     MQTT        `yaml:"mqtt"`
	Influxdb Influxdb    `yaml:"influxdb"`
	Apps     []model.App `yaml:"apps"`
}

type Web struct {
	Addr string `yaml:"addr"`
}

type MQTT struct {
	Url      string `yaml:"url"`
	ClientId string `yaml:"client_id"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Influxdb struct {
	Url         string `yaml:"url"`
	Org         string `yaml:"org"`
	Bucket      string `yaml:"bucket"`
	Token       string `yaml:"token"`
	Measurement string `yaml:"measurement"`
}

// Load 加载
func Load(filename string) error {
	// 如果没有文件，则使用默认信息创建
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		//config.MQTT.Url = "unix://[" + url.PathEscape(filepath.Join(os.TempDir(), "iot-master.sock")) + "]"
		//config.MQTT.Url = "unix://" + url.PathEscape(filepath.Join(os.TempDir(), "iot-master.sock"))
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
		err = d.Decode(&config)
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

	err = e.Encode(&config)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
