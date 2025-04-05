package updater

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/mphartzheim/f1viewer/loader"
)

// Endpoint represents an API endpoint with a name, URL, and a parser function.
type Endpoint struct {
	Name  string
	URL   string
	Parse func([]byte) (any, error)
}

// EndpointResult holds the outcome of fetching an endpoint.
type EndpointResult struct {
	Name     string
	Data     any
	Err      error
	Duration time.Duration
	Hash     string
	Raw      []byte
}

// FetchEndpoint retrieves, hashes, and parses the data from an endpoint.
func FetchEndpoint(ep Endpoint) EndpointResult {
	start := time.Now()
	raw, err := loader.FetchData(ep.URL)
	duration := time.Since(start)
	if err != nil {
		return EndpointResult{Name: ep.Name, Err: err, Duration: duration}
	}

	// Compute SHA‑256 hash of the raw data.
	hashBytes := sha256.Sum256(raw)
	hashString := hex.EncodeToString(hashBytes[:])

	parsed, err := ep.Parse(raw)
	return EndpointResult{
		Name:     ep.Name,
		Data:     parsed,
		Err:      err,
		Duration: duration,
		Hash:     hashString,
		Raw:      raw,
	}
}

// LoadEndpoints concurrently fetches data for all endpoints and prints update information.
// It compares each endpoint's hash with the previously stored value in lastHashes.
func LoadEndpoints(endpoints []Endpoint, lastHashes map[string]string) {
	var wg sync.WaitGroup
	resultsChan := make(chan EndpointResult, len(endpoints))

	for _, ep := range endpoints {
		wg.Add(1)
		go func(ep Endpoint) {
			defer wg.Done()
			resultsChan <- FetchEndpoint(ep)
		}(ep)
	}

	wg.Wait()
	close(resultsChan)

	for res := range resultsChan {
		if res.Err != nil {
			fmt.Printf("%s: error: %v (took %s)\n", res.Name, res.Err, res.Duration)
			continue
		}
		if last, exists := lastHashes[res.Name]; exists && last == res.Hash {
			fmt.Printf("%s: no change (took %s)\n", res.Name, res.Duration)
		} else {
			lastHashes[res.Name] = res.Hash
			fmt.Printf("%s: updated (took %s, hash: %s)\n", res.Name, res.Duration, res.Hash)
			// Additional logic for updating a UI or data cache can be inserted here.
		}
	}
}
