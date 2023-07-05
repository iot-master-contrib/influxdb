package internal

import (
	"encoding/json"
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/iot-master-contrib/influxdb/influx"
	"github.com/zgwit/iot-master/v3/payload"
	"github.com/zgwit/iot-master/v3/pkg/mqtt"
	"strings"
	"time"
)

func SubscribeProperty() {
	//订阅消息
	mqtt.SubscribeJson("up/property/+/+", func(topic string, values map[string]any) {
		topics := strings.Split(topic, "/")
		pid := topics[2]
		id := topics[3]

		tm := time.Now()
		influx.Insert(pid, id, values, tm)
		//写入
		//writeApi.Flush()
	})
}

func SubscribePropertyOld(client paho.Client) {
	//订阅消息
	client.Subscribe("up/property/+/+", 0, func(client paho.Client, message paho.Message) {
		names := strings.Split(message.Topic(), "/")
		pid := names[2]
		//id := names[3]

		var prop payload.DevicePropertyUp
		err := json.Unmarshal(message.Payload(), &prop)
		if err != nil {
			//log
			return
		}

		//解析设备数据
		//tm := time.UnixMilli(prop.Timestamp)
		var tm time.Time
		if prop.Time > 0 {
			tm = time.UnixMilli(prop.Time)
		} else {
			tm = time.Now()
		}

		if prop.Properties != nil {
			for _, v := range prop.Properties {
				ts := tm
				if v.Time > 0 {
					ts = time.UnixMilli(v.Time)
				}
				fields := map[string]any{v.Name: v.Value}
				influx.Insert(pid, prop.Id, fields, ts)
			}
		}

		//解析子设备数据
		if prop.Devices != nil {
			for _, dev := range prop.Devices {
				ts := tm
				if dev.Time > 0 {
					ts = time.UnixMilli(dev.Time)
				}
				for _, v := range prop.Properties {
					tss := ts
					if v.Time > 0 {
						tss = time.UnixMilli(v.Time)
					}
					fields := map[string]any{v.Name: v.Value}
					influx.Insert(pid, prop.Id, fields, tss)
				}
			}
		}

		//写入
		//writeApi.Flush()
	})
}
