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

	// transformedUrlData represents a formatted urlData
	transformedUrlData struct {
		ID        uint   `json:"id"`
		Url     string `json:"url"`
		CrawlTime   int `json:"crawlTime"`
		WaitTime    int    `json:"waitTime"`
		Threshold int `json:"threshold"`
	}

	// testingData
	TestingData struct {
		gorm.Model
		UrlId int
		RunId int
		//AttemptsDone int   //////For later
		AttemptNumber int
		Health string
	}
)