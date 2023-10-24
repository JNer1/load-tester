package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

func sendRequest(url string, payload []byte, ch chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	// resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	// if err != nil {
	//     fmt.Println("Request failed:", err)
	// } else {
	//     fmt.Println("Response:", resp.Status)
	// }
	time.Sleep(time.Millisecond * 100)
	ch <- string(payload)
}

func main() {
	startTime := time.Now()
	numConnections := 10

	ch := make(chan string)
	var wg sync.WaitGroup

	for i := 0; i < numConnections; i++ {
		payload := []byte(fmt.Sprintf(`{"WG Number": %s}`, strconv.Itoa(i)))

		wg.Add(1)
		go sendRequest("", payload, ch, &wg)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var results []string

	for res := range ch {
		results = append(results, res)
	}

	fmt.Printf("Items received: %v\n", len(results))
	fmt.Printf("This took: %s\n", time.Since(startTime))

}
