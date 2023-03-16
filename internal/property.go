package internal

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/iot-master-contribe/influxdb/influx"
	"github.com/zgwit/iot-master/v3/model"
	"strings"
	"time"
)

func SubscribeProperty(client mqtt.Client) {
	//订阅消息
	client.Subscribe("up/property/+/+", 0, func(client mqtt.Client, message mqtt.Message) {
		names := strings.Split(message.Topic(), "/")
		pid := names[2]
		//id := names[3]

		var prop model.PayloadPropertyUp
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
				} else if !v.Time.IsZero() {
					ts = v.Time
				}
				fields := map[string]any{v.Name: v.Value}
				influx.Insert(pid, prop.Id, fields, ts)
			}
		}

		//解析子设备数据
		if prop.Devices != nil {
			for _, dev := range prop.Devices {
				ts := tm
				if dev.Timestamp > 0 {
					ts = time.UnixMilli(dev.Timestamp)
				} else if !dev.Time.IsZero() {
					ts = dev.Time
				}
				for _, v := range prop.Properties {
					tss := ts
					if v.Timestamp > 0 {
						tss = time.UnixMilli(v.Timestamp)
					} else if !v.Time.IsZero() {
						tss = v.Time
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
