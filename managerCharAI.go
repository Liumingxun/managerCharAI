package managerCharAI

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Lorebook struct {
	Name              string                 `json:"name"`
	Description       string                 `json:"description"`
	ScanDepth         int                    `json:"scan_depth"`
	TokenBudget       int                    `json:"token_budget"`
	RecursiveScanning bool                   `json:"recursive_scanning"`
	Extensions        map[string]interface{} `json:"extensions"`
	Entries           []struct {
		Keys           []string               `json:"keys"`
		Content        string                 `json:"content"`
		Extensions     map[string]interface{} `json:"extensions"`
		Enabled        bool                   `json:"enabled"`
		InsertionOrder int                    `json:"insertion_order"`
		CaseSensitive  *bool                  `json:"case_sensitive"`
		UseRegex       bool                   `json:"use_regex"`
		Constant       *bool                  `json:"constant"`
		Name           *string                `json:"name"`
		Priority       *int                   `json:"priority"`
		Id             *string                `json:"id"`
		Comment        *string                `json:"comment"`
		Selective      *bool                  `json:"selective"`
		SecondaryKeys  *[]string              `json:"secondary_keys"`
		Position       *string                `json:"position"`
	} `json:"entries"`
}

type CharacterCardV2 struct {
	Name                    string                 `json:"name"`
	Description             string                 `json:"description"`
	Tags                    []string               `json:"tags"`
	Creator                 string                 `json:"creator"`
	CharacterVersion        string                 `json:"character_version"`
	MesExample              string                 `json:"mes_example"`
	Extensions              map[string]interface{} `json:"extensions"`
	SystemPrompt            string                 `json:"system_prompt"`
	PostHistoryInstructions string                 `json:"post_history_instructions"`
	FirstMes                string                 `json:"first_mes"`
	AlternateGreetings      []string               `json:"alternate_greetings"`
	Personality             string                 `json:"personality"`
	Scenario                string                 `json:"scenario"`
	CreatorNotes            string                 `json:"creator_notes"`
	CharacterBook           *Lorebook              `json:"character_book"`
	Assets                  []struct {
		Type string `json:"type"`
		Uri  string `json:"uri"`
		Name string `json:"name"`
		Ext  string `json:"ext"`
	} `json:"assets"`
}

type CharacterCardV3 struct {
	Name                    string                 `json:"name"`
	Description             string                 `json:"description"`
	Tags                    []string               `json:"tags"`
	Creator                 string                 `json:"creator"`
	CharacterVersion        string                 `json:"character_version"`
	MesExample              string                 `json:"mes_example"`
	Extensions              map[string]interface{} `json:"extensions"`
	SystemPrompt            string                 `json:"system_prompt"`
	PostHistoryInstructions string                 `json:"post_history_instructions"`
	FirstMes                string                 `json:"first_mes"`
	AlternateGreetings      []string               `json:"alternate_greetings"`
	Personality             string                 `json:"personality"`
	Scenario                string                 `json:"scenario"`
	CreatorNotes            string                 `json:"creator_notes"`
	CharacterBook           *Lorebook              `json:"character_book"`
	Assets                  []struct {
		Type string `json:"type"`
		Uri  string `json:"uri"`
		Name string `json:"name"`
		Ext  string `json:"ext"`
	} `json:"assets"`
	Nickname                 string            `json:"nickname"`
	CreatorNotesMultilingual map[string]string `json:"creator_notes_multilingual"`
	Source                   []string          `json:"source"`
	GroupOnlyGreetings       []string          `json:"group_only_greetings"`
	CreationDate             *int              `json:"creation_date"`
	ModificationDate         *int              `json:"modification_date"`
}

// ReadPNG reads a PNG file and extracts the Character Card base64 metadata
func ReadPNG(file string) (string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Extract base64 data from PNG
	base64Data, err := extractBase64FromPNG(data)
	if err != nil {
		return "", fmt.Errorf("failed to extract base64: %w", err)
	}

	return base64Data, nil
}

// ReadPNGAsV3 reads a PNG file and extracts the Character Card V3 metadata as a struct
func ReadPNGAsV3(file string) (*CharacterCardV3, error) {
	base64Data, err := ReadPNG(file)
	if err != nil {
		return nil, err
	}

	// Decode base64
	jsonData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	// Parse JSON
	var card CharacterCardV3
	if err := json.Unmarshal(jsonData, &card); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &card, nil
}

// extractBase64FromPNG extracts the base64-encoded character data from a PNG file
func extractBase64FromPNG(data []byte) (string, error) {
	// Verify PNG signature
	if len(data) < 8 || !bytes.Equal(data[:8], []byte{137, 80, 78, 71, 13, 10, 26, 10}) {
		return "", errors.New("not a valid PNG file")
	}

	// Parse PNG chunks looking for text chunks with keyword "chara"
	pos := 8 // Skip PNG signature

	for pos < len(data)-8 {
		// Read chunk length (4 bytes, big-endian)
		if pos+4 > len(data) {
			break
		}
		chunkLen := binary.BigEndian.Uint32(data[pos : pos+4])
		pos += 4

		// Read chunk type (4 bytes)
		if pos+4 > len(data) {
			break
		}
		chunkType := string(data[pos : pos+4])
		pos += 4

		// Check if we have enough data for chunk data
		if pos+int(chunkLen) > len(data) {
			break
		}

		// Check for tEXt or zTXt chunks with keyword "chara"
		if chunkType == "tEXt" {
			chunkData := data[pos : pos+int(chunkLen)]
			// tEXt format: keyword\0text
			nullPos := bytes.IndexByte(chunkData, 0)
			if nullPos != -1 && string(chunkData[:nullPos]) == "chara" {
				return string(chunkData[nullPos+1:]), nil
			}
		}

		// Move to next chunk (skip chunk data + 4 bytes CRC)
		pos += int(chunkLen) + 4
	}

	return "", errors.New("no character card metadata found in PNG chunks")
}

// findTextChunk searches for a tEXt chunk with the specified keyword
func findTextChunk(data []byte, keyword string) (string, bool) {
	pos := 0
	for pos < len(data)-8 {
		// Check if we have enough data for chunk length
		if pos+4 > len(data) {
			break
		}

		// Read chunk length (big-endian)
		chunkLen := binary.BigEndian.Uint32(data[pos : pos+4])
		pos += 4

		// Check if we have enough data for chunk type
		if pos+4 > len(data) {
			break
		}

		// Read chunk type
		chunkType := string(data[pos : pos+4])
		pos += 4

		// Check if we have enough data for chunk data
		if pos+int(chunkLen) > len(data) {
			break
		}

		if chunkType == "tEXt" {
			chunkData := data[pos : pos+int(chunkLen)]
			// tEXt format: keyword\0text
			nullPos := bytes.IndexByte(chunkData, 0)
			if nullPos != -1 {
				kw := string(chunkData[:nullPos])
				if kw == keyword {
					return string(chunkData[nullPos+1:]), true
				}
			}
		}

		// Move to next chunk (skip chunk data + 4 bytes CRC)
		pos += int(chunkLen) + 4
	}

	return "", false
}

// extractBase64Direct extracts base64 data when it's appended directly after IEND
func extractBase64Direct(data []byte) string {
	// Look for the fixed part of the footer/trailer that marks the end of the base64 data
	// Footer pattern: [4 variable bytes] 00 00 00 00 49 45 4E 44 AE 42 60 82
	// The fixed part is: 00 00 00 00 49 45 4E 44 AE 42 60 82
	fixedFooter := []byte{0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82}
	footerIdx := bytes.Index(data, fixedFooter)

	// If footer found, trim the data to exclude it and the 4 variable bytes before it
	if footerIdx != -1 {
		// Go back 4 bytes to include the variable prefix
		if footerIdx >= 4 {
			data = data[:footerIdx-4]
		} else {
			data = data[:footerIdx]
		}
	}

	// Remove any leading/trailing whitespace and null bytes
	data = bytes.TrimSpace(data)
	data = bytes.Trim(data, "\x00")

	// Look for the start of base64 data
	// Character Card V3 base64 typically starts with 'eyJ' (which is '{' in JSON)
	startIdx := bytes.Index(data, []byte("eyJ"))
	if startIdx == -1 {
		return ""
	}

	// Extract from that point to the end, removing non-base64 characters
	base64Data := data[startIdx:]

	// Clean up: keep only valid base64 characters
	var cleaned bytes.Buffer
	for _, b := range base64Data {
		if (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || (b >= '0' && b <= '9') || b == '+' || b == '/' || b == '=' {
			cleaned.WriteByte(b)
		}
	}

	return cleaned.String()
}
