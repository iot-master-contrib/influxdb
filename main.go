package main

import (
	"embed"
	"encoding/json"
	"github.com/iot-master-contribe/influxdb/api"
	"github.com/iot-master-contribe/influxdb/config"
	"github.com/iot-master-contribe/influxdb/influx"
	"github.com/iot-master-contribe/influxdb/internal"
	"github.com/zgwit/iot-master/v3/model"
	"github.com/zgwit/iot-master/v3/pkg/banner"
	"github.com/zgwit/iot-master/v3/pkg/build"
	"github.com/zgwit/iot-master/v3/pkg/log"
	"github.com/zgwit/iot-master/v3/pkg/mqtt"
	"github.com/zgwit/iot-master/v3/pkg/web"
	"net/http"
)

//go:embed all:app/influx
var wwwFiles embed.FS

// @title 历史数据库接口文档
// @version 1.0 版本
// @description API文档
// @BasePath /app/influx/api/
// @query.collection.format multi
func main() {
	banner.Print("iot-master-plugin:influx")
	build.Print()

	config.Load()

	err := log.Open()
	if err != nil {
		log.Fatal(err)
	}

	influx.Open()

	//MQTT总线
	err = mqtt.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer mqtt.Close()

	//注册应用
	payload, _ := json.Marshal(model.App{
		Id:   "influx",
		Name: "Influxdb",
		Entries: []model.AppEntry{{
			Path: "app/influx/influx",
			Name: "Influxdb",
		}},
		Type:    "tcp",
		Address: "http://localhost" + web.GetOptions().Addr,
	})
	_ = mqtt.Publish("master/register", payload, false, 0)

	internal.SubscribeProperty(mqtt.Client)

	app := web.CreateEngine()

	//注册前端接口
	api.RegisterRoutes(app.Group("/app/influx/api"))

	//注册接口文档
	web.RegisterSwaggerDocs(app.Group("/app/influx"))

	//前端静态文件
	app.RegisterFS(http.FS(wwwFiles), "", "app/influx/index.html")

	//监听HTTP
	app.Serve()
}
