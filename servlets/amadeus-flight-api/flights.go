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

	FlightDatesTool = ToolDescription{
		Name: "am-flights-dates",
		Description: "Find the cheapest flight dates from an origin to a destination. " +
			"The API provides list of flight options with dates and prices, and allows you to order by price, departure date or duration. " +
			"The results can be displayed as a table. Users should be offered to display such a table",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"origin":      prop("string", "city/airport IATA code from which the traveler will depart, e.g. BOS for Boston"),
				"destination": prop("string", "city/airport IATA code to which the traveler is going, e.g. PAR for Paris"),
				"departureDate": prop("string", "the date on which the traveler will depart from the origin to go to the destination. "+
					"Dates are specified in the ISO 8601 YYYY-MM-DD format, e.g. 2017-12-25. Ranges are specified with a comma and are inclusive, e.g. 2017-12-25,2017-12-30"),
				"oneWay":   prop("boolean", "one way or round trip. if this parameter is set to true, only one-way flights are considered. If this parameter is not set or set to false, only round-trip flights are considered"),
				"duration": prop("string", "exact duration or range of durations of the travel, in days. This parameter must not be set if oneWay is true. Ranges are specified with a comma and are inclusive, e.g. 2,8"),
				"nonStop":  prop("boolean", "if this parameter is set to true, only flights going from the origin to the destination with no stop in-between are considered. Default value: false"),
				"maxPrice": prop("integer", "the maximum price of the flight offers. The value should be a positive number with no decimals"),
				"viewBy": prop("string", "view the flight dates by DATE, DURATION, or WEEK. View by DATE (default when oneWay is true) to get the cheapest flight dates for every departure date in the given range. "+
					"View by DURATION (default when oneWay is false) to get the cheapest flight dates for every departure date and for every duration in the given ranges. "+
					"View by WEEK to get the cheapest flight destination for every week in the given range of departure dates. Note that specifying a detailed view but large ranges may result in a huge number of flight dates being returned. "+
					"For some very large numbers of flight dates, the API may refuse to provide a response."+
					"Available values : DATE, DURATION, WEEK"),
			},
			"required": []string{"origin", "destination", "departureDate", "maxPrice", "viewBy"},
		},
	}

	FlightInspirationTool = ToolDescription{
		Name: "am-flights-inspiration",
		Description: "Find the cheapest destinations where you can fly to." +
			"The Flight Inspiration Search API provides a list of destinations from a given city that is ordered by price and can be filtered by departure date or maximum price" +
			"The results can be displayed as a table. Users should be offered to display such a table",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"origin": prop("string", "IATA code of the city from which the flight will depart"),
				"departureDate": prop("string", "The date, or range of dates, on which the flight will depart from the origin. "+
					"Dates are specified in the ISO 8601 YYYY-MM-DD format, e.g. 2017-12-25. Ranges are specified with a comma and are inclusive. "+
					"Departure date can not be more than 180 days in the future."),
				"oneWay": prop("boolean", "if this parameter is set to true, only one-way flights are considered. "+
					"If this parameter is not set or set to false, only round-trip flights are considered"),
				"duration": prop("string", "Exact duration or range of durations of the travel, in days. "+
					"This parameter must not be set if oneWay is true. Ranges are specified with a comma and are inclusive, e.g. 2,8. "+
					"Duration can not be lower than 1 days or higher than 15 days"),
				"nonStop": prop("boolean", "if this parameter is set to true, only flights going from the origin to the destination with no stop in-between are considered. "+
					"Default value: false"),
				"maxPrice": prop("integer", "defines the price limit for each offer returned. The value should be a positive number, without decimals"),
				"viewBy": prop("string", "view the flight destinations by DATE, DESTINATION, DURATION, WEEK, or COUNTRY. "+
					"View by DATE (default when oneWay is true) to get the cheapest flight destination for every departure date in the given range. "+
					"View by DURATION (default when oneWay is false) to get the cheapest flight destination for every departure date and for every duration in the given ranges. "+
					"View by WEEK to get the cheapest flight destination for every week in the given range of departure dates. "+
					"View by COUNTRY to get the cheapest flight destination by country. Note that specifying a detailed view but large ranges may result in a huge number of flight destinations being returned. "+
					"For some very large numbers of flight destinations, the API may refuse to provide a response"),
			},
		},
	}

	FlightsTools = []ToolDescription{
		FlightOfferSearchTool, FlightDatesTool, FlightInspirationTool,
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
	q.Set("nonStop", fmt.Sprint(args["nonStop"]))
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
		return callToolError("error while searching for flight offers")
	}

	return callToolSuccess(string(resp.Body()))
}

// https://developers.amadeus.com/self-service/category/flights/api-doc/flight-cheapest-date-search/api-reference
func flightDates(args map[string]any) CallToolResult {
	const endpoint = "/v1/shopping/flight-dates"

	q := url.Values{}

	// origin IATA
	q.Add("origin", args["origin"].(string))
	// destination IATA
	q.Add("destination", args["destination"].(string))
	// departure date
	q.Add("departureDate", args["departureDate"].(string))
	// one way or round trip
	if oneWay, hasOneWay := args["oneWay"]; hasOneWay {
		q.Add("oneWay", fmt.Sprint(oneWay))
	}
	// exact duration or range of durations of the travel, in days.
	// This parameter must not be set if oneWay is true.
	// Ranges are specified with a comma and are inclusive, e.g. 2,8
	if duration, hasDuration := args["duration"]; hasDuration {
		q.Add("duration", fmt.Sprint(duration))
	}
	// direct flights only
	if nonStop, hasNonStop := args["nonStop"]; hasNonStop {
		q.Add("nonStop", fmt.Sprint(nonStop))
	}
	// max price
	q.Add("maxPrice", fmt.Sprint(args["maxPrice"]))
	// view the flight dates by DATE, DURATION, or WEEK.
	// View by DATE (default when oneWay is true) to get the cheapest flight dates for every departure date in the given range.
	// View by DURATION (default when oneWay is false) to get the cheapest flight dates for every departure date and
	// for every duration in the given ranges. View by WEEK to get the cheapest flight destination for every week
	// in the given range of departure dates.
	// Note that specifying a detailed view but large ranges may result in a huge number of flight dates being returned.
	// For some very large numbers of flight dates, the API may refuse to provide a response
	//
	// Available values : DATE, DURATION, WEEK
	q.Add("viewBy", args["viewBy"].(string))

	req := pdk.NewHTTPRequest(pdk.MethodGet, config.baseUrl+endpoint+"?"+q.Encode())
	req.SetHeader("Authorization", "Bearer "+config.token)
	req.SetHeader("Accept", "application/json")
	resp := req.Send()

	if resp.Status() != 200 {
		return callToolError("error while searching for flight dates " + string(resp.Body()))
	}

	return callToolSuccess(string(resp.Body()))
}

// https://developers.amadeus.com/self-service/category/flights/api-doc/flight-inspiration-search/api-reference
func flightInspiration(args map[string]any) CallToolResult {
	const endpoint = "/v1/shopping/flight-destinations"

	q := url.Values{}

	// origin IATA
	q.Add("origin", args["origin"].(string))
	// departure date
	q.Add("departureDate", args["departureDate"].(string))
	// one way or round trip
	if oneWay, hasOneWay := args["oneWay"]; hasOneWay {
		q.Add("oneWay", fmt.Sprint(oneWay))
	}
	// exact duration or range of durations of the travel, in days.
	// This parameter must not be set if oneWay is true.
	// Ranges are specified with a comma and are inclusive, e.g. 2,8
	if duration, hasDuration := args["duration"]; hasDuration {
		q.Add("duration", fmt.Sprint(duration))
	}
	// direct flights only
	if nonStop, hasNonStop := args["nonStop"]; hasNonStop {
		q.Add("nonStop", fmt.Sprint(nonStop))
	}
	// max price
	q.Add("maxPrice", fmt.Sprint(args["maxPrice"]))
	// view the flight destinations by DATE, DESTINATION, DURATION, WEEK, or COUNTRY.
	// View by DATE (default when oneWay is true) to get the cheapest flight destination for every departure date in the given range.
	// View by DURATION (default when oneWay is false) to get the cheapest flight destination for every departure date and
	// for every duration in the given ranges. View by WEEK to get the cheapest flight destination for every week
	// in the given range of departure dates. View by COUNTRY to get the cheapest flight destination by country.
	// Note that specifying a detailed view but large ranges may result in a huge number of flight destinations being returned.
	// For some very large numbers of flight destinations, the API may refuse to provide a response
	//
	// Available values : DATE, DURATION, WEEK
	q.Add("viewBy", args["viewBy"].(string))

	req := pdk.NewHTTPRequest(pdk.MethodGet, config.baseUrl+endpoint+"?"+q.Encode())
	req.SetHeader("Authorization", "Bearer "+config.token)
	req.SetHeader("Accept", "application/json")
	resp := req.Send()

	if resp.Status() != 200 {
		return callToolError("error while searching for flight inspiration " + string(resp.Body()))
	}

	return callToolSuccess(string(resp.Body()))
}
