package main

import (
	"HealthMonitor/controllers"
	"HealthMonitor/databases"
	"HealthMonitor/routes"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/robfig/cron"
)


func init() {

	databases.InitDB() //Database connection

	///For running the tests after a specific duration periodically
	c := cron.New()
	c.AddFunc("*/10 * * * *",controllers.TestingFunc)
	c.Start()

}



func main(){

	fmt.Println("Hiii")
	routes.InitRoutes()

}


/*
// create records of urlData
func createRecords(c *gin.Context) {

	var ex []urlData

	c.Bind(&ex) /////For taking the json data sent through POST request into the variable ex

	for i := 0; i < len(ex); i++ {
		var count int
		var y urlData
		db.Model(&urlData{}).Where("url = ?", ex[i].Url).Count(&count)   //////Checking for unique records

		if count == 0 {
			db.Save(&ex[i])       ////////Saving the record into table
			c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "urlData record inserted successfully!", "resourceId": ex[i].ID})
			fmt.Println("Inserted a record")
		} else {
			//////// Updating the records
			db.Where("url = ?",ex[i].Url).First(&y)
			y.CrawlTime=ex[i].CrawlTime
			y.WaitTime=ex[i].WaitTime
			y.Threshold=ex[i].Threshold
			db.Save(&y)
			c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Updated urlData record!", "resourceId": y.ID})
		}
	}
}

// Fetch records of urlData
func fetchRecords(c *gin.Context) {
	var dat []urlData
	var _dat []transformedUrlData

	db.Find(&dat)

	if len(dat) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}

	//transforms the todos for building a good response
	for _, item := range dat {
		_dat = append(_dat, transformedUrlData{ID: item.ID, Url: item.Url, CrawlTime: item.CrawlTime, WaitTime: item.WaitTime, Threshold: item.Threshold})
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _dat})
}

//// Testing all the urls
func testingFunc() {
	var urls []urlData
	var x []testingData
	db.Find(&urls)
	db.Last(&x)
	k:=0
	fmt.Println(urls)
	if len(x)!=0{
		k=x[0].RunId
	}
	for _, item := range urls{
		go testingUsingGo(item, k)

	}

}

func testingUsingGo(item urlData,k int){
	for i:=0;i<item.Threshold;i++{
		timeout := time.Duration(item.CrawlTime) * time.Second
		client := http.Client{
			Timeout: timeout,
		}
		resp, err := client.Get(item.Url)

		fmt.Println(resp, err)
		if err == nil {

			if resp.StatusCode == 200 {
				x := testingData{
					UrlId:         int(item.ID),
					RunId:k+1,
					AttemptNumber: i + 1,
					Health:        "Good",
				}
				db.Save(&x)
				break
			} else {
				x := testingData{
					UrlId:         int(item.ID),
					RunId:k+1,
					AttemptNumber: i + 1,
					Health:        "Bad",
				}
				db.Save(&x)
				/////Wait for item.WaitTime then next iteration
				time.Sleep(time.Duration(item.WaitTime)*time.Second)

			}
		} else{
			x := testingData{
				UrlId:         int(item.ID),
				RunId:k+1,
				AttemptNumber: i + 1,
				Health:        "Bad",
			}
			db.Save(&x)
			/////Wait for item.WaitTime then next iteration
			time.Sleep(time.Duration(item.WaitTime)*time.Second)
		}
	}
}


*/
