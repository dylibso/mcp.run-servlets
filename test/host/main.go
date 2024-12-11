package main

import (
	"encoding/json"
	"fmt"

	pdk "github.com/extism/go-pdk"
)

//go:wasmimport xtp:test/harness mock_input
func mock_input() uint64

//go:export config_get
func config_get(key uint64) uint64 {
	input := getMockInput()

	keyMem := pdk.FindMemory(key)
	k := string(keyMem.ReadBytes())

	emptyString, err := pdk.AllocateJSON("")
	if err != nil {
		panic(err)
	}

	v, ok := input.Config[k]
	if !ok {
		pdk.Log(pdk.LogDebug, fmt.Sprintf("config_get: key %s not found", k))
		return emptyString.Offset()
	}

	pdk.Log(pdk.LogDebug, fmt.Sprintf("config_get: key %s found => %s", k, v))

	output, err := json.Marshal(v)
	if err != nil {
		pdk.Log(pdk.LogError, fmt.Sprintf("config_get: failed to marshal value %s", v))
		return emptyString.Offset()
	}

	valMem := pdk.AllocateBytes(output)
	return valMem.Offset()
}

func main() {}

func getMockInput() *MockInput {
	offs := mock_input()

	if offs == 0 {
		pdk.Log(pdk.LogDebug, "getMockInput: no mock input found, using default")
		return &MockInput{
			Config: map[string]string{},
		}
	}

	inputMem := pdk.FindMemory(offs)
	var input MockInput
	err := json.Unmarshal(inputMem.ReadBytes(), &input)
	if err != nil {
		panic(err)
	}

	return &input
}

type MockInput struct {
	Config map[string]string
}
