package main

import (
	"encoding/json"
	"fmt"

	"github.com/extism/go-pdk"
)

var (
	CreateIssueTool = ToolDescription{
		Name:        "create-issue",
		Description: "Create an issue on a GitHub repository",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"owner": map[string]interface{}{
					"type":        "string",
					"description": "The owner of the repository",
				},
				"repo": map[string]interface{}{
					"type":        "string",
					"description": "The repository name",
				},
				"title": map[string]interface{}{
					"type":        "string",
					"description": "The title of the issue",
				},
				"body": map[string]interface{}{
					"type":        "string",
					"description": "The body of the issue",
				},
			},
			"required": []string{"owner", "repo", "title", "body"},
		},
	}
	GetIssueTool = ToolDescription{
		Name:        "get-issue",
		Description: "Get an issue from a GitHub repository",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"owner": map[string]interface{}{
					"type":        "string",
					"description": "The owner of the repository",
				},
				"repo": map[string]interface{}{
					"type":        "string",
					"description": "The repository name",
				},
				"issue": map[string]interface{}{
					"type":        "integer",
					"description": "The issue number",
				},
			},
			"required": []string{"owner", "repo", "issue"},
		},
	}
	AddIssueCommentTool = ToolDescription{
		Name:        "add-issue-comment",
		Description: "Add a comment to an issue in a GitHub repository",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"owner": map[string]interface{}{
					"type":        "string",
					"description": "The owner of the repository",
				},
				"repo": map[string]interface{}{
					"type":        "string",
					"description": "The repository name",
				},
				"issue": map[string]interface{}{
					"type":        "integer",
					"description": "The issue number",
				},
				"body": map[string]interface{}{
					"type":        "string",
					"description": "The body of the comment",
				},
			},
			"required": []string{"owner", "repo", "issue", "body"},
		},
	}
	IssueTools = []ToolDescription{
		CreateIssueTool,
		GetIssueTool,
		AddIssueCommentTool,
	}
)

func issueCreate(apiKey string, owner, repo, title, body string) (CallToolResult, error) {
	url := fmt.Sprint("https://api.github.com/repos/", owner, "/", repo, "/issues")
	pdk.Log(pdk.LogDebug, fmt.Sprint("Adding comment: ", url))

	req := pdk.NewHTTPRequest(pdk.MethodPost, url)
	req.SetHeader("Authorization", fmt.Sprint("token ", apiKey))
	req.SetHeader("Accept", "application/vnd.github.v3+json")
	req.SetHeader("User-Agent", "github-mcpx-servlet")
	req.SetHeader("Content-Type", "application/json")

	res, err := json.Marshal(map[string]string{
		"title": title,
		"body":  body,
	})

	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprint("Failed to create issue: ", err)),
			}},
		}, nil
	}

	req.SetBody([]byte(res))
	resp := req.Send()

	if resp.Status() != 201 {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprint("Failed to create issue: ", resp.Status(), " ", string(resp.Body()))),
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

func issueGet(apiKey string, owner, repo string, issue int) (CallToolResult, error) {
	url := fmt.Sprint("https://api.github.com/repos/", owner, "/", repo, "/issues/", issue)
	pdk.Log(pdk.LogDebug, fmt.Sprint("Getting issue: ", url))

	req := pdk.NewHTTPRequest(pdk.MethodGet, url)
	req.SetHeader("Authorization", fmt.Sprint("token ", apiKey))
	req.SetHeader("Accept", "application/vnd.github.v3+json")
	req.SetHeader("User-Agent", "github-mcpx-servlet")
	resp := req.Send()
	if resp.Status() != 200 {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprint("Failed to get issue: ", resp.Status())),
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

func issueAddComment(apiKey string, owner, repo string, issue int, comment string) (CallToolResult, error) {
	url := fmt.Sprint("https://api.github.com/repos/", owner, "/", repo, "/issues/", issue, "/comments")
	pdk.Log(pdk.LogDebug, fmt.Sprint("Adding comment: ", url))

	req := pdk.NewHTTPRequest(pdk.MethodPost, url)
	req.SetHeader("Authorization", fmt.Sprint("token ", apiKey))
	req.SetHeader("Accept", "application/vnd.github.v3+json")
	req.SetHeader("User-Agent", "github-mcpx-servlet")
	req.SetHeader("Content-Type", "application/json")

	res, err := json.Marshal(map[string]string{
		"body": comment,
	})

	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprint("Failed to create issue: ", err)),
			}},
		}, nil
	}

	req.SetBody([]byte(res))
	resp := req.Send()

	if resp.Status() != 201 {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprint("Failed to add comment: ", resp.Status())),
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
