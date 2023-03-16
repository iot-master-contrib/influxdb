package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/iot-master-contribe/influxdb/api"
	"github.com/iot-master-contribe/influxdb/config"
	"log"
	"net/http"
)

func OpenWeb() {
	gin.SetMode(gin.ReleaseMode)

	//GIN初始化
	app := gin.Default()
	//app.Use(gin.Logger())

	api.RegisterRoutes(app.Group("/app/history/api"))

	log.Println("Web服务启动", config.Config.Web)
	server := &http.Server{
		Addr:    config.Config.Web.Addr,
		Handler: app,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Web服务启动错误", err)
	}
}
