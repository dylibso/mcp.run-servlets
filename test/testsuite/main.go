package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
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
	case "eval-py":
		testEvalPy()
	case "fetch":
		testFetch()
	case "fetch-image":
		testFetchImage()
	case "gif-search":
		testTenor()
	case "coulomb_force":
		testCoulombForce()
	case "induced_emf":
		testInducedEMF()
	case "magnetic_field":
		testMagneticField()
	case "lorentz_force":
		testLorentzForce()
	case "cyclotron_frequency":
		testCyclotronFrequency()
	case "electric_potential_energy":
		testElectricPotentialEnergy()
	case "magnetic_flux":
		testMagneticFlux()
	case "capacitor_energy":
		testCapacitorEnergy()
	case "solenoid_inductance":
		testSolenoidInductance()
	case "rc_time_constant":
		testRCTimeConstant()
	case "validate":
	default:
		xtptest.Assert(fmt.Sprintf("Unable to test '%s'", name), false, "Unknown tool")
	}
}

func testJSONValidator() {
	xtptest.Group("test json-schema validate tool", func() {
		pdk.Log(pdk.LogDebug, "Testing JSON Schema validate tool")

		arguments := map[string]interface{}{
			"schema": map[string]interface{}{
				"$schema": "http://json-schema.org/draft-07/schema#",
				"title":   "Person",
				"type":    "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						"type":        "string",
						"description": "The person's full name",
					},
					"age": map[string]interface{}{
						"type":        "integer",
						"minimum":     0,
						"description": "The person's age in years",
					},
					"isStudent": map[string]interface{}{
						"type":        "boolean",
						"description": "Whether the person is a student",
					},
				},
				"required": []string{"name", "age"},
			},
			"document": map[string]interface{}{
				"name":      "John Doe",
				"age":       25,
				"isStudent": true,
			},
		}

		input := CallToolRequest{
			Method: nil,
			Params: Params{
				Arguments: &arguments,
				Name:      "validate",
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
		xtptest.AssertEq("Validation should pass", *result.Content[0].Text, "true")
	})
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

func testEvalPy() {
	xtptest.Group("test eval-py tool", func() {
		pdk.Log(pdk.LogDebug, "Testing eval-py tool")

		arguments := map[string]interface{}{"code": "print(1 + 1)"}
		input := CallToolRequest{
			Method: nil,
			Params: Params{
				Arguments: &arguments,
				Name:      "eval-py",
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
		xtptest.AssertEq("Content text should be '2\n'", *result.Content[0].Text, "2\n")
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

func testCoulombForce() {
	xtptest.Group("test coulomb_force tool", func() {
		pdk.Log(pdk.LogDebug, "Testing Coulomb force calculation")

		arguments := map[string]interface{}{
			"charge1": 1e-6, // 1 microcoulomb
			"position1": map[string]interface{}{
				"x": 0.0,
				"y": 0.0,
				"z": 0.0,
			},
			"charge2": -1e-6, // -1 microcoulomb
			"position2": map[string]interface{}{
				"x": 1.0,
				"y": 0.0,
				"z": 0.0,
			},
		}

		input := CallToolRequest{
			Method: nil,
			Params: Params{
				Arguments: &arguments,
				Name:      "coulomb_force",
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

		hasErrored := result.IsError != nil && *result.IsError
		xtptest.AssertEq("Tool call should not have errored", hasErrored, false)
		xtptest.AssertEq("Tool call should have one content item", len(result.Content), 1)
		xtptest.AssertEq("Content type should be text", result.Content[0].Type, ContentTypeText)

		// Parse the JSON response
		var response struct {
			Force struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
				Z float64 `json:"z"`
			} `json:"force"`
			Magnitude float64 `json:"magnitude"`
			Unit      string  `json:"unit"`
		}
		err = json.Unmarshal([]byte(*result.Content[0].Text), &response)
		if err != nil {
			xtptest.Assert("Failed to parse response JSON", false, err.Error())
		}

		xtptest.AssertEq("Unit should be Newtons", response.Unit, "Newtons")
		xtptest.Assert("Force magnitude should be non-zero", response.Magnitude > 0,
			fmt.Sprintf("Got magnitude: %f", response.Magnitude))
		xtptest.Assert("Force should be attractive (negative x component)",
			response.Force.X < 0, fmt.Sprintf("Got x component: %f", response.Force.X))
	})
}

func testInducedEMF() {
	xtptest.Group("test induced_emf tool", func() {
		pdk.Log(pdk.LogDebug, "Testing induced EMF calculation")

		arguments := map[string]interface{}{
			"fluxChange":   0.5, // 0.5 Weber
			"timeInterval": 0.1, // 0.1 seconds
			"turns":        100, // 100 turns
		}

		input := CallToolRequest{
			Method: nil,
			Params: Params{
				Arguments: &arguments,
				Name:      "induced_emf",
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

		hasErrored := result.IsError != nil && *result.IsError
		xtptest.AssertEq("Tool call should not have errored", hasErrored, false)
		xtptest.AssertEq("Tool call should have one content item", len(result.Content), 1)
		xtptest.AssertEq("Content type should be text", result.Content[0].Type, ContentTypeText)

		// Parse the JSON response
		var response struct {
			EMF  float64 `json:"emf"`
			Unit string  `json:"unit"`
		}
		err = json.Unmarshal([]byte(*result.Content[0].Text), &response)
		if err != nil {
			xtptest.Assert("Failed to parse response JSON", false, err.Error())
		}

		xtptest.AssertEq("Unit should be Volts", response.Unit, "Volts")
		expectedEMF := -500.0 // -N * ΔΦ/Δt = -(100 * 0.5/0.1) = -500
		xtptest.Assert("EMF should match expected value",
			math.Abs(response.EMF-expectedEMF) < 0.001,
			fmt.Sprintf("Got EMF: %f, expected: %f", response.EMF, expectedEMF))
	})
}
func testMagneticField() {
	xtptest.Group("test magnetic_field tool", func() {
		pdk.Log(pdk.LogDebug, "Testing magnetic field calculation")

		arguments := map[string]interface{}{
			"current": 1.0, // 1 Ampere
			"wirePath": map[string]interface{}{
				"x": 0.0,
				"y": 0.0,
				"z": 1.0,
			},
			"observationPoint": map[string]interface{}{
				"x": 1.0,
				"y": 0.0,
				"z": 0.0,
			},
		}

		input := CallToolRequest{
			Method: nil,
			Params: Params{
				Arguments: &arguments,
				Name:      "magnetic_field",
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

		hasErrored := result.IsError != nil && *result.IsError
		xtptest.AssertEq("Tool call should not have errored", hasErrored, false)
		xtptest.AssertEq("Tool call should have one content item", len(result.Content), 1)
		xtptest.AssertEq("Content type should be text", result.Content[0].Type, ContentTypeText)

		var response struct {
			Field struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
				Z float64 `json:"z"`
			} `json:"field"`
			Magnitude float64 `json:"magnitude"`
			Unit      string  `json:"unit"`
		}
		err = json.Unmarshal([]byte(*result.Content[0].Text), &response)
		if err != nil {
			xtptest.Assert("Failed to parse response JSON", false, err.Error())
		}

		xtptest.AssertEq("Unit should be Tesla", response.Unit, "Tesla")
		xtptest.Assert("Field magnitude should be non-zero", response.Magnitude > 0,
			fmt.Sprintf("Got magnitude: %f", response.Magnitude))
	})
}

func testLorentzForce() {
	xtptest.Group("test lorentz_force tool", func() {
		pdk.Log(pdk.LogDebug, "Testing Lorentz force calculation")

		arguments := map[string]interface{}{
			"charge": 1.6e-19, // electron charge
			"velocity": map[string]interface{}{
				"x": 1000.0,
				"y": 0.0,
				"z": 0.0,
			},
			"electricField": map[string]interface{}{
				"x": 100.0,
				"y": 0.0,
				"z": 0.0,
			},
			"magneticField": map[string]interface{}{
				"x": 0.0,
				"y": 0.0,
				"z": 1.0,
			},
		}

		input := CallToolRequest{
			Method: nil,
			Params: Params{
				Arguments: &arguments,
				Name:      "lorentz_force",
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

		hasErrored := result.IsError != nil && *result.IsError
		xtptest.AssertEq("Tool call should not have errored", hasErrored, false)
		xtptest.AssertEq("Tool call should have one content item", len(result.Content), 1)
		xtptest.AssertEq("Content type should be text", result.Content[0].Type, ContentTypeText)

		var response struct {
			Force struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
				Z float64 `json:"z"`
			} `json:"force"`
			Magnitude float64 `json:"magnitude"`
			Unit      string  `json:"unit"`
		}
		err = json.Unmarshal([]byte(*result.Content[0].Text), &response)
		if err != nil {
			xtptest.Assert("Failed to parse response JSON", false, err.Error())
		}

		xtptest.AssertEq("Unit should be Newtons", response.Unit, "Newtons")
		xtptest.Assert("Force magnitude should be non-zero", response.Magnitude > 0,
			fmt.Sprintf("Got magnitude: %f", response.Magnitude))
	})
}

func testCyclotronFrequency() {
	xtptest.Group("test cyclotron_frequency tool", func() {
		pdk.Log(pdk.LogDebug, "Testing cyclotron frequency calculation")

		arguments := map[string]interface{}{
			"charge":        1.6e-19, // electron charge
			"magneticField": 1.0,     // 1 Tesla
			"mass":          9.1e-31, // electron mass
		}

		input := CallToolRequest{
			Method: nil,
			Params: Params{
				Arguments: &arguments,
				Name:      "cyclotron_frequency",
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

		hasErrored := result.IsError != nil && *result.IsError
		xtptest.AssertEq("Tool call should not have errored", hasErrored, false)
		xtptest.AssertEq("Tool call should have one content item", len(result.Content), 1)
		xtptest.AssertEq("Content type should be text", result.Content[0].Type, ContentTypeText)

		var response struct {
			Frequency float64 `json:"frequency"`
			Unit      string  `json:"unit"`
		}
		err = json.Unmarshal([]byte(*result.Content[0].Text), &response)
		if err != nil {
			xtptest.Assert("Failed to parse response JSON", false, err.Error())
		}

		xtptest.AssertEq("Unit should be Hertz", response.Unit, "Hertz")
		xtptest.Assert("Frequency should be positive", response.Frequency > 0,
			fmt.Sprintf("Got frequency: %f", response.Frequency))
	})
}

func testElectricPotentialEnergy() {
	xtptest.Group("test electric_potential_energy tool", func() {
		pdk.Log(pdk.LogDebug, "Testing electric potential energy calculation")

		arguments := map[string]interface{}{
			"charge1":  1e-6,  // 1 microcoulomb
			"charge2":  -1e-6, // -1 microcoulomb
			"distance": 1.0,   // 1 meter
		}

		input := CallToolRequest{
			Method: nil,
			Params: Params{
				Arguments: &arguments,
				Name:      "electric_potential_energy",
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

		hasErrored := result.IsError != nil && *result.IsError
		xtptest.AssertEq("Tool call should not have errored", hasErrored, false)
		xtptest.AssertEq("Tool call should have one content item", len(result.Content), 1)
		xtptest.AssertEq("Content type should be text", result.Content[0].Type, ContentTypeText)

		var response struct {
			Energy float64 `json:"energy"`
			Unit   string  `json:"unit"`
		}
		err = json.Unmarshal([]byte(*result.Content[0].Text), &response)
		if err != nil {
			xtptest.Assert("Failed to parse response JSON", false, err.Error())
		}

		xtptest.AssertEq("Unit should be Joules", response.Unit, "Joules")
		xtptest.Assert("Energy should be negative for opposite charges",
			response.Energy < 0, fmt.Sprintf("Got energy: %f", response.Energy))
	})
}

func testMagneticFlux() {
	xtptest.Group("test magnetic_flux tool", func() {
		pdk.Log(pdk.LogDebug, "Testing magnetic flux calculation")

		arguments := map[string]interface{}{
			"magneticField": map[string]interface{}{
				"x": 0.0,
				"y": 0.0,
				"z": 1.0, // 1 Tesla
			},
			"area":  1.0, // 1 square meter
			"angle": 0.0, // 0 radians (parallel)
		}

		input := CallToolRequest{
			Method: nil,
			Params: Params{
				Arguments: &arguments,
				Name:      "magnetic_flux",
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

		hasErrored := result.IsError != nil && *result.IsError
		xtptest.AssertEq("Tool call should not have errored", hasErrored, false)
		xtptest.AssertEq("Tool call should have one content item", len(result.Content), 1)
		xtptest.AssertEq("Content type should be text", result.Content[0].Type, ContentTypeText)

		var response struct {
			Flux float64 `json:"flux"`
			Unit string  `json:"unit"`
		}
		err = json.Unmarshal([]byte(*result.Content[0].Text), &response)
		if err != nil {
			xtptest.Assert("Failed to parse response JSON", false, err.Error())
		}

		xtptest.AssertEq("Unit should be Weber", response.Unit, "Weber")
		expectedFlux := 1.0 // B * A * cos(0) = 1 * 1 * 1 = 1
		xtptest.Assert("Flux should match expected value",
			math.Abs(response.Flux-expectedFlux) < 0.001,
			fmt.Sprintf("Got flux: %f, expected: %f", response.Flux, expectedFlux))
	})
}

func testCapacitorEnergy() {
	xtptest.Group("test capacitor_energy tool", func() {
		pdk.Log(pdk.LogDebug, "Testing capacitor energy calculation")

		arguments := map[string]interface{}{
			"capacitance": 1e-6, // 1 microfarad
			"voltage":     12.0, // 12 volts
		}

		input := CallToolRequest{
			Method: nil,
			Params: Params{
				Arguments: &arguments,
				Name:      "capacitor_energy",
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

		hasErrored := result.IsError != nil && *result.IsError
		xtptest.AssertEq("Tool call should not have errored", hasErrored, false)
		xtptest.AssertEq("Tool call should have one content item", len(result.Content), 1)
		xtptest.AssertEq("Content type should be text", result.Content[0].Type, ContentTypeText)

		var response struct {
			Energy float64 `json:"energy"`
			Unit   string  `json:"unit"`
		}
		err = json.Unmarshal([]byte(*result.Content[0].Text), &response)
		if err != nil {
			xtptest.Assert("Failed to parse response JSON", false, err.Error())
		}

		xtptest.AssertEq("Unit should be Joules", response.Unit, "Joules")
		expectedEnergy := 0.5 * 1e-6 * 12.0 * 12.0
		xtptest.Assert("Energy should match expected value",
			math.Abs(response.Energy-expectedEnergy) < 0.001,
			fmt.Sprintf("Got energy: %f, expected: %f", response.Energy, expectedEnergy))
	})
}

func testSolenoidInductance() {
	xtptest.Group("test solenoid_inductance tool", func() {
		pdk.Log(pdk.LogDebug, "Testing solenoid inductance calculation")

		arguments := map[string]interface{}{
			"turns":  1000,   // 1000 turns
			"length": 0.1,    // 10 cm
			"area":   0.0001, // 1 cm²
		}

		input := CallToolRequest{
			Method: nil,
			Params: Params{
				Arguments: &arguments,
				Name:      "solenoid_inductance",
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

		hasErrored := result.IsError != nil && *result.IsError
		xtptest.AssertEq("Tool call should not have errored", hasErrored, false)
		xtptest.AssertEq("Tool call should have one content item", len(result.Content), 1)
		xtptest.AssertEq("Content type should be text", result.Content[0].Type, ContentTypeText)
		var response struct {
			Inductance float64 `json:"inductance"`
			Unit       string  `json:"unit"`
		}

		err = json.Unmarshal([]byte(*result.Content[0].Text), &response)
		if err != nil {
			xtptest.Assert("Failed to parse response JSON", false, err.Error())
		}

		xtptest.AssertEq("Unit should be Henry", response.Unit, "Henry")

		expectedInductance := (4 * math.Pi * 1e-7 * 1000 * 1000 * 0.0001) / 0.1 // μ₀N²A/l
		xtptest.Assert("Inductance should match expected value",
			math.Abs(response.Inductance-expectedInductance) < 0.001,
			fmt.Sprintf("Got inductance: %f, expected: %f", response.Inductance, expectedInductance))
	})
}

func testRCTimeConstant() {
	xtptest.Group("test rc_time_constant tool", func() {
		pdk.Log(pdk.LogDebug, "Testing RC time constant calculation")

		arguments := map[string]interface{}{
			"resistance":  1000.0, // 1 kΩ
			"capacitance": 1e-6,   // 1 μF
		}

		input := CallToolRequest{
			Method: nil,
			Params: Params{
				Arguments: &arguments,
				Name:      "rc_time_constant",
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

		hasErrored := result.IsError != nil && *result.IsError
		xtptest.AssertEq("Tool call should not have errored", hasErrored, false)
		xtptest.AssertEq("Tool call should have one content item", len(result.Content), 1)
		xtptest.AssertEq("Content type should be text", result.Content[0].Type, ContentTypeText)

		var response struct {
			TimeConstant float64 `json:"timeConstant"`
			Unit         string  `json:"unit"`
		}

		err = json.Unmarshal([]byte(*result.Content[0].Text), &response)
		if err != nil {
			xtptest.Assert("Failed to parse response JSON", false, err.Error())
		}

		xtptest.AssertEq("Unit should be Seconds", response.Unit, "Seconds")

		expectedTimeConstant := 0.001 // RC = 1000 * 1e-6 = 0.001
		xtptest.Assert("Time constant should match expected value",
			math.Abs(response.TimeConstant-expectedTimeConstant) < 0.001,
			fmt.Sprintf("Got time constant: %f, expected: %f", response.TimeConstant, expectedTimeConstant))
	})
}

func main() {}
