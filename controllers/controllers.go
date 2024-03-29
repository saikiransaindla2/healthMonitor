package controllers

import (
	"HealthMonitor/databases"
	"HealthMonitor/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var wg sync.WaitGroup
// create records of UrlData
func CreateRecords(c *gin.Context) {

	var ex []models.UrlData

	c.Bind(&ex) /////For taking the json data sent through POST request into the variable ex

	wg.Add(len(ex))

	//For every single data record, update or create a record in the table UrlData
	for i := 0; i < len(ex); i++ {

		go AddRecords(&ex[i], c)
		//var count int
		//var y models.UrlData
		//databases.Db.Model(&models.UrlData{}).Where("url = ?", ex[i].Url).Count(&count)   //////Checking for unique records
		//
		//if count == 0 {
		//	databases.Db.Save(&ex[i])       ////////Saving the record into table
		//	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "urlData record inserted successfully!", "resourceId": ex[i].ID})
		//	fmt.Println("Inserted a record")
		//} else {
		//	//////// Updating the record
		//	databases.Db.Where("url = ?",ex[i].Url).First(&y)
		//	y.CrawlTime=ex[i].CrawlTime
		//	y.WaitTime=ex[i].WaitTime
		//	y.Threshold=ex[i].Threshold
		//	databases.Db.Save(&y)
		//	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Updated urlData record!", "resourceId": y.ID})
		//}
	}
	wg.Wait()
}

// For creating or updating records from the data in JSON file in the table UrlData
func ReadFileData(c *gin.Context){
	var ex []models.UrlData
	p, _ := ioutil.ReadFile(c.Query("path"))
	err:=json.Unmarshal(p, &ex)
	if err == nil {
		wg.Add(len(ex))
		//For every single data record, update or create a record in the table UrlData
		for i := 0; i < len(ex); i++ {
			go AddRecords(&ex[i], c)
		}
		wg.Wait()
	} else{
		panic("failed to read the JSON file")
	}
}


//Function used in Go routines for saving in database
func AddRecords(x *models.UrlData, c *gin.Context) {
	var count int
	var y models.UrlData
	databases.Db.Model(&models.UrlData{}).Where("url = ?", (*x).Url).Count(&count)   //////Checking for unique records

	if count == 0 {
		databases.Db.Save(x)       ////////Saving the record into table
		c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "urlData record inserted successfully!", "resourceId": x.ID})
		fmt.Println("Inserted a record")
	} else {
		//////// Updating the record
		databases.Db.Where("url = ?",(*x).Url).First(&y)
		y.CrawlTime=(*x).CrawlTime
		y.WaitTime=(*x).WaitTime
		y.Threshold=(*x).Threshold
		databases.Db.Save(&y)
		c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Updated urlData record!", "resourceId": y.ID})
	}
	defer wg.Done()
	///Add defer
}

// Fetch all the records of UrlData
func FetchRecords(c *gin.Context) {

	var dat []models.UrlData
	databases.Db.Find(&dat)

	if len(dat) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No record found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": dat})
}

//// Testing all the urls
func TestingFunc() {
	var urls []models.UrlData
	var x models.TestingData
	databases.Db.Find(&urls)
	databases.Db.Last(&x)   ///For checking the previous runId
	k:=0
	fmt.Println(urls)
	if x.RunId!=0{
		k=x.RunId
	}
	for _, item := range urls{
		go TestingUsingGo(item, k)

	}

}

func TestingUsingGo(item models.UrlData,k int){

	for i:=0;i<item.Threshold;i++{   ////Tries threshold number of times before giving up
		timeout := time.Duration(item.CrawlTime) * time.Second
		client := http.Client{
			Timeout: timeout,
		}
		resp, err := client.Get(item.Url)

		fmt.Println(resp, err)
		if err == nil {

			if resp.StatusCode == 200 {
				x := models.TestingData{
					UrlId:         int(item.ID),
					RunId:k+1,
					AttemptNumber: i + 1,
					Health:        "Good",
				}
				databases.Db.Save(&x)
				break
			} else {
				x := models.TestingData{
					UrlId:         int(item.ID),
					RunId:k+1,
					AttemptNumber: i + 1,
					Health:        "Bad",
				}
				databases.Db.Save(&x)
				/////Wait for item.WaitTime then next iteration
				time.Sleep(time.Duration(item.WaitTime)*time.Second)

			}
		} else{
			x := models.TestingData{
				UrlId:         int(item.ID),
				RunId:k+1,
				AttemptNumber: i + 1,
				Health:        "Bad",
			}
			databases.Db.Save(&x)
			/////Wait for item.WaitTime then next iteration
			time.Sleep(time.Duration(item.WaitTime)*time.Second)
		}
	}
}

// Shows the records of certain RunID
func FetchTestData(c *gin.Context) {
	var dat []models.TestingData
	runId := c.Param("id")

	databases.Db.Where("run_id = ?", runId).Find(&dat)

	if len(dat) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No record found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": dat})
}

