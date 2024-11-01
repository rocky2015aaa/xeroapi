package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func callXeroApi(method, xeroApiUrl, accessToken, tenantID string, payload []byte) (map[string]interface{}, error) {
	// Create a new HTTP request
	httpReq, err := http.NewRequest(method, xeroApiUrl, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	// Set the necessary headers
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	httpReq.Header.Set("Xero-Tenant-Id", tenantID)
	httpReq.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Handle the response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON into a map or struct
	var responseData map[string]interface{}
	if err := json.Unmarshal(body, &responseData); err != nil {
		return nil, err
	}
	return responseData, nil
}
