package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/extism/go-pdk"
)

var (
	GetFileContentsTool = ToolDescription{
		Name:        "gh-get-file-contents",
		Description: "Get the contents of a file or a directory in a GitHub repository",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"owner":  prop("string", "The owner of the repository"),
				"repo":   prop("string", "The repository name"),
				"path":   prop("string", "The path of the file"),
				"branch": prop("string", "(optional string): Branch to get contents from"),
			},
			"required": []string{"owner", "repo", "path"},
		},
	}
	CreateOrUpdateFileTool = ToolDescription{
		Name:        "gh-create-or-update-file",
		Description: "Create or update a file in a GitHub repository",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"owner":   prop("string", "The owner of the repository"),
				"repo":    prop("string", "The repository name"),
				"path":    prop("string", "The path of the file"),
				"content": prop("string", "The content of the file"),
				"message": prop("string", "The commit message"),
				"branch":  prop("string", "The branch name"),
				"sha":     prop("string", "(optional) The sha of the file, required for updates"),
			},
			"required": []string{"owner", "repo", "path", "content", "message", "branch"},
		},
	}
	PushFilesTool = ToolDescription{
		Name:        "gh-push-files",
		Description: "Push files to a GitHub repository",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"owner":   prop("string", "The owner of the repository"),
				"repo":    prop("string", "The repository name"),
				"branch":  prop("string", "The branch name to push to"),
				"message": prop("string", "The commit message"),
				"files": SchemaProperty{
					Type:        "array",
					Description: "Array of files to push",
					Items: &schema{
						"type": "object",
						"properties": props{
							"path":    prop("string", "The path of the file"),
							"content": prop("string", "The content of the file"),
						},
					},
				},
			},
		},
	}
	FileTools = []ToolDescription{
		GetFileContentsTool,
		CreateOrUpdateFileTool,
		PushFilesTool,
	}
)

type FileCreate struct {
	Content string  `json:"content"`
	Message string  `json:"message"`
	Branch  string  `json:"branch"`
	Sha     *string `json:"sha,omitempty"`
}

func fileCreateFromArgs(args map[string]interface{}) FileCreate {
	file := FileCreate{}
	if content, ok := args["content"].(string); ok {
		b64c := base64.StdEncoding.EncodeToString([]byte(content))
		file.Content = b64c
	}
	if message, ok := args["message"].(string); ok {
		file.Message = message
	}
	if branch, ok := args["branch"].(string); ok {
		file.Branch = branch
	}
	if sha, ok := args["sha"].(string); ok {
		file.Sha = some(sha)
	}
	return file
}

func filesCreateOrUpdate(apiKey string, owner string, repo string, path string, file FileCreate) (CallToolResult, error) {
	if file.Sha == nil {
		uc, err := filesGetContentsInternal(apiKey, owner, repo, path, &file.Branch)
		if err != nil {
			pdk.Log(pdk.LogDebug, "File does not exist, creating it")
		} else if !uc.isArray {
			sha := uc.FileContent.Sha
			file.Sha = &sha
		}
	}

	url := fmt.Sprint("https://api.github.com/repos/", owner, "/", repo, "/contents/", path)
	req := pdk.NewHTTPRequest(pdk.MethodPut, url)
	req.SetHeader("Authorization", fmt.Sprint("token ", apiKey))
	req.SetHeader("Accept", "application/vnd.github.v3+json")
	req.SetHeader("User-Agent", "github-mcpx-servlet")

	res, err := json.Marshal(file)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprint("Failed to marshal file: ", err)),
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
				Text: some(fmt.Sprint("Failed to create or update file: ", resp.Status(), " ", string(resp.Body()))),
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

type UnionContent struct {
	isArray           bool
	FileContent       FileContent
	DirectoryContents []DirectoryContent
}

type FileContent struct {
	Type        string `json:"type"`
	Encoding    string `json:"encoding"`
	Size        int    `json:"size"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	Content     string `json:"content"`
	Sha         string `json:"sha"`
	Url         string `json:"url"`
	GitUrl      string `json:"git_url"`
	HtmlUrl     string `json:"html_url"`
	DownloadUrl string `json:"download_url"`
}

type DirectoryContent struct {
	Type        string  `json:"type"`
	Size        int     `json:"size"`
	Name        string  `json:"name"`
	Path        string  `json:"path"`
	Sha         string  `json:"sha"`
	Url         string  `json:"url"`
	GitUrl      string  `json:"git_url"`
	HtmlUrl     string  `json:"html_url"`
	DownloadUrl *string `json:"download_url"`
}

func filesGetContents(apiKey string, owner string, repo string, path string, branch *string) CallToolResult {
	res, err := filesGetContentsInternal(apiKey, owner, repo, path, branch)
	if err == nil {
		var v []byte
		if res.isArray {
			v, err = json.Marshal(res.DirectoryContents)
		} else {
			v, err = json.Marshal(res.FileContent)
		}
		if err == nil {
			return CallToolResult{
				Content: []Content{{
					Type: ContentTypeText,
					Text: some(string(v)),
				}},
			}
		}
	}
	return CallToolResult{
		IsError: some(true),
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(err.Error()),
		}},
	}
}

func filesGetContentsInternal(apiKey string, owner string, repo string, path string, branch *string) (UnionContent, error) {
	u := fmt.Sprint("https://api.github.com/repos/", owner, "/", repo, "/contents/", path)

	params := url.Values{}
	if branch != nil {
		params.Add("ref", *branch)
	}
	u = fmt.Sprint(u, "?", params.Encode())

	req := pdk.NewHTTPRequest(pdk.MethodGet, u)
	req.SetHeader("Authorization", fmt.Sprint("token ", apiKey))
	req.SetHeader("Accept", "application/vnd.github.v3+json")
	req.SetHeader("User-Agent", "github-mcpx-servlet")

	resp := req.Send()
	if resp.Status() != 200 {
		return UnionContent{}, fmt.Errorf("Failed to get file contents: %d %s (%s)", resp.Status(), string(resp.Body()), u)
	}

	// attempt to parse this as a file
	uc := UnionContent{}
	fc := &uc.FileContent
	if err := json.Unmarshal(resp.Body(), fc); err == nil {
		base64.StdEncoding.DecodeString(fc.Content)
		// replace it with the decoded content
		fc.Content = string(fc.Content)
		return uc, nil
	} else {
		// if it's not a file, try to parse it as a directory
		d := []DirectoryContent{}
		if err := json.Unmarshal(resp.Body(), &d); err != nil {
			return UnionContent{}, fmt.Errorf("Failed to unmarshal directory contents: %w", err)
		}
		uc.DirectoryContents = d
		uc.isArray = true
		return uc, nil
	}
}

type FileOperation struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

func filePushFromArgs(args map[string]interface{}) []FileOperation {
	files := []FileOperation{}
	if f, ok := args["files"].([]interface{}); ok {
		for _, file := range f {
			if file, ok := file.(map[string]interface{}); ok {
				files = append(files, FileOperation{
					Path:    file["path"].(string),
					Content: file["content"].(string),
				})
			}
		}
	}
	return files
}

func filesPush(apiKey, owner, repo, branch, message string, files []FileOperation) CallToolResult {
	url := fmt.Sprint("https://api.github.com/repos/", owner, "/", repo, "/git/refs/heads/", branch)
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
				Text: some(fmt.Sprint("Failed to get branch: ", resp.Status())),
			}},
		}
	}

	ref := RefSchema{}
	json.Unmarshal(resp.Body(), &ref)

	commitSha := ref.Object.Sha
	if tree, err := createTree(apiKey, owner, repo, files, commitSha); err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprint("Failed to create tree: ", err)),
			}},
		}
	} else if commit, err := createCommit(apiKey, owner, repo, message, tree.Sha, []string{commitSha}); err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprint("Failed to create commit: ", err)),
			}},
		}
	} else {
		return updateRef(apiKey, owner, repo, "heads/"+branch, commit.Sha)
	}
}

type TreeSchema struct {
	BaseTree  string      `json:"base_tree,omitempty"`
	Tree      []TreeEntry `json:"tree"`
	Truncated bool        `json:"truncated,omitempty"`
	Url       string      `json:"url,omitempty"`
	Sha       string      `json:"sha,omitempty"`
}
type TreeEntry struct {
	Path    string `json:"path"`
	Mode    string `json:"mode"`
	Type    string `json:"type"`
	Content string `json:"content,omitempty"`
	Size    int    `json:"size,omitempty"`
	Sha     string `json:"sha,omitempty"`
	Url     string `json:"url,omitempty"`
}

func createTree(apiKey, owner, repo string, files []FileOperation, baseTree string) (TreeSchema, error) {
	tree := TreeSchema{
		BaseTree: baseTree,
		Tree:     []TreeEntry{},
	}

	for _, file := range files {
		tree.Tree = append(tree.Tree, TreeEntry{
			Path: file.Path, Mode: "100644", Type: "blob", Content: file.Content})
	}

	url := fmt.Sprint("https://api.github.com/repos/", owner, "/", repo, "/git/trees")
	req := pdk.NewHTTPRequest(pdk.MethodPost, url)
	req.SetHeader("Authorization", fmt.Sprint("token ", apiKey))
	req.SetHeader("Accept", "application/vnd.github.v3+json")
	req.SetHeader("User-Agent", "github-mcpx-servlet")
	req.SetHeader("Content-Type", "application/json")

	res, err := json.Marshal(tree)

	if err != nil {
		return TreeSchema{}, fmt.Errorf("Failed to marshal tree: %w", err)
	}
	req.SetBody(res)

	resp := req.Send()
	if resp.Status() != 201 {
		return TreeSchema{}, fmt.Errorf("Failed to create tree: %d %s", resp.Status(), string(resp.Body()))
	}

	ts := TreeSchema{}
	err = json.Unmarshal(resp.Body(), &res)
	return ts, err
}

type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Date  string `json:"date"`
}

type Commit struct {
	Sha       string `json:"sha"`
	NodeID    string `json:"node_id"`
	Url       string `json:"url"`
	Author    Author `json:"author"`
	Committer Author `json:"committer"`
	Message   string `json:"message"`
	Tree      []struct {
		Sha string `json:"sha"`
		Url string `json:"url"`
	} `json:"tree"`
	Parents []struct {
		Sha string `json:"sha"`
		Url string `json:"url"`
	} `json:"parents"`
}

func createCommit(apiKey, owner, repo, message, tree string, parents []string) (Commit, error) {
	commit := map[string]interface{}{
		"message": message,
		"tree":    tree,
		"parents": parents,
	}

	url := fmt.Sprint("https://api.github.com/repos/", owner, "/", repo, "/git/commits")
	req := pdk.NewHTTPRequest(pdk.MethodPost, url)
	req.SetHeader("Authorization", fmt.Sprint("token ", apiKey))
	req.SetHeader("Accept", "application/vnd.github.v3+json")
	req.SetHeader("User-Agent", "github-mcpx-servlet")
	req.SetHeader("Content-Type", "application/json")

	res, _ := json.Marshal(commit)
	req.SetBody(res)

	resp := req.Send()
	if resp.Status() != 201 {
		return Commit{}, fmt.Errorf("Failed to create commit: %d %s", resp.Status(), string(resp.Body()))
	}

	cs := Commit{}
	json.Unmarshal(resp.Body(), &cs)
	return cs, nil
}

func updateRef(apiKey, owner, repo, ref, sha string) CallToolResult {
	url := fmt.Sprint("https://api.github.com/repos/", owner, "/", repo, "/git/refs/", ref)
	req := pdk.NewHTTPRequest(pdk.MethodPatch, url)
	req.SetHeader("Authorization", fmt.Sprint("token ", apiKey))
	req.SetHeader("Accept", "application/vnd.github.v3+json")
	req.SetHeader("User-Agent", "github-mcpx-servlet")
	req.SetHeader("Content-Type", "application/json")

	res, _ := json.Marshal(map[string]any{"sha": sha, "force": true})
	req.SetBody(res)

	resp := req.Send()
	if resp.Status() != 200 {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprint("Failed to update ref: ", resp.Status(), " ", string(resp.Body()))),
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
