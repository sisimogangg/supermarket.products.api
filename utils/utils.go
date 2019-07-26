package utils

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

// Message creates a message wrapper to send over http
func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

// Respond sends a message over http
func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// EnableCors configures cors for the application
func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// GetRequest gets data from server
func GetRequest(ctx context.Context, URL string) (*io.ReadCloser, error) {

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		err = errors.New(res.Status)
		log.Fatal(err)
		return nil, err
	}

	return &res.Body, err
}
