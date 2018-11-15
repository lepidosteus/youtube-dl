package main

import (
	"os"
	"testing"
)

func TestGetVideoInfo(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cmd", "youtube-dl", "https://www.youtube.com/watch?v=WUhWZbKvLsQ"}
	tt := []struct {
		name     string
		value    string
		expected string
	}{
		{name: "test output value", value: "y3LB4zXafpY", expected: ""},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// returnedBody, err := getVideoInfo(tc.value)

			if err != nil {
				t.Errorf("error during request: %v", err)
			}
			if returnedBody != tc.expected {
				t.Errorf("expected %v but get %v", tc.expected, returnedBody)
			}
		})
	}
}
