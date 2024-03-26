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
	defer file.Close()
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
		resultsFile.WriteString(url + ", " + "Error " + strconv.Itoa(response.StatusCode) + ", " + responseURL + "\r\n")
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
