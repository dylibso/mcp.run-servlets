// main.go
package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Called when the tool is invoked
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
	case "google-calendar-login":
		return handleLogin(argsMap)
	case "list_events", "create_event", "update_event", "delete_event", "get_event", "list_calendars":
		// Try to get access token from args
		accessToken, ok := argsMap["access_token"].(string)
		if !ok {
			// If no access token, check if we have device code to try one-time token fetch
			deviceCode, ok := argsMap["device_code"].(string)
			if !ok {
				return CallToolResult{}, errors.New("either access_token or device_code is required")
			}

			// Try to get token once
			auth, err := NewAuthManager()
			if err != nil {
				return CallToolResult{}, fmt.Errorf("failed to create auth manager: %v", err)
			}

			tokenResp, err := auth.PollForToken(deviceCode)
			if err != nil {
				return CallToolResult{}, fmt.Errorf("failed to get token: %v", err)
			}

			if tokenResp.Error != "" || tokenResp.AccessToken == "" {
				return CallToolResult{}, fmt.Errorf("authorization still pending, try again with device_code")
			}

			accessToken = tokenResp.AccessToken
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
		case "delete_event":
			return client.DeleteEvent(argsMap)
		case "get_event":
			return client.GetEvent(argsMap)
		case "list_calendars":
			return client.ListCalendars()
		}
	default:
		return CallToolResult{}, fmt.Errorf("unknown tool: %s", input.Params.Name)
	}

	return CallToolResult{}, nil
}

func handleLogin(args map[string]interface{}) (CallToolResult, error) {
	auth, err := NewAuthManager()
	if err != nil {
		return CallToolResult{}, fmt.Errorf("failed to create auth manager: %v", err)
	}

	// If we have a device code, poll for tokens
	if deviceCode, ok := args["device_code"].(string); ok {
		tokenResp, err := auth.PollForToken(deviceCode)
		if err != nil {
			return CallToolResult{}, fmt.Errorf("failed to poll for token: %v", err)
		}

		respJSON, err := json.Marshal(tokenResp)
		if err != nil {
			return CallToolResult{}, fmt.Errorf("failed to marshal response: %v", err)
		}

		return CallToolResult{
			Content: []Content{{
				Type: ContentTypeText,
				Text: ptr(string(respJSON)),
			}},
		}, nil
	}

	// Otherwise start a new flow
	deviceCode, err := auth.StartDeviceFlow()
	if err != nil {
		return CallToolResult{}, fmt.Errorf("failed to start device flow: %v", err)
	}

	respJSON, err := json.Marshal(deviceCode)
	if err != nil {
		return CallToolResult{}, fmt.Errorf("failed to marshal response: %v", err)
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: ptr(string(respJSON)),
		}},
	}, nil
}

// Describe implements the tool description
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
				Name:        "get_freebusy",
				Description: "Get free/busy information for a set of calendars",
				InputSchema: map[string]interface{}{
					"type": "object",
					"required": []string{
						"time_min",
						"time_max",
					},
					"anyOf": []map[string]interface{}{
						{
							"required": []string{"access_token"},
						},
						{
							"required": []string{"device_code"},
						},
					},
					"properties": map[string]interface{}{
						"access_token": map[string]interface{}{
							"type":        "string",
							"description": "Access token from google-calendar-login",
						},
						"device_code": map[string]interface{}{
							"type":        "string",
							"description": "Device code if access token is not yet available",
						},
						"time_min": map[string]interface{}{
							"type":        "string",
							"description": "Start of the interval to look up availability (RFC3339)",
						},
						"time_max": map[string]interface{}{
							"type":        "string",
							"description": "End of the interval to look up availability (RFC3339)",
						},
						"calendar_ids": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"type": "string",
							},
							"description": "List of calendar IDs to check (defaults to primary calendar)",
							"default":     []string{"primary"},
						},
					},
				},
			},
			{
				Name:        "google-calendar-login",
				Description: "Start the Google Calendar OAuth2 device flow authentication and poll for tokens",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"device_code": map[string]interface{}{
							"type":        "string",
							"description": "Device code from previous response, if continuing auth flow",
						},
						"verification_url": map[string]interface{}{
							"type":        "string",
							"description": "Verification URL from previous response, if continuing auth flow",
						},
						"user_code": map[string]interface{}{
							"type":        "string",
							"description": "User code from previous response, if continuing auth flow",
						},
					},
				},
			},
			{
				Name:        "list_events",
				Description: "List events from a calendar with optional time range and search criteria",
				InputSchema: map[string]interface{}{
					"type": "object",
					"anyOf": []map[string]interface{}{
						{
							"required": []string{"access_token"},
						},
						{
							"required": []string{"device_code"},
						},
					},
					"properties": map[string]interface{}{
						"access_token": map[string]interface{}{
							"type":        "string",
							"description": "Access token from google-calendar-login",
						},
						"device_code": map[string]interface{}{
							"type":        "string",
							"description": "Device code if access token is not yet available",
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
						"show_deleted": map[string]interface{}{
							"type":        "boolean",
							"description": "Whether to include deleted events",
							"default":     false,
						},
					},
				},
			},
			{
				Name:        "get_event",
				Description: "Get a single calendar event by ID",
				InputSchema: map[string]interface{}{
					"type": "object",
					"required": []string{
						"event_id",
					},
					"anyOf": []map[string]interface{}{
						{
							"required": []string{"access_token"},
						},
						{
							"required": []string{"device_code"},
						},
					},
					"properties": map[string]interface{}{
						"access_token": map[string]interface{}{
							"type":        "string",
							"description": "Access token from google-calendar-login",
						},
						"device_code": map[string]interface{}{
							"type":        "string",
							"description": "Device code if access token is not yet available",
						},
						"calendar_id": map[string]interface{}{
							"type":        "string",
							"description": "Calendar ID (use 'primary' for primary calendar)",
							"default":     "primary",
						},
						"event_id": map[string]interface{}{
							"type":        "string",
							"description": "Event ID to retrieve",
						},
					},
				},
			},
			{
				Name:        "list_calendars",
				Description: "List all calendars in the user's calendar list",
				InputSchema: map[string]interface{}{
					"type": "object",
					"anyOf": []map[string]interface{}{
						{
							"required": []string{"access_token"},
						},
						{
							"required": []string{"device_code"},
						},
					},
					"properties": map[string]interface{}{
						"access_token": map[string]interface{}{
							"type":        "string",
							"description": "Access token from google-calendar-login",
						},
						"device_code": map[string]interface{}{
							"type":        "string",
							"description": "Device code if access token is not yet available",
						},
						"max_results": map[string]interface{}{
							"type":        "number",
							"description": "Maximum number of calendars to return",
							"default":     100,
						},
					},
				},
			},
			{
				Name:        "google-calendar-login",
				Description: "Start the Google Calendar OAuth2 device flow authentication and poll for tokens",
				InputSchema: map[string]interface{}{
					"type":     "object",
					"required": []string{"device_code"},
					"properties": map[string]interface{}{
						"device_code": map[string]interface{}{
							"type":        "string",
							"description": "Device code from previous response, if continuing auth flow",
						},
						"verification_url": map[string]interface{}{
							"type":        "string",
							"description": "Verification URL from previous response, if continuing auth flow",
						},
						"user_code": map[string]interface{}{
							"type":        "string",
							"description": "User code from previous response, if continuing auth flow",
						},
					},
				},
			},
			{
				Name:        "list_events",
				Description: "List events from a calendar with optional time range and search criteria",
				InputSchema: map[string]interface{}{
					"type": "object",
					"anyOf": []map[string]interface{}{
						{
							"required": []string{"access_token"},
						},
						{
							"required": []string{"device_code"},
						},
					},
					"properties": map[string]interface{}{
						"access_token": map[string]interface{}{
							"type":        "string",
							"description": "Access token from google-calendar-login",
						},
						"device_code": map[string]interface{}{
							"type":        "string",
							"description": "Device code if access token is not yet available",
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
						"show_deleted": map[string]interface{}{
							"type":        "boolean",
							"description": "Whether to include deleted events",
							"default":     false,
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
				Name:        "delete_event",
				Description: "Delete a calendar event",
				InputSchema: map[string]interface{}{
					"type": "object",
					"required": []string{
						"event_id",
					},
					"properties": map[string]interface{}{
						"calendar_id": map[string]interface{}{
							"type":        "string",
							"description": "Calendar ID (use 'primary' for primary calendar)",
							"default":     "primary",
						},
						"event_id": map[string]interface{}{
							"type":        "string",
							"description": "Event ID to delete",
						},
					},
				},
			},
			{
				Name:        "get_event",
				Description: "Get a single calendar event by ID",
				InputSchema: map[string]interface{}{
					"type": "object",
					"required": []string{
						"event_id",
					},
					"properties": map[string]interface{}{
						"calendar_id": map[string]interface{}{
							"type":        "string",
							"description": "Calendar ID (use 'primary' for primary calendar)",
							"default":     "primary",
						},
						"event_id": map[string]interface{}{
							"type":        "string",
							"description": "Event ID to retrieve",
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
