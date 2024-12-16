package main

import "encoding/json"

type ToolDescription struct {
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
	Name        string      `json:"name"`
}

func parseToolDescription(data []byte) (ToolDescription, error) {
	var res ToolDescription
	err := json.Unmarshal(data, &res)
	return res, err
}

type CallToolRequest struct {
	Method *string `json:"method,omitempty"`
	Params Params  `json:"params"`
}

func (c *CallToolRequest) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

type Params struct {
	Arguments *map[string]interface{} `json:"arguments,omitempty"`
	Name      string                  `json:"name"`
}

func parseCallToolResult(data []byte) (CallToolResult, error) {
	var res CallToolResult
	err := json.Unmarshal(data, &res)
	return res, err
}

type CallToolResult struct {
	Content []Content `json:"content"`
	IsError *bool     `json:"isError,omitempty"`
}

type ContentType string

const (
	ContentTypeText     ContentType = "text"
	ContentTypeImage    ContentType = "image"
	ContentTypeResource ContentType = "resource"
)

type Content struct {
	Data     *string           `json:"data,omitempty"`
	MimeType *string           `json:"mimeType,omitempty"`
	Resource *ResourceContents `json:"resource,omitempty"`
	Text     *string           `json:"text,omitempty"`
	Type     ContentType       `json:"type"`
}

type ResourceContents struct {
	Blob     *string `json:"blob,omitempty"`
	MimeType *string `json:"mimeType,omitempty"`
	Text     *string `json:"text,omitempty"`
	Uri      string  `json:"uri"`
}
