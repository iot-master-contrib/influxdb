package influxdb

import (
	"github.com/zgwit/iot-master/v4/history"
	"github.com/zgwit/iot-master/v4/pkg/config"
)

var h *Historian

func Open() error {
	h = &Historian{
		Url:    config.GetString(MODULE, "url"),
		Org:    config.GetString(MODULE, "org"),
		Bucket: config.GetString(MODULE, "bucket"),
		Token:  config.GetString(MODULE, "token"),
	}

	h.Open()

	//注册
	history.Register(h)

	return nil
}

func Close() error {

	h.Close()

	return nil
}
