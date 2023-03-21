package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/iot-master-contribe/influxdb/influx"
	"github.com/xuri/excelize/v2"
	"github.com/zgwit/iot-master/v3/pkg/curd"
)

func rateRouter(app *gin.RouterGroup) {

	app.GET("/:pid/:id/:name", func(ctx *gin.Context) {
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

	app.GET("/:pid/:id/:name/excel", func(ctx *gin.Context) {
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

		if len(values) == 0 {
			curd.Fail(ctx, "无记录")
			return
		}

		//创建文件
		excel := excelize.NewFile()
		defer excel.Close()

		index, err := excel.NewSheet(key)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Set value of a cell.
		_ = excel.SetCellValue(key, "A1", "time")
		_ = excel.SetCellValue(key, "B1", "value")

		for k, v := range values {
			_ = excel.SetCellValue(key, fmt.Sprintf("A%d", k+1), v.Time)
			_ = excel.SetCellValue(key, fmt.Sprintf("B%d", k+1), v.Value)
		}

		// Set active sheet of the workbook.
		excel.SetActiveSheet(index)

		filename := pid + "-" + key + "-" + values[0].Time.Format("20060102150405") + "-" + values[len(values)-1].Time.Format("20060102150405") + ".xlsx"

		//下载头
		ctx.Header("Content-Type", "application/octet-stream")
		ctx.Header("Content-Disposition", "attachment; filename="+filename+".xlsx") // 用来指定下载下来的文件名
		ctx.Header("Content-Transfer-Encoding", "binary")

		_ = excel.Write(ctx.Writer)
	})

}
