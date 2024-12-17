// Note: run `go doc -all` in this package to see all of the types and functions available.
// ./pdk.gen.go contains the domain types from the host where your plugin will run.
package main

import (
	"fmt"

	"github.com/extism/go-pdk"
)

// Called when the tool is invoked.
// If you support multiple tools, you must switch on the input.params.name to detect which tool is being called.
// The name will match one of the tool names returned from "describe".
// It takes CallToolRequest as input (The incoming tool request from the LLM)
// And returns CallToolResult (The servlet's response to the given tool call)
func Call(input CallToolRequest) (CallToolResult, error) {
	apiKey, ok := pdk.GetConfig("api-key")
	if !ok {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("No api-key configured"),
			}},
		}, nil
	}
	args := input.Params.Arguments.(map[string]interface{})
	pdk.Log(pdk.LogDebug, fmt.Sprint("Args: ", args))
	switch input.Params.Name {
	case GetIssueTool.Name:
		owner, _ := args["owner"].(string)
		repo, _ := args["repo"].(string)
		issue, _ := args["issue"].(float64)
		return issueGet(apiKey, owner, repo, int(issue))
	case AddIssueCommentTool.Name:
		owner, _ := args["owner"].(string)
		repo, _ := args["repo"].(string)
		issue, _ := args["issue"].(float64)
		body, _ := args["body"].(string)
		return issueAddComment(apiKey, owner, repo, int(issue), body)
	case CreateIssueTool.Name:
		owner, _ := args["owner"].(string)
		repo, _ := args["repo"].(string)
		data := issueFromArgs(args)
		return issueCreate(apiKey, owner, repo, data)
	case UpdateIssueTool.Name:
		owner, _ := args["owner"].(string)
		repo, _ := args["repo"].(string)
		issue, _ := args["issue"].(float64)
		data := issueFromArgs(args)
		return issueUpdate(apiKey, owner, repo, int(issue), data)
	default:
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Unknown tool " + input.Params.Name),
			}},
		}, nil
	}

}

// Called by mcpx to understand how and why to use this tool.
// Note: Your servlet configs will not be set when this function is called,
// so do not rely on config in this function
// And returns ListToolsResult (The tools' descriptions, supporting multiple tools from a single servlet.)
func Describe() (ListToolsResult, error) {
	return ListToolsResult{
		Tools: IssueTools,
	}, nil
}

func some[T any](t T) *T {
	return &t
}
