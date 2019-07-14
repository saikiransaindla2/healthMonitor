package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
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
		v1.POST("/", createTodo)
		v1.GET("/", fetchAllTodo)
		v1.GET("/testing", testingFunc)
		//v1.GET("/:id", fetchSingleTodo)
		//v1.PUT("/:id", updateTodo)
		//v1.DELETE("/:id", deleteTodo)
	}
	router.Run()



}

// createTodo add a new todo
func createTodo(c *gin.Context) {

	var ex []urlData
	c.Bind(&ex)

	for  i:=0;i<len(ex);i++ {
		db.Save(&ex[i])
		c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Todo item created successfully!", "resourceId": ex[i].ID})
	}
	fmt.Println("Inserted a record")
	}

// fetchAllTodo fetch all todos
func fetchAllTodo(c *gin.Context) {
	var todos []urlData
	var _todos []transformedUrlData

	db.Find(&todos)

	if len(todos) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}

	//transforms the todos for building a good response
	for _, item := range todos {
		_todos = append(_todos, transformedUrlData{ID: item.ID, Url: item.Url, CrawlTime: item.CrawlTime, WaitTime: item.WaitTime, Threshold: item.Threshold})
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _todos})
}

func testingFunc(cn *gin.Context) {
	var urls []urlData

	db.Find(&urls)
	fmt.Println(urls)
	//CrawlTime:=500
	for _, item := range urls{
		fmt.Println(item.Threshold)
	//
		go testingUsingGo(item)
		/*
		go func(){

			for i:=0;i<item.Threshold;i++{
				timeout := time.Duration(item.CrawlTime) * time.Millisecond
				client := http.Client{
					Timeout: timeout,
				}
				resp, err := client.Get(item.Url)

				fmt.Println(resp, err)
				if err == nil {

					if resp.StatusCode == 200 {
						x := testingData{
							UrlId:         int(item.ID),
							AttemptNumber: i + 1,
							Health:        "Good",
						}
						db.Save(&x)
						break
					} else {
						x := testingData{
							UrlId:         int(item.ID),
							AttemptNumber: i + 1,
							Health:        "Bad",
						}
						db.Save(&x)
						/////Wait for item.WaitTime then next iteration
						time.Sleep(time.Duration(item.WaitTime)*time.Millisecond)
					}
				} else{
					x := testingData{
						UrlId:         int(item.ID),
						AttemptNumber: i + 1,
						Health:        "Bad",
					}
					db.Save(&x)
					/////Wait for item.WaitTime then next iteration
					time.Sleep(time.Duration(item.WaitTime)*time.Millisecond)
				}
			}

		}() */
	//
	//
	}

}

func testingUsingGo(item urlData){
	for i:=0;i<item.Threshold;i++{
		timeout := time.Duration(item.CrawlTime) * time.Millisecond
		client := http.Client{
			Timeout: timeout,
		}
		resp, err := client.Get(item.Url)

		fmt.Println(resp, err)
		if err == nil {

			if resp.StatusCode == 200 {
				x := testingData{
					UrlId:         int(item.ID),
					AttemptNumber: i + 1,
					Health:        "Good",
				}
				db.Save(&x)
				break
			} else {
				x := testingData{
					UrlId:         int(item.ID),
					AttemptNumber: i + 1,
					Health:        "Bad",
				}
				db.Save(&x)
				/////Wait for item.WaitTime then next iteration
				time.Sleep(time.Duration(item.WaitTime)*time.Millisecond)
			}
		} else{
			x := testingData{
				UrlId:         int(item.ID),
				AttemptNumber: i + 1,
				Health:        "Bad",
			}
			db.Save(&x)
			/////Wait for item.WaitTime then next iteration
			time.Sleep(time.Duration(item.WaitTime)*time.Millisecond)
		}
	}
}

/*
// fetchSingleTodo fetch a single todo
func fetchSingleTodo(c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")

	db.First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}

	completed := false
	if todo.Completed == 1 {
		completed = true
	} else {
		completed = false
	}

	_todo := transformedTodo{ID: todo.ID, Title: todo.Title, Completed: completed}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _todo})
}

// updateTodo update a todo
func updateTodo(c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")

	db.First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}

	db.Model(&todo).Update("title", c.PostForm("title"))
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	db.Model(&todo).Update("completed", completed)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo updated successfully!"})
}

// deleteTodo remove a todo
func deleteTodo(c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")

	db.First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}

	db.Delete(&todo)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo deleted successfully!"})
}
*/