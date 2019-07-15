package controllers

import (

	"HealthMonitor/databases"
	"HealthMonitor/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// create records of urlData
func CreateRecords(c *gin.Context) {

	var ex []models.UrlData

	c.Bind(&ex) /////For taking the json data sent through POST request into the variable ex

	for i := 0; i < len(ex); i++ {
		var count int
		var y models.UrlData
		databases.Db.Model(&models.UrlData{}).Where("url = ?", ex[i].Url).Count(&count)   //////Checking for unique records

		if count == 0 {
			databases.Db.Save(&ex[i])       ////////Saving the record into table
			c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "urlData record inserted successfully!", "resourceId": ex[i].ID})
			fmt.Println("Inserted a record")
		} else {
			//////// Updating the records
			databases.Db.Where("url = ?",ex[i].Url).First(&y)
			y.CrawlTime=ex[i].CrawlTime
			y.WaitTime=ex[i].WaitTime
			y.Threshold=ex[i].Threshold
			databases.Db.Save(&y)
			c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Updated urlData record!", "resourceId": y.ID})
		}
	}
}

// Fetch records of urlData
func FetchRecords(c *gin.Context) {
	var dat []models.UrlData

//fmt.Println(dat)
	databases.Db.Find(&dat)



	if len(dat) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}


	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": dat})
}

//// Testing all the urls
func TestingFunc() {
	var urls []models.UrlData
	var x []models.TestingData
	databases.Db.Find(&urls)
	databases.Db.Last(&x)
	k:=0
	fmt.Println(urls)
	if len(x)!=0{
		k=x[0].RunId
	}
	for _, item := range urls{
		go TestingUsingGo(item, k)

	}

}

func TestingUsingGo(item models.UrlData,k int){
	for i:=0;i<item.Threshold;i++{
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

// fetchSingleTodo fetch a single todo
func FetchSingleTodo(c *gin.Context) {
	var todo []models.TestingData
	runId := c.Param("id")

	databases.Db.Where("run_id = ?", runId).Find(&todo)

	if len(todo) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No record found!"})
		return
	}


	//_todo := transformedTodo{ID: todo.ID, Title: todo.Title, Completed: completed}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": todo})
}