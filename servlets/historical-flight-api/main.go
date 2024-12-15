// Note: run `go doc -all` in this package to see all of the types and functions available.
// ./pdk.gen.go contains the domain types from the host where your plugin will run.
package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/extism/go-pdk"
)

const (
	apiOpenSkyBaseUrl   = "https://opensky-network.org/api"
	apiOpenSkyArrival   = "/flights/arrival"
	apiOpenSkyDeparture = "/flights/departure"

	apiAdsbdbBaseUrl  = "https://api.adsbdb.com/v0"
	apiAdsbdbCallsign = "/callsign"
	apiAdsbdbAircraft = "/aircraft"
)

var (
	basicAuthToken string
)

// Called when the tool is invoked.
// It takes CallToolRequest as input (The incoming tool request from the LLM)
// And returns CallToolResult (The servlet's response to the given tool call)
func Call(input CallToolRequest) (res CallToolResult, err error) {
	if err = loadBasicAuthToken(); err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(err.Error()),
			}},
		}, nil
	}

	args := input.Params.Arguments.(map[string]interface{})

	requestType, _ := args["requestType"].(string)

	var result string
	switch requestType {
	case "arrival", "departure":
		var airport = args["airport"].(string)
		var begin = args["begin"].(string)
		var end = args["end"].(string)

		// send the request, get response back (can check status on response via res.Status())
		result, err = flightInfo(requestType, airport, begin, end)
		if err != nil {
			return
		}
		return CallToolResult{
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(result),
			}},
		}, nil
	default:
		var icao24 = args["icao24"].(string)
		var callsign = args["callsign"].(string)
		bytes := aircraft(icao24, callsign)
		imgData := fetchImage(bytes)

		return CallToolResult{
			Content: []Content{{
				Type:     ContentTypeImage,
				MimeType: some("image/jpeg"),
				Data:     some(base64.StdEncoding.EncodeToString(imgData)),
			}, {
				Type: ContentTypeText,
				Text: some(string(bytes)),
			}}}, nil
	}
}

func loadBasicAuthToken() error {
	if basicAuthToken != "" {
		return nil
	}

	user, uok := pdk.GetConfig("username")
	pass, pok := pdk.GetConfig("password")
	if !uok || !pok {
		return errors.New("username or password not set")
	}
	basicAuthToken = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, pass)))

	auth := user + ":" + pass
	basicAuthToken = "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	return nil
}

func flightInfo(requestType string, airport string, begin string, end string) (string, error) {
	query := &url.Values{}
	query.Add("airport", airport)
	query.Add("begin", begin)
	query.Add("end", end)

	var path string
	switch requestType {
	case "arrival":
		path = apiOpenSkyArrival
	case "departure":
		path = apiOpenSkyDeparture
	default:
		return "", errors.New("Invalid type")
	}

	req := pdk.NewHTTPRequest(pdk.MethodGet, apiOpenSkyBaseUrl+path+"?"+query.Encode())
	req.SetHeader("Authorization", basicAuthToken)

	res := req.Send()
	return string(res.Body()), nil
}

func aircraft(modeS, callsign string) []byte {
	req := pdk.NewHTTPRequest(pdk.MethodGet, apiAdsbdbBaseUrl+apiAdsbdbAircraft+"/"+modeS+"?callsign="+callsign)
	res := req.Send()
	return res.Body()
}

func fetchImage(aircraftResponse []byte) []byte {
	jsonData := map[string]interface{}{}
	json.Unmarshal(aircraftResponse, &jsonData)
	if response, ok := jsonData["response"].(map[string]interface{}); ok {
		if aircraft, ok := response["aircraft"].(map[string]interface{}); ok {
			if urlPhoto, ok := aircraft["url_photo"].(string); ok {
				data := pdk.NewHTTPRequest(pdk.MethodGet, urlPhoto).Send()
				return data.Body()
			}
		}
	}
	return nil
}

func Describe() (ListToolsResult, error) {
	return ListToolsResult{Tools: []ToolDescription{{
		Name:        "historical-flight-api",
		Description: "Get the flight arrivals and departures for a given airport by ICAO identifier within a given time range; or get the details and picture of a flight by callsign and ICAO24 hex code.",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"requestType": map[string]interface{}{
					"type":        "string",
					"description": "The type of the request, 'departure', 'arrival' or 'aircraft'",
				},
			},
			"required": []string{"requestType"},
			"oneOf": []interface{}{
				map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"airport": map[string]interface{}{
							"type":        "string",
							"description": "The ICAO identifier of the airport",
						},
						"begin": map[string]interface{}{
							"type":        "string",
							"description": "The start of the time range as a UNIX timestamp in UTC",
						},
						"end": map[string]interface{}{
							"type":        "string",
							"description": "The end of the time range as a UNIX timestamp in UTC",
						},
					},
					"required": []string{"airport", "begin", "end"},
				},
				map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"callsign": map[string]interface{}{
							"type":        "string",
							"description": "The callsign of the flight",
						},
						"icao24": map[string]interface{}{
							"type":        "string",
							"description": "The aircraft as ICAO24 hex code",
						},
					},
					"required": []string{"callsign", "icao24"},
				},
			},
		},
	}}}, nil
}

// box the value to return a nil-able reference
func some[T any](t T) *T {
	return &t
}
