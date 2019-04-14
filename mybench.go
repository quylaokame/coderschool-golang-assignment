package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"io"
	"io/ioutil"
	"time"
)

type responseInfo struct {
	status int
	bytes int64
	duration time.Duration
}

type summaryInfo struct {
	requested int64
	responsed int64
}

func main() {
	fmt.Println("chay roi")
	requests := flag.Int64("n", 1, "Number of requests to perform")
	concurrency := flag.Int64("c", 1, "Number of multiple requests to make at a time")
	
	fmt.Println(requests, concurrency)
	flag.Parse()

	if flag.NArg() == 0 || *requests == 0|| *requests < *concurrency {
		flag.PrintDefaults()
		os.Exit(-1)
	}
	
	link := flag.Arg(0)
	c := make(chan responseInfo)
	sum := summaryInfo{}
	
	for i := int64(0); i < *concurrency; i++{
		sum.requested++
		go checkLink(link,c)
	}

	for response := range c {
		if sum.requested < *requests{
			sum.requested++
			go checkLink(link,c)
		}

		sum.responsed++
		fmt.Println(response)
		if sum.responsed == sum.requested {

		}
	}
}

func checkLink(link string, c chan responseInfo)  {
	start := time.Now()
	res, err := http.Get(link)
	if err != nil {
		panic(err)
	}
	read, _ := io.Copy(ioutil.Discard, res.Body)

	c <- responseInfo{
		status: res.StatusCode,
		bytes: read,
		duration: time.Now().Sub(start),
	}
} 