package main

import (
	"fmt"
<<<<<<< HEAD
=======
	"io"
>>>>>>> 2800bf52dbc8f99c4467731458a1623e5bc4462c
	"net/http"
	"os"
	"testing"
)

<<<<<<< HEAD
func ErrChecker(t *testing.T, ErrMsg string, err error) {
	if err != nil {
		t.Fatalf("%s: %v", ErrMsg, err)
	}
}

func Equals(t *testing.T, myanswer, expected string) {
	if myanswer != expected {
		t.Errorf("Expected %s but get %s", myanswer, expected)
	}
}
func NotEquals(t *testing.T, myanswer, expected string) {
	if myanswer == expected {
		t.Errorf("Expected %s but get %s", myanswer, expected)
	}
}

=======
>>>>>>> 2800bf52dbc8f99c4467731458a1623e5bc4462c
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
<<<<<<< HEAD
			ErrChecker(t, "Error Getting video Info", err)
			NotEquals(t, returnedBody, tc.expected)
=======
			if err != nil {
				t.Errorf("error during request: %v", err)
			}
			if returnedBody != tc.expected {
				t.Errorf("expected %v but get %v", tc.expected, returnedBody)
			}
>>>>>>> 2800bf52dbc8f99c4467731458a1623e5bc4462c
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

<<<<<<< HEAD
			_, err = getAudioData(resp)
			ErrChecker(t, "Error getting audio data", err)

			// normally should not be problem streaming data to file
			// file, _ := os.Create("./testAudioDataFile")
			// io.Copy(file, responseAudio.Body)
=======
			responseAudio, err := getAudioData(resp)
			if err != nil {
				t.Errorf("error getting Audio Data: %v", err)
			}

			// normally should not be problem streaming data to file
			file, _ := os.Create("./testAudioDataFile")
			io.Copy(file, responseAudio.Body)
>>>>>>> 2800bf52dbc8f99c4467731458a1623e5bc4462c
		})
	}
}
