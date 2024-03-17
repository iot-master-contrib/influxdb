package influx

import (
	"github.com/zgwit/iot-master/v4/pkg/config"
)

const MODULE = "influxdb"

func init() {
	config.Register(MODULE, "url", "")
	config.Register(MODULE, "org", "")
	config.Register(MODULE, "bucket", "")
	config.Register(MODULE, "token", "")
}
