package graphapi

import (
	"encoding/json"
	"strings"
)

// Slot represents a connection point within a GraphNode.
// It holds various properties that define the behavior and appearance
// of the connection, such as the name, type, associated widget, and more.
type Slot struct {
	Name       string     `json:"name"` // The name of the slot
	CustomType int        `json:"-"`
	Node       *GraphNode `json:"-"`                // The node the slot belongs to
	Type       string     `json:"-"`                // The type of the data the slot accepts (handled by custom unmarshaler)
	Link       int        `json:"link,omitempty"`   // Index of the link for an input slot
	Links      *[]int     `json:"links,omitempty"`  // Array of links for output slots
	Widget     *Widget    `json:"widget,omitempty"` // Collection of widgets that allow setting properties
	Shape      *int       `json:"shape,omitempty"`
	SlotIndex  *int       `json:"slot_index,omitempty"` // Index of the Slot in relation to other Slots
	Property   Property   `json:"-"`                    // non-null for inputs that are exported widgets
}

// UnmarshalJSON implements custom JSON unmarshaling for Slot
func (s *Slot) UnmarshalJSON(data []byte) error {
	// Define a temporary struct that matches the JSON structure
	type TempSlot struct {
		Name      string      `json:"name"`
		Type      interface{} `json:"type"` // Can be string or array
		Link      int         `json:"link,omitempty"`
		Links     *[]int      `json:"links,omitempty"`
		Widget    *Widget     `json:"widget,omitempty"`
		Shape     *int        `json:"shape,omitempty"`
		SlotIndex *int        `json:"slot_index,omitempty"`
	}

	var temp TempSlot
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Copy all fields except Type
	s.Name = temp.Name
	s.Link = temp.Link
	s.Links = temp.Links
	s.Widget = temp.Widget
	s.Shape = temp.Shape
	s.SlotIndex = temp.SlotIndex

	// Handle Type field - convert array to pipe-separated string if needed
	switch v := temp.Type.(type) {
	case string:
		s.Type = v
	case []interface{}:
		// Convert array to pipe-separated string
		var types []string
		for _, item := range v {
			if str, ok := item.(string); ok {
				types = append(types, str)
			}
		}
		s.Type = strings.Join(types, "|")
	case nil:
		s.Type = ""
	default:
		// Fallback: convert to string
		if typeBytes, err := json.Marshal(v); err == nil {
			s.Type = string(typeBytes)
		}
	}

	return nil
}
