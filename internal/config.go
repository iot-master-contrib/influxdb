package internal

import (
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
}

type Config struct {
	Web      Web      `yaml:"web"`
	MQTT     MQTT     `yaml:"mqtt"`
	Influxdb Influxdb `yaml:"influxdb"`
}

type Web struct {
	Addr string `yaml:"addr"`
}

type MQTT struct {
	Broker   string `yaml:"broker"`
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

func init() {
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
