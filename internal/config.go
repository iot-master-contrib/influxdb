package internal

import (
	"github.com/zgwit/iot-master/v3/model"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var config = Config{
	Web: Web{
		Addr: ":8088",
	},
	MQTT: MQTT{
		Type:     "unix",
		Url:      "iot-master.sock",
		ClientId: "",
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
	Type     string `yaml:"type"` //tcp, unix
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

func init2() {
	app, _ := filepath.Abs(os.Args[0])
	ext := filepath.Ext(os.Args[0])
	//替换后缀名.exe为.yaml
	cfg := strings.TrimSuffix(app, ext) + ".yaml"

	err := Load(cfg)
	if err != nil {
		log.Fatal(err)
	}
}

// Load 加载
func Load(filename string) error {
	// 如果没有文件，则使用默认信息创建
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		config.MQTT.Type = "unix"
		config.MQTT.Url = filepath.Join(os.TempDir(), "iot-master.sock")
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
