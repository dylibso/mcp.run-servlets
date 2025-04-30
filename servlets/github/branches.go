package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/extism/go-pdk"
)

var (
	CreateBranchTool = ToolDescription{
		Name:        "gh-create-branch",
		Description: "Create a branch in a GitHub repository",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"owner":       prop("string", "The owner of the repository"),
				"repo":        prop("string", "The repository name"),
				"branch":      prop("string", "The branch name"),
				"from_branch": prop("string", "Source branch (defaults to `main` if not provided)"),
			},
			"required": []string{"owner", "repo", "branch", "from_branch"},
		},
	}
	ListPullRequestsTool = ToolDescription{
		Name:        "gh-list-pull-requests",
		Description: "Lists pull requests in a specified repository. Supports different response formats via accept parameter.",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"owner":     prop("string", "The account owner of the repository. The name is not case sensitive."),
				"repo":      prop("string", "The name of the repository without the .git extension. The name is not case sensitive."),
				"state":     prop("string", "Either open, closed, or all to filter by state."),
				"head":      prop("string", "Filter pulls by head user or head organization and branch name in the format of user:ref-name or organization:ref-name."),
				"base":      prop("string", "Filter pulls by base branch name. Example: gh-pages"),
				"sort":      prop("string", "What to sort results by. Can be one of: created, updated, popularity, long-running"),
				"direction": prop("string", "The direction of the sort. Default: desc when sort is created or not specified, otherwise asc"),
				"per_page":  prop("integer", "The number of results per page (max 100)"),
				"page":      prop("integer", "The page number of the results to fetch"),
				"accept":    prop("string", "Response format: raw (default), text, html, or full. Raw returns body, text returns body_text, html returns body_html, full returns all."),
			},
			"required": []string{"owner", "repo"},
		},
	}
	CreatePullRequestTool = ToolDescription{
		Name:        "gh-create-pull-request",
		Description: "Create a pull request in a GitHub repository",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"owner":                 prop("string", "The owner of the repository"),
				"repo":                  prop("string", "The repository name"),
				"title":                 prop("string", "The title of the pull request"),
				"body":                  prop("string", "The body of the pull request"),
				"head":                  prop("string", "The branch you want to merge into the base branch"),
				"base":                  prop("string", "The branch you want to merge into"),
				"draft":                 prop("boolean", "Create as draft (optional)"),
				"maintainer_can_modify": prop("boolean", "Allow maintainers to modify the pull request"),
			},
			"required": []string{"owner", "repo", "title", "body", "head", "base"},
		},
	}
)

var BranchTools = []ToolDescription{
	CreateBranchTool,
	ListPullRequestsTool,
	CreatePullRequestTool,
}

type RefObjectSchema struct {
	Sha  string `json:"sha"`
	Type string `json:"type"`
	URL  string `json:"url"`
}
type RefSchema struct {
	Ref    string          `json:"ref"`
	NodeID string          `json:"node_id"`
	URL    string          `json:"url"`
	Object RefObjectSchema `json:"object"`
}

func branchCreate(apiKey, owner, repo, branch string, fromBranch *string) CallToolResult {
	from := "main"
	if fromBranch != nil {
		from = *fromBranch
	}
	sha, err := branchGetSha(apiKey, owner, repo, from)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to get sha for branch %s: %s", from, err)),
			}},
		}
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/refs", owner, repo)
	req := pdk.NewHTTPRequest(pdk.MethodPost, url)
	req.SetHeader("Authorization", fmt.Sprintf("token %s", apiKey))
	req.SetHeader("Content-Type", "application/json")
	req.SetHeader("Accept", "application/vnd.github.v3+json")
	req.SetHeader("User-Agent", "github-mcpx-servlet")

	data := map[string]interface{}{
		"ref": fmt.Sprintf("refs/heads/%s", branch),
		"sha": sha,
	}
	res, err := json.Marshal(data)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to marshal branch data: %s", err)),
			}},
		}
	}

	req.SetBody([]byte(res))
	resp := req.Send()
	if resp.Status() != 201 {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to create branch: %d %s", resp.Status(), string(resp.Body()))),
			}},
		}
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body())),
		}},
	}
}

type PullRequestSchema struct {
	Title               string `json:"title"`
	Body                string `json:"body"`
	Head                string `json:"head"`
	Base                string `json:"base"`
	Draft               bool   `json:"draft"`
	MaintainerCanModify bool   `json:"maintainer_can_modify"`
}

func branchPullRequestSchemaFromArgs(args map[string]interface{}) PullRequestSchema {
	prs := PullRequestSchema{
		Title: args["title"].(string),
		Body:  args["body"].(string),
		Head:  args["head"].(string),
		Base:  args["base"].(string),
	}
	if draft, ok := args["draft"].(bool); ok {
		prs.Draft = draft
	}
	if canModify, ok := args["maintainer_can_modify"].(bool); ok {
		prs.MaintainerCanModify = canModify
	}
	return prs
}

func pullRequestList(apiKey string, owner, repo string, args map[string]interface{}) (CallToolResult, error) {
	baseURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls", owner, repo)
	params := make([]string, 0)

	// Handle state parameter
	if state, ok := args["state"].(string); ok && state != "" {
		switch state {
		case "open", "closed", "all":
			params = append(params, fmt.Sprintf("state=%s", state))
		}
	} else {
		params = append(params, "state=open") // Default value
	}

	// Handle head parameter (user:ref-name or organization:ref-name format)
	if head, ok := args["head"].(string); ok && head != "" {
		params = append(params, fmt.Sprintf("head=%s", head))
	}

	// Handle base parameter
	if base, ok := args["base"].(string); ok && base != "" {
		params = append(params, fmt.Sprintf("base=%s", base))
	}

	// Handle sort parameter
	sort := "created" // Default value
	if sortArg, ok := args["sort"].(string); ok && sortArg != "" {
		switch sortArg {
		case "created", "updated", "popularity", "long-running":
			sort = sortArg
		}
	}
	params = append(params, fmt.Sprintf("sort=%s", sort))

	// Handle direction parameter
	direction := "desc" // Default for created or unspecified sort
	if sort != "created" {
		direction = "asc" // Default for other sort types
	}
	if dirArg, ok := args["direction"].(string); ok {
		switch dirArg {
		case "asc", "desc":
			direction = dirArg
		}
	}
	params = append(params, fmt.Sprintf("direction=%s", direction))

	// Handle pagination
	perPage := 30 // Default value
	if perPageArg, ok := args["per_page"].(float64); ok {
		if perPageArg > 100 {
			perPage = 100 // Max value
		} else if perPageArg > 0 {
			perPage = int(perPageArg)
		}
	}
	params = append(params, fmt.Sprintf("per_page=%d", perPage))

	page := 1 // Default value
	if pageArg, ok := args["page"].(float64); ok && pageArg > 0 {
		page = int(pageArg)
	}
	params = append(params, fmt.Sprintf("page=%d", page))

	// Build final URL
	url := fmt.Sprintf("%s?%s", baseURL, strings.Join(params, "&"))
	pdk.Log(pdk.LogDebug, fmt.Sprint("Listing pull requests: ", url))

	// Make request
	req := pdk.NewHTTPRequest(pdk.MethodGet, url)
	req.SetHeader("Authorization", fmt.Sprint("token ", apiKey))

	// Handle Accept header based on requested format
	acceptHeader := "application/vnd.github+json" // Default recommended header
	if format, ok := args["accept"].(string); ok {
		switch format {
		case "raw":
			acceptHeader = "application/vnd.github.raw+json"
		case "text":
			acceptHeader = "application/vnd.github.text+json"
		case "html":
			acceptHeader = "application/vnd.github.html+json"
		case "full":
			acceptHeader = "application/vnd.github.full+json"
		}
	}
	req.SetHeader("Accept", acceptHeader)
	req.SetHeader("User-Agent", "github-mcpx-servlet")

	resp := req.Send()

	// Handle response status codes
	switch resp.Status() {
	case 200:
		return CallToolResult{
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(string(resp.Body())),
			}},
		}, nil
	case 304:
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Not modified"),
			}},
		}, nil
	case 422:
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Validation failed, or the endpoint has been spammed."),
			}},
		}, nil
	default:
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Request failed with status %d: %s", resp.Status(), string(resp.Body()))),
			}},
		}, nil
	}
}

func branchCreatePullRequest(apiKey, owner, repo string, pr PullRequestSchema) CallToolResult {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls", owner, repo)
	req := pdk.NewHTTPRequest(pdk.MethodPost, url)
	req.SetHeader("Authorization", fmt.Sprintf("token %s", apiKey))
	req.SetHeader("Accept", "application/vnd.github.v3+json")
	req.SetHeader("User-Agent", "github-mcpx-servlet")
	req.SetHeader("Content-Type", "application/json")

	res, err := json.Marshal(pr)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to marshal pull request data: %s", err)),
			}},
		}
	}

	req.SetBody([]byte(res))
	resp := req.Send()
	if resp.Status() != 201 {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to create pull request: %d %s", resp.Status(), string(resp.Body()))),
			}},
		}
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body())),
		}},
	}
}

func branchGetSha(apiKey, owner, repo, ref string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/refs/heads/%s", owner, repo, ref)
	req := pdk.NewHTTPRequest(pdk.MethodGet, url)
	req.SetHeader("Authorization", fmt.Sprintf("token %s", apiKey))
	req.SetHeader("Accept", "application/vnd.github.v3+json")
	req.SetHeader("User-Agent", "github-mcpx-servlet")

	resp := req.Send()
	if resp.Status() != 200 {
		return "", fmt.Errorf("Failed to get main branch sha: %d", resp.Status())
	}

	var refDetail RefSchema
	json.Unmarshal(resp.Body(), &refDetail)
	return refDetail.Object.Sha, nil
}
