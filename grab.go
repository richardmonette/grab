package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"regexp"
	"flag"
	"strconv"
)

type urlSrc struct {
	url string
	index int
}

func URLWorker(id int, urlSrcsChan <-chan urlSrc, resultsChan chan<- bool) {
	for urlSrc := range urlSrcsChan {
		fmt.Println("grabbing", urlSrc.url)
		resp, err := http.Get(urlSrc.url)
		if err != nil {
			fmt.Println(err)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		path := "./" + strconv.Itoa(urlSrc.index) + ".mp3"
		fmt.Println("writing", path)
		ioutil.WriteFile(path, body, 0644)
		resultsChan <- true
	}
}

func main() {

	urlPtr := flag.String("url", "", "URL of page to query for HTTP response")
	urlRegexPtr := flag.String("urlRegex", "\"http.*?\"", "Regex for finding URLs in HTTP response")
	contentRegexPtr := flag.String("contentRegex", ".*mp3.*", "Regex for filtering URLs in HTTP response")
	numWorkersPtr := flag.Int("numWorkers", 4, "Number of workers to spawn")

	flag.Parse()

	resp, err := http.Get(*urlPtr)

	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
	}

	strBody := string(body)

	urlRegex, _ := regexp.Compile(*urlRegexPtr)

	urls := urlRegex.FindAllString(strBody, -1)

	for index, element := range urls {
		urls[index] = element[1:len(element)-1]
	}

	contentRegex , _ := regexp.Compile(*contentRegexPtr)

	urlSrcsChan := make(chan urlSrc, 64)
	resultsChan := make(chan bool, 64)

	for workerId := 0; workerId < *numWorkersPtr; workerId++ {
		go URLWorker(workerId, urlSrcsChan, resultsChan)
	}

	workUnits := 0
	
	for index, element := range urls {
		if contentRegex.MatchString(element) {
			urlSrcsChan <- urlSrc{element, index}
			workUnits++
		}
	}
	close(urlSrcsChan)

	for workUnit := 0; workUnit < workUnits; workUnit++ {
		<-resultsChan
	}

}
