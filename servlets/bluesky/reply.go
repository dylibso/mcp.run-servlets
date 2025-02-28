package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	pdk "github.com/extism/go-pdk"
)

func reply(args map[string]any) (CallToolResult, error) {
	if err := loadConfig(); err != nil {
		return callToolError(fmt.Sprintf("failed to load config: %s", err.Error())), nil
	}
	if err := loginSession(); err != nil {
		return callToolError(fmt.Sprintf("failed to login: %s", err.Error())), nil
	}

	if text, ok := args["text"].(string); !ok {
		return callToolError("missing text argument"), nil
	} else {
		replyTo, ok := args["reply_to"].(string)
		if !ok {
			return callToolError("missing reply_to argument"), nil
		}
		parentUri := webUriToAT(replyTo)
		refs, err := getReplyRefs(parentUri)
		if err != nil {
			return callToolError(fmt.Sprintf("failed to get reply refs: %s", err.Error())), nil
		}
		return doPost(text, refs)
	}
}

// webUriToAT converts each URL to an AT URI as follows:
// - https://bsky.app/profile/<DID>/post/<RKEY>
// - at://<DID>/app.bsky.feed.post/<RKEY>
func webUriToAT(uri string) string {
	if strings.HasPrefix(uri, "https://bsky.app/profile/") {
		parts := strings.Split(uri, "/")
		return fmt.Sprintf("at://%s/app.bsky.feed.post/%s", parts[3], parts[5])
	}
	return uri
}

type UriParts struct {
	Repo       string `json:"repo"`
	Collection string `json:"collection"`
	Rkey       string `json:"rkey"`
}

type ReplyValue struct {
	Reply Reply `json:"reply"`
}

type RecordResponse struct {
	URI   string     `json:"uri"`
	CID   string     `json:"cid"`
	Value ReplyValue `json:"value"`
}

// GetReplyRefs resolves the parent record and gets reply references
func getReplyRefs(parentURI string) (*Reply, error) {
	uriParts, err := parseURI(parentURI)
	if err != nil {
		return nil, err
	}

	// Get parent record
	parentResp := pdk.NewHTTPRequest(pdk.MethodGet, fmt.Sprintf(
		"https://bsky.social/xrpc/com.atproto.repo.getRecord?repo=%s&collection=%s&rkey=%s",
		uriParts.Repo, uriParts.Collection, uriParts.Rkey,
	)).Send()

	if parentResp.Status() != http.StatusOK {
		return nil, fmt.Errorf("failed to get parent record: %d, %s", parentResp.Status(), string(parentResp.Body()))
	}

	var parent RecordResponse
	if err := json.Unmarshal(parentResp.Body(), &parent); err != nil {
		return nil, err
	}

	var root RecordResponse

	// Check if parent has a reply reference
	if parent.Value.Reply.Root.URI != "" {
		// Parent is a reply, so get the root post
		rootURI := parent.Value.Reply.Root.URI
		rootParts, err := parseURI(rootURI)
		if err != nil {
			return nil, err
		}

		rootResp := pdk.NewHTTPRequest(pdk.MethodGet, fmt.Sprintf(
			"https://bsky.social/xrpc/com.atproto.repo.getRecord?repo=%s&collection=%s&rkey=%s",
			rootParts.Repo, rootParts.Collection, rootParts.Rkey,
		)).Send()

		if rootResp.Status() != http.StatusOK {
			return nil, fmt.Errorf("failed to get root record: %d, %s", rootResp.Status(), string(rootResp.Body()))
		}

		if err := json.Unmarshal(rootResp.Body(), &root); err != nil {
			return nil, err
		}
	} else {
		// The parent record is a top-level post, so it is also the root
		root = parent
	}

	return &Reply{
		Root: Post{
			URI: root.URI,
			CID: root.CID,
		},
		Parent: Post{
			URI: parent.URI,
			CID: parent.CID,
		},
	}, nil
}

// ParseURI parses a URI and returns the parts
func parseURI(uri string) (UriParts, error) {
	if strings.HasPrefix(uri, "at://") {
		parts := strings.Split(uri, "/")
		if len(parts) < 5 {
			return UriParts{}, fmt.Errorf("invalid at:// URI format: %s", uri)
		}
		return UriParts{
			Repo:       parts[2],
			Collection: parts[3],
			Rkey:       parts[4],
		}, nil
	} else if strings.HasPrefix(uri, "https://bsky.app/") {
		parts := strings.Split(uri, "/")
		if len(parts) < 7 {
			return UriParts{}, fmt.Errorf("invalid bsky.app URI format: %s", uri)
		}

		repo := parts[4]
		collection := parts[5]
		rkey := parts[6]

		// Map collection names
		switch collection {
		case "post":
			collection = "app.bsky.feed.post"
		case "lists":
			collection = "app.bsky.graph.list"
		case "feed":
			collection = "app.bsky.feed.generator"
		}

		return UriParts{
			Repo:       repo,
			Collection: collection,
			Rkey:       rkey,
		}, nil
	}

	return UriParts{}, fmt.Errorf("unhandled URI format: %s", uri)
}
