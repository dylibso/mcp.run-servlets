package main

import (
	"fmt"
	"strings"

	"github.com/extism/go-pdk"
)

// MetaClient represents a client for Meta Ads APIs
type MetaClient struct {
	token   string
	baseURL string
}

// NewMetaClient creates a new Meta Ads API client
func NewMetaClient(token string) *MetaClient {
	return &MetaClient{
		token:   token,
		baseURL: "https://graph.facebook.com/v23.0",
	}
}

// Response represents an API response
type Response struct {
	StatusCode int
	Body       []byte
}

// MakeRequest sends an HTTP request to the Meta API
func (c *MetaClient) MakeRequest(method pdk.HTTPMethod, path string, query map[string][]string, body []byte) (Response, error) {
	return c.MakeRequestWithHeaders(method, path, query, body, nil)
}

// MakeRequestWithHeaders sends an HTTP request with custom headers to the Meta API
func (c *MetaClient) MakeRequestWithHeaders(method pdk.HTTPMethod, path string, query map[string][]string, body []byte, headers map[string]string) (Response, error) {
	// Construct the URL
	url := fmt.Sprintf("%s/%s", c.baseURL, path)

	// Add access token to query parameters
	if query == nil {
		query = make(map[string][]string)
	}
	query["access_token"] = []string{c.token}

	if len(query) > 0 {
		parts := []string{}
		for key, values := range query {
			for _, value := range values {
				parts = append(parts, fmt.Sprintf("%s=%s", key, value))
			}
		}
		url = fmt.Sprintf("%s?%s", url, strings.Join(parts, "&"))
	}

	req := pdk.NewHTTPRequest(method, url)
	req.SetHeader("Accept", "application/json")
	req.SetHeader("Content-Type", "application/json")

	// Add custom headers if provided
	if headers != nil {
		for key, value := range headers {
			req.SetHeader(key, value)
		}
	}

	if body != nil {
		req.SetBody(body)
	}

	// Send the request
	pdk.Log(pdk.LogDebug, "Sending request...")
	resp := req.Send()
	statusCode := resp.Status()
	responseBody := resp.Body()

	pdk.Log(pdk.LogDebug, fmt.Sprintf("Received response with status: %d", statusCode))

	// Handle error status codes
	if statusCode >= 400 {
		errorMsg := fmt.Sprintf("Request failed with status %d: %s", statusCode, string(responseBody))
		pdk.Log(pdk.LogError, errorMsg)
		return Response{
			StatusCode: int(statusCode),
			Body:       resp.Body(),
		}, fmt.Errorf(errorMsg)
	}

	// Return successful response
	return Response{
		StatusCode: int(statusCode),
		Body:       resp.Body(),
	}, nil
}

// MakeFormRequest sends a form-encoded POST request to the Meta API
func (c *MetaClient) MakeFormRequest(path string, formData map[string]string) (Response, error) {
	// Construct the URL
	url := fmt.Sprintf("%s/%s", c.baseURL, path)

	// Add access token to form data
	if formData == nil {
		formData = make(map[string]string)
	}
	formData["access_token"] = c.token

	// Convert form data to URL-encoded string
	parts := []string{}
	for key, value := range formData {
		parts = append(parts, fmt.Sprintf("%s=%s", key, value))
	}
	formBody := strings.Join(parts, "&")

	req := pdk.NewHTTPRequest(pdk.MethodPost, url)
	req.SetHeader("Accept", "application/json")
	req.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	req.SetBody([]byte(formBody))

	// Send the request
	pdk.Log(pdk.LogDebug, "Sending form request...")
	resp := req.Send()
	statusCode := resp.Status()
	responseBody := resp.Body()

	pdk.Log(pdk.LogDebug, fmt.Sprintf("Received response with status: %d", statusCode))

	// Handle error status codes
	if statusCode >= 400 {
		errorMsg := fmt.Sprintf("Request failed with status %d: %s", statusCode, string(responseBody))
		pdk.Log(pdk.LogError, errorMsg)
		return Response{
			StatusCode: int(statusCode),
			Body:       resp.Body(),
		}, fmt.Errorf(errorMsg)
	}

	// Return successful response
	return Response{
		StatusCode: int(statusCode),
		Body:       resp.Body(),
	}, nil
}
