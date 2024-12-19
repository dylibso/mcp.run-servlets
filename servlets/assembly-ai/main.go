package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

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

	// Get API key
	apiKey, ok := pdk.GetConfig("ASSEMBLYAI_API_KEY")
	if !ok {
		return CallToolResult{}, fmt.Errorf("ASSEMBLYAI_API_KEY config must be set")
	}

	var audioData []byte
	var err error

	// Check if base64 input is provided
	if base64Input, ok := args["audio_base64"].(string); ok {
		// Decode base64 input
		audioData, err = base64.StdEncoding.DecodeString(base64Input)
		if err != nil {
			return CallToolResult{}, fmt.Errorf("failed to decode base64 audio: %w", err)
		}
		pdk.Log(pdk.LogInfo, fmt.Sprintf("Using base64 input data (%d bytes)", len(audioData)))
	} else if fileInput, ok := args["audio_path"].(string); ok {
		// Use file path input
		fullPath := fileInput
		if !filepath.IsAbs(fileInput) {
			fullPath = filepath.Join("/", fileInput)
		}

		pdk.Log(pdk.LogDebug, fmt.Sprintf("Transcribing audio file: %s", truncateString(fullPath, 100)))

		// Read the audio file
		audioData, err = os.ReadFile(fullPath)
		if err != nil {
			return CallToolResult{}, fmt.Errorf("failed to read audio file: %w. Trying to read from %s", err, fullPath)
		}
		pdk.Log(pdk.LogInfo, fmt.Sprintf("Read %d bytes from input file", len(audioData)))
	} else {
		return CallToolResult{}, fmt.Errorf("either audio_path or audio_base64 must be provided")
	}

	// Create client and transcribe
	transcript, err := transcribeAudio(apiKey, audioData)
	if err != nil {
		return CallToolResult{}, fmt.Errorf("transcription failed: %w", err)
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: &transcript.Text,
		}},
	}, nil
}

// Helper function to truncate string with ellipsis
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// Called by mcpx to understand how and why to use this tool.
func Describe() (ListToolsResult, error) {
	return ListToolsResult{
		Tools: []ToolDescription{
			{
				Name:        "transcribe",
				Description: "Transcribes an audio file using AssemblyAI and returns the transcription as text. Supports common audio formats (mp3, wav, FLAC, AAC, M4A, etc). Audio can be provided either as a file path or base64-encoded string. Either audio_path or audio_base64 must be provided.",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"audio_path": map[string]interface{}{
							"type":        "string",
							"description": "Path to the input audio file. The file paths must either be absolute or relative to the directory this servlet has access to. This servlet understands the followin root directories: /, /home/, and /tmp",
						},
						"audio_base64": map[string]interface{}{
							"type":        "string",
							"description": "Base64-encoded audio data",
						},
					},
				},
			},
		},
	}, nil
}
