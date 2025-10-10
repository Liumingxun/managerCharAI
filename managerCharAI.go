package managerCharAI

import (
	"bytes"
	"encoding/binary"
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
