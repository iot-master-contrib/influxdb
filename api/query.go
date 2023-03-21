package api

import (
	"github.com/gin-gonic/gin"
	"github.com/iot-master-contribe/influxdb/influx"
	"github.com/zgwit/iot-master/v3/pkg/curd"
)

func queryRouter(app *gin.RouterGroup) {

	app.GET("/query/:pid/:id/:name", func(ctx *gin.Context) {
		pid := ctx.Param("pid")
		id := ctx.Param("id")
		key := ctx.Param("name")

		start := ctx.DefaultQuery("start", "-5h")
		end := ctx.DefaultQuery("end", "0h")
		window := ctx.DefaultQuery("window", "10m")
		fn := ctx.DefaultQuery("fn", "mean") //last

		values, err := influx.Query(pid, id, key, start, end, window, fn)
		if err != nil {
			curd.Error(ctx, err)
			return
		}

		curd.OK(ctx, values)
	})

	app.GET("/rate/:pid/:id/:name", func(ctx *gin.Context) {
		pid := ctx.Param("pid")
		id := ctx.Param("id")
		key := ctx.Param("name")

		start := ctx.DefaultQuery("start", "-5h")
		end := ctx.DefaultQuery("end", "0h")
		window := ctx.DefaultQuery("window", "10m")

		values, err := influx.Rate(pid, id, key, start, end, window)
		if err != nil {
			curd.Error(ctx, err)
			return
		}

		curd.OK(ctx, values)
	})

}
