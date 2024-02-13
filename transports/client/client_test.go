package client

import (
	"testing"
)

func TestGetAbsoluteValue(t *testing.T) {
	// Test cases
	testCases := []struct {
		input    int64
		expected uint64
	}{
		{10, 10},
		{-10, 10},
		{0, 0},
	}

	for _, testCase := range testCases {
		result := getAbsoluteValue(testCase.input)
		if result != testCase.expected {
			t.Errorf("Input: %d, Expected: %d, Got: %d", testCase.input, testCase.expected, result)
		}
	}
}
