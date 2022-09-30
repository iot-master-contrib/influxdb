package main

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/influxdata/influxdb-client-go/api/write"
	"log"
	"strings"
	"time"
)

func OpenMQTT() error {

	//物联大师 主连接
	opts := mqtt.NewClientOptions()

	opts.AddBroker(config.Broker)
	opts.SetClientID(config.ClientId)
	opts.SetUsername(config.Username)
	opts.SetPassword(config.Password)

	client := mqtt.NewClient(opts)
	token := client.Connect()
	token.Wait()
	log.Println(token.Error())

	//TODO 使用iot-master model.Service
	payload, _ := json.Marshal(map[string]string{
		"name": "history",
		"addr": "http://localhost:8088",
	})
	client.Publish("/service/register", 0, false, payload)

	//订阅消息
	client.Subscribe("/device/+/values", 0, func(client mqtt.Client, message mqtt.Message) {
		id := strings.Split(message.Topic(), "/")[2]
		var values map[string]interface{}
		_ = json.Unmarshal(message.Payload(), &values)
		writeApi.WritePoint(write.NewPoint("wattmeter", map[string]string{"id": id}, values, time.Now()))
	})

	return nil
}
