// utils/json.go

package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ParseJSON reads and decodes JSON from an http.Request into the provided struct.
// It safely reads and restores the body so it can be reused later.
func ParseJSON(r *http.Request, v interface{}) error {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}

	// Optional: log for debugging
	// fmt.Println("Raw JSON:", string(bodyBytes))

	// Restore the body for future reads (important!)
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Decode into the target struct
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return fmt.Errorf("invalid json: %w", err)
	}

	return nil
}
