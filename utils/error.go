package utils

import "fmt"

// HTTPError stores an error
type HTTPError struct {
	Status   int
	ErrorMsg string
}

func (h *HTTPError) Error() string {
	return fmt.Sprintf("Error processing the request. Server code: %v. Error: %v ", h.Status, h.ErrorMsg)
}
