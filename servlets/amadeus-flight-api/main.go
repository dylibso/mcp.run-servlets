// Note: run `go doc -all` in this package to see all of the types and functions available.
// ./pdk.gen.go contains the domain types from the host where your plugin will run.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/extism/go-pdk"
)

var (
	config = Config{}
)

// Called when the tool is invoked.
// If you support multiple tools, you must switch on the input.params.name to detect which tool is being called.
// It takes CallToolRequest as input (The incoming tool request from the LLM)
// And returns CallToolResult (The servlet's response to the given tool call)
func Call(input CallToolRequest) (CallToolResult, error) {
	if err := loadConfig(); err != nil {
		return callToolError(err.Error()), nil
	}
	refreshToken()
	args := input.Params.Arguments.(map[string]any)
	switch input.Params.Name {
	case FlightOfferSearchTool.Name:
		return flightOfferSearch(args), nil
	case FlightDatesTool.Name:
		return flightDates(args), nil
	case FlightInspirationTool.Name:
		return flightInspiration(args), nil
	default:
		return CallToolResult{}, fmt.Errorf("unknown tool name %q", input.Params.Name)
	}
}

// Called by mcpx to understand how and why to use this tool.
// Note: Your servlet configs will not be set when this function is called,
// so do not rely on config in this function
// And returns ListToolsResult (The tools' descriptions, supporting multiple tools from a single servlet.)
func Describe() (ListToolsResult, error) {
	return ListToolsResult{
		Tools: FlightsTools,
	}, nil
}

func loadConfig() error {
	if config.apiKey != "" && config.apiSecret != "" {
		return nil
	}

	var okU, okK, okS bool
	config.baseUrl, okU = pdk.GetConfig("base-url")
	config.apiKey, okK = pdk.GetConfig("api-key")
	config.apiSecret, okS = pdk.GetConfig("api-secret")
	if !okU || !okK || !okS {
		return errors.New("missing required configuration")
	}

	return nil
}

func refreshToken() {
	if time.Now().UTC().Unix() > config.expiration {
		req := pdk.NewHTTPRequest(pdk.MethodPost, config.baseUrl+"/v1/security/oauth2/token")
		req.SetHeader("Content-Type", "application/x-www-form-urlencoded")
		form := "grant_type=client_credentials&client_id=" + config.apiKey + "&client_secret=" + config.apiSecret
		pdk.Log(pdk.LogDebug, form)
		req.SetBody([]byte(form))
		resp := req.Send()
		if resp.Status() == 200 {
			var res map[string]any
			json.Unmarshal(resp.Body(), &res)
			config.token = res["access_token"].(string)
			config.expiration = time.Now().UTC().Unix() + int64(res["expires_in"].(float64))
		}
	}
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

func arrprop(tpe, description, itemstpe string) SchemaProperty {
	items := schema{"type": itemstpe}
	return SchemaProperty{Type: tpe, Description: description, Items: &items}
}

type schema = map[string]any
type props = map[string]SchemaProperty

func callToolSuccess(msg string) (res CallToolResult) {
	res.Content = []Content{{Type: ContentTypeText, Text: some(msg)}}
	return
}

func callToolError(msg string) (res CallToolResult) {
	res.IsError = some(true)
	res.Content = []Content{{Type: ContentTypeText, Text: some(msg)}}
	return
}

type Config struct {
	// grant_type=client_credentials&client_id={client_id}&client_secret={client_secret}
	baseUrl    string
	apiKey     string
	apiSecret  string
	token      string
	expiration int64
}
