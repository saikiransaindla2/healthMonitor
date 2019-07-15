package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/robfig/cron"
	"net/http"
	"time"
)

type (
	// urlData describes a urlData type
	urlData struct {
		gorm.Model
		Url     string `json:"url"`
		CrawlTime   int `json:"crawlTime"`
		WaitTime    int    `json:"waitTime"`
		Threshold int `json:"threshold"`
	}

	// transformedUrlData represents a formatted urlData
	transformedUrlData struct {
		ID        uint   `json:"id"`
		Url     string `json:"url"`
		CrawlTime   int `json:"crawlTime"`
		WaitTime    int    `json:"waitTime"`
		Threshold int `json:"threshold"`
	}

	// testingData
	testingData struct {
		gorm.Model
		UrlId int
		RunId int
		//AttemptsDone int   //////For later
		AttemptNumber int
		Health string
	}
)

//Not needed
/*
type Websites struct {
	Websites []Website `json:"websites"`
	////Websites is variable name => can be different name too
}

type Website struct {
	Url   string `json:"url"`
	CrawlTime   int `json:"crawlTime"`
	WaitTime    int    `json:"waitTime"`
	Threshold int `json:"threshold"`
}

 */
// From https://tutorialedge.net/golang/parsing-json-with-golang/

var db *gorm.DB
func init() {
	//open a db connection
	var err error
	db, err = gorm.Open("mysql", "root:@/health_monitor?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	} else {
		fmt.Println("Connected to database")
	}
	//Migrate the schema
	db.AutoMigrate(&urlData{})
	db.AutoMigrate(&testingData{})

	c := cron.New()
	c.AddFunc("*/10 * * * *",testingFunc)
	c.Start()

}



func main(){

	fmt.Println("Hiii")


	///For taking data from a json file
	/* jsonFile, err := os.Open("websites.json")

	if err!=nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully opened websites.json")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var webs Websites

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &webs)

	// we iterate through every user within our users array and
	// print out the user Type, their name, and their facebook url
	// as just an example
	for i := 0; i < len(webs.Websites); i++ {
		fmt.Println("Url: " + webs.Websites[i].Url)
		fmt.Println("Crawl Time: ", webs.Websites[i].CrawlTime)
		fmt.Println("Wait Time: ", webs.Websites[i].WaitTime)
		fmt.Println("Threshold: " , webs.Websites[i].Threshold)
		todo := urlData{Url: webs.Websites[i].Url, CrawlTime: webs.Websites[i].CrawlTime, WaitTime: webs.Websites[i].WaitTime, Threshold: webs.Websites[i].Threshold}
		db.Save(&todo)
	} */

	router := gin.Default()
	v1 := router.Group("/HealthMonitor")
	{
		v1.POST("/", createRecords)
		v1.GET("/", fetchRecords)
		//v1.GET("/testing", testingFunc)
		v1.GET("/:id", fetchSingleTodo)
		//v1.PUT("/:id", updateTodo)
		//v1.DELETE("/:id", deleteTodo)
	}
	router.Run()

}

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

// fetchSingleTodo fetch a single todo
func fetchSingleTodo(c *gin.Context) {
	var todo []testingData
	runId := c.Param("id")

	db.Where("run_id = ?", runId).Find(&todo)

	if len(todo) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No record found!"})
		return
	}


	//_todo := transformedTodo{ID: todo.ID, Title: todo.Title, Completed: completed}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": todo})
}