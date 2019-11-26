package gui

import "testing"

func TestFormatDisasm(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"data",
			"00:0000: OP_DATA_3 deadbeef",
			" 0000   deadbeef"},
		{"op",
			"00:0000: OP_1",
			" 0000   OP_1"},
	}

	for _, test := range tests {
		actual := formatDisasm(test.input)
		if actual != test.expected {
			t.Errorf("expected: %s, actual: %s", test.expected, actual)
		}
	}
}
