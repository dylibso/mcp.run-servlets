package main

import (
	"net/url"
	"strconv"

	"github.com/extism/go-pdk"
)

var (
	WebSearchTool = ToolDescription{
		Name: "brave-web-search",
		Description: "Queries Brave Search and get back search results from the web. " +
			"The results contains image and text URLs that can be fetched using the fetch tool.",
		InputSchema: schema{
			"type":     "object",
			"required": []string{"q"},
			"properties": map[string]SchemaProperty{
				"q": prop("string", "The user’s search query term. Query can not be empty. Maximum of 400 characters and 50 words in the query."),
				"country": prop("string", "The search query country, where the results come from. "+
					"The country string is limited to 2 character country codes of supported countries. Default: US"),
				"search_lang": prop("string", "The search language preference. "+
					"The 2 or more character language code for which the search results are provided. Default: en"),
				"ui_lang": prop("string", "User interface language preferred in response. "+
					"Usually of the format ‘<language_code>-<country_code>’ (RFC 9110). Default: en-US"),
				"count": prop("number", "The number of search results returned in response. The maximum is 20. "+
					"The actual number delivered may be less than requested. Combine this parameter with offset to paginate search results. Default: 20"),
				"offset": prop("number", "The zero based offset that indicates number of search results per page (count) to skip before returning the result. "+
					"The maximum is 9. The actual number delivered may be less than requested based on the query. "+
					"In order to paginate results use this parameter together with count. Default: 0"),
				"safesearch": prop("string", "Filters search results for adult content. "+
					"off: No filtering is done. moderate: Filters explicit content, like images and videos but allows adult domains in the search results. "+
					"strict: Drops all adult content from search results. Default: moderate"),
				"freshness": prop("string", "Filters search results by when they were discovered. "+
					"pd: Discovered within the last 24 hours. pw: Discovered within the last 7 Days. pm: Discovered within the last 31 Days. "+
					"py: Discovered within the last 365 Days… YYYY-MM-DDtoYYYY-MM-DD: timeframe is also supported by specifying the date range e.g. 2022-04-01to2022-07-30."),
				"text_decorations": prop("bool", "Whether display strings (e.g. result snippets) should include decoration markers (e.g. highlighting characters). Default: true"),
				"spellcheck":       prop("bool", "Whether to spellcheck provided query. If the spellchecker is enabled, the modified query is always used for search. The modified query can be found in altered key from the query response model. Default: true"),
				"result_filter": prop("string", "A comma delimited string of result types to include in the search response. "+
					"Not specifying this parameter will return back all result types in search response where data is available and a plan with the corresponding option is subscribed. "+
					"The response always includes query and type to identify any query modifications and response type respectively. "+
					"Av ailable result filter values are: discussions, faq, infobox, news, query, summarizer, videos, web, locations"),
				"goggles_id": prop("string", "Goggles act as a custom re-ranking on top of Brave’s search index. For more details, refer to the Goggles repository."),
				"units": prop("string", "The measurement units. If not provided, units are derived from search country. "+
					"Possible values are: metric: The standardized measurement system, imperial: The British Imperial system of units."),
				"extra_snippets": prop("bool", "A snippet is an excerpt from a page you get as a result of the query, and extra_snippets allow you to get up to 5 additional, alternative excerpts. "+
					"Only available under Free AI, Base AI, Pro AI, Base Data, Pro Data and Custom plans. Default: true"),
				"summary": prop("bool", "This parameter enables summary key generation in web search results. This is required for summarizer to be enabled."),
			},
		},
	}
)

func callWebSearch(args args) (CallToolResult, error) {
	params := url.Values{}
	if q, ok := args.string("q"); ok {
		params.Set("q", q)
	} else {
		return callToolError("missing required argument q"), nil
	}
	if country, ok := args.string("country"); ok {
		params.Set("country", country)
	}

	if searchLang, ok := args.string("search_lang"); ok {
		params.Set("search_lang", searchLang)
	}

	if uiLang, ok := args.string("ui_lang"); ok {
		params.Set("ui_lang", uiLang)
	}

	if count, ok := args.number("count"); ok {
		params.Set("count", strconv.Itoa(int(count)))
	}

	if offset, ok := args.number("offset"); ok {
		params.Set("offset", strconv.Itoa(int(offset)))
	}

	if safesearch, ok := args.string("safesearch"); ok {
		params.Set("safesearch", safesearch)
	}

	if freshness, ok := args.string("freshness"); ok {
		params.Set("freshness", freshness)
	}

	if textDecorations, ok := args.bool("text_decorations"); ok {
		params.Set("text_decorations", strconv.FormatBool(textDecorations))
	}

	if spellcheck, ok := args.bool("spellcheck"); ok {
		params.Set("spellcheck", strconv.FormatBool(spellcheck))
	}

	if resultFilter, ok := args.string("result_filter"); ok {
		params.Set("result_filter", resultFilter)
	}

	if gogglesID, ok := args.string("goggles_id"); ok {
		params.Set("goggles_id", gogglesID)
	}

	if units, ok := args.string("units"); ok {
		params.Set("units", units)
	}

	if extraSnippets, ok := args.bool("extra_snippets"); ok {
		params.Set("extra_snippets", strconv.FormatBool(extraSnippets))
	} else {
		params.Set("extra_snippets", "true")
	}

	if summary, ok := args.bool("summary"); ok {
		params.Set("summary", strconv.FormatBool(summary))
	}

	// Call the Brave Search API
	u := "https://api.search.brave.com/res/v1/web/search?" + params.Encode()
	req := pdk.NewHTTPRequest(pdk.MethodGet, u)
	req.SetHeader("X-Subscription-Token", apiKey)
	res := req.Send()

	pdk.Log(pdk.LogDebug, apiKey)

	return callToolTextSuccess(res.Body()), nil
}
