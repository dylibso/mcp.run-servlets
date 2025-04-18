package main

import (
	"fmt"

	"net/url"

	"github.com/extism/go-pdk"
)

var (
	ListCommitsTool = ToolDescription{
		Name:        "gh-list-commits",
		Description: "List commits in a GitHub repository",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"owner":     prop("string", "The owner of the repository"),
				"repo":      prop("string", "The repository name"),
				"sha":       prop("string", "SHA or branch to start listing commits from. Default: the repository’s default branch (usually main)."),
				"path":      prop("string", "Source branch (defaults to `main` if not provided)"),
				"author":    prop("string", "GitHub username or email address to use to filter by commit author."),
				"committer": prop("string", "GitHub username or email address to use to filter by commit committer."),
				"since":     prop("string", "Only show results that were last updated after the given time. This is a timestamp in ISO 8601 format: YYYY-MM-DDTHH:MM:SSZ. Due to limitations of Git, timestamps must be between 1970-01-01 and 2099-12-31 (inclusive) or unexpected results may be returned."),
				"until":     prop("string", "Only commits before this date will be returned. This is a timestamp in ISO 8601 format: YYYY-MM-DDTHH:MM:SSZ. Due to limitations of Git, timestamps must be between 1970-01-01 and 2099-12-31 (inclusive) or unexpected results may be returned."),
				"per_page":  prop("integer", "Results per page (max 100). Defaults to 30."),
				"page":      prop("integer", "Page number of the results to fetch. Defaults to 1."),
			},
			"required": []string{"owner", "repo"},
		},
	}
	GetCommitTool = ToolDescription{
		Name:        "gh-get-commit",
		Description: "Returns the contents of a single commit reference",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"owner":    prop("string", "The owner of the repository"),
				"repo":     prop("string", "The repository name"),
				"ref":      prop("string", "The commit reference. Can be a commit SHA, branch name (heads/BRANCH_NAME), or tag name (tags/TAG_NAME). For more information, see 'Git References' in the Git documentation."),
				"per_page": prop("integer", "Results per page (max 100). Defaults to 30."),
				"page":     prop("integer", "Page number of the results to fetch. Defaults to 1."),
			},
			"required": []string{"owner", "repo", "ref"},
		},
	}
)

var CommitTools = []ToolDescription{
	ListCommitsTool,
	GetCommitTool,
}

func commitList(apiKey, owner, repo string, args map[string]interface{}) CallToolResult {
	q := url.Values{}
	if sha, ok := args["sha"].(string); ok && sha != "" {
		q.Add("sha", sha)
	}
	if path, ok := args["path"].(string); ok && path != "" {
		q.Add("path", path)
	}
	if author, ok := args["author"].(string); ok && author != "" {
		q.Add("author", author)
	}
	if committer, ok := args["committer"].(string); ok && committer != "" {
		q.Add("committer", committer)
	}
	if since, ok := args["since"].(string); ok && since != "" {
		q.Add("since", since)
	}
	if until, ok := args["until"].(string); ok && until != "" {
		q.Add("until", until)
	}
	if perPage, ok := args["per_page"].(int); ok && perPage > 0 {
		q.Add("per_page", fmt.Sprintf("%d", perPage))
	}
	if page, ok := args["page"].(int); ok && page > 0 {
		q.Add("page", fmt.Sprintf("%d", page))
	}

	u := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits?%s", owner, repo, q.Encode())
	req := pdk.NewHTTPRequest(pdk.MethodGet, u)
	req.SetHeader("Authorization", fmt.Sprintf("token %s", apiKey))
	req.SetHeader("Accept", "application/vnd.github.v3+json")
	req.SetHeader("User-Agent", "github-mcpx-servlet")
	resp := req.Send()
	switch resp.Status() {
	case 200:
		return CallToolResult{
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(string(resp.Body())),
			}},
		}
	default:
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Request failed with status %d: %s", resp.Status(), string(resp.Body()))),
			}},
		}
	}
}

func commitGet(apiKey, owner, repo, ref string, args map[string]interface{}) CallToolResult {
	q := url.Values{}
	if perPage, ok := args["per_page"].(int); ok && perPage > 0 {
		q.Add("per_page", fmt.Sprintf("%d", perPage))
	}
	if page, ok := args["page"].(int); ok && page > 0 {
		q.Add("page", fmt.Sprintf("%d", page))
	}

	u := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/%s?%s", owner, repo, ref, q.Encode())
	req := pdk.NewHTTPRequest(pdk.MethodGet, u)
	req.SetHeader("Authorization", fmt.Sprintf("tokexn %s", apiKey))
	req.SetHeader("Accept", "application/vnd.github.v3+json")
	req.SetHeader("User-Agent", "github-mcpx-servlet")

	resp := req.Send()
	switch resp.Status() {
	case 200:
		return CallToolResult{
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(string(resp.Body())),
			}},
		}
	default:
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Request failed with status %d: %s", resp.Status(), string(resp.Body()))),
			}},
		}
	}
}
