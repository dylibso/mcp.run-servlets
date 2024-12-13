package main

import (
	"encoding/base64"
	"fmt"
	"regexp"

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
	default:
		xtptest.Assert(fmt.Sprintf("Unable to test '%s'", name), false, "Unknown tool")
	}
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
