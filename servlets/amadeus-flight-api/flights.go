package main

import (
	"fmt"
	"net/url"

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
	q.Set("returnDate", args["returnDate"].(string))
	q.Set("adults", fmt.Sprint(args["adults"]))
	if children, hasChildren := args["children"]; hasChildren {
		q.Set("children", fmt.Sprint(children))
	}
	if infants, hasInfants := args["infants"]; hasInfants {
		q.Set("infants", fmt.Sprint(infants))
	}
	q.Set("travelClass", args["travelClass"].(string))
	if nonStop, hasNonStop := args["nonStop"]; hasNonStop {
		q.Set("nonStop", fmt.Sprint(nonStop))
	}
	q.Set("currencyCode", args["currencyCode"].(string))
	if maxPrice, hasMaxPrice := args["maxPrice"]; hasMaxPrice {
		q.Set("maxPrice", fmt.Sprint(maxPrice))
	}
	if max, hasMax := args["max"]; hasMax {
		q.Set("max", fmt.Sprint(max))
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
