package main

import (
	"fmt"
	"regexp"

	xtptest "github.com/dylibso/xtp-test-go"
	"github.com/extism/go-pdk"
)

//go:export test
func test() int32 {
	xtptest.Group("basic tests", func() {
		output := xtptest.CallBytes("describe", []byte{})
		description, err := parseToolDescription(output)
		if err != nil {
			xtptest.Assert("Failed to parse tool description", false, err.Error())
			return
		}

		xtptest.AssertGte("name should be non-empty", len(description.Name), 1)
		xtptest.AssertGte("description should be non-empty", len(description.Description), 1)
		xtptest.AssertNe("inputSchema should be non-nil", description.InputSchema, nil)

		testToolCall(description.Name)
	})

	return 0
}

func testToolCall(name string) {
	switch name {
	case "greet":
		testGreet()

	default:
		xtptest.Assert(fmt.Sprintf("Unable to test %s", name), false, "Unknown tool")
	}
}

func testGreet() {
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
}

func main() {}
