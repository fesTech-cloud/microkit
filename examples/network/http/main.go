package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/festech-cloud/microkit/adapters/http"
	"github.com/festech-cloud/microkit/network"
)

func main() {
	client := http.NewClient(10 * time.Second)
	defer client.Close()

	ctx := context.Background()

	// HTTP GET with retry
	resp, err := client.Get(ctx, "https://httpbin.org/get",
		network.WithHeader("User-Agent", "microkit/1.0"),
		network.WithRetry(3, 100*time.Millisecond, 2*time.Second, 2.0),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Response: %d\n", resp.StatusCode)
}