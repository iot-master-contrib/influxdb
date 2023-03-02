package internal

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/zgwit/iot-master/v3/model"
	"strings"
	"time"
)

func OpenMQTT() error {
	//调试MQTT
	//mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	//mqtt.CRITICAL = log.New(os.Stdout, "[CRIT] ", 0)
	//mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)
	//mqtt.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)

	//物联大师 主连接
	opts := mqtt.NewClientOptions()

	opts.AddBroker(config.MQTT.Url)
	opts.SetClientID(config.MQTT.ClientId)
	opts.SetUsername(config.MQTT.Username)
	opts.SetPassword(config.MQTT.Password)

	opts.SetOnConnectHandler(func(client mqtt.Client) {
		//注册应用
		for _, v := range config.Apps {
			payload, _ := json.Marshal(v)
			client.Publish("master/register", 0, false, payload)
		}
	})


	client := mqtt.NewClient(opts)
	token := client.Connect()
	token.Wait()
	err := token.Error()
	if err != nil {
		return err
	}


	//订阅消息
	client.Subscribe("up/property/+/+", 0, func(client mqtt.Client, message mqtt.Message) {
		names := strings.Split(message.Topic(), "/")
		pid := names[2]
		//id := names[3]

		var prop model.UpPropertyPayload
		err := json.Unmarshal(message.Payload(), &prop)
		if err != nil {
			//log
			return
		}

		//解析设备数据
		//tm := time.UnixMilli(prop.Timestamp)
		var tm time.Time
		if prop.Timestamp > 0 {
			tm = time.UnixMilli(prop.Timestamp)
		} else {
			tm = time.Now()
		}
		if prop.Properties != nil {
			for _, v := range prop.Properties {
				ts := tm
				if v.Timestamp > 0 {
					ts = time.UnixMilli(v.Timestamp)
				}
				tags := make(map[string]string)
				tags["id"] = prop.Id
				fields := make(map[string]interface{})
				fields[v.Name] = v.Value
				pt := influxdb2.NewPoint(pid, tags, fields, ts)
				writeApi.WritePoint(pt)
			}
		}

		//解析子设备数据
		if prop.Devices != nil {
			for _, dev := range prop.Devices {
				tm := time.UnixMilli(dev.Timestamp)
				for _, v := range prop.Properties {
					ts := tm
					if v.Timestamp > 0 {
						ts = time.UnixMilli(v.Timestamp)
					}
					tags := make(map[string]string)
					tags["id"] = dev.Id
					fields := make(map[string]interface{})
					fields[v.Name] = v.Value
					pt := influxdb2.NewPoint(pid, tags, fields, ts)
					writeApi.WritePoint(pt)
				}
			}
		}

		//写入
		//writeApi.Flush()
	})

	return nil
}
