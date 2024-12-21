package main

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	pdk "github.com/extism/go-pdk"
)

var (
	FlightOfferSearchTool = ToolDescription{
		Name: "am-flights-offer-search",
		Description: "Search for flight offers. For each itinerary, the API provides a list of flight offers with prices, " +
			"fare details, airline names, baggage allowances and departure terminals. " +
			"The results can be displayed as a table. Users should be offered to display such a table",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"originLocationCode":      prop("string", "city/airport IATA code from which the traveler will depart, e.g. BOS for Boston"),
				"destinationLocationCode": prop("string", "city/airport IATA code to which the traveler is going, e.g. PAR for Paris"),
				"departureDate":           prop("string", "the date on which the traveler will depart from the origin to go to the destination. Dates are specified in the ISO 8601 YYYY-MM-DD format, e.g. 2017-12-25"),
				"returnDate": prop("string", "the date on which the traveler will depart from the destination to return to the origin. "+
					"If this parameter is not specified, only one-way itineraries are found. "+
					"If this parameter is specified, only round-trip itineraries are found. Dates are specified in the ISO 8601 YYYY-MM-DD format, e.g. 2018-02-28"),
				"adults": prop("integer", "the number of adult travelers (age 12 or older on date of departure)."+
					"The total number of seated travelers (adult and children) cannot exceed 9."),
				"children": prop("integer", "(optional) The number of children"),
				"infants":  prop("integer", "(optional) The number of infants"),
				"travelClass": prop("string", "most of the flight time should be spent in a cabin of this quality or higher. "+
					"The accepted travel class is economy, premium economy, business or first class. If no travel class is specified, "+
					"the search considers any travel class. Available values : ECONOMY, PREMIUM_ECONOMY, BUSINESS, FIRST"),
				"nonStop":      prop("boolean", "If set to true, the search will find only flights going from the origin to the destination with no stop in between"),
				"currencyCode": prop("string", "The preferred currency for the flight offers. Currency is specified in the ISO 4217 format, e.g. EUR for Euro"),
				"maxPrice":     prop("integer", "The maximum price of the flight offers. If specified, the value should be a positive number with no decimals"),
				"max":          prop("integer", "The maximum number of offers to return (default: 10)"),
			},
			"required": []string{"originLocationCode", "destinationLocationCode", "departureDate", "adults", "travelClass", "currencyCode"},
		},
	}

	FlightsTools = []ToolDescription{
		FlightOfferSearchTool,
	}
)

// https://developers.amadeus.com/self-service/category/flights/api-doc/flight-offers-search
func flightOfferSearch(args map[string]any) CallToolResult {
	const endpoint = "/v2/shopping/flight-offers"

	q := url.Values{}
	q.Set("originLocationCode", args["originLocationCode"].(string))
	q.Set("destinationLocationCode", args["destinationLocationCode"].(string))
	q.Set("departureDate", args["departureDate"].(string))

	if returnDate, hasReturnDate := getString(args, "returnDate"); hasReturnDate {
		q.Set("returnDate", returnDate)
	}

	if adults, hasAdults := getNumber(args, "adults"); hasAdults {
		q.Set("adults", fmt.Sprint(adults))
	}

	children, _ := getNumber(args, "children")
	q.Set("children", fmt.Sprint(children))

	infants, _ := getNumber(args, "infants")
	q.Set("infants", fmt.Sprint(infants))

	if travelClass, hasTravelClass := getString(args, "travelClass"); hasTravelClass {
		q.Set("travelClass", strings.ToUpper(travelClass))
	} else {
		q.Set("travelClass", "ECONOMY")
	}

	nonStop, _ := getBool(args, "nonStop")

	q.Set("nonStop", fmt.Sprint(nonStop))

	if currencyCode, hasCurrencyCode := getString(args, "currencyCode"); hasCurrencyCode {
		q.Set("currencyCode", currencyCode)
	} else {
		q.Set("currencyCode", "USD")
	}

	if maxPrice, hasMaxPrice := getString(args, "maxPrice"); hasMaxPrice {
		q.Set("maxPrice", fmt.Sprint(maxPrice))
	}
	if max, hasMax := getNumber(args, "max"); hasMax {
		q.Set("max", fmt.Sprint(max))
	} else {
		q.Set("max", fmt.Sprint(10))
	}

	req := pdk.NewHTTPRequest(pdk.MethodGet, config.baseUrl+endpoint+"?"+q.Encode())
	req.SetHeader("Authorization", "Bearer "+config.token)
	req.SetHeader("Accept", "application/json")
	resp := req.Send()

	if resp.Status() != 200 {
		return callToolError("error while searching for flight offers " + string(resp.Body()))
	}

	return callToolSuccess(string(resp.Body()))
}

func getNumber(args map[string]any, key string) (float64, bool) {
	if n, ok := args[key]; ok && n != nil {
		if i, ok := n.(float64); ok {
			return i, true
		}
		if i, ok := n.(int); ok {
			return float64(i), true
		}
		if i, ok := n.(string); ok {
			if f, err := strconv.ParseFloat(i, 64); err == nil {
				return f, true
			}
		}
	}
	return 0, false
}

func getBool(args map[string]any, key string) (bool, bool) {
	if n, ok := args[key]; ok && n != nil {
		if i, ok := n.(bool); ok {
			return i, true
		}
		if i, ok := n.(string); ok {
			if f, err := strconv.ParseBool(i); err == nil {
				return f, true
			}
		}
	}
	return false, false
}

func getString(args map[string]any, key string) (string, bool) {
	if s, ok := args[key]; ok && s != nil {
		if ss, ok := s.(string); ok {
			ss = strings.TrimSpace(ss)
			if ss == "" {
				return "", false
			}
			return ss, ok
		}
	}
	return "", false
}
