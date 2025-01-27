// main.go
package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/extism/go-pdk"
)

func Call(input CallToolRequest) (CallToolResult, error) {
	args := input.Params.Arguments
	if args == nil {
		return CallToolResult{}, errors.New("Arguments must be provided")
	}

	argsMap, ok := args.(map[string]interface{})
	if !ok {
		return CallToolResult{}, errors.New("Invalid arguments format")
	}

	switch input.Params.Name {
	case "login-initiate":
		auth, err := NewAuthManager()
		if err != nil {
			return CallToolResult{}, fmt.Errorf("failed to create auth manager: %v", err)
		}

		pdk.Log(pdk.LogInfo, "Starting device flow")

		deviceCode, err := auth.StartDeviceFlow()
		if err != nil {
			return CallToolResult{}, fmt.Errorf("failed to start device flow: %v", err)
		}

		respJSON, err := json.Marshal(deviceCode)
		if err != nil {
			return CallToolResult{}, fmt.Errorf("failed to marshal response: %v", err)
		}

		return CallToolResult{
			Content: []Content{{Type: ContentTypeText, Text: ptr(string(respJSON))}},
		}, nil

	case "login-complete":
		deviceCode, ok := argsMap["device_code"].(string)
		if !ok {
			return CallToolResult{}, errors.New("device_code is required")
		}

		auth, err := NewAuthManager()
		if err != nil {
			return CallToolResult{}, fmt.Errorf("failed to create auth manager: %v", err)
		}

		tokenResp, err := auth.PollForToken(deviceCode)
		if err != nil {
			return CallToolResult{}, fmt.Errorf("failed to get token: %v", err)
		}

		if tokenResp.Error == "authorization_pending" {
			return CallToolResult{}, fmt.Errorf("authorization_pending: user has not yet completed authorization")
		}

		respJSON, err := json.Marshal(tokenResp)
		if err != nil {
			return CallToolResult{}, fmt.Errorf("failed to marshal response: %v", err)
		}

		return CallToolResult{
			Content: []Content{{Type: ContentTypeText, Text: ptr(string(respJSON))}},
		}, nil

	default:
		// All other operations require an access token
		accessToken, ok := argsMap["access_token"].(string)
		if !ok {
			return CallToolResult{}, errors.New("access_token is required. Call login-initiate and login-complete first")
		}

		// Create calendar client with access token
		client := NewCalendarClient(accessToken)

		// Call appropriate method
		switch input.Params.Name {
		case "list_events":
			return client.ListEvents(argsMap)
		case "create_event":
			return client.CreateEvent(argsMap)
		case "update_event":
			return client.UpdateEvent(argsMap)
		case "list_calendars":
			return client.ListCalendars()
		default:
			return CallToolResult{}, fmt.Errorf("unknown tool: %s", input.Params.Name)
		}
	}
}
func Describe() (ListToolsResult, error) {
	return ListToolsResult{
		Tools: []ToolDescription{
			{
				Name:        "login-initiate",
				Description: "Start the Google Calendar OAuth2 device flow authentication process",
				InputSchema: map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{},
				},
			},
			{
				Name:        "login-complete",
				Description: "Complete the OAuth2 device flow and get access token",
				InputSchema: map[string]interface{}{
					"type":     "object",
					"required": []string{"device_code"},
					"properties": map[string]interface{}{
						"device_code": map[string]interface{}{
							"type":        "string",
							"description": "Device code from login-initiate response",
						},
					},
				},
			},
			{
				Name:        "list_events",
				Description: "List events from a calendar with optional time range",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"access_token": map[string]interface{}{
							"type":        "string",
							"description": "Access token from login-complete. Either access_token or device_code must be provided",
						},
						"device_code": map[string]interface{}{
							"type":        "string",
							"description": "Device code if access token is not yet available. Either access_token or device_code must be provided",
						},
						"calendar_id": map[string]interface{}{
							"type":        "string",
							"description": "Calendar ID (use 'primary' for primary calendar)",
							"default":     "primary",
						},
						"time_min": map[string]interface{}{
							"type":        "string",
							"description": "Start time in RFC3339 format (default: now)",
						},
						"time_max": map[string]interface{}{
							"type":        "string",
							"description": "End time in RFC3339 format",
						},
						"max_results": map[string]interface{}{
							"type":        "number",
							"description": "Maximum number of events to return",
							"default":     10,
						},
					},
				},
			},
			{
				Name:        "create_event",
				Description: "Create a new calendar event",
				InputSchema: map[string]interface{}{
					"type": "object",
					"required": []string{
						"summary",
						"start_time",
						"end_time",
					},
					"properties": map[string]interface{}{
						"access_token": map[string]interface{}{
							"type":        "string",
							"description": "Access token from login-complete. Either access_token or device_code must be provided",
						},
						"device_code": map[string]interface{}{
							"type":        "string",
							"description": "Device code if access token is not yet available. Either access_token or device_code must be provided",
						},
						"calendar_id": map[string]interface{}{
							"type":        "string",
							"description": "Calendar ID (use 'primary' for primary calendar)",
							"default":     "primary",
						},
						"summary": map[string]interface{}{
							"type":        "string",
							"description": "Event title",
						},
						"description": map[string]interface{}{
							"type":        "string",
							"description": "Event description",
						},
						"location": map[string]interface{}{
							"type":        "string",
							"description": "Event location",
						},
						"start_time": map[string]interface{}{
							"type":        "string",
							"description": "Start time in RFC3339 format",
						},
						"end_time": map[string]interface{}{
							"type":        "string",
							"description": "End time in RFC3339 format",
						},
						"attendees": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"type": "string",
							},
							"description": "List of attendee email addresses",
						},
					},
				},
			},
			{
				Name:        "update_event",
				Description: "Update an existing calendar event",
				InputSchema: map[string]interface{}{
					"type": "object",
					"required": []string{
						"event_id",
					},
					"properties": map[string]interface{}{
						"access_token": map[string]interface{}{
							"type":        "string",
							"description": "Access token from login-complete. Either access_token or device_code must be provided",
						},
						"device_code": map[string]interface{}{
							"type":        "string",
							"description": "Device code if access token is not yet available. Either access_token or device_code must be provided",
						},
						"calendar_id": map[string]interface{}{
							"type":        "string",
							"description": "Calendar ID (use 'primary' for primary calendar)",
							"default":     "primary",
						},
						"event_id": map[string]interface{}{
							"type":        "string",
							"description": "Event ID to update",
						},
						"summary": map[string]interface{}{
							"type":        "string",
							"description": "New event title",
						},
						"description": map[string]interface{}{
							"type":        "string",
							"description": "New event description",
						},
						"location": map[string]interface{}{
							"type":        "string",
							"description": "New event location",
						},
						"start_time": map[string]interface{}{
							"type":        "string",
							"description": "New start time in RFC3339 format",
						},
						"end_time": map[string]interface{}{
							"type":        "string",
							"description": "New end time in RFC3339 format",
						},
						"attendees": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"type": "string",
							},
							"description": "New list of attendee email addresses",
						},
					},
				},
			},
			{
				Name:        "list_calendars",
				Description: "List all calendars in the user's calendar list",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"access_token": map[string]interface{}{
							"type":        "string",
							"description": "Access token from login-complete. Either access_token or device_code must be provided",
						},
						"device_code": map[string]interface{}{
							"type":        "string",
							"description": "Device code if access token is not yet available. Either access_token or device_code must be provided",
						},
						"max_results": map[string]interface{}{
							"type":        "number",
							"description": "Maximum number of calendars to return",
							"default":     100,
						},
					},
				},
			},
		},
	}, nil
}
