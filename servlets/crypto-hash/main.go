// Note: run `go doc -all` in this package to see all of the types and functions available.
// ./pdk.gen.go contains the domain types from the host where your plugin will run.
package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Called when the tool is invoked.
// If you support multiple tools, you must switch on the input.params.name to detect which tool is being called.
// The name will match one of the tool names returned from "describe".
// It takes CallToolRequest as input (The incoming tool request from the LLM)
// And returns CallToolResult (The servlet's response to the given tool call)
func Call(input CallToolRequest) (CallToolResult, error) {
	args := input.Params.Arguments
	if args == nil {
		return CallToolResult{}, errors.New("Arguments must be provided")
	}

	argsMap := args.(map[string]interface{})

	switch input.Params.Name {
	case "md5":
		return runMD5(argsMap)

	case "sha1":
		return runSHA1(argsMap)

	case "sha2":
		return runSHA256(argsMap)

	case "sha256":
		return runSHA256(argsMap)

	case "bcrypt":
		return runBcrypt(argsMap)

	default:
		return CallToolResult{}, errors.New("Unknown tool")
	}
}

func runSHA1(argsMap map[string]interface{}) (CallToolResult, error) {
	text, ok := argsMap["text"].(string)
	if !ok {
		return CallToolResult{}, errors.New("text must be provided")
	}

	h1 := sha1.New()
	h1.Write([]byte(text))
	hash := hex.EncodeToString(h1.Sum(nil))

	return CallToolResult{
		Content: []Content{
			{
				Type: ContentTypeText,
				Text: &hash,
			},
		},
	}, nil
}

func runSHA256(argsMap map[string]interface{}) (CallToolResult, error) {
	text, ok := argsMap["text"].(string)
	if !ok {
		return CallToolResult{}, errors.New("text must be provided")
	}

	h256 := sha256.New()
	h256.Write([]byte(text))
	hash := hex.EncodeToString(h256.Sum(nil))

	return CallToolResult{
		Content: []Content{
			{
				Type: ContentTypeText,
				Text: &hash,
			},
		},
	}, nil
}

func runMD5(args map[string]interface{}) (CallToolResult, error) {
	text, ok := args["text"].(string)
	if !ok {
		return CallToolResult{}, errors.New("text must be provided")
	}

	hash := fmt.Sprintf("%x", md5.Sum([]byte(text)))
	return CallToolResult{
		Content: []Content{
			{
				Type: ContentTypeText,
				Text: &hash,
			},
		},
	}, nil
}

func runBcrypt(args map[string]interface{}) (CallToolResult, error) {
	pwd, ok := args["password"].(string)
	if !ok {
		return CallToolResult{}, errors.New("password must be provided")
	}

	// Get cost parameter with default value of 10 if not provided
	cost := 10
	if costVal, exists := args["cost"].(float64); exists {
		cost = int(costVal)
		if cost < 4 || cost > 31 {
			return CallToolResult{}, errors.New("cost must be between 4 and 31")
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), cost)
	if err != nil {
		return CallToolResult{}, fmt.Errorf("bcrypt error: %v", err)
	}

	hashStr := string(hash)
	return CallToolResult{
		Content: []Content{
			{
				Type: ContentTypeText,
				Text: &hashStr,
			},
		},
	}, nil
}

// Called by mcpx to understand how and why to use this tool.
// Note: Your servlet configs will not be set when this function is called,
// so do not rely on config in this function
// And returns ListToolsResult (The tools' descriptions, supporting multiple tools from a single servlet.)
func Describe() (ListToolsResult, error) {
	return ListToolsResult{
		Tools: []ToolDescription{
			{
				Name:        "md5",
				Description: "Hash a string using MD5 (Note: MD5 is not cryptographically secure, use for checksums only)",
				InputSchema: map[string]interface{}{
					"type":     "object",
					"required": []string{"text"},
					"properties": map[string]interface{}{
						"text": map[string]interface{}{
							"type":        "string",
							"description": "the text to hash",
						},
					},
				},
			},
			{
				Name:        "sha1",
				Description: "Hash a string using SHA-1 (Note: SHA-1 is not recommended for new applications)",
				InputSchema: map[string]interface{}{
					"type":     "object",
					"required": []string{"text"},
					"properties": map[string]interface{}{
						"text": map[string]interface{}{
							"type":        "string",
							"description": "the text to hash",
						},
					},
				},
			},
			{
				Name:        "sha256",
				Description: "Hash a string using SHA-256 (also called SHA2, cryptographically secure)",
				InputSchema: map[string]interface{}{
					"type":     "object",
					"required": []string{"text"},
					"properties": map[string]interface{}{
						"text": map[string]interface{}{
							"type":        "string",
							"description": "the text to hash",
						},
					},
				},
			},
			{
				Name:        "bcrypt",
				Description: "Hash a string using bcrypt (recommended for passwords)",
				InputSchema: map[string]interface{}{
					"type":     "object",
					"required": []string{"password"},
					"properties": map[string]interface{}{
						"password": map[string]interface{}{
							"type":        "string",
							"description": "the password to hash",
						},
						"cost": map[string]interface{}{
							"type":        "integer",
							"description": "cost factor for bcrypt (4-31, default: 10)",
							"minimum":     4,
							"maximum":     31,
							"default":     10,
						},
					},
				},
			},
		},
	}, nil
}
