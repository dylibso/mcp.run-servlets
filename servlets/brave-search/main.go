// Note: run `go doc -all` in this package to see all of the types and functions available.
// ./pdk.gen.go contains the domain types from the host where your plugin will run.
package main

import (
	"encoding/base64"
	"strconv"
	"strings"

	pdk "github.com/extism/go-pdk"
)

var (
	apiKey string
)

func loadKeys() {
	if apiKey != "" { // already loaded
		return
	}
	k, kerr := pdk.GetConfig("api-key")
	if !kerr {
		panic("missing required configuration")
	}
	apiKey = k
}

// Called when the tool is invoked.
// If you support multiple tools, you must switch on the input.params.name to detect which tool is being called.
// It takes CallToolRequest as input (The incoming tool request from the LLM)
// And returns CallToolResult (The servlet's response to the given tool call)
func Call(input CallToolRequest) (CallToolResult, error) {
	loadKeys()
	args := args{input.Params.Arguments.(map[string]any)}
	switch input.Params.Name {
	case WebSearchTool.Name:
		return callWebSearch(args)
	case ImageSearchTool.Name:
		return callImageSearch(args)
	default:
		return callToolError("unknown tool " + input.Params.Name), nil
	}
}

func Describe() (ListToolsResult, error) {
	return ListToolsResult{
		Tools: []ToolDescription{
			WebSearchTool,
			ImageSearchTool,
		}}, nil

}

func some[T any](t T) *T {
	return &t
}

type SchemaProperty struct {
	Type        string  `json:"type"`
	Description string  `json:"description,omitempty"`
	Items       *schema `json:"items,omitempty"`
}

func prop(tpe, description string) SchemaProperty {
	return SchemaProperty{Type: tpe, Description: description}
}

func arrprop(tpe, description string, items schema) SchemaProperty {
	return SchemaProperty{Type: tpe, Description: description, Items: &items}
}

type schema = map[string]any
type props = map[string]SchemaProperty
type args struct {
	args map[string]any
}

func callToolTextSuccess(bytes []byte) (res CallToolResult) {
	res.Content = []Content{{Type: ContentTypeText, Text: some(string(bytes))}}
	return
}

func callToolImageSuccess(bytes []byte) (res CallToolResult) {
	b64s := base64.StdEncoding.EncodeToString(bytes)
	res.Content = []Content{{Type: ContentTypeImage, Data: some(b64s), MimeType: some("image/png")}}
	return
}

func callToolError(msg string) (res CallToolResult) {
	res.IsError = some(true)
	res.Content = []Content{{Type: ContentTypeText, Text: some(msg)}}
	return
}

func (a args) number(key string) (float64, bool) {
	if n, ok := a.args[key]; ok && n != nil {
		if i, ok := n.(float64); ok {
			return i, true
		}
		if i, ok := n.(int); ok {
			return float64(i), true
		}
		if i, ok := n.(string); ok {
			if f, err := strconv.ParseFloat(i, 64); err == nil {
				return f, true
			}
		}
	}
	return 0, false
}

func (a args) bool(key string) (bool, bool) {
	if n, ok := a.args[key]; ok && n != nil {
		if i, ok := n.(bool); ok {
			return i, true
		}
		if i, ok := n.(string); ok {
			if f, err := strconv.ParseBool(i); err == nil {
				return f, true
			}
		}
	}
	return false, false
}

func (a args) string(key string) (string, bool) {
	if s, ok := a.args[key]; ok && s != nil {
		if ss, ok := s.(string); ok {
			ss = strings.TrimSpace(ss)
			if ss == "" || ss == "null" {
				return "", false
			}
			return ss, ok
		}
	}
	return "", false
}
