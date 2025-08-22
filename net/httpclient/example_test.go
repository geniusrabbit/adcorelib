package httpclient_test

import (
	"fmt"
	"io"
	"log"

	"github.com/geniusrabbit/adcorelib/net/httpclient/stdhttpclient"
)

// Example demonstrates the basic usage of the httpclient package.
// It shows the two-phase request pattern: create request, then execute.
func Example() {
	// Create a new HTTP client driver
	driver := stdhttpclient.NewDriver()

	// Phase 1: Create a request
	req, err := driver.Request("GET", "https://httpbin.org/get", nil)
	if err != nil {
		log.Fatal(err)
	}

	// Modify the request by setting headers
	req.SetHeader("User-Agent", "httpclient-example/1.0")
	req.SetHeader("Accept", "application/json")

	// Phase 2: Execute the request
	resp, err := driver.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Close()

	// Process the response
	fmt.Printf("Status Code: %d\n", resp.StatusCode())

	// Read response body
	body, err := io.ReadAll(resp.Body())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Response Body: %s\n", string(body))
}
