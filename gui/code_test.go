package gui

import "testing"

func TestFormatDisasm(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		indentation int
		expected    string
	}{
		{name: "data",
			input:       "00:0000: OP_DATA_3 deadbeef",
			indentation: 0,
			expected:    " 0000   deadbeef"},

		{name: "op",
			input:       "00:0000: OP_1",
			indentation: 0,
			expected:    " 0000   OP_1"},

		{name: "with indentation",
			input:       "00:0000: OP_1",
			indentation: 4,
			expected:    " 0000       OP_1"},
	}

	for _, test := range tests {
		actual := formatDisasm(test.input, &test.indentation, 2)
		if actual != test.expected {
			t.Errorf("expected: %s, actual: %s", test.expected, actual)
		}
	}
}
