package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	requestsPerSecond := 100
	totalRequests := 100000

	wg.Add(totalRequests)

	limiter := time.Tick(time.Second / time.Duration(requestsPerSecond))

	for i := 0; i < totalRequests; i++ {
		go func() {
			<-limiter

			resp, err := http.Get("http://localhost:8080")
			if err != nil {
				fmt.Printf("Error making request: %s\n", err)
				wg.Done()
				return
			}
			defer resp.Body.Close()

			// Read and discard the response body
			_, err = io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Error reading response: %s\n", err)
			}

			wg.Done()
		}()
	}

	// Wait for all goroutines to complete
	wg.Wait()

	fmt.Println("All requests completed.")
}
