package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const filePath = "./resources/urls.txt"
const resultFileName = "./resources/testedResults.csv"
const resultFilePermissions = 0666

func main() {
	urls := readUrlsFromFile()
	resultsUrls := pingUrls(urls)
	urls = nil
	fmt.Println("All urls have been reached, around " + strconv.Itoa(len((resultsUrls))) + " entries")
	fmt.Println("Starting to write csv...")
	resultsFile := resultsFile(resultFileName, resultFilePermissions)
	resultsFile.WriteString("URLs, hasRedirect, Actual Redirect" + "\r\n")
	resultsFile.WriteString(strings.Join(resultsUrls, "\n"))
	resultsFile.Close()
	time.Sleep(5 * time.Second)

}

func readUrlsFromFile() []string {
	file, err := os.Open(filePath)
	var urls []string
	if err != nil {
		fmt.Println("There was an error while trying to read the requested file:")
		fmt.Println("File Path: ", filePath)
		fmt.Println("Error: ", err)
	}
	buf := []byte{}
	scanner := bufio.NewScanner(file)
	scanner.Buffer(buf, 2048*1024)
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}
	file.Close()
	return urls
}
func pingUrls(urls []string) []string {
	fmt.Println(strconv.Itoa(len(urls)) + " Entries have been loaded, calling them now...")
	var urlsResults []string
	for index, url := range urls {
		fmt.Println(index)
		trimmedUrl := strings.TrimSpace(url)
		response, error := http.Get(strings.TrimSpace(url))
		if error != nil {
			urlsResults = append(urlsResults, trimmedUrl+", "+"ERROR ON CALL ")
		} else if error == nil && response.StatusCode == http.StatusOK || response.StatusCode == http.StatusUnauthorized {
			fmt.Println(strconv.Itoa(response.StatusCode) + " Response on " + url)
			responseURL := response.Request.URL.String()
			hasRedirect := strconv.FormatBool(strings.Compare(trimmedUrl, responseURL) != 0)
			urlsResults = append(urlsResults, trimmedUrl+", "+hasRedirect+", "+responseURL)
		} else {
			responseURL := response.Request.URL.String()
			fmt.Println("Error Response on " + response.Request.URL.String())
			urlsResults = append(urlsResults, trimmedUrl+", "+"Error "+strconv.Itoa(response.StatusCode)+", "+responseURL)
		}
	}
	return urlsResults
}

func resultsFile(fileName string, filePermission fs.FileMode) os.File {
	os.Remove(fileName)
	resultsFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, filePermission)
	if err != nil {
		fmt.Println("There was an error while trying to create log file on path: ", filePermission)
		fmt.Println("Error: ", err)
	}
	return *resultsFile
}
