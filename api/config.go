package api

import (
	"github.com/gin-gonic/gin"
	"github.com/zgwit/iot-master/v3/pkg/curd"
	"influxdb/influx"
)

// @Summary 查询WEB配置
// @Schemes
// @Description 查询WEB配置
// @Tags config
// @Accept json
// @Produce json
// @Success 200 {object} ReplyData[influx.Options] 返回WEB配置
// @Router /config/influxdb [get]
func configGetInfluxdb(ctx *gin.Context) {
	curd.OK(ctx, influx.GetOptions())
}

// @Summary 修改WEB配置
// @Schemes
// @Description 修改WEB配置
// @Tags config
// @Param cfg body influx.Options true "WEB配置"
// @Accept json
// @Produce json
// @Success 200 {object} ReplyData[int]
// @Router /config/influxdb [post]
func configSetInfluxdb(ctx *gin.Context) {
	var conf influx.Options
	err := ctx.BindJSON(&conf)
	if err != nil {
		curd.Error(ctx, err)
		return
	}
	influx.SetOptions(conf)
	err = influx.Store()
	if err != nil {
		curd.Error(ctx, err)
		return
	}
	curd.OK(ctx, nil)
}

func configRouter(app *gin.RouterGroup) {

	app.POST("/influxdb", configSetInfluxdb)
	app.GET("/influxdb", configGetInfluxdb)

}
