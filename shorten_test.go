package main

import (
	"fmt"
	"testing"
)

func TestShorten(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		wantLen  int
		wantErr  bool
	}{
		{
			name:    "Regular URL",
			input:   "https://www.example.com",
			wantLen: 8,
			wantErr: false,
		},
		{
			name:    "Empty string",
			input:   "",
			wantLen: 8,
			wantErr: false,
		},
		{
			name:    "Long URL",
			input:   "https://www.example.com/very/long/path/with/many/segments",
			wantLen: 8,
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := shorten(tc.input)
			
			if (err != nil) != tc.wantErr {
				t.Errorf("shorten() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			
			if len(got) != tc.wantLen {
				fmt.Println(got)
				t.Errorf("shorten() got length = %v, want length %v", len(got), tc.wantLen)
			}
			
			// Test idempotency - same input should give same output
			second, _ := shorten(tc.input)
			if got != second {
				t.Errorf("shorten() not idempotent: first = %v, second = %v", got, second)
			}
		})
	}
}