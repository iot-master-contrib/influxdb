package influxdb

import (
	"embed"
	"encoding/json"
	"github.com/iot-master-contrib/influxdb/api"
	_ "github.com/iot-master-contrib/influxdb/docs"
	"github.com/iot-master-contrib/influxdb/internal"
	"github.com/zgwit/iot-master/v3/model"
	"github.com/zgwit/iot-master/v3/pkg/mqtt"
	"github.com/zgwit/iot-master/v3/pkg/web"
	"net/http"
)

func App() *model.App {
	return &model.App{
		Id:   "influxdb",
		Name: "InfluxDB",
		Icon: "/app/influxdb/assets/influxdb.svg",
		Entries: []model.AppEntry{{
			Path: "app/influxdb/history",
			Name: "历史",
		}, {
			Path: "app/influxdb/setting",
			Name: "配置",
		}},
		Type:    "tcp",
		Address: "http://localhost" + web.GetOptions().Addr,
	}
}

//go:embed all:app/influxdb
var wwwFiles embed.FS

// @title 历史数据库接口文档
// @version 1.0 版本
// @description API文档
// @BasePath /app/influxdb/api/
// @query.collection.format multi
func main() {
}

func Startup(app *web.Engine) error {
	internal.SubscribeProperty(mqtt.Client)

	//注册前端接口
	api.RegisterRoutes(app.Group("/app/influxdb/api"))

	//注册接口文档
	web.RegisterSwaggerDocs(app.Group("/app/influxdb"), "influxdb")

	return nil
}

func Register() error {
	payload, _ := json.Marshal(App())
	return mqtt.Publish("master/register", payload, false, 0)
}

func Static(fs *web.FileSystem) {
	//前端静态文件
	fs.Put("/app/influxdb", http.FS(wwwFiles), "", "app/influxdb/index.html")
}

func Shutdown() error {

	//只关闭Web就行了，其他通过defer关闭

	return nil
}
