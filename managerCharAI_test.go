package managerCharAI_test

import (
	"encoding/base64"
	"encoding/json"
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
			t.Logf("First 100 chars of base64: %s", base64Data[:min(100, len(base64Data))])

			// Decode base64 to JSON
			jsonData, err := base64.StdEncoding.DecodeString(base64Data)
			if err != nil {
				t.Errorf("Failed to decode base64: %v", err)
				return
			}

			t.Logf("Decoded JSON length: %d bytes", len(jsonData))

			// Pretty print JSON
			var jsonObj interface{}
			if err := json.Unmarshal(jsonData, &jsonObj); err != nil {
				t.Errorf("Failed to parse JSON: %v", err)
				t.Logf("Raw JSON (first 500 chars): %s", string(jsonData[:min(500, len(jsonData))]))
				return
			}

			prettyJSON, err := json.MarshalIndent(jsonObj, "", "  ")
			if err != nil {
				t.Errorf("Failed to pretty print JSON: %v", err)
				return
			}

			t.Logf("Parsed JSON:\n%s", string(prettyJSON))
		})
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestReadPNGAsCard(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		wantErr bool
	}{
		{
			name:    "valid character card struct",
			file:    "test.png",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			card, err := managerCharAI.ReadPNGAsCard(tt.file)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("ReadPNGAsCard() failed: %v", err)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ReadPNGAsCard() succeeded unexpectedly")
			}
			if card == nil {
				t.Fatal("ReadPNGAsCard() returned nil card")
			}

			// Validate card structure
			t.Logf("Character Name: %s", card.Name)
			t.Logf("Spec: %s v%s", card.Spec, card.SpecVersion)
			t.Logf("Description (first 100 chars): %s", card.Description[:min(100, len(card.Description))])
			t.Logf("Tags: %v", card.Tags)
			t.Logf("Data.Name: %s", card.Data.Name)
			t.Logf("Data.Creator: %s", card.Data.Creator)
			t.Logf("Data.Tags: %v", card.Data.Tags)
			t.Logf("Alternate Greetings: %d", len(card.Data.AlternateGreetings))

			// Basic validation
			if card.Name == "" && card.Data.Name == "" {
				t.Error("Both Name and Data.Name are empty")
			}
			if card.Spec == "" {
				t.Error("Spec is empty")
			}
			if card.SpecVersion == "" {
				t.Error("SpecVersion is empty")
			}
		})
	}
}
