package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const filePath = "./resources/urls.txt"
const resultFileName = "testedResults.txt"
const resultFilePermissions = 0666

func main() {
	resultsFile := resultsFile(resultFileName, resultFilePermissions)
	urlList := readUrlsFromFile()
	for _, url := range urlList {
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

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println("There was an issue while reading the current line: ", line)
			fmt.Println("The error was: ", err)
		}

		urls = append(urls, line)
	}
	file.Close()
	return urls
}

func ping(url string, resultsFile os.File) {
	response, error := http.Get(strings.TrimSpace(url))
	if error != nil {
		fmt.Printf("Error on ping method", error)
	}
	if error == nil && response.StatusCode == http.StatusOK {
		fmt.Printf("OK Response on " + url + "/n")
		responseURL := response.Request.URL.String()
		hasRedirect := strconv.FormatBool(strings.Compare(url, responseURL) == 0)
		fmt.Println("url: " + url)
		fmt.Println("responseUrl: " + responseURL)
		fmt.Println("status: " + strconv.Itoa(response.StatusCode))
		fmt.Println("redirect: " + hasRedirect)
		resultsFile.WriteString(responseURL + ", " + hasRedirect + "/n")
	} else {
		fmt.Printf("Error Response on " + response.Request.URL.String() + "/n")
		responseURL := response.Request.URL.String()
		resultsFile.WriteString(responseURL + ", " + "Error " + strconv.Itoa(response.StatusCode) + "/n")
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
