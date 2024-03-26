package main

import (
	"log"
	"net/http"
	"sync"
	"time"
)

var urls = []string{
	"http://netflix.com",
	"http://google.com",
	"http://github.com",
}

var wg sync.WaitGroup

func main() {
	wg.Add(len(urls))
	benchmark(false)
	benchmark(true)
}

func callUrls(urls []string) {
	for _, url := range urls {
		ping(url)
	}
}

func ping(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error while trying to call %s, Error: %v", url, err)
	}
	log.Printf("Response status from %s was: %s", url, resp.Status)
}

func callUrlsConcurrently(urls []string, channel chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	pingConcurrently(urls, channel)
}

func pingConcurrently(urls []string, channel chan string) {
	go ping(urls)
}

func benchmark(concurrently bool) {
	if concurrently {
		start := time.Now()
		log.Printf("Calling concurrently? : %b", concurrently)
		callUrlsConcurrently(urls, &wg)
		end := time.Now()
		elapsed := end.Sub(start)
		log.Printf("Calling concurrently, the elapsed time was: %s", elapsed)
		return
	}

	start := time.Now()
	log.Printf("Calling concurrently? : %b", concurrently)
	callUrls(urls)
	end := time.Now()
	elapsed := end.Sub(start)
	log.Printf("Calling in order, the elapsed time was: %s", elapsed)
}
