package routes

import (
	"HealthMonitor/controllers"
	"github.com/gin-gonic/gin"
)

func InitRoutes(){

	router := gin.Default()
	v1 := router.Group("/HealthMonitor")
	{
		v1.POST("/", controllers.CreateRecords)
		v1.GET("/", controllers.FetchRecords)
		v1.GET("/:id", controllers.FetchSingleTodo)
	}

	router.Run()
}
