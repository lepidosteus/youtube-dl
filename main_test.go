package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
)

func TestGetVideoInfo(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cmd", "youtube-dl", "https://www.youtube.com/watch?v=iwyXbD1Rn7g"}
	tt := []struct {
		name     string
		value    string
		expected string
	}{
		{name: "test output value", value: "iwyXbD1Rn7g", expected: ""},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			returnedBody, err := getVideoInfo(tc.value)
			if err != nil {
				t.Errorf("error during request: %v", err)
			}
			if returnedBody != tc.expected {
				t.Errorf("expected %v but get %v", tc.expected, returnedBody)
			}
		})
	}
}

func assertType(t *testing.T, a, b interface{}) {
	_, ok := a.(*http.Response)
	if !ok {
		t.Errorf("expected type (%v) but get (%v)", a, b)
	}
}

func TestGetAudioData(t *testing.T) {
	// oldArgs := os.Args
	// defer func() { os.Args = oldArgs }()
	// os.Args = []string{"cmd", "youtube-dl", "https://www.youtube.com/watch?v=WUhWZbKvLsQ"}
	tt := []struct {
		name     string
		id       string
		expected interface{}
	}{
		{name: "Test Audio Data", id: "iwyXbD1Rn7g", expected: ""},
	}

	for _, tc := range tt {
		fmt.Println("Running TestGetAudioData")
		t.Run(tc.name, func(t *testing.T) {
			resp, err := getVideoInfo(tc.id)
			if err != nil {
				t.Errorf("error during request: %v", err)
			}

			responseAudio, err := getAudioData(resp)
			if err != nil {
				t.Errorf("error getting Audio Data: %v", err)
			}

			// normally should not be problem streaming data to file
			file, _ := os.Create("./testAudioDataFile")
			io.Copy(file, responseAudio.Body)
		})
	}
}
