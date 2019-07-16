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
		v1.GET("/fetch/:id", controllers.FetchTestData)
		v1.GET("/readFileData/", controllers.ReadFileData) //For reading data from JSON file.
		//NOTE:  Add path of the file while calling /readFileData
	}

	router.Run()
}
