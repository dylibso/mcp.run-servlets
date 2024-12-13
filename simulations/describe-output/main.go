package main

import (
	"encoding/json"
	"fmt"
	"strings"

	xtptest "github.com/dylibso/xtp-test-go"
)

type ListToolsResult struct {
	Tools []ToolDescription `json:"tools"`
}

type ToolDescription struct {
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
	Name        string      `json:"name"`
}

func validateProperty(toolName, propName string, prop interface{}) error {
	propMap, ok := prop.(map[string]interface{})
	if !ok {
		return fmt.Errorf("tool %s: property %s must be an object", toolName, propName)
	}

	// Check for required property fields
	requiredFields := []string{"type", "description"}
	for _, field := range requiredFields {
		val, exists := propMap[field]
		if !exists {
			return fmt.Errorf("tool %s: property %s missing required field %s", toolName, propName, field)
		}

		strVal, ok := val.(string)
		if !ok || strVal == "" {
			return fmt.Errorf("tool %s: property %s.%s must be a non-empty string", toolName, propName, field)
		}
	}

	return nil
}

func validateToolDescription(tool ToolDescription) error {
	// Validate Name
	if tool.Name == "" {
		return fmt.Errorf("tool name is required")
	}
	if strings.TrimSpace(tool.Name) != tool.Name {
		return fmt.Errorf("tool name cannot have leading or trailing whitespace: %q", tool.Name)
	}

	// Validate Description
	if tool.Description == "" {
		return fmt.Errorf("tool %s: description is required", tool.Name)
	}
	if strings.TrimSpace(tool.Description) != tool.Description {
		return fmt.Errorf("tool %s: description cannot have leading or trailing whitespace", tool.Name)
	}

	// Validate InputSchema structure
	schema, ok := tool.InputSchema.(map[string]interface{})
	if !ok {
		return fmt.Errorf("tool %s: inputSchema must be an object", tool.Name)
	}

	// Validate schema type exists
	typeVal, hasType := schema["type"]
	if !hasType {
		return fmt.Errorf("tool %s: missing inputSchema.type", tool.Name)
	}
	typeStr, ok := typeVal.(string)
	if !ok || typeStr == "" {
		return fmt.Errorf("tool %s: inputSchema.type must be a non-empty string", tool.Name)
	}

	// If type is object, validate properties
	if typeStr == "object" {
		// Validate properties
		props, hasProps := schema["properties"]
		if !hasProps {
			return fmt.Errorf("tool %s: object schema missing properties", tool.Name)
		}
		properties, ok := props.(map[string]interface{})
		if !ok {
			return fmt.Errorf("tool %s: inputSchema.properties must be an object", tool.Name)
		}
		if len(properties) == 0 {
			return fmt.Errorf("tool %s: inputSchema.properties cannot be empty", tool.Name)
		}

		// Validate each property
		for propName, prop := range properties {
			if err := validateProperty(tool.Name, propName, prop); err != nil {
				return err
			}
		}

		// Validate required fields if present
		if reqVal, hasRequired := schema["required"]; hasRequired {
			required, ok := reqVal.([]interface{})
			if !ok {
				return fmt.Errorf("tool %s: inputSchema.required must be an array", tool.Name)
			}

			var missingFields []string
			seenFields := make(map[string]bool)

			for _, field := range required {
				fieldStr, ok := field.(string)
				if !ok {
					return fmt.Errorf("tool %s: required field names must be strings", tool.Name)
				}

				// Check for duplicate required fields
				if seenFields[fieldStr] {
					return fmt.Errorf("tool %s: duplicate required field: %s", tool.Name, fieldStr)
				}
				seenFields[fieldStr] = true

				// Check if field exists in properties
				if _, exists := properties[fieldStr]; !exists {
					missingFields = append(missingFields, fieldStr)
				}
			}

			if len(missingFields) > 0 {
				return fmt.Errorf("tool %s: required fields missing from properties: %s",
					tool.Name,
					strings.Join(missingFields, ", "))
			}
		}
	}

	return nil
}

//go:export test
func test() int32 {
	output := xtptest.CallBytes("describe", nil)
	xtptest.AssertNe("describe returned output", string(output), "")

	var result ListToolsResult
	if err := json.Unmarshal(output, &result); err != nil {
		xtptest.Assert("describe returns valid JSON", false, fmt.Sprintf("invalid JSON: %v", err))
		return 0
	}

	if len(result.Tools) == 0 {
		xtptest.Assert("describe provides at least one tool", false, "describe must provide at least one tool")
		return 0
	}

	xtptest.Group("validate tool descriptions", func() {
		for i, tool := range result.Tools {
			testName := fmt.Sprintf("tool[%d](%s) is valid", i, tool.Name)
			if err := validateToolDescription(tool); err != nil {
				xtptest.Assert(testName, false, err.Error())
				return
			}
			xtptest.Assert(testName, true, "tool validation passed")
		}
	})

	return 0
}

func main() {}
