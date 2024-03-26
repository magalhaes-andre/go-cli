package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const filePath = "./resources/urls.txt"
const resultFileName = "testedResults.txt"
const resultFilePermissions = 0666

func main() {
	readUrlsFromFile()
}

func readUrlsFromFile() {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("There was an error while trying to read the requested file:")
		fmt.Println("File Path: ", filePath)
		fmt.Println("Error: ", err)
	}
	defer file.Close()
	buf := []byte{}
	scanner := bufio.NewScanner(file)
	scanner.Buffer(buf, 2048*1024)
	for scanner.Scan() {
		ping(scanner.Text())
	}
}

func ping(url string) {
	trimmedUrl := strings.TrimSpace(url)
	response, error := http.Get(strings.TrimSpace(url))
	resultsFile, _ := os.Open(resultFileName)
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
	resultsFile.Close()
}

/*
func resultsFile(fileName string, filePermission fs.FileMode) os.File {
	os.Remove(fileName)
	resultsFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, filePermission)
	if err != nil {
		fmt.Println("There was an error while trying to create log file on path: ", filePermission)
		fmt.Println("Error: ", err)
	}
	return *resultsFile
}
*/
