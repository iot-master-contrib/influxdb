package influxdb

import (
	"github.com/iot-master-contrib/influxdb/influx"
	"github.com/zgwit/iot-master/v4/history"
	"github.com/zgwit/iot-master/v4/pkg/config"
)

var h *influx.Historian

func Open() error {
	h = &influx.Historian{
		Url:    config.GetString(influx.MODULE, "url"),
		Org:    config.GetString(influx.MODULE, "org"),
		Bucket: config.GetString(influx.MODULE, "bucket"),
		Token:  config.GetString(influx.MODULE, "token"),
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
