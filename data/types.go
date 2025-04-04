package data

// DataType represents a single record from the JSON endpoint.
type DataType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	// Add additional fields as needed.
}
