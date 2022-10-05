package internal

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"strings"
)

func OpenMQTT() error {

	//物联大师 主连接
	opts := mqtt.NewClientOptions()

	opts.AddBroker(config.MQTT.Broker)
	opts.SetClientID(config.MQTT.ClientId)
	opts.SetUsername(config.MQTT.Username)
	opts.SetPassword(config.MQTT.Password)

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
		//writeApi.WritePoint(write.NewPoint("wattmeter", map[string]string{"id": id}, values, time.Now()))
		Insert(id, values)
	})

	return nil
}
