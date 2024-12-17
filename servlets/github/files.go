package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/extism/go-pdk"
)

var (
	GetFileContentsTool = ToolDescription{
		Name:        "get-file-contents",
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
		Name:        "create-or-update-file",
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
				"sha":     prop("string", "(optional) The sha of the file, for updates"),
			},
			"required": []string{"owner", "repo", "path", "content", "message", "branch"},
		},
	}
	PushFilesTool = ToolDescription{
		Name:        "push-files",
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

type File struct {
	Content string  `json:"content"`
	Message string  `json:"message"`
	Branch  string  `json:"branch"`
	Sha     *string `json:"sha,omitempty"`
}

func fileFromArgs(args map[string]interface{}) File {
	file := File{}
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

func filesCreateOrUpdate(apiKey string, owner string, repo string, path string, file File) (CallToolResult, error) {
	if file.Sha == nil {
		res, isArray := filesGetContents(apiKey, owner, repo, path, &file.Branch)
		if res.IsError != nil && *res.IsError {
			pdk.Log(pdk.LogDebug, "File does not exist, creating it")
		} else if isArray {
			f := map[string]interface{}{}
			json.Unmarshal([]byte(*res.Content[0].Text), &f)
			sha := f["sha"].(string)
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

func filesGetContents(apiKey string, owner string, repo string, path string, branch *string) (res CallToolResult, isArray bool) {
	// url := `https://api.github.com/repos/${owner}/${repo}/contents/${path}`;
	url := fmt.Sprint("https://api.github.com/repos/", owner, "/", repo, "/contents/", path)

	if branch != nil {
		url += `?ref=${branch}`
	}

	req := pdk.NewHTTPRequest(pdk.MethodPost, url)
	req.SetHeader("Authorization", fmt.Sprint("token ", apiKey))
	req.SetHeader("Accept", "application/vnd.github.v3+json")
	req.SetHeader("User-Agent", "github-mcpx-servlet")

	resp := req.Send()
	if resp.Status() != 200 {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprint("Failed to get file contents: ", resp.Status())),
			}},
		}, false
	}

	// attempt to parse this as a file
	type content struct {
		Content string `json:"content"`
	}
	var c content
	if err := json.Unmarshal(resp.Body(), &c); err == nil {
		base64.StdEncoding.DecodeString(c.Content)
		// replace it with the decoded content
		c.Content = string(c.Content)
		res, _ := json.Marshal(c)
		return CallToolResult{
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(string(res)),
			}},
		}, false
	} else {
		// otherwise just return the result
		return CallToolResult{
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(string(resp.Body())),
			}},
		}, true
	}
}
