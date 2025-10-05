package main

import (
	"encoding/json"
	"log"
)

// Task 1: Struct to JSON (Marshaling)

// Define a Go struct and convert an instance of it into a JSON byte array.
type ToolDefinition struct {
	Name           string   `json:"name"`
	Description    string   `json:"description,omitempty"`
	RequiredInputs []string `json:"required_inputs"`
}

func main() {
	// Create an instance of ToolDefinition
	tool := ToolDefinition{
		Name:           "calculate_area",
		Description:    "Calculated the area of a Circle",
		RequiredInputs: []string{"radius"},
	}
	// Task 1: Marshal the 'tool' struct into JSON bytes (and handle the error!)
	// Hint: Use json.Marshal()
	tool_json, error := json.Marshal(tool)
	if error != nil {
		log.Fatalf("Couldn't Marshal %v", error)
	}
	// fmt.Printf("Tools in JSON %s\n", tool_json)
	// Task 2: Unmarshal the JSON bytes back into a new struct instance (and handle the error!)
	// Hint: Use json.Unmarshal()
	var new_tool ToolDefinition
	error = json.Unmarshal(tool_json, &new_tool)
	if error != nil {
		log.Fatalf("Couldn't Unmarshal, %v", error)
	}

	// fmt.Printf("New Tool: %T\n", new_tool)
	// fmt.Printf("New Tool Name: %s\n", new_tool.Name)
	// fmt.Printf("New Tool Description: %s\n", new_tool.Description)
	// fmt.Printf("New Tool Required Inputs: %v\n", new_tool.RequiredInputs)
	// ... continue with Task 3 and 4

	// Task 3: Using Struct Tags

	// Modify the struct to use json:"..." struct tags to enforce specific JSON key names (e.g., using snake_case in JSON for a CamelCase field in Go).
	// Task 4: Handling Dynamic/Untyped Data #############DONE

	// Write code to parse a JSON object where you don't know the exact structure (i.e., treating it as a map[string]interface{}).
	// Yeah , so here I knew the structure was ToolDefinition. But what if I don't know that. Good Point.

	// Let's take this for example
	// When the JSON structure is unknown or dynamic, we unmarshal into a map.
	rawJSON := []byte(`{"method":"execute", "id": 123, "params": {"tool_name": "add", "arguments": [5, 2]}}`)
	// WHEN THE JSON STRUCUTRE IS UNKNOWN WE UNMARSHALL IT INTO A MAP
	// THE TARGET IS A MAP MAPPING STRINGS TO AN EMPTY INTERFACE (interface{})
	var genericData map[string]interface{}
	error = json.Unmarshal(rawJSON, &genericData)
	if error != nil {
		log.Fatalf("Coundn't umarshall %v of type %t", genericData, genericData)
	}
	// fmt.Printf("JSON rawJSON %v", genericData)
	// fmt.Printf("%t", genericData["params"])
	// fmt.Printf("%t", genericData)
	// Now we need to explicitly convert the values to their expected types

	// This is not a good way to do it. Go will let it pass. But if the types are not as expected, it will panic at runtime.
	// So we need to be very careful while doing this.
	// Also note that JSON numbers are float64 by default in Go.
	// params := genericData["params"].(map[string]interface{})
	// tool_name := params["tool_name"].(string)
	// arguments := params["arguments"].([]interface{})
	// arg1 := arguments[0].(float64) // JSON numbers are float64 by default
	// arg2 := arguments[1].(float64)
	// fmt.Printf("Tool Name: %s, Arguments: %v, %v\n", tool_name, arg1, arg2)
	params, err := genericData["params"].(map[string]interface{})
	if err != nil {
		log.Fatal("Couldn't convert params to map[string]interface{}")
	}
}
