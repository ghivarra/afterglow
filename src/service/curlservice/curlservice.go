package curlservice

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v3/log"
)

func CreateNewRequest(ctx context.Context, httpMethod string, url string, payloads io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, httpMethod, url, payloads)
	if err != nil {
		log.Error("failed to open curl request")
	}

	return req, err
}

func ValidateResponseBody(response *http.Response, dataInterface any) (int, error) {
	// parse body
	decoder := json.NewDecoder(response.Body)
	if dataInterface != nil {
		if err := decoder.Decode(dataInterface); err != nil {
			return response.StatusCode, fmt.Errorf("failed to parse body. Error: %v", err)
		}
	}

	return response.StatusCode, nil
}
