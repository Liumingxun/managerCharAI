package managerCharAI_test

import (
	"testing"

	"github.com/jonathanhecl/managerCharAI"
)

func TestReadPNG(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		wantErr bool
	}{
		{
			name:    "valid character card",
			file:    "test.png",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base64Data, err := managerCharAI.ReadPNG(tt.file)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("ReadPNG() failed: %v", err)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ReadPNG() succeeded unexpectedly")
			}
			if base64Data == "" {
				t.Fatal("ReadPNG() returned empty base64 data")
			}

			// Basic validation - should start with base64 characters
			if len(base64Data) < 10 {
				t.Error("Base64 data is too short")
			}

			t.Logf("Successfully extracted base64 data, length: %d bytes", len(base64Data))
			t.Logf("First 50 chars: %s", base64Data[:min(50, len(base64Data))])
		})
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
