package schemaelement

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSchemaPositiveCases(t *testing.T) {
	tests := []struct {
		inputValues []string
		expected    string
		msg         string
	}{
		{
			inputValues: []string{"a", "b", "c", "d"},
			expected:    "string",
			msg:         "ElmType was not a string with all strings.",
		},
		{
			inputValues: []string{"a", "1", "c", "d"},
			expected:    "string",
			msg:         "ElmType was not a string with most strings.",
		},
		{
			inputValues: []string{"0", "1", "2", "3"},
			expected:    "int",
			msg:         "ElmType was not a int with all int values.",
		},
		{
			inputValues: []string{"0.5", "1.2", "2.9", "3.88"},
			expected:    "float",
			msg:         "ElmType was not a float with all float values.",
		},
		{
			inputValues: []string{"true", "false", "true", "true"},
			expected:    "bool",
			msg:         "ElmType was not a bool with all bool values.",
		},
	}

	for _, test := range tests {
		se, err := NewSchemaElement("test")

		require.NoError(t, err)
		assert.Equal(t, "test", se.Name)

		for _, iv := range test.inputValues {
			se.UpdateType(iv)
		}

		assert.Equal(t, test.expected, se.ElmType)
	}

}
