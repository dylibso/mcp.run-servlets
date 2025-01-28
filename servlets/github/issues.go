package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/extism/go-pdk"
)

var (
	ListIssuesTool = ToolDescription{
		Name:        "gh-list-issues",
		Description: "List issues from a GitHub repository",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"owner":     prop("string", "The owner of the repository"),
				"repo":      prop("string", "The repository name"),
				"filter":    prop("string", "Filter by assigned, created, mentioned, subscribed, repos, all"),
				"state":     prop("string", "The state of the issues (open, closed, all)"),
				"labels":    prop("string", "A list of comma separated label names (e.g. bug,ui,@high)"),
				"sort":      prop("string", "Sort field (created, updated, comments)"),
				"direction": prop("string", "Sort direction (asc or desc)"),
				"since":     prop("string", "ISO 8601 timestamp (YYYY-MM-DDTHH:MM:SSZ)"),
				"collab":    prop("boolean", "Filter by issues that are collaborated on"),
				"orgs":      prop("boolean", "Filter by organization issues"),
				"owned":     prop("boolean", "Filter by owned issues"),
				"pulls":     prop("boolean", "Include pull requests in results"),
				"per_page":  prop("integer", "Number of results per page (max 100)"),
				"page":      prop("integer", "Page number for pagination"),
			},
			"required": []string{"owner", "repo"},
		},
	}
	CreateIssueTool = ToolDescription{
		Name:        "gh-create-issue",
		Description: "Create an issue on a GitHub repository",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"owner":     prop("string", "The owner of the repository"),
				"repo":      prop("string", "The repository name"),
				"title":     prop("string", "The title of the issue"),
				"body":      prop("string", "The body of the issue"),
				"state":     prop("string", "The state of the issue"),
				"assignees": arrprop("string", "The assignees of the issue", "string"),
				"milestone": prop("integer", "The milestone of the issue"),
			},
			"required": []string{"owner", "repo", "title", "body"},
		},
	}
	GetIssueTool = ToolDescription{
		Name:        "gh-get-issue",
		Description: "Get an issue from a GitHub repository",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"owner": prop("string", "The owner of the repository"),
				"repo":  prop("string", "The repository name"),
				"issue": prop("integer", "The issue number"),
			},
			"required": []string{"owner", "repo", "issue"},
		},
	}
	AddIssueCommentTool = ToolDescription{
		Name:        "gh-add-issue-comment",
		Description: "Add a comment to an issue in a GitHub repository",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"owner": prop("string", "The owner of the repository"),
				"repo":  prop("string", "The repository name"),
				"issue": prop("integer", "The issue number"),
				"body":  prop("string", "The body of the issue"),
			},
			"required": []string{"owner", "repo", "issue", "body"},
		},
	}
	UpdateIssueTool = ToolDescription{
		Name:        "gh-update-issue",
		Description: "Update an issue in a GitHub repository",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"owner":     prop("string", "The owner of the repository"),
				"repo":      prop("string", "The repository name"),
				"issue":     prop("integer", "The issue number"),
				"title":     prop("string", "The title of the issue"),
				"body":      prop("string", "The body of the issue"),
				"state":     prop("string", "The state of the issue"),
				"assignees": arrprop("string", "The assignees of the issue", "string"),
				"milestone": prop("integer", "The milestone of the issue"),
			},
			"required": []string{"owner", "repo", "issue"},
		},
	}
	IssueTools = []ToolDescription{
		ListIssuesTool,
		CreateIssueTool,
		GetIssueTool,
		UpdateIssueTool,
		AddIssueCommentTool,
	}
)

type Issue struct {
	Title     string   `json:"title,omitempty"`
	Body      string   `json:"body,omitempty"`
	Assignees []string `json:"assignees,omitempty"`
	Milestone int      `json:"milestone,omitempty"`
	Labels    []string `json:"labels,omitempty"`
}

func issueList(apiKey string, owner, repo string, args map[string]interface{}) (CallToolResult, error) {
	baseURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", owner, repo)
	params := make([]string, 0)

	// String parameters
	stringParams := map[string]string{
		"filter":    "assigned", // Default value
		"state":     "open",     // Default value
		"labels":    "",
		"sort":      "created", // Default value
		"direction": "desc",    // Default value
		"since":     "",
	}

	for key := range stringParams {
		if value, ok := args[key].(string); ok && value != "" {
			params = append(params, fmt.Sprintf("%s=%s", key, value))
		} else if stringParams[key] != "" {
			// Add default value if one exists
			params = append(params, fmt.Sprintf("%s=%s", key, stringParams[key]))
		}
	}

	// Boolean parameters
	boolParams := []string{"collab", "orgs", "owned", "pulls"}
	for _, param := range boolParams {
		if value, ok := args[param].(bool); ok {
			params = append(params, fmt.Sprintf("%s=%t", param, value))
		}
	}

	// Pagination parameters
	perPage := 30 // Default value
	if value, ok := args["per_page"].(float64); ok {
		if value > 100 {
			perPage = 100 // Max value
		} else if value > 0 {
			perPage = int(value)
		}
	}
	params = append(params, fmt.Sprintf("per_page=%d", perPage))

	page := 1 // Default value
	if value, ok := args["page"].(float64); ok && value > 0 {
		page = int(value)
	}
	params = append(params, fmt.Sprintf("page=%d", page))

	// Build final URL
	url := baseURL
	if len(params) > 0 {
		url = fmt.Sprintf("%s?%s", baseURL, strings.Join(params, "&"))
	}

	pdk.Log(pdk.LogDebug, fmt.Sprint("Listing issues: ", url))

	// Make request
	req := pdk.NewHTTPRequest(pdk.MethodGet, url)
	req.SetHeader("Authorization", fmt.Sprint("token ", apiKey))
	req.SetHeader("Accept", "application/vnd.github+json")
	req.SetHeader("User-Agent", "github-mcpx-servlet")

	resp := req.Send()
	if resp.Status() != 200 {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to list issues: %d %s", resp.Status(), string(resp.Body()))),
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

func issueFromArgs(args map[string]interface{}) Issue {
	data := Issue{}
	if title, ok := args["title"].(string); ok {
		data.Title = title
	}
	if body, ok := args["body"].(string); ok {
		data.Body = body
	}
	if assignees, ok := args["assignees"].([]interface{}); ok {
		for _, a := range assignees {
			data.Assignees = append(data.Assignees, a.(string))
		}
	}
	if milestone, ok := args["milestone"].(float64); ok {
		data.Milestone = int(milestone)
	}
	if labels, ok := args["labels"].([]interface{}); ok {
		for _, l := range labels {
			data.Labels = append(data.Labels, l.(string))
		}
	}
	return data
}

func issueCreate(apiKey string, owner, repo string, data Issue) (CallToolResult, error) {
	url := fmt.Sprint("https://api.github.com/repos/", owner, "/", repo, "/issues")
	pdk.Log(pdk.LogDebug, fmt.Sprint("Adding comment: ", url))

	req := pdk.NewHTTPRequest(pdk.MethodPost, url)
	req.SetHeader("Authorization", fmt.Sprint("token ", apiKey))
	req.SetHeader("Accept", "application/vnd.github.v3+json")
	req.SetHeader("User-Agent", "github-mcpx-servlet")
	req.SetHeader("Content-Type", "application/json")

	res, err := json.Marshal(data)

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

func issueUpdate(apiKey string, owner, repo string, issue int, data Issue) (CallToolResult, error) {
	url := fmt.Sprint("https://api.github.com/repos/", owner, "/", repo, "/issues/", issue)
	pdk.Log(pdk.LogDebug, fmt.Sprint("Getting issue: ", url))

	req := pdk.NewHTTPRequest(pdk.MethodPatch, url)
	req.SetHeader("Authorization", fmt.Sprint("token ", apiKey))
	req.SetHeader("Accept", "application/vnd.github.v3+json")
	req.SetHeader("User-Agent", "github-mcpx-servlet")
	req.SetHeader("Content-Type", "application/json")

	res, err := json.Marshal(data)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprint("Failed to update issue: ", err)),
			}},
		}, nil
	}

	req.SetBody([]byte(res))
	resp := req.Send()
	if resp.Status() != 200 {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprint("Failed to update issue: ", resp.Status())),
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
