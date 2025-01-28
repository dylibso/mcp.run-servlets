package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/extism/go-pdk"
)

var (
	GetRepositoryContributorsTool = ToolDescription{
		Name:        "gh-get-repo-contributors",
		Description: "Get the list of contributors for a GitHub repository, including their contributions count and profile details",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"owner":    prop("string", "The owner of the repository"),
				"repo":     prop("string", "The repository name"),
				"per_page": prop("integer", "Number of results per page (max 100)"),
				"page":     prop("integer", "Page number for pagination"),
			},
			"required": []string{"owner", "repo"},
		},
	}
	GetRepositoryCollaboratorsTool = ToolDescription{
		Name:        "gh-get-repo-collaborators",
		Description: "Get the list of collaborators for a GitHub repository, including their permissions and profile details",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"owner":    prop("string", "The owner of the repository"),
				"repo":     prop("string", "The repository name"),
				"per_page": prop("integer", "Number of results per page (max 100)"),
				"page":     prop("integer", "Page number for pagination"),
			},
			"required": []string{"owner", "repo"},
		},
	}
	GetRepositoryDetailsTool = ToolDescription{
		Name:        "gh-get-repo-details",
		Description: "Get detailed information about a GitHub repository, including stars, forks, issues, and more",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"owner": prop("string", "The owner of the repository"),
				"repo":  prop("string", "The repository name"),
			},
			"required": []string{"owner", "repo"},
		},
	}
	ListReposTool = ToolDescription{
		Name:        "gh-list-repos",
		Description: "List repositories for a GitHub user or organization",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"username":  prop("string", "The GitHub username or organization name"),
				"type":      prop("string", "The type of repositories to list (all, owner, member)"),
				"sort":      prop("string", "The sort field (created, updated, pushed, full_name)"),
				"direction": prop("string", "The sort direction (asc or desc)"),
				"per_page":  prop("integer", "Number of results per page (max 100)"),
				"page":      prop("integer", "Page number for pagination"),
			},
			"required": []string{"username"},
		},
	}
	RepoTools = []ToolDescription{
		GetRepositoryContributorsTool,
		GetRepositoryCollaboratorsTool,
		GetRepositoryDetailsTool,
		ListReposTool,
	}
)

type Contributor struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
	Contributions     int    `json:"contributions"`
}

func reposGetContributors(apiKey string, owner, repo string, args map[string]interface{}) (CallToolResult, error) {
	baseURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contributors", owner, repo)
	params := make([]string, 0)

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

	pdk.Log(pdk.LogDebug, fmt.Sprint("Fetching contributors: ", url))

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
				Text: some(fmt.Sprintf("Failed to fetch contributors: %d %s", resp.Status(), string(resp.Body()))),
			}},
		}, nil
	}

	// Parse the response
	var contributors []Contributor
	if err := json.Unmarshal(resp.Body(), &contributors); err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to parse contributors: %s", err)),
			}},
		}, nil
	}

	// Marshal the response
	responseJSON, err := json.Marshal(contributors)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to marshal response: %s", err)),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(responseJSON)),
		}},
	}, nil
}

type Collaborator struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
	Permissions       struct {
		Admin bool `json:"admin"`
		Push  bool `json:"push"`
		Pull  bool `json:"pull"`
	} `json:"permissions"`
}

func reposGetCollaborators(apiKey string, owner, repo string, args map[string]interface{}) (CallToolResult, error) {
	baseURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/collaborators", owner, repo)
	params := make([]string, 0)

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

	pdk.Log(pdk.LogDebug, fmt.Sprint("Fetching collaborators: ", url))

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
				Text: some(fmt.Sprintf("Failed to fetch collaborators: %d %s", resp.Status(), string(resp.Body()))),
			}},
		}, nil
	}

	// Parse the response
	var collaborators []Collaborator
	if err := json.Unmarshal(resp.Body(), &collaborators); err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to parse collaborators: %s", err)),
			}},
		}, nil
	}

	// Marshal the response
	responseJSON, err := json.Marshal(collaborators)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to marshal response: %s", err)),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(responseJSON)),
		}},
	}, nil
}

type RepositoryDetails struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	Owner       struct {
		Login string `json:"login"`
	} `json:"owner"`
	HTMLURL       string `json:"html_url"`
	Stargazers    int    `json:"stargazers_count"`
	Watchers      int    `json:"watchers_count"`
	Forks         int    `json:"forks_count"`
	OpenIssues    int    `json:"open_issues_count"`
	DefaultBranch string `json:"default_branch"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	PushedAt      string `json:"pushed_at"`
}

func reposGetDetails(apiKey string, owner, repo string) (CallToolResult, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
	pdk.Log(pdk.LogDebug, fmt.Sprint("Fetching repository details: ", url))

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
				Text: some(fmt.Sprintf("Failed to fetch repository details: %d %s", resp.Status(), string(resp.Body()))),
			}},
		}, nil
	}

	var repoDetails RepositoryDetails
	if err := json.Unmarshal(resp.Body(), &repoDetails); err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to parse repository details: %s", err)),
			}},
		}, nil
	}

	responseJSON, err := json.Marshal(repoDetails)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to marshal response: %s", err)),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(responseJSON)),
		}},
	}, nil
}

func reposList(apiKey string, username string, args map[string]interface{}) (CallToolResult, error) {
	baseURL := fmt.Sprintf("https://api.github.com/users/%s/repos", username)
	params := make([]string, 0)

	// Optional parameters
	if value, ok := args["type"].(string); ok && value != "" {
		params = append(params, fmt.Sprintf("type=%s", value))
	}
	if value, ok := args["sort"].(string); ok && value != "" {
		params = append(params, fmt.Sprintf("sort=%s", value))
	}
	if value, ok := args["direction"].(string); ok && value != "" {
		params = append(params, fmt.Sprintf("direction=%s", value))
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

	pdk.Log(pdk.LogDebug, fmt.Sprint("Fetching repositories: ", url))

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
				Text: some(fmt.Sprintf("Failed to fetch repositories: %d %s", resp.Status(), string(resp.Body()))),
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
