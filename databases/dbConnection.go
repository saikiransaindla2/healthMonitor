package databases

import (
	"HealthMonitor/models"
	"fmt"
	"github.com/jinzhu/gorm"
)

var Db *gorm.DB

func InitDB(){
	var err error
	Db, err = gorm.Open("mysql", "root:@/health_monitor?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	} else {
		fmt.Println("Connected to database")
	}

	//Migrate the schema
	Db.AutoMigrate(&models.UrlData{})
	Db.AutoMigrate(&models.TestingData{})
}