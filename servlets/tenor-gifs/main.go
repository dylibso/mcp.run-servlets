package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	pdk "github.com/extism/go-pdk"
)

// TenorClient handles Tenor API interactions
type TenorClient struct {
	apiKey string
}

// TenorResponse represents the Tenor API search response
type TenorResponse struct {
	Results []struct {
		ID           string                    `json:"id"`
		Title        string                    `json:"title"`
		Description  string                    `json:"content_description"`
		ItemURL      string                    `json:"itemurl"`
		MediaFormats map[string]TenorMediaItem `json:"media_formats"`
	} `json:"results"`
	Next string `json:"next"`
}

// TenorMediaItem represents a single media format in Tenor's response
type TenorMediaItem struct {
	URL      string  `json:"url"`
	Duration float64 `json:"duration"`
	Preview  string  `json:"preview"`
	Dims     []int   `json:"dims"`
	Size     int     `json:"size"`
}

// NewTenorClient creates a new Tenor API client
func NewTenorClient() (*TenorClient, error) {
	// Default API key embedded at build time
	defaultAPIKey := "AIzaSyBPEvZcufj7hklpDLUmgz2MtLyQY8XHOio"

	// Try to get API key from config first (allows users to override)
	apiKey, ok := pdk.GetConfig("API_KEY")
	if !ok || apiKey == "" {
		// Fall back to build-time key
		apiKey = defaultAPIKey
		if apiKey == "TENOR_API_KEY_PLACEHOLDER" {
			return nil, errors.New("no API key found in config and default key not set during build")
		}
	}

	return &TenorClient{
		apiKey: apiKey,
	}, nil
}

// Helper function to fetch and encode image to base64
func (c *TenorClient) fetchAndEncodeGif(url string) (string, error) {
	req := pdk.NewHTTPRequest(pdk.MethodGet, url)
	res := req.Send()

	if res.Status() != 200 {
		return "", fmt.Errorf("failed to fetch GIF: %d - %s", res.Status(), string(res.Body()))
	}

	// Encode to base64
	return base64.StdEncoding.EncodeToString(res.Body()), nil
}

func (c *TenorClient) searchGifs(query string, limit int) (*TenorResponse, error) {
	// Build the search URL
	baseURL := "https://tenor.googleapis.com/v2/search"
	params := url.Values{}
	params.Add("q", query)
	params.Add("key", c.apiKey)
	params.Add("client_key", "tenor_servlet")
	params.Add("limit", fmt.Sprintf("%d", limit))
	params.Add("media_filter", "tinygif,gif") // Request both preview and full sizes

	fullURL := baseURL + "?" + params.Encode()
	pdk.Log(pdk.LogInfo, "Making request to: "+fullURL)

	// Create request using PDK
	req := pdk.NewHTTPRequest(pdk.MethodGet, fullURL)

	// Send request
	res := req.Send()

	// Check for rate limiting or other errors
	if res.Status() == 429 {
		return nil, errors.New("rate limit exceeded - please configure your own API key")
	}
	if res.Status() != 200 {
		return nil, fmt.Errorf("Tenor API error: %d - Response: %s", res.Status(), string(res.Body()))
	}

	// Parse response
	var response TenorResponse
	if err := json.Unmarshal(res.Body(), &response); err != nil {
		return nil, fmt.Errorf("failed to parse Tenor response: %v", err)
	}

	return &response, nil
}

// Call handles all tool requests
func Call(input CallToolRequest) (CallToolResult, error) {
	if input.Params.Arguments == nil {
		return CallToolResult{}, errors.New("arguments must be provided")
	}

	args, ok := input.Params.Arguments.(map[string]interface{})
	if !ok {
		return CallToolResult{}, errors.New("invalid arguments format")
	}

	client, err := NewTenorClient()
	if err != nil {
		return CallToolResult{}, fmt.Errorf("failed to create Tenor client: %v", err)
	}

	switch input.Params.Name {
	case "gif-search":
		query, ok := args["query"].(string)
		if !ok {
			return CallToolResult{}, errors.New("query parameter required")
		}

		limit := 1 // Default limit
		if limitVal, ok := args["limit"].(float64); ok {
			limit = int(limitVal)
		}

		response, err := client.searchGifs(query, limit)
		if err != nil {
			return CallToolResult{}, err
		}

		// Format response with proper attribution
		var content []Content

		// Add header text
		headerText := fmt.Sprintf("Search results for \"%s\" (Powered by Tenor):", query)
		content = append(content, Content{
			Type: ContentTypeText,
			Text: &headerText,
		})

		// Add each GIF result
		for _, result := range response.Results {
			if mediaItem, ok := result.MediaFormats["tinygif"]; ok {
				// Fetch and encode the GIF
				encodedGif, err := client.fetchAndEncodeGif(mediaItem.URL)
				if err != nil {
					// Log the error but continue with other results
					pdk.Log(pdk.LogWarn, fmt.Sprintf("Failed to fetch GIF: %v", err))
					continue
				}

				if len(strings.TrimSpace(result.Title)) > 0 {
					content = append(content, Content{
						Type: ContentTypeText,
						Text: ptr(result.Title),
					})
				}

				content = append(content, Content{
					Type:     ContentTypeImage,
					Data:     &encodedGif,
					MimeType: ptr("image/gif"),
				})

				content = append(content, Content{
					Type: ContentTypeText,
					Text: ptr(result.Description),
				})

				content = append(content, Content{
					Type: ContentTypeText,
					Text: &result.ItemURL,
				})
			}
		}

		return CallToolResult{
			Content: content,
		}, nil

	default:
		return CallToolResult{}, fmt.Errorf("unknown tool: %s", input.Params.Name)
	}
}

// Describe implements the tool description
func Describe() (ListToolsResult, error) {
	return ListToolsResult{
		Tools: []ToolDescription{
			{
				Name:        "gif-search",
				Description: "Search for GIFs on Tenor",
				InputSchema: map[string]interface{}{
					"type":     "object",
					"required": []string{"query"},
					"properties": map[string]interface{}{
						"query": map[string]interface{}{
							"type":        "string",
							"description": "Search query for GIFs",
						},
						"limit": map[string]interface{}{
							"type":        "number",
							"description": "Maximum number of results to return (default: 1)",
							"minimum":     1,
							"maximum":     3,
						},
					},
				},
			},
		},
	}, nil
}

// Helper function for creating string pointers
func ptr(s string) *string {
	return &s
}
