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
	case ListIssuesTool.Name:
		owner, _ := args["owner"].(string)
		repo, _ := args["repo"].(string)
		return issueList(apiKey, owner, repo, args)
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

	case GetFileContentsTool.Name:
		owner, _ := args["owner"].(string)
		repo, _ := args["repo"].(string)
		path, _ := args["path"].(string)
		branch, _ := args["branch"].(string)
		res := filesGetContents(apiKey, owner, repo, path, &branch)
		return res, nil
	case CreateOrUpdateFileTool.Name:
		owner, _ := args["owner"].(string)
		repo, _ := args["repo"].(string)
		path, _ := args["path"].(string)
		file := fileCreateFromArgs(args)
		return filesCreateOrUpdate(apiKey, owner, repo, path, file)

	case CreateBranchTool.Name:
		owner, _ := args["owner"].(string)
		repo, _ := args["repo"].(string)
		from, _ := args["branch"].(string)
		var maybeBranch *string
		if branch, ok := args["from_branch"].(string); ok {
			maybeBranch = &branch
		}
		return branchCreate(apiKey, owner, repo, from, maybeBranch), nil

	case ListCommitsTool.Name:
		owner, _ := args["owner"].(string)
		repo, _ := args["repo"].(string)
		return commitList(apiKey, owner, repo, args), nil
	case GetCommitTool.Name:
		owner, _ := args["owner"].(string)
		repo, _ := args["repo"].(string)
		ref, _ := args["ref"].(string)
		return commitGet(apiKey, owner, repo, ref, args), nil

	case ListPullRequestsTool.Name:
		owner, _ := args["owner"].(string)
		repo, _ := args["repo"].(string)
		return pullRequestList(apiKey, owner, repo, args)

	case CreatePullRequestTool.Name:
		owner, _ := args["owner"].(string)
		repo, _ := args["repo"].(string)
		pr := branchPullRequestSchemaFromArgs(args)
		return branchCreatePullRequest(apiKey, owner, repo, pr), nil

	case PushFilesTool.Name:
		owner, _ := args["owner"].(string)
		repo, _ := args["repo"].(string)
		branch, _ := args["branch"].(string)
		message, _ := args["message"].(string)
		files := filePushFromArgs(args)
		return filesPush(apiKey, owner, repo, branch, message, files), nil

	case ListReposTool.Name:
		owner, _ := args["owner"].(string)
		return reposList(apiKey, owner, args)

	case GetRepositoryCollaboratorsTool.Name:
		owner, _ := args["owner"].(string)
		repo, _ := args["repo"].(string)
		return reposGetCollaborators(apiKey, owner, repo, args)

	case GetRepositoryContributorsTool.Name:
		owner, _ := args["owner"].(string)
		repo, _ := args["repo"].(string)
		return reposGetContributors(apiKey, owner, repo, args)

	case GetRepositoryDetailsTool.Name:
		owner, _ := args["owner"].(string)
		repo, _ := args["repo"].(string)
		return reposGetDetails(apiKey, owner, repo)

	case CreateGistTool.Name:
		description, _ := args["description"].(string)
		files, _ := args["files"].(map[string]any)
		return gistCreate(apiKey, description, files), nil

	case GetGistTool.Name:
		gistId, _ := args["gist_id"].(string)
		return gistGet(apiKey, gistId), nil

	case UpdateGistTool.Name:
		gistId, _ := args["gist_id"].(string)
		description, _ := args["description"].(string)
		files, _ := args["files"].(map[string]any)
		return gistUpdate(apiKey, gistId, description, files), nil

	case DeleteGistTool.Name:
		gistId, _ := args["gist_id"].(string)
		return gistDelete(apiKey, gistId), nil

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

func Describe() (ListToolsResult, error) {
	toolsets := [][]ToolDescription{
		IssueTools,
		FileTools,
		BranchTools,
		CommitTools,
		RepoTools,
		GistTools,
	}

	tools := []ToolDescription{}

	for _, toolset := range toolsets {
		tools = append(tools, toolset...)
	}

	return ListToolsResult{
		Tools: tools,
	}, nil
}

func some[T any](t T) *T {
	return &t
}

type SchemaProperty struct {
	Type                 string  `json:"type"`
	Description          string  `json:"description,omitempty"`
	AdditionalProperties *schema `json:"additionalProperties,omitempty"`
	Items                *schema `json:"items,omitempty"`
}

func prop(tpe, description string) SchemaProperty {
	return SchemaProperty{Type: tpe, Description: description}
}

func arrprop(tpe, description, itemstpe string) SchemaProperty {
	items := schema{"type": itemstpe}
	return SchemaProperty{Type: tpe, Description: description, Items: &items}
}

type schema = map[string]interface{}
type props = map[string]SchemaProperty
