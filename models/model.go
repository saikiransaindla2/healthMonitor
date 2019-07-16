package models

import(
	"github.com/jinzhu/gorm"
)

type (
	// urlData describes a urlData type
	UrlData struct {
		gorm.Model
		Url     string `json:"url"`
		CrawlTime   int `json:"crawlTime"`
		WaitTime    int    `json:"waitTime"`
		Threshold int `json:"threshold"`
	}

	// testingData for storing the testing data of all the urls with their testing time and health for every test run
	TestingData struct {
		gorm.Model
		UrlId int
		RunId int
		//AttemptsDone int   //////For later
		AttemptNumber int
		Health string
	}
)