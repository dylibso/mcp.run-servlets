package main

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"

	xtptest "github.com/dylibso/xtp-test-go"
	"github.com/extism/go-pdk"
)

//go:export test
func test() int32 {
	output := xtptest.CallBytes("describe", []byte{})
	list, err := parseToolListResult(output)
	if err != nil {
		xtptest.Assert("Failed to parse tool list", false, err.Error())
		return 1
	}

	for _, tool := range list.Tools {
		if len(tool.Name) == 0 {
			xtptest.Assert("Name should be non-empty", false, "Name is empty")
		} else {
			xtptest.Assert("Successfully parsed tool description", true, tool.Name)
		}

		testToolCall(tool.Name)
	}

	return 0
}

func testToolCall(name string) {
	switch name {
	case "greet":
		testGreet()
	case "qr-code":
		testQRCode()
	case "currency-converter":
		testCurrencyConverter()
	case "eval_js":
		testEvalJS()
	case "fetch":
		testFetch()
	case "fetch-image":
		testFetchImage()
	case "gif-search":
		testTenor()
	default:
		xtptest.Assert(fmt.Sprintf("Unable to test '%s'", name), false, "Unknown tool")
	}
}

func testTenor() {
	xtptest.Group("test tenor gif-search tool", func() {
		pdk.Log(pdk.LogDebug, "Testing Tenor GIF search tool")

		arguments := map[string]interface{}{
			"query": "happy",
			"limit": 2,
		}

		input := CallToolRequest{
			Method: nil,
			Params: Params{
				Arguments: &arguments,
				Name:      "gif-search",
			},
		}

		inputBytes, err := input.Marshal()
		if err != nil {
			xtptest.Assert("Failed to marshal input", false, err.Error())
		}

		output := xtptest.CallBytes("call", inputBytes)
		result, err := parseCallToolResult(output)
		if err != nil {
			xtptest.Assert("Failed to parse tool call result", false, err.Error())
		}

		pdk.Log(pdk.LogDebug, fmt.Sprintf("Tool call result: %v", string(output)))

		hasErrored := result.IsError != nil && *result.IsError
		xtptest.AssertEq("Tool call should not have errored", hasErrored, false)

		// Should have at least 3 content items (header + GIF + attribution for first result)
		xtptest.Assert("Tool call should have multiple content items", len(result.Content) >= 3, fmt.Sprintf("Got %d content items", len(result.Content)))

		// First content should be text with "Powered by Tenor"
		xtptest.AssertEq("First content should be text", result.Content[0].Type, ContentTypeText)
		xtptest.Assert("Header should contain Tenor attribution",
			strings.Contains(*result.Content[0].Text, "Powered by Tenor"),
			*result.Content[0].Text)

		// Second content should be an image
		xtptest.AssertEq("Second content should be image", result.Content[1].Type, ContentTypeImage)
		xtptest.AssertEq("Image should have GIF mime type", *result.Content[1].MimeType, "image/gif")

		// Validate base64 GIF data
		imageData, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(*result.Content[1].Data, "data:image/gif;base64,"))
		if err != nil {
			xtptest.Assert("Failed to decode GIF data", false, err.Error())
			return
		}

		// Check GIF magic number
		if len(imageData) < 6 {
			xtptest.Assert("GIF data too short", false, fmt.Sprintf("Got %d bytes", len(imageData)))
			return
		}
		xtptest.AssertEq("Image should start with GIF magic number", string(imageData[:6]), "GIF89a")

		// Third content should be text with attribution
		xtptest.AssertEq("Third content should be text", result.Content[2].Type, ContentTypeText)
	})
}

func testFetch() {
	xtptest.Group("test fetch tool", func() {
		pdk.Log(pdk.LogDebug, "Testing fetch tool")

		arguments := map[string]interface{}{"url": "https://getxtp.com"}
		input := CallToolRequest{
			Method: nil,
			Params: Params{
				Arguments: &arguments,
				Name:      "fetch",
			},
		}

		inputBytes, err := input.Marshal()
		if err != nil {
			xtptest.Assert("Failed to marshal input", false, err.Error())
		}

		output := xtptest.CallBytes("call", inputBytes)
		result, err := parseCallToolResult(output)
		if err != nil {
			xtptest.Assert("Failed to parse tool call result", false, err.Error())
		}

		pdk.Log(pdk.LogDebug, fmt.Sprintf("Tool call result: %v", string(output)))

		hasErrored := result.IsError != nil && *result.IsError

		xtptest.AssertEq("Tool call should not have errored", hasErrored, false)
		xtptest.AssertEq("Tool call should have one content item", len(result.Content), 1)
		xtptest.AssertEq("Content type should be text", result.Content[0].Type, ContentTypeText)

		// test if the content contains the string "xtp", case-insensitive
		regex := regexp.MustCompile(`(?i)xtp`)
		xtptest.Assert("Content text should contain 'xtp'", regex.MatchString(*result.Content[0].Text), *result.Content[0].Text)
	})
}

func testFetchImage() {
	xtptest.Group("test fetch-image tool", func() {
		pdk.Log(pdk.LogDebug, "Testing fetch image tool")

		arguments := map[string]interface{}{"url": "https://httpbin.org/image/jpeg", "mime-type": "image/jpeg"}
		input := CallToolRequest{
			Method: nil,
			Params: Params{
				Arguments: &arguments,
				Name:      "fetch-image",
			},
		}

		inputBytes, err := input.Marshal()
		if err != nil {
			xtptest.Assert("Failed to marshal input", false, err.Error())
		}

		output := xtptest.CallBytes("call", inputBytes)
		result, err := parseCallToolResult(output)
		if err != nil {
			xtptest.Assert("Failed to parse tool call result", false, err.Error())
		}

		pdk.Log(pdk.LogDebug, fmt.Sprintf("Tool call result: %v", string(output)))

		hasErrored := result.IsError != nil && *result.IsError

		xtptest.AssertEq("Tool call should not have errored", hasErrored, false)
		xtptest.AssertEq("Tool call should have one content item", len(result.Content), 1)
		xtptest.AssertEq("Content type should be text", result.Content[0].Type, ContentTypeImage)

		imageData, err := base64.StdEncoding.DecodeString(*result.Content[0].Data)
		if err != nil {
			xtptest.Assert("Failed to decode image data", false, err.Error())
			return
		}

		if len(imageData) == 0 {
			xtptest.Assert("Image data should be greater than 0", false, "Image data is empty")
			return
		}

		// test the image for the magic number for JPEG
		xtptest.AssertEq("Image data should start with JPEG magic number", string(imageData[:2]), "\xFF\xD8")
	})
}

func testEvalJS() {
	xtptest.Group("test eval-js tool", func() {
		pdk.Log(pdk.LogDebug, "Testing eval-js tool")

		arguments := map[string]interface{}{"code": "1 + 1"}
		input := CallToolRequest{
			Method: nil,
			Params: Params{
				Arguments: &arguments,
				Name:      "eval-js",
			},
		}

		inputBytes, err := input.Marshal()
		if err != nil {
			xtptest.Assert("Failed to marshal input", false, err.Error())
		}

		output := xtptest.CallBytes("call", inputBytes)
		result, err := parseCallToolResult(output)
		if err != nil {
			xtptest.Assert("Failed to parse tool call result", false, err.Error())
		}

		pdk.Log(pdk.LogDebug, fmt.Sprintf("Tool call result: %v", string(output)))

		hasErrored := result.IsError != nil && *result.IsError

		xtptest.AssertEq("Tool call should not have errored", hasErrored, false)
		xtptest.AssertEq("Tool call should have one content item", len(result.Content), 1)
		xtptest.AssertEq("Content type should be text", result.Content[0].Type, ContentTypeText)
		xtptest.AssertEq("Content text should be '2'", *result.Content[0].Text, "2")
	})
}

func testGreet() {
	xtptest.Group("test greet tool", func() {
		pdk.Log(pdk.LogDebug, "Testing greet tool")

		arguments := map[string]interface{}{"name": "Steve"}
		input := CallToolRequest{
			Method: nil,
			Params: Params{
				Arguments: &arguments,
				Name:      "greet",
			},
		}

		inputBytes, err := input.Marshal()
		if err != nil {
			xtptest.Assert("Failed to marshal input", false, err.Error())
		}

		output := xtptest.CallBytes("call", inputBytes)
		result, err := parseCallToolResult(output)
		if err != nil {
			xtptest.Assert("Failed to parse tool call result", false, err.Error())
		}

		pdk.Log(pdk.LogDebug, fmt.Sprintf("Tool call result: %v", string(output)))

		hasErrored := result.IsError != nil && *result.IsError

		xtptest.AssertEq("Tool call should not have errored", hasErrored, false)
		xtptest.AssertEq("Tool call should have one content item", len(result.Content), 1)
		xtptest.AssertEq("Content type should be text", result.Content[0].Type, ContentTypeText)

		xtptest.AssertEq("Content text should be `Hello Steve!!!`", "Hello Steve!!!", *result.Content[0].Text)
	})
}

func testQRCode() {
	xtptest.Group("test qr-code tool", func() {
		pdk.Log(pdk.LogDebug, "Testing QR Code tool")

		arguments := map[string]interface{}{
			"data":  "hello, world",
			"ecc":   1,
			"width": 200,
		}

		input := CallToolRequest{
			Method: nil,
			Params: Params{
				Arguments: &arguments,
				Name:      "qr-code",
			},
		}

		inputBytes, err := input.Marshal()
		if err != nil {
			xtptest.Assert("Failed to marshal input", false, err.Error())
		}

		output := xtptest.CallBytes("call", inputBytes)
		result, err := parseCallToolResult(output)
		if err != nil {
			xtptest.Assert(fmt.Sprintf("Failed to parse tool call result: %s", err.Error()), false, string(output))
		}

		pdk.Log(pdk.LogDebug, fmt.Sprintf("Tool call result: %v", string(output)))

		hasErrored := result.IsError != nil && *result.IsError

		xtptest.AssertEq("Tool call should not have errored", hasErrored, false)
		xtptest.AssertEq("Tool call should have one content item", len(result.Content), 1)
		xtptest.AssertEq("Content type should be text", result.Content[0].Type, ContentTypeImage)
		xtptest.AssertEq("MimeType should be image/png", *result.Content[0].MimeType, "image/png")
	})
}

func testCurrencyConverter() {
	xtptest.Group("test currency-converter tool", func() {
		pdk.Log(pdk.LogDebug, "Testing QR Code tool")

		arguments := map[string]interface{}{
			"amount": 50.0,
			"from":   "USD",
			"to":     "EUR",
		}

		input := CallToolRequest{
			Method: nil,
			Params: Params{
				Arguments: &arguments,
				Name:      "currency-converter",
			},
		}

		inputBytes, err := input.Marshal()
		if err != nil {
			xtptest.Assert("Failed to marshal input", false, err.Error())
		}

		output := xtptest.CallBytes("call", inputBytes)
		result, err := parseCallToolResult(output)
		if err != nil {
			xtptest.Assert(fmt.Sprintf("Failed to parse tool call result: %s", err.Error()), false, string(output))
		}

		pdk.Log(pdk.LogDebug, fmt.Sprintf("Tool call result: %v", string(output)))

		hasErrored := result.IsError != nil && *result.IsError

		xtptest.AssertEq("Tool call should not have errored", hasErrored, false)

		convertedAmount := 0.0
		_, err = fmt.Sscanf(*result.Content[0].Text, "%f", &convertedAmount)
		if err != nil {
			xtptest.Assert("Failed to parse converted amount", false, err.Error())
		}

		pdk.Log(pdk.LogDebug, fmt.Sprintf("Converted amount: %f", convertedAmount))

		xtptest.Assert("Converted amount should be greater than 0", convertedAmount > 0, fmt.Sprintf("Converted amount: %f", convertedAmount))
	})
}

func main() {}
