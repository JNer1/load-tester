package main

import (
	"fmt"
	"net/http"
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
	numConnections := 10
	failedRequests := 0

	ch := make(chan int64)
	var wg sync.WaitGroup

	for i := 0; i < numConnections; i++ {
		payload := []byte(fmt.Sprintf(`{"WG Number": %s}`, strconv.Itoa(i)))

		wg.Add(1)
		go sendRequest("http://localhost:3000", payload, ch, &wg, &failedRequests)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var total int64 = 0

	for res := range ch {
		total += res
	}

	averageResponseTime := total / int64(numConnections)

	fmt.Printf("Average Response Time: %v\n", averageResponseTime)
	fmt.Printf("Failed Requests: %v\n", failedRequests)
}
