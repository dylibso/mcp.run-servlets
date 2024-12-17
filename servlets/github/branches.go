package main

import (
	"encoding/json"
	"fmt"

	"github.com/extism/go-pdk"
)

var (
	CreateBranchTool = ToolDescription{
		Name:        "create-branch",
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
	CreatePullRequestTool = ToolDescription{
		Name:        "create-pull-request",
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
