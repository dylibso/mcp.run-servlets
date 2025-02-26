package main

import (
	"fmt"
	"net/http"
	"net/url"

	pdk "github.com/extism/go-pdk"
)

func search(args map[string]any) (CallToolResult, error) {
	if err := loadConfig(); err != nil {
		return callToolError(fmt.Sprintf("failed to load config: %s", err.Error())), nil
	}
	if err := loginSession(); err != nil {
		return callToolError(fmt.Sprintf("failed to login: %s", err.Error())), nil
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
	if facets, ok := args["facets"]; ok {
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
		return callToolError(fmt.Sprintf("failed to search: %d\n %s", resp.Status(), string(resp.Body()))), nil
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
