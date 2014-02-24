package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"regexp"
	"flag"
	"strconv"
)

func main() {

	urlPtr := flag.String("url", "", "URL of page to query for HTTP response")
	urlRegexPtr := flag.String("urlRegex", "\"http.*?\"", "Regex for finding URLs in HTTP response")
	contentRegexPtr := flag.String("contentRegex", ".*mp3.*", "Regex for filtering URLs in HTTP response")

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

	for index, element := range urls {
		if contentRegex.MatchString(element) {
			fmt.Println("grabbing ", element)
			resp, err := http.Get(element)
			if err != nil {
				fmt.Println(err)
			}
			body, err = ioutil.ReadAll(resp.Body)
			ioutil.WriteFile("./" + strconv.Itoa(index) + ".mp3", body, 0644)
		}
	}

}
