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

	// If type is object, validate properties and required fields
	if typeStr == "object" {
		properties, err := parseProperties(schema)
		if err != nil {
			return fmt.Errorf("tool %s: %w", tool.Name, err)
		}

		required, err := parseRequiredFields(schema)
		if err != nil {
			return fmt.Errorf("tool %s: %w", tool.Name, err)
		}

		// If properties is nil, required should also be empty
		if properties == nil && len(required) > 0 {
			return fmt.Errorf("tool %s: cannot have required fields without properties", tool.Name)
		}

		if properties != nil {
			// Validate each property
			for propName, prop := range properties {
				if err := validateProperty(tool.Name, propName, prop); err != nil {
					return err
				}
			}
			var missingFields []string
			for _, field := range required {
				if _, exists := properties[field]; !exists {
					missingFields = append(missingFields, field)
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

// parseProperties returns the properties map if it exists and is valid
func parseProperties(schema map[string]interface{}) (map[string]interface{}, error) {
	props, exists := schema["properties"]
	if !exists || props == nil {
		return nil, nil
	}

	properties, ok := props.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("properties must be an object")
	}

	return properties, nil
}

// parseRequiredFields returns the required fields if they exist and are valid
func parseRequiredFields(schema map[string]interface{}) ([]string, error) {
	req, exists := schema["required"]
	if !exists || req == nil {
		return []string{}, nil
	}

	reqArray, ok := req.([]interface{})
	if !ok {
		return []string{}, fmt.Errorf("required fields must be an array")
	}

	required := make([]string, 0, len(reqArray))
	seenFields := make(map[string]bool)

	for _, field := range reqArray {
		fieldStr, ok := field.(string)
		if !ok {
			return []string{}, fmt.Errorf("required field names must be strings")
		}

		// Check for duplicate required fields
		if seenFields[fieldStr] {
			return []string{}, fmt.Errorf("duplicate required field: %s", fieldStr)
		}
		seenFields[fieldStr] = true

		required = append(required, fieldStr)
	}

	return required, nil
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
