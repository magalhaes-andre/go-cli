package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const filePath = "./resources/urls.txt"
const resultFileName = "./resources/testedResults.txt"
const resultFilePermissions = 0666

func main() {
	resultsFile := resultsFile(resultFileName, resultFilePermissions)
	urls := readUrlsFromFile()
	for index, url := range urls {
		fmt.Println(index)
		ping(url, resultsFile)
	}

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

func ping(url string, resultsFile os.File) {
	trimmedUrl := strings.TrimSpace(url)
	response, error := http.Get(strings.TrimSpace(url))
	if error == nil && response.StatusCode == http.StatusOK || response.StatusCode == http.StatusUnauthorized {
		fmt.Println(strconv.Itoa(response.StatusCode) + " Response on " + url + "\r\n")
		responseURL := response.Request.URL.String()
		hasRedirect := strconv.FormatBool(strings.Compare(trimmedUrl, responseURL) != 0)
		resultsFile.WriteString(trimmedUrl + ", " + hasRedirect + ", " + responseURL + "\r\n")
	} else {
		responseURL := response.Request.URL.String()
		fmt.Println("Error Response on " + response.Request.URL.String() + "\r\n")
		resultsFile.WriteString(trimmedUrl + ", " + "Error " + strconv.Itoa(response.StatusCode) + ", " + responseURL + "\r\n")
	}
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
