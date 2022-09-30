package main

var config = Config{
	Web:    ":8088",
	Broker: ":1843",
}

type Config struct {
	Web      string
	Broker   string
	ClientId string
	Username string
	Password string
}
