package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

func sendRequest(url string, payload []byte, ch chan<- int64, wg *sync.WaitGroup, failedRequests *int) {
	defer wg.Done()

	reqStartTime := time.Now()

	_, err := http.Get(url)
	if err != nil {
		fmt.Println("Request Failed: ", err)
		*failedRequests++
		return
	}

	ch <- time.Since((reqStartTime)).Milliseconds()
}

func main() {

	numConnections := flag.Int("c", 1, "Number of connections")
	url := flag.String("u", "http://localhost:3000", "API endpoint to be tested")
	flag.Parse()

	failedRequests := 0

	if *numConnections <= 0 {
		fmt.Println("Must use a positive integer and not zero")
		os.Exit(1)
	}

	ch := make(chan int64)
	var wg sync.WaitGroup

	for i := 0; i < *numConnections; i++ {
		payload := []byte(fmt.Sprintf(`{"WG Number": %s}`, strconv.Itoa(i)))

		wg.Add(1)
		go sendRequest(*url, payload, ch, &wg, &failedRequests)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var total int64 = 0

	for res := range ch {
		total += res
	}

	averageResponseTime := total / int64(*numConnections)

	fmt.Printf("Average Response Time: %v\n", averageResponseTime)
	fmt.Printf("Failed Requests: %v\n", failedRequests)
}
