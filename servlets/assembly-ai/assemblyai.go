package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/extism/go-pdk"
)

// Client manages communication with AssemblyAI API
type Client struct {
	apiKey  string
	baseURL string
}

// NewClient creates a new AssemblyAI client
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: "https://api.assemblyai.com/v2",
	}
}

// TranscriptStatus represents the status of a transcript
type TranscriptStatus string

const (
	TranscriptStatusQueued     TranscriptStatus = "queued"
	TranscriptStatusProcessing TranscriptStatus = "processing"
	TranscriptStatusCompleted  TranscriptStatus = "completed"
	TranscriptStatusError      TranscriptStatus = "error"
)

// Transcript represents an AssemblyAI transcript
type Transcript struct {
	ID       *string          `json:"id"`
	Status   TranscriptStatus `json:"status"`
	Text     string           `json:"text"`
	AudioURL string           `json:"audio_url"`
	Error    string           `json:"error,omitempty"`
}

// Upload uploads an audio file and returns the URL
func (c *Client) Upload(data []byte) (string, error) {
	req := pdk.NewHTTPRequest(pdk.MethodPost, c.baseURL+"/upload")
	req.SetHeader("Authorization", c.apiKey)
	req.SetHeader("Content-Type", "application/octet-stream")
	req.SetBody(data)

	resp := req.Send()
	if resp.Status() != 200 {
		return "", fmt.Errorf("upload failed with status %d", resp.Status())
	}

	var result struct {
		UploadURL string `json:"upload_url"`
	}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return "", fmt.Errorf("failed to parse upload response: %w", err)
	}

	return result.UploadURL, nil
}

// SubmitTranscript submits a URL for transcription
func (c *Client) SubmitTranscript(audioURL string) (*Transcript, error) {
	req := pdk.NewHTTPRequest(pdk.MethodPost, c.baseURL+"/transcript")
	req.SetHeader("Authorization", c.apiKey)
	req.SetHeader("Content-Type", "application/json")

	params := struct {
		AudioURL string `json:"audio_url"`
	}{
		AudioURL: audioURL,
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req.SetBody(body)
	resp := req.Send()

	if resp.Status() != 200 {
		return nil, fmt.Errorf("transcript submission failed with status %d", resp.Status())
	}

	var transcript Transcript
	if err := json.Unmarshal(resp.Body(), &transcript); err != nil {
		return nil, fmt.Errorf("failed to parse transcript response: %w", err)
	}

	return &transcript, nil
}

// GetTranscript gets the status and result of a transcript
func (c *Client) GetTranscript(transcriptID string) (*Transcript, error) {
	req := pdk.NewHTTPRequest(pdk.MethodGet, fmt.Sprintf("%s/transcript/%s", c.baseURL, transcriptID))
	req.SetHeader("Authorization", c.apiKey)

	resp := req.Send()
	if resp.Status() != 200 {
		return nil, fmt.Errorf("get transcript failed with status %d", resp.Status())
	}

	var transcript Transcript
	if err := json.Unmarshal(resp.Body(), &transcript); err != nil {
		return nil, fmt.Errorf("failed to parse transcript response: %w", err)
	}

	return &transcript, nil
}

// Example usage function to showcase how to use the client
func transcribeAudio(apiKey string, audioData []byte) (*Transcript, error) {
	client := NewClient(apiKey)

	// Upload the audio file
	uploadURL, err := client.Upload(audioData)
	if err != nil {
		return nil, fmt.Errorf("upload failed: %w", err)
	}

	//uploadURL := "https://cdn.assemblyai.com/upload/1e8a7b90-7b3b-4b7b-8b3b-7b3b4b7b8b3b"

	// Submit for transcription
	transcript, err := client.SubmitTranscript(uploadURL)
	if err != nil {
		return nil, fmt.Errorf("submission failed: %w", err)
	}

	// Poll until completion
	for transcript.Status != TranscriptStatusCompleted && transcript.Status != TranscriptStatusError {
		pdk.Log(pdk.LogInfo, "Checking status")
		transcript, err = client.GetTranscript(*transcript.ID)
		if err != nil {
			return nil, fmt.Errorf("status check failed: %w", err)
		}

		if transcript.Status == TranscriptStatusError {
			return nil, fmt.Errorf("transcription failed: %s", transcript.Error)
		}

		// Sleep for a few seconds before next poll
		// Note: In a real implementation, you'd want to use backoff
		pdk.Log(pdk.LogInfo, fmt.Sprintf("Status: %s", transcript.Status))

		time.Sleep(100)
	}

	pdk.Log(pdk.LogInfo, "Transcription completed")

	return transcript, nil
}
