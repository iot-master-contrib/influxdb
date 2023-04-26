package internal

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/zgwit/iot-master/v3/payload"
	"influxdb/influx"
	"strings"
	"time"
)

func SubscribeProperty(client mqtt.Client) {
	//订阅消息
	client.Subscribe("up/property/+/+", 0, func(client mqtt.Client, message mqtt.Message) {
		topics := strings.Split(message.Topic(), "/")
		pid := topics[2]
		id := topics[3]

		var properties map[string]interface{}
		err := json.Unmarshal(message.Payload(), &properties)
		if err != nil {
			//log
			return
		}

		tm := time.Now()
		influx.Insert(pid, id, properties, tm)
		//写入
		//writeApi.Flush()
	})
}

func SubscribePropertyOld(client mqtt.Client) {
	//订阅消息
	client.Subscribe("up/property/+/+", 0, func(client mqtt.Client, message mqtt.Message) {
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
					//} else if !v.Time.IsZero() {
					ts = time.Time(v.Time)
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
					//} else if !dev.Time.IsZero() {
					ts = time.Time(dev.Time)
				}
				for _, v := range prop.Properties {
					tss := ts
					if v.Timestamp > 0 {
						tss = time.UnixMilli(v.Timestamp)
						//} else if !v.Time.IsZero() {
						tss = time.Time(v.Time)
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
