// Note: run `go doc -all` in this package to see all of the types and functions available.
// ./pdk.gen.go contains the domain types from the host where your plugin will run.
package main

import (
	"encoding/json"
)

var (
	BASE_URL       string
	HANDLE         string
	PASSWORD       string
	currentSession Session
)

// Called when the tool is invoked.
// If you support multiple tools, you must switch on the input.params.name to detect which tool is being called.
// It takes CallToolRequest as input (The incoming tool request from the LLM)
// And returns CallToolResult (The servlet's response to the given tool call)
func Call(input CallToolRequest) (CallToolResult, error) {
	switch input.Params.Name {
	case "post":
		return post(input.Params.Arguments.(map[string]any))
	case "reply":
		return reply(input.Params.Arguments.(map[string]any))
	case "latest_mentions":
		return latestMentions(input.Params.Arguments.(map[string]any))
	case "search":
		return search(input.Params.Arguments.(map[string]any))
	case "get_thread":
		return getThread(input.Params.Arguments.(map[string]any))
	default:
		return CallToolResult{IsError: some(true), Content: []Content{{
			Type: ContentTypeText,
			Text: some("Unknown tool " + input.Params.Name),
		}}}, nil
	}
}

// Called by mcpx to understand how and why to use this tool.
// Note: Your servlet configs will not be set when this function is called,
// so do not rely on config in this function
// And returns ListToolsResult (The tools' descriptions, supporting multiple tools from a single servlet.)
func Describe() (res ListToolsResult, err error) {
	err = json.Unmarshal([]byte(`
			{
				"tools":[
					{
						"name": "post",
						"description": "Post a message to your feed",
						"inputSchema": {
							"type": "object",
							"properties": {
								"text": {
									"type": "string",
									"description": "The text of the post"
								}
							},
							"required": ["text"]
						}
					},
					{
						"name": "reply",
						"description": "Reply to a message for a given at:// URI",
						"inputSchema": {
							"type": "object",
							"properties": {
								"text": {
									"type": "string",
									"description": "The text of the post"
								},
								"reply_to": {
									"type": "string",
									"description": "The at:// URI of the post we are replying to. Additionally, if the URI is a web URI (e.g. https://bsky.app/profile/<DID>/post/<RKEY>), it will be converted to an AT URI (e.g. at://<DID>/app.bsky.feed.post/<RKEY>)."
								}
							},
							"required": ["text"]
						}
					},
					{
						"name": "get_thread",
						"description": "Get a thread from a given AT-URI. Use this to find replies to an existing post. For instance, this is useful to make sure a bot is not replying to the same message twice.",
						"inputSchema": {
							"type": "object",
							"required": ["uri"],
							"properties": {
								"uri": {
									"type": "string",
									"description": "AT-URI to post record (e.g. as returned by the search)"
								},
								"depth": {
									"type": "integer",
									"description": "How many levels of reply depth should be included in response. (default: 6, limit <=1000)"
								},
								"parentHeight": {
									"type": "integer",
									"description": "How many levels of parent (and grandparent, etc) post to include. (default: 80, limit <=1000)"
								}
							}
						}
					},
					{
						"name": "latest_mentions",
						"description": "Search for your latest mentions on Bluesky.",
						"inputSchema": {
							"type": "object",
							"properties": {
								"within": {
									"type": "string",
									"description": "Find mentions within now and the given duration. A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as \"300ms\", \"-1.5h\" or \"2h45m\". Valid time units are \"ns\", \"us\" (or \"µs\"), \"ms\", \"s\", \"m\", \"h\". default: 5m"
								},
								"limit": {
									"type": "integer",
									"description": "Maximum number of posts to return (>=1, <=100) default: 25"
								}
							}
						}
					}
				]
			}
		`), &res)

	return
}

func callToolError(msg string) CallToolResult {
	return CallToolResult{
		IsError: some(true),
		Content: []Content{{
			Text: some(msg),
			Type: ContentTypeText,
		}},
	}
}
