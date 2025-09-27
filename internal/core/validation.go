package core

import (
	"strings"
)

// ValidationError aggregates per-field validation failures so callers can
// present inline messages in the UI. It satisfies the error interface.
type ValidationError struct {
	Fields map[string][]string `json:"fields"`
}

// Add records a message for the given field key.
func (v *ValidationError) Add(field, message string) {
	if v.Fields == nil {
		v.Fields = make(map[string][]string)
	}
	v.Fields[field] = append(v.Fields[field], message)
}

// Merge combines another ValidationError into this one.
func (v *ValidationError) Merge(other *ValidationError) {
	if other == nil || len(other.Fields) == 0 {
		return
	}
	if v.Fields == nil {
		v.Fields = make(map[string][]string)
	}
	for k, messages := range other.Fields {
		v.Fields[k] = append(v.Fields[k], messages...)
	}
}

// Error implements the error interface.
func (v *ValidationError) Error() string {
	if v == nil || len(v.Fields) == 0 {
		return ""
	}
	var parts []string
	for field, messages := range v.Fields {
		parts = append(parts, field+": "+strings.Join(messages, ", "))
	}
	return strings.Join(parts, "; ")
}

// Empty reports whether there are no validation errors recorded.
func (v *ValidationError) Empty() bool {
	return v == nil || len(v.Fields) == 0
}
