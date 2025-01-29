// calendar.go
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	pdk "github.com/extism/go-pdk"
)

type CalendarClient struct {
	baseURL    string
	authHeader string
}

// EventData represents a calendar event
type EventData struct {
	Summary     string          `json:"summary"`
	Description string          `json:"description,omitempty"`
	Location    string          `json:"location,omitempty"`
	Start       *EventDateTime  `json:"start"`
	End         *EventDateTime  `json:"end"`
	Attendees   []EventAttendee `json:"attendees,omitempty"`
}

type EventDateTime struct {
	DateTime string `json:"dateTime,omitempty"`
	Date     string `json:"date,omitempty"`
	TimeZone string `json:"timeZone,omitempty"`
}

type EventAttendee struct {
	Email string `json:"email"`
}

func NewCalendarClient(accessToken string) *CalendarClient {
	return &CalendarClient{
		baseURL:    "https://www.googleapis.com/calendar/v3",
		authHeader: fmt.Sprintf("Bearer %s", accessToken),
	}
}

func (c *CalendarClient) makeRequest(method pdk.HTTPMethod, endpoint string, body []byte) ([]byte, error) {
	url := c.baseURL + endpoint
	pdk.Log(pdk.LogInfo, fmt.Sprintf("Making %s request to %s", method, url))

	req := pdk.NewHTTPRequest(method, url)
	req.SetHeader("Authorization", c.authHeader)
	req.SetHeader("Content-Type", "application/json")

	if len(body) > 0 {
		req.SetBody(body)
	}

	res := req.Send()

	if res.Status() < 200 || res.Status() >= 300 {
		return nil, fmt.Errorf("Calendar API error %d: %s", res.Status(), string(res.Body()))
	}

	return res.Body(), nil
}

func (c *CalendarClient) ListEvents(args map[string]interface{}) (CallToolResult, error) {
	calendarID := getStringArg(args, "calendar_id", "primary")
	maxResults := getIntArg(args, "max_results", 10)
	timeMin := time.Now().Format(time.RFC3339)
	if tm, ok := args["time_min"].(string); ok {
		timeMin = tm
	}

	queryParams := []string{
		fmt.Sprintf("maxResults=%d", maxResults),
		fmt.Sprintf("timeMin=%s", timeMin),
		"singleEvents=true",
		"orderBy=startTime",
	}

	if timeMax, ok := args["time_max"].(string); ok {
		queryParams = append(queryParams, fmt.Sprintf("timeMax=%s", timeMax))
	}

	endpoint := fmt.Sprintf("/calendars/%s/events?%s", calendarID, strings.Join(queryParams, "&"))
	resp, err := c.makeRequest(pdk.MethodGet, endpoint, nil)
	if err != nil {
		return CallToolResult{}, err
	}

	return CallToolResult{
		Content: []Content{{Type: ContentTypeText, Text: ptr(string(resp))}},
	}, nil
}

func (c *CalendarClient) CreateEvent(args map[string]interface{}) (CallToolResult, error) {
	calendarID := getStringArg(args, "calendar_id", "primary")

	event := EventData{
		Summary:     getStringArg(args, "summary", ""),
		Description: getStringArg(args, "description", ""),
		Location:    getStringArg(args, "location", ""),
		Start:       &EventDateTime{DateTime: getStringArg(args, "start_time", "")},
		End:         &EventDateTime{DateTime: getStringArg(args, "end_time", "")},
	}

	if attendees, ok := args["attendees"].([]interface{}); ok {
		event.Attendees = make([]EventAttendee, len(attendees))
		for i, a := range attendees {
			if email, ok := a.(string); ok {
				event.Attendees[i] = EventAttendee{Email: email}
			}
		}
	}

	body, err := json.Marshal(event)
	if err != nil {
		return CallToolResult{}, err
	}

	resp, err := c.makeRequest(pdk.MethodPost, fmt.Sprintf("/calendars/%s/events", calendarID), body)
	if err != nil {
		return CallToolResult{}, err
	}

	return CallToolResult{
		Content: []Content{{Type: ContentTypeText, Text: ptr(string(resp))}},
	}, nil
}

func (c *CalendarClient) UpdateEvent(args map[string]interface{}) (CallToolResult, error) {
	calendarID := getStringArg(args, "calendar_id", "primary")
	eventID := getStringArg(args, "event_id", "")
	if eventID == "" {
		return CallToolResult{}, errors.New("event_id is required")
	}

	event := EventData{
		Summary:     getStringArg(args, "summary", ""),
		Description: getStringArg(args, "description", ""),
		Location:    getStringArg(args, "location", ""),
	}

	if startTime := getStringArg(args, "start_time", ""); startTime != "" {
		event.Start = &EventDateTime{DateTime: startTime}
	}
	if endTime := getStringArg(args, "end_time", ""); endTime != "" {
		event.End = &EventDateTime{DateTime: endTime}
	}

	if attendees, ok := args["attendees"].([]interface{}); ok {
		event.Attendees = make([]EventAttendee, len(attendees))
		for i, a := range attendees {
			if email, ok := a.(string); ok {
				event.Attendees[i] = EventAttendee{Email: email}
			}
		}
	}

	body, err := json.Marshal(event)
	if err != nil {
		return CallToolResult{}, err
	}

	resp, err := c.makeRequest(pdk.MethodPut, fmt.Sprintf("/calendars/%s/events/%s", calendarID, eventID), body)
	if err != nil {
		return CallToolResult{}, err
	}

	return CallToolResult{
		Content: []Content{{Type: ContentTypeText, Text: ptr(string(resp))}},
	}, nil
}

func (c *CalendarClient) ListCalendars() (CallToolResult, error) {
	resp, err := c.makeRequest(pdk.MethodGet, "/users/me/calendarList", nil)
	if err != nil {
		return CallToolResult{}, err
	}

	return CallToolResult{
		Content: []Content{{Type: ContentTypeText, Text: ptr(string(resp))}},
	}, nil
}

// Helper functions
func getStringArg(args map[string]interface{}, key, defaultValue string) string {
	if val, ok := args[key].(string); ok {
		return val
	}
	return defaultValue
}

func getIntArg(args map[string]interface{}, key string, defaultValue int) int {
	if val, ok := args[key].(float64); ok {
		return int(val)
	}
	return defaultValue
}

func ptr(s string) *string {
	return &s
}
