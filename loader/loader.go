package loader

import (
	"fmt"
	"io"
	"net/http"
)

// FetchData retrieves JSON data from the provided URL and reports the number of bytes downloaded.
func FetchData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Report how many bytes were downloaded.
	fmt.Printf("Fetched %d bytes from %s\n", len(data), url)
	return data, nil
}
