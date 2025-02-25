// Note: run `go doc -all` in this package to see all of the types and functions available.
// ./pdk.gen.go contains the domain types from the host where your plugin will run.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/extism/go-pdk"
)

var (
	BASE_URL       string
	HANDLE         string
	PASSWORD       string
	currentSession Session
)

// Called when the tool is invoked.
// If you support multiple tools, you must switch on the input.params.name to detect which tool is being called.
// It takes CallToolRequest as input (The incoming tool request from the LLM)
// And returns CallToolResult (The servlet's response to the given tool call)
func Call(input CallToolRequest) (CallToolResult, error) {
	switch input.Params.Name {
	case "post":
		return post(input.Params.Arguments.(map[string]any))
	case "search":
		return search(input.Params.Arguments.(map[string]any))
	default:
		return CallToolResult{IsError: some(true), Content: []Content{{
			Type: ContentTypeText,
			Text: some("Unknown tool " + input.Params.Name),
		}}}, nil
	}
}

// Called by mcpx to understand how and why to use this tool.
// Note: Your servlet configs will not be set when this function is called,
// so do not rely on config in this function
// And returns ListToolsResult (The tools' descriptions, supporting multiple tools from a single servlet.)
func Describe() (res ListToolsResult, err error) {
	err = json.Unmarshal([]byte(`
			{
				"tools":[
					{
						"name": "post",
						"description": "Post a message to your feed",
						"inputSchema": {
							"type": "object",
							"properties": {
								"text": {
									"type": "string",
									"description": "The text of the post"
								}
							},
							"required": ["text"]
						}
					},
					{
						"name": "search",
						"description": "Search for posts on Bluesky",
						"inputSchema": {
							"type": "object",
							"properties": {
								"q": {
									"type": "string",
									"description": "Search query"
								},
								"sort": {
									"type": "string",
									"description": "Sort order"
								},
								"since": {
									"type": "string",
									"description": "Timestamp of the last seen message"
								},
								"until": {
									"type": "string",
									"description": "Timestamp of the first seen message"
								},
								"mentions": {
									"type": "array",
									"items": {
										"type": "string"
									},
									"description": "List of mentions"
								},
								"author": {
									"type": "string",
									"description": "Filter to posts by the given account. Handles are resolved to DID before query-time."
								},
								"lang": {
									"type": "string",
									"description": "Filter to posts in the given language"
								},
								"domain": {
									"type": "string",
									"description": "Filter to posts with URLs (facet links or embeds) linking to the given domain (hostname). Server may apply hostname normalization."
								},
								"url": {
									"type": "string",
									"description": "Filter to posts with links (facet links or embeds) linking to the given URL. Server may apply URL normalization."
								},
								"tag": {
									"type": "array",
									"items": {
										"type": "string"
									},
									"description": "Filter to posts with the given tag (hashtag), based on rich-text facet or tag field. Do not include the hash (#) prefix. Multiple tags can be specified, with 'AND' matching. (<=640 characters)"
								},
								"limit": {
									"type": "integer",
									"description": "Maximum number of posts to return (>=1, <=100) default: 25"
								},
								"cursor": {
									"type": "string",
									"description": "Cursor for pagination"
								}
							}
						}
					}
				]
			}
		`), &res)

	return
}

type LoginPayload struct {
	Handle   string `json:"identifier"`
	Password string `json:"password"`
}

type Session struct {
	AccessJwt string `json:"accessJwt"`
	DID       string `json:"did"`
}

func loadConfig() error {
	BASE_URL, _ = pdk.GetConfig("BASE_URL") // default https://bsky.social
	HANDLE, _ = pdk.GetConfig("HANDLE")
	PASSWORD, _ = pdk.GetConfig("APP_PASSWORD")

	if BASE_URL == "" {
		BASE_URL = "https://bsky.social"
	}
	if HANDLE == "" {
		return errors.New("handle is required")
	}
	if PASSWORD == "" {
		return errors.New("password is required")
	}
	return nil
}

func loginSession() error {
	url := BASE_URL + "/xrpc/com.atproto.server.createSession"
	req := pdk.NewHTTPRequest(pdk.MethodPost, url)
	req.SetHeader("Content-Type", "application/json")
	loginPayload := LoginPayload{
		Handle:   HANDLE,
		Password: PASSWORD,
	}
	jsonBytes, err := json.Marshal(&loginPayload)
	if err != nil {
		return err
	}
	req.SetBody(jsonBytes)
	resp := req.Send()
	if resp.Status() != http.StatusOK {
		return fmt.Errorf("failed to login: %d, %s", resp.Status(), string(resp.Body()))
	}
	body := resp.Body()
	if err := json.Unmarshal(body, &currentSession); err != nil {
		return err
	}
	pdk.Log(pdk.LogInfo, "logged in")
	return nil
}

func post(args map[string]any) (CallToolResult, error) {
	if err := loadConfig(); err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("failed to load config: %s", err.Error())),
			}},
		}, err
	}
	if err := loginSession(); err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("failed to login: %s", err.Error())),
			}},
		}, err
	}
	url := BASE_URL + "/xrpc/com.atproto.repo.createRecord"
	req := pdk.NewHTTPRequest(pdk.MethodPost, url)
	req.SetHeader("Content-Type", "application/json")
	req.SetHeader("Authorization", "Bearer "+currentSession.AccessJwt)
	pdk.Log(pdk.LogInfo, currentSession.AccessJwt)
	jsonBytes, err := json.Marshal(map[string]interface{}{
		"repo":       currentSession.DID,
		"collection": "app.bsky.feed.post",
		"record": map[string]string{
			"$type":     "app.bsky.feed.post",
			"text":      args["text"].(string),
			"createdAt": time.Now().Format(time.RFC3339),
		},
	})
	if err != nil {
		return CallToolResult{
			IsError: some(true),
		}, err
	}
	req.SetBody(jsonBytes)
	resp := req.Send()
	if resp.Status() != http.StatusOK {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("failed to post: %d\n %s", resp.Status(), string(resp.Body()))),
			}},
		}, nil
	}
	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body())),
		}},
	}, nil
}

func search(args map[string]any) (CallToolResult, error) {
	if err := loadConfig(); err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("failed to load config: %s", err.Error())),
			}},
		}, err
	}
	if err := loginSession(); err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("failed to login: %s", err.Error())),
			}},
		}, err
	}

	q := url.Values{}
	if qq, ok := args["q"].(string); ok {
		q.Add("q", qq)
	}
	if sort, ok := args["sort"].(string); ok {
		q.Add("sort", sort)
	}
	if since, ok := args["since"].(string); ok {
		q.Add("since", since)
	}
	if until, ok := args["until"].(string); ok {
		q.Add("until", until)
	}
	if facets, ok := args["facets"].([]interface{}); ok {
		q.Add("facets", fmt.Sprintf("%v", facets))
	}
	if mentions, ok := args["mentions"].([]string); ok {
		q.Add("mentions", fmt.Sprintf("%v", mentions))
	}
	if author, ok := args["author"].(string); ok {
		q.Add("author", author)
	}
	if domain, ok := args["domain"].(string); ok {
		q.Add("domain", domain)
	}
	if lang, ok := args["lang"].(string); ok {
		q.Add("lang", lang)
	}
	if url, ok := args["url"].(string); ok {
		q.Add("url", url)
	}
	if tag, ok := args["tag"].(string); ok {
		q.Add("tag", tag)
	}
	if limit, ok := args["limit"].(int); ok {
		q.Add("limit", fmt.Sprintf("%d", limit))
	}
	if cursor, ok := args["cursor"].(string); ok {
		q.Add("cursor", cursor)
	}

	url := BASE_URL + "/xrpc/app.bsky.feed.searchPosts?" + q.Encode()
	req := pdk.NewHTTPRequest(pdk.MethodGet, url)
	req.SetHeader("Content-Type", "application/json")
	req.SetHeader("Authorization", "Bearer "+currentSession.AccessJwt)

	resp := req.Send()
	if resp.Status() != http.StatusOK {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("failed to search: %d\n %s", resp.Status(), string(resp.Body()))),
			}},
		}, nil
	}
	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body())),
		}},
	}, nil
}

func some[T any](t T) *T {
	return &t
}
