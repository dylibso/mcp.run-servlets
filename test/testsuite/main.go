package main

import (
	"fmt"
	"regexp"

	xtptest "github.com/dylibso/xtp-test-go"
	"github.com/extism/go-pdk"
)

//go:export test
func test() int32 {
	output := xtptest.CallBytes("describe", []byte{})
	description, err := parseToolDescription(output)
	if err != nil {
		xtptest.Assert("Failed to parse tool description", false, err.Error())
		return 1
	}

	if len(description.Name) == 0 {
		xtptest.Assert("Name should be non-empty", false, "Name is empty")
	} else {
		xtptest.Assert("Successfully parsed tool description", true, description.Name)
	}

	testToolCall(description.Name)

	return 0
}

func testToolCall(name string) {
	switch name {
	case "greet":
		testGreet()
	case "qr-code":
		testQRCode()
	default:
		xtptest.Assert(fmt.Sprintf("Unable to test '%s'", name), false, "Unknown tool")
	}
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

		regex := regexp.MustCompile(`^(.)+ says: Hello, Steve!$`)
		xtptest.Assert("Content text should match `X says: Hello, Steve!`", regex.MatchString(*result.Content[0].Text), *result.Content[0].Text)
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
			xtptest.Assert("Failed to parse tool call result", false, err.Error())
		}

		pdk.Log(pdk.LogDebug, fmt.Sprintf("Tool call result: %v", string(output)))

		hasErrored := result.IsError != nil && *result.IsError

		xtptest.AssertEq("Tool call should not have errored", hasErrored, false)
	})
}

func main() {}