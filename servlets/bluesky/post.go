package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/extism/go-pdk"
)

func post(args map[string]any) (CallToolResult, error) {
	if err := loadConfig(); err != nil {
		return callToolError(fmt.Sprintf("failed to load config: %s", err.Error())), nil
	}
	if err := loginSession(); err != nil {
		return callToolError(fmt.Sprintf("failed to login: %s", err.Error())), nil
	}

	if text, ok := args["text"].(string); !ok {
		return callToolError("missing text argument"), nil
	} else {
		return doPost(text, nil)
	}
}

func doPost(text string, reply *Reply) (CallToolResult, error) {
	facets, err := parseFacets(text)
	if err != nil {
		return callToolError(fmt.Sprintf("failed to parse facets: %s", err.Error())), err
	}

	url := BASE_URL + "/xrpc/com.atproto.repo.createRecord"
	req := pdk.NewHTTPRequest(pdk.MethodPost, url)
	req.SetHeader("Content-Type", "application/json")
	req.SetHeader("Authorization", "Bearer "+currentSession.AccessJwt)
	jsonBytes, err := json.Marshal(map[string]any{
		"repo":       currentSession.DID,
		"collection": "app.bsky.feed.post",
		"record": Record{
			Type:      "app.bsky.feed.post",
			Text:      text,
			CreatedAt: time.Now().Format(time.RFC3339),
			Facets:    facets,
			Reply:     reply,
		},
	})
	if err != nil {
		return callToolError(err.Error()), nil
	}
	req.SetBody(jsonBytes)
	resp := req.Send()
	if resp.Status() != http.StatusOK {
		return callToolError(fmt.Sprintf("failed to post: %d\n %s", resp.Status(), string(resp.Body()))), nil
	}
	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body())),
		}},
	}, nil
}

func getThread(args map[string]any) (CallToolResult, error) {
	if err := loadConfig(); err != nil {
		return callToolError(fmt.Sprintf("failed to load config: %s", err.Error())), nil
	}
	if err := loginSession(); err != nil {
		return callToolError(fmt.Sprintf("failed to login: %s", err.Error())), nil
	}
	uri, ok := args["uri"].(string)
	// allow web URIs and conver them to AT URIs
	uri = webUriToAT(uri)
	if !ok {
		return callToolError("missing uri"), nil
	}
	depth, ok := args["depth"].(int)
	if !ok {
		depth = 6
	}
	parentHeight, ok := args["parentHeight"].(int)
	if !ok {
		parentHeight = 80
	}
	q := url.Values{}
	q.Set("uri", uri)
	q.Set("depth", fmt.Sprintf("%d", depth))
	q.Set("parentHeight", fmt.Sprintf("%d", parentHeight))

	url := fmt.Sprintf("%s/xrpc/app.bsky.feed.getPostThread?%s", BASE_URL, q.Encode())
	req := pdk.NewHTTPRequest(pdk.MethodGet, url)
	req.SetHeader("Content-Type", "application/json")
	req.SetHeader("Authorization", "Bearer "+currentSession.AccessJwt)
	resp := req.Send()
	if resp.Status() != http.StatusOK {
		return callToolError(fmt.Sprintf("failed to get thread: %d\n %s", resp.Status(), string(resp.Body()))), nil
	}
	// Fetch the thread using the provided URI
	// Implement the logic to fetch the thread here
	return CallToolResult{Content: []Content{{
		Type: ContentTypeText,
		Text: some(string(resp.Body())),
	}}}, nil
}

type Span struct {
	Start  int    `json:"start"`
	End    int    `json:"end"`
	Handle string `json:"handle,omitempty"`
	URL    string `json:"url,omitempty"`
}

func parseMentions(text string) []Span {
	spans := []Span{}
	// regex based on: https://atproto.com/specs/handle#handle-identifier-syntax
	// mentionRegex := regexp.MustCompile(`[$|\W](@([a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)`)
	mentionRegex := regexp.MustCompile(`(@[a-zA-Z0-9.-]+)`)
	textBytes := []byte(text)

	matches := mentionRegex.FindAllSubmatchIndex(textBytes, -1)
	for _, match := range matches {
		start := match[2]
		end := match[3]
		handle := string(textBytes[start+1 : end]) // +1 to skip the @ character

		spans = append(spans, Span{
			Start:  start,
			End:    end,
			Handle: handle,
		})
	}

	return spans
}

func parseURLs(text string) []Span {
	spans := []Span{}
	// partial/naive URL regex based on: https://stackoverflow.com/a/3809435
	// tweaked to disallow some training punctuation
	textBytes := []byte(text)

	urlRegex := regexp.MustCompile(`(https?://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}[a-zA-Z0-9/\-._?=&%]*)`)
	matches := urlRegex.FindAllSubmatchIndex(textBytes, -1)
	for _, match := range matches {
		start := match[2]
		end := match[3]
		url := string(textBytes[start:end])

		spans = append(spans, Span{
			Start: start,
			End:   end,
			URL:   url,
		})
	}

	return spans
}

type ResolveHandleResponse struct {
	Did string `json:"did"`
}

type Feature struct {
	Type string `json:"$type"`
	Did  string `json:"did,omitempty"`
	URI  string `json:"uri,omitempty"`
}

type Index struct {
	ByteStart int `json:"byteStart"`
	ByteEnd   int `json:"byteEnd"`
}

type Facet struct {
	Index    Index     `json:"index"`
	Features []Feature `json:"features"`
}

// ParseFacets parses facets from text and resolves handles to DIDs
func parseFacets(text string) ([]Facet, error) {
	facets := []Facet{}

	// Process mentions
	mentions := parseMentions(text)
	for _, m := range mentions {
		// Create HTTP request to resolve handle
		req := pdk.NewHTTPRequest(pdk.MethodGet, "https://bsky.social/xrpc/com.atproto.identity.resolveHandle?handle="+m.Handle)
		resp := req.Send()

		// If the handle can't be resolved, just skip it!
		// It will be rendered as text in the post instead of a link
		if resp.Status() == 400 {
			continue
		}

		// Parse response
		var resolveResp ResolveHandleResponse
		if err := json.Unmarshal(resp.Body(), &resolveResp); err != nil {
			continue
		}

		facets = append(facets, Facet{
			Index: Index{
				ByteStart: m.Start,
				ByteEnd:   m.End,
			},
			Features: []Feature{
				{
					Type: "app.bsky.richtext.facet#mention",
					Did:  resolveResp.Did,
				},
			},
		})
	}

	// Process URLs
	urls := parseURLs(text)
	for _, u := range urls {
		facets = append(facets, Facet{
			Index: Index{
				ByteStart: u.Start,
				ByteEnd:   u.End,
			},
			Features: []Feature{
				{
					Type: "app.bsky.richtext.facet#link",
					// NOTE: URI ("I") not URL ("L")
					URI: u.URL,
				},
			},
		})
	}

	return facets, nil
}
