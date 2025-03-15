package main

import (
	"encoding/json"
	"fmt"
	"github.com/extism/go-pdk"
)

var (
	CreateGistTool = ToolDescription{
		Name:        "gh-create-gist",
		Description: "Create a GitHub Gist",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"description": prop("string", "Description of the gist"),
				"files": SchemaProperty{
					Type:        "object",
					Description: "Files contained in the gist.",
					AdditionalProperties: &schema{
						"type": "object",
						"properties": schema{
							"content": schema{
								"type":        "string",
								"description": "Content of the file",
							},
						},
						"required": []string{"content"},
					},
				},
			},
			"required": []string{"files"},
		},
	}
	GetGistTool = ToolDescription{
		Name:        "gh-get-gist",
		Description: "Gets a specified gist.",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"gist_id": prop("string", "The unique identifier of the gist."),
			},
			"required": []string{"gist_id"},
		},
	}
	UpdateGistTool = ToolDescription{
		Name:        "gh-update-gist",
		Description: "Lists pull requests in a specified repository. Supports different response formats via accept parameter.",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"gist_id":     prop("string", "The unique identifier of the gist."),
				"description": prop("string", "Description of the gist"),
				"files": SchemaProperty{
					Type:        "object",
					Description: "Files contained in the gist.",
					AdditionalProperties: &schema{
						"type": "object",
						"properties": schema{
							"content": schema{
								"type":        "string",
								"description": "Content of the file",
							},
						},
						"required": []string{"content"},
					},
				},
			},
			"required": []string{"gist_id"},
		},
	}
	DeleteGistTool = ToolDescription{
		Name:        "gh-delete-gist",
		Description: "Delete a specified gist.",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"gist_id": prop("string", "The unique identifier of the gist."),
			},
			"required": []string{"gist_id"},
		},
	}
)

var GistTools = []ToolDescription{
	CreateGistTool,
	GetGistTool,
	UpdateGistTool,
	DeleteGistTool,
}

func gistCreate(apiKey, description string, files map[string]any) CallToolResult {
	url := "https://api.github.com/gists"
	req := pdk.NewHTTPRequest(pdk.MethodPost, url)
	req.SetHeader("Authorization", fmt.Sprintf("token %s", apiKey))
	req.SetHeader("Content-Type", "application/json")
	req.SetHeader("Accept", "application/vnd.github+json")
	req.SetHeader("User-Agent", "github-mcpx-servlet")

	data := map[string]any{
		"description": description,
		"files":       files,
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
	req.SetBody(res)
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

func gistUpdate(apiKey, gistId, description string, files map[string]any) CallToolResult {
	url := fmt.Sprintf("https://api.github.com/gists/%s", gistId)
	req := pdk.NewHTTPRequest(pdk.MethodPatch, url)
	req.SetHeader("Authorization", fmt.Sprintf("token %s", apiKey))
	req.SetHeader("Content-Type", "application/json")
	req.SetHeader("Accept", "application/vnd.github+json")
	req.SetHeader("User-Agent", "github-mcpx-servlet")

	data := map[string]any{
		"description": description,
		"files":       files,
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
	req.SetBody(res)
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

func gistGet(apiKey, gistId string) CallToolResult {
	url := fmt.Sprintf("https://api.github.com/gists/%s", gistId)
	req := pdk.NewHTTPRequest(pdk.MethodGet, url)
	req.SetHeader("Authorization", fmt.Sprintf("token %s", apiKey))
	req.SetHeader("Content-Type", "application/json")
	req.SetHeader("Accept", "application/vnd.github+json")
	req.SetHeader("User-Agent", "github-mcpx-servlet")

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

func gistDelete(apiKey, gistId string) CallToolResult {
	url := fmt.Sprintf("https://api.github.com/gists/%s", gistId)
	req := pdk.NewHTTPRequest(pdk.MethodDelete, url)
	req.SetHeader("Authorization", fmt.Sprintf("token %s", apiKey))
	req.SetHeader("Content-Type", "application/json")
	req.SetHeader("Accept", "application/vnd.github+json")
	req.SetHeader("User-Agent", "github-mcpx-servlet")

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
