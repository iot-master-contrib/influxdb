package api

import "github.com/gin-gonic/gin"

func RegisterRoutes(app *gin.RouterGroup) {

	queryRouter(app.Group("/query"))

	rateRouter(app.Group("/rate"))

	configRouter(app.Group("/config"))
}
