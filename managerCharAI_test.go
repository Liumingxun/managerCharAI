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
	err = managerCharAI.WritePNG(outputFile, imageBase64, metadataBase64)
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

func TestReadString(t *testing.T) {
	// Test JSON string
	testJSON := `{
		"name": "String Test Character",
		"spec": "chara_card_v3",
		"spec_version": "3.0",
		"description": "A test character from JSON string",
		"tags": ["test", "string", "parsing"],
		"data": {
			"name": "String Test Character",
			"creator": "StringTester",
			"description": "Testing character parsing from JSON string",
			"first_mes": "Hello! I was parsed from a JSON string.",
			"personality": "String-based and efficient",
			"scenario": "String parsing test environment",
			"tags": ["test", "string", "golang"],
			"character_version": "1.0",
			"creator_notes": "This character was created by parsing a JSON string",
			"system_prompt": "You are a helpful test character from string"
		}
	}`

	card, err := managerCharAI.ReadString(testJSON)
	if err != nil {
		t.Fatalf("ReadString() failed: %v", err)
	}

	if card == nil {
		t.Fatal("ReadString() returned nil card")
	}

	// Validate parsed data
	if card.Name != "String Test Character" {
		t.Errorf("Name mismatch. Expected: String Test Character, Got: %s", card.Name)
	}

	if card.Spec != "chara_card_v3" {
		t.Errorf("Spec mismatch. Expected: chara_card_v3, Got: %s", card.Spec)
	}

	if card.Data.Creator != "StringTester" {
		t.Errorf("Creator mismatch. Expected: StringTester, Got: %s", card.Data.Creator)
	}

	if len(card.Tags) != 3 {
		t.Errorf("Tags count mismatch. Expected: 3, Got: %d", len(card.Tags))
	}

	t.Logf("Successfully parsed character from string: %s", card.Name)
	t.Logf("Creator: %s", card.Data.Creator)
	t.Logf("Tags: %v", card.Tags)
}

func TestReadJSON(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		wantErr bool
	}{
		{
			name:    "valid JSON file",
			file:    "test.json",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			card, err := managerCharAI.ReadJSON(tt.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if card == nil {
					t.Fatal("ReadJSON() returned nil card")
				}

				t.Logf("Character Name: %s", card.Name)
				t.Logf("Spec: %s v%s", card.Spec, card.SpecVersion)
				t.Logf("Description (first 100 chars): %s", card.Description[:min(100, len(card.Description))])

				if card.Name == "" && card.Data.Name == "" {
					t.Error("Both Name and Data.Name are empty")
				}
			}
		})
	}
}

func TestWriteJSON(t *testing.T) {
	// Create a test card
	card := &managerCharAI.CharacterCard{
		Name:        "Test Character",
		Description: "A test character for JSON writing",
		Spec:        "chara_card_v3",
		SpecVersion: "3.0",
		Tags:        []string{"test", "json", "example"},
		CreateDate:  "2025-10-10",
		Fav:         false,
		Data: managerCharAI.CharacterData{
			Name:             "Test Character",
			Creator:          "TestCreator",
			Description:      "Detailed test description for demonstration purposes",
			FirstMes:         "Hello from JSON! This is a test character card.",
			Personality:      "Friendly, helpful, and enthusiastic",
			Scenario:         "Test scenario in a development environment",
			Tags:             []string{"test", "json", "golang"},
			CharacterVersion: "1.0",
			CreatorNotes:     "This is a test character card created by the managerCharAI library",
			SystemPrompt:     "You are a helpful test character",
			Extensions: managerCharAI.Extensions{
				Talkativeness: "0.8",
				Fav:           false,
				World:         "Test World",
			},
		},
	}

	outputFile := "test_output.json"
	err := managerCharAI.WriteJSON(outputFile, card)
	if err != nil {
		t.Fatalf("WriteJSON() failed: %v", err)
	}
	// NO eliminar el archivo para poder inspeccionarlo
	// defer os.Remove(outputFile)

	t.Logf("Successfully created %s (file saved to disk for inspection)", outputFile)

	// Read back and verify
	readCard, err := managerCharAI.ReadJSON(outputFile)
	if err != nil {
		t.Fatalf("Failed to read back JSON: %v", err)
	}

	if readCard.Name != card.Name {
		t.Errorf("Name mismatch. Expected: %s, Got: %s", card.Name, readCard.Name)
	}

	if readCard.Spec != card.Spec {
		t.Errorf("Spec mismatch. Expected: %s, Got: %s", card.Spec, readCard.Spec)
	}

	if readCard.Data.Creator != card.Data.Creator {
		t.Errorf("Creator mismatch. Expected: %s, Got: %s", card.Data.Creator, readCard.Data.Creator)
	}

	t.Logf("Successfully verified JSON content")
	t.Logf("File location: %s", outputFile)
	t.Logf("Character: %s by %s", readCard.Name, readCard.Data.Creator)
}

func TestJSONToPNG(t *testing.T) {
	// Read JSON
	card, err := managerCharAI.ReadJSON("test.json")
	if err != nil {
		t.Fatalf("Failed to read JSON: %v", err)
	}

	t.Logf("Loaded character: %s", card.Name)

	// Read base image
	imageData, err := os.ReadFile("clean.png")
	if err != nil {
		t.Fatalf("Failed to read image: %v", err)
	}
	imageBase64 := base64.StdEncoding.EncodeToString(imageData)

	// Create PNG with character card
	outputFile := "test_json_to_png.png"
	err = managerCharAI.WritePNGFromCard(outputFile, imageBase64, card)
	if err != nil {
		t.Fatalf("WritePNGFromCard() failed: %v", err)
	}
	defer os.Remove(outputFile)

	t.Logf("Successfully created PNG from JSON: %s", outputFile)

	// Verify by reading back
	readCard, err := managerCharAI.ReadPNGAsCard(outputFile)
	if err != nil {
		t.Fatalf("Failed to read back PNG: %v", err)
	}

	if readCard.Name != card.Name {
		t.Errorf("Name mismatch after JSON->PNG conversion")
	}

	t.Logf("Successfully verified JSON to PNG conversion")
}

func TestCharacterCard_SaveJSON(t *testing.T) {
	// Create a character card
	card := &managerCharAI.CharacterCard{
		Name:        "Method Test Character",
		Description: "Testing SaveJSON method",
		Spec:        "chara_card_v3",
		SpecVersion: "3.0",
		Tags:        []string{"test", "method"},
		Data: managerCharAI.CharacterData{
			Name:        "Method Test Character",
			Creator:     "MethodTester",
			Description: "Testing the SaveJSON method",
			FirstMes:    "Hello! I was saved using the SaveJSON method.",
			Personality: "Methodical and organized",
			Scenario:    "Testing environment",
			Tags:        []string{"test", "method"},
		},
	}

	outputFile := "test_method_output.json"
	err := card.SaveJSON(outputFile)
	if err != nil {
		t.Fatalf("SaveJSON() method failed: %v", err)
	}
	// Keep file for inspection
	// defer os.Remove(outputFile)

	t.Logf("Successfully saved using method: %s", outputFile)

	// Verify by reading back
	readCard, err := managerCharAI.ReadJSON(outputFile)
	if err != nil {
		t.Fatalf("Failed to read back: %v", err)
	}

	if readCard.Name != card.Name {
		t.Errorf("Name mismatch. Expected: %s, Got: %s", card.Name, readCard.Name)
	}

	t.Logf("Method SaveJSON() works correctly!")
}

func TestCharacterCard_SavePNG(t *testing.T) {
	// Create a character card
	card := &managerCharAI.CharacterCard{
		Name:        "PNG Method Test",
		Description: "Testing SavePNG method",
		Spec:        "chara_card_v3",
		SpecVersion: "3.0",
		Tags:        []string{"test", "png", "method"},
		Data: managerCharAI.CharacterData{
			Name:        "PNG Method Test",
			Creator:     "PNGMethodTester",
			Description: "Testing the SavePNG method",
			FirstMes:    "Hello! I was saved using the SavePNG method.",
			Personality: "Visual and embedded",
			Scenario:    "PNG testing environment",
			Tags:        []string{"test", "png", "method"},
		},
	}

	// Read base image
	imageData, err := os.ReadFile("clean.png")
	if err != nil {
		t.Fatalf("Failed to read base image: %v", err)
	}
	imageBase64 := base64.StdEncoding.EncodeToString(imageData)

	outputFile := "test_method_output.png"
	err = card.SavePNG(outputFile, imageBase64)
	if err != nil {
		t.Fatalf("SavePNG() method failed: %v", err)
	}
	// Keep file for inspection
	// defer os.Remove(outputFile)

	t.Logf("Successfully saved using method: %s", outputFile)

	// Verify by reading back
	readCard, err := managerCharAI.ReadPNGAsCard(outputFile)
	if err != nil {
		t.Fatalf("Failed to read back: %v", err)
	}

	if readCard.Name != card.Name {
		t.Errorf("Name mismatch. Expected: %s, Got: %s", card.Name, readCard.Name)
	}

	if readCard.Data.Creator != card.Data.Creator {
		t.Errorf("Creator mismatch. Expected: %s, Got: %s", card.Data.Creator, readCard.Data.Creator)
	}

	t.Logf("Method SavePNG() works correctly!")
}
