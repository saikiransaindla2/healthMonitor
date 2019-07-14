package main

import (
	"fmt"
	"net/http"
	"time"

	//"time"
)

func main(){
	CrawlTime:=500
	timeout := time.Duration(CrawlTime) * time.Millisecond
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get("https://www.google.com/")

	fmt.Println(resp, err)
}