package main

import (
	"context"
	"fmt"
	"log"
	"time"

	microhttp "github.com/festus/microkit/adapters/http"
	"github.com/festus/microkit/network"
)

func main() {
	client := microhttp.NewClient(10 * time.Second)
	defer client.Close()

	ctx := context.Background()

	// Example 1: HTTP GET with retry on errors
	fmt.Println("Example 1: Basic retry on errors")
	resp, err := client.Get(ctx, "https://httpbin.org/get",
		network.WithHeader("User-Agent", "microkit/1.0"),
		network.WithRetry(3, 100*time.Millisecond, 2*time.Second, 2.0),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Response: %d\n\n", resp.StatusCode)

	// Example 2: Retry on specific status codes (429, 503)
	fmt.Println("Example 2: Retry on status codes 429 and 503")
	resp2, err := client.Get(ctx, "https://httpbin.org/status/200",
		network.WithRetryOnStatus(3, 100*time.Millisecond, 2*time.Second, 2.0, 429, 503),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Response: %d\n", resp2.StatusCode)
}