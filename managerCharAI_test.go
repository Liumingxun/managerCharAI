package managerCharAI_test

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"os"
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

func TestWritePNG(t *testing.T) {
	// Read an existing PNG and remove existing metadata to get a clean base
	imageData, err := os.ReadFile("clean.png")
	if err != nil {
		t.Fatalf("Failed to read clean.png: %v", err)
	}

	// Remove existing tEXt chunks to get clean PNG
	cleanImageData := removeTextChunks(imageData)
	imageBase64 := base64.StdEncoding.EncodeToString(cleanImageData)

	// Create simple test metadata
	testMetadata := `{"name":"Test Character","spec":"chara_card_v3","spec_version":"3.0","description":"Test description","tags":["test1","test2","test3"]}`
	metadataBase64 := base64.StdEncoding.EncodeToString([]byte(testMetadata))

	// Write new PNG
	outputFile := "test_output.png"
	err = managerCharAI.WritePNG(imageBase64, metadataBase64, outputFile)
	if err != nil {
		t.Fatalf("WritePNG() failed: %v", err)
	}
	// No eliminar el archivo para que puedas verlo en disco
	// defer os.Remove(outputFile)

	t.Logf("Successfully created %s (file saved to disk)", outputFile)

	// Read back and verify
	readBase64, err := managerCharAI.ReadPNG(outputFile)
	if err != nil {
		t.Fatalf("Failed to read back the created PNG: %v", err)
	}

	if readBase64 != metadataBase64 {
		t.Errorf("Metadata mismatch.\nExpected: %s\nGot: %s", metadataBase64, readBase64)
		// Decode both to see the difference
		expectedJSON, _ := base64.StdEncoding.DecodeString(metadataBase64)
		gotJSON, _ := base64.StdEncoding.DecodeString(readBase64)
		t.Logf("Expected JSON: %s", string(expectedJSON))
		t.Logf("Got JSON: %s", string(gotJSON))
	} else {
		t.Logf("Metadata verification successful!")
	}

	// Decode and verify JSON
	decodedJSON, err := base64.StdEncoding.DecodeString(readBase64)
	if err != nil {
		t.Fatalf("Failed to decode base64 from created PNG: %v", err)
	}

	t.Logf("Verified metadata content: %s", string(decodedJSON))
}

// removeTextChunks removes all tEXt chunks from a PNG to get a clean base image
func removeTextChunks(data []byte) []byte {
	if len(data) < 8 {
		return data
	}

	var result []byte
	result = append(result, data[:8]...) // Keep PNG signature

	pos := 8
	for pos < len(data)-8 {
		if pos+4 > len(data) {
			break
		}

		chunkLen := binary.BigEndian.Uint32(data[pos : pos+4])
		if pos+8+int(chunkLen) > len(data) {
			break
		}

		chunkType := string(data[pos+4 : pos+8])

		// Copy all chunks except tEXt
		if chunkType != "tEXt" {
			chunkSize := 4 + 4 + int(chunkLen) + 4 // length + type + data + CRC
			if pos+chunkSize <= len(data) {
				result = append(result, data[pos:pos+chunkSize]...)
			}
		}

		pos += 4 + 4 + int(chunkLen) + 4
	}

	return result
}
