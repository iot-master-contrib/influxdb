package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/iot-master-contribe/influxdb/config"
	"github.com/iot-master-contribe/influxdb/influx"
	"log"
	"net/http"
)

func replyList(ctx *gin.Context, data interface{}, total int64) {
	ctx.JSON(http.StatusOK, gin.H{"data": data, "total": total})
}

func replyOk(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

func replyFail(ctx *gin.Context, err string) {
	ctx.JSON(http.StatusOK, gin.H{"error": err})
}

func replyError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusOK, gin.H{"error": err.Error()})
}

func OpenWeb() {
	//gin.SetMode(gin.ReleaseMode)

	//GIN初始化
	app := gin.Default()

	app.Use(gin.Logger())

	app.GET("/$history/:pid/:id/:name", func(ctx *gin.Context) {
		pid := ctx.Param("pid")
		id := ctx.Param("id")
		key := ctx.Param("name")
		start := ctx.DefaultQuery("start", "-5h")
		end := ctx.DefaultQuery("end", "0h")
		window := ctx.DefaultQuery("window", "10m")
		fn := ctx.DefaultQuery("fn", "mean") //last

		values, err := influx.Query(pid, id, key, start, end, window, fn)
		if err != nil {
			replyError(ctx, err)
			return
		}

		replyOk(ctx, values)
	})

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
