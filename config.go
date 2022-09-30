package main

var config = Config{
	Web: Web{
		Addr: ":8088",
	},
}

type Config struct {
	Web      Web      `json:"web"`
	MQTT     MQTT     `json:"mqtt"`
	Influxdb Influxdb `json:"influxdb"`
}

type Web struct {
	Addr string `json:"addr"`
}

type MQTT struct {
	Broker   string `json:"broker"`
	ClientId string `json:"client_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Influxdb struct {
	Url         string `json:"url"`
	Org         string `json:"org"`
	Bucket      string `json:"bucket"`
	Token       string `json:"token"`
	Measurement string `json:"measurement"`
}
