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

// DepthPrompt contains depth prompt configuration
type DepthPrompt struct {
	Depth  int    `json:"depth"`
	Prompt string `json:"prompt"`
	Role   string `json:"role"`
}

// Extensions contains additional metadata and settings
type Extensions struct {
	DepthPrompt   DepthPrompt `json:"depth_prompt"`
	Fav           bool        `json:"fav"`
	Talkativeness string      `json:"talkativeness"`
	World         string      `json:"world"`
}

// CharacterData contains the core character information
type CharacterData struct {
	AlternateGreetings      []string   `json:"alternate_greetings"`
	CharacterVersion        string     `json:"character_version"`
	Creator                 string     `json:"creator"`
	CreatorNotes            string     `json:"creator_notes"`
	Description             string     `json:"description"`
	Extensions              Extensions `json:"extensions"`
	FirstMes                string     `json:"first_mes"`
	GroupOnlyGreetings      []string   `json:"group_only_greetings"`
	MesExample              string     `json:"mes_example"`
	Name                    string     `json:"name"`
	Personality             string     `json:"personality"`
	PostHistoryInstructions string     `json:"post_history_instructions"`
	Scenario                string     `json:"scenario"`
	SystemPrompt            string     `json:"system_prompt"`
	Tags                    []string   `json:"tags"`
}

// CharacterCard represents the complete character card structure
type CharacterCard struct {
	Avatar         string        `json:"avatar"`
	Chat           string        `json:"chat"`
	CreateDate     string        `json:"create_date"`
	CreatorComment string        `json:"creatorcomment"`
	Data           CharacterData `json:"data"`
	Description    string        `json:"description"`
	Fav            bool          `json:"fav"`
	FirstMes       string        `json:"first_mes"`
	MesExample     string        `json:"mes_example"`
	Name           string        `json:"name"`
	Personality    string        `json:"personality"`
	Scenario       string        `json:"scenario"`
	Spec           string        `json:"spec"`
	SpecVersion    string        `json:"spec_version"`
	Tags           []string      `json:"tags"`
	Talkativeness  string        `json:"talkativeness"`
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

// ReadPNGAsCard reads a PNG file and extracts the Character Card as a struct
func ReadPNGAsCard(file string) (*CharacterCard, error) {
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
	var card CharacterCard
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
