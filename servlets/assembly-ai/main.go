package main

import (
	"encoding/base64"
	"fmt"

	"github.com/extism/go-pdk"
)

// Called when the tool is invoked.
func Call(input CallToolRequest) (CallToolResult, error) {
	if input.Params.Arguments == nil {
		return CallToolResult{}, fmt.Errorf("arguments must be provided")
	}

	switch input.Params.Name {
	case "transcribe":
		return handleTranscribe(input)
	default:
		return CallToolResult{}, fmt.Errorf("unknown tool: %s", input.Params.Name)
	}
}

func handleTranscribe(input CallToolRequest) (CallToolResult, error) {
	args := input.Params.Arguments.(map[string]interface{})

	// Get base64 encoded audio data
	audioBase64, ok := args["audio"].(string)
	if !ok {
		return CallToolResult{}, fmt.Errorf("audio parameter must be provided as base64 string")
	}

	audioData, err := base64.StdEncoding.DecodeString(audioBase64)
	if err != nil {
		return CallToolResult{}, fmt.Errorf("failed to decode audio data: %w", err)
	}

	apiKey, ok := pdk.GetConfig("ASSEMBLYAI_API_KEY")
	if !ok {
		return CallToolResult{}, fmt.Errorf("ASSEMBLYAI_API_KEY config must be set!")
	}

	// Create client and transcribe
	transcript, err := transcribeAudio(apiKey, audioData)
	if err != nil {
		return CallToolResult{}, fmt.Errorf("transcription failed: %w", err)
	}

	pdk.Log(pdk.LogInfo, fmt.Sprintf("Transcript: %s", transcript.Text))

	// Return the transcribed text
	text := transcript.Text
	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: &text,
		}},
	}, nil
}

// Called by mcpx to understand how and why to use this tool.
func Describe() (ListToolsResult, error) {
	return ListToolsResult{
		Tools: []ToolDescription{
			{
				Name:        "transcribe",
				Description: "Transcribes an audio file using AssemblyAI. Supports mp3, wav, and other common audio formats.",
				InputSchema: map[string]interface{}{
					"type":     "object",
					"required": []string{"audio"},
					"properties": map[string]interface{}{
						"audio": map[string]interface{}{
							"type":        "string",
							"description": "Base64 encoded audio file content",
						},
					},
				},
			},
		},
	}, nil
}
