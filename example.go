package main

import (
	"fmt"
	"net/http"
	"time"
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
// From https://tutorialedge.net/golang/parsing-json-with-golang/
*/

func main(){
	CrawlTime:=500
	timeout := time.Duration(CrawlTime) * time.Millisecond
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get("https://www.google.com/")

	fmt.Println(resp, err)

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


}