package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"sort"
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
	var responseTimes []int

	for res := range ch {
		total += res
		responseTimes = append(responseTimes, int(res))
	}

	sort.Ints(responseTimes)
	numberOfResponses := len(responseTimes)

	index99 := int(0.99 * float64(numberOfResponses-1))
	index95 := int(0.95 * float64(numberOfResponses-1))

	p99 := responseTimes[index99]
	p95 := responseTimes[index95]

	averageResponseTime := total / int64(*numConnections)

	fmt.Printf("Average Response Time: %v\n", averageResponseTime)
	fmt.Printf("99: %v\n", p99)
	fmt.Printf("95: %v\n", p95)
	fmt.Printf("Failed Requests: %v\n", failedRequests)
}
