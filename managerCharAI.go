package managerCharAI

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
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

// ToJSON converts the CharacterCard to a JSON string
func (c *CharacterCard) ToJSON() (string, error) {
	jsonData, err := json.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return string(jsonData), nil
}

// Format returns the specification format of the character card
func (c *CharacterCard) Format() string {
	return c.Spec
}

// SaveJSON saves the CharacterCard to a JSON file
// outputFile: path where to save the JSON file
func (c *CharacterCard) SaveJSON(outputFile string) error {
	return WriteJSON(outputFile, c)
}

// SavePNG saves the CharacterCard to a PNG file with embedded metadata
// outputFile: path where to save the PNG file
// imageBase64: base64-encoded PNG image to use as the base
func (c *CharacterCard) SavePNG(outputFile, imageBase64 string) error {
	return WritePNGFromCard(outputFile, imageBase64, c)
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

// ReadJSON reads a JSON file and parses it into a CharacterCard struct
func ReadJSON(file string) (*CharacterCard, error) {
	jsonData, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON file: %w", err)
	}

	var card CharacterCard
	if err := json.Unmarshal(jsonData, &card); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &card, nil
}

// ReadString reads a JSON string and parses it into a CharacterCard struct
func ReadString(jsonStr string) (*CharacterCard, error) {
	var card CharacterCard
	if err := json.Unmarshal([]byte(jsonStr), &card); err != nil {
		return nil, fmt.Errorf("failed to parse JSON string: %w", err)
	}

	return &card, nil
}

// WriteJSON writes a CharacterCard struct to a JSON file
// outputFile: path where to save the JSON file
// card: CharacterCard struct to write
func WriteJSON(outputFile string, card *CharacterCard) error {
	jsonData, err := json.MarshalIndent(card, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(outputFile, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	return nil
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

// WritePNG creates a PNG file with embedded character card metadata
// outputFile: path to save the output PNG file
// imageBase64: base64-encoded PNG image
// metadataBase64: base64-encoded character card JSON
func WritePNG(outputFile, imageBase64, metadataBase64 string) error {
	// Decode the image from base64
	imageData, err := base64.StdEncoding.DecodeString(imageBase64)
	if err != nil {
		return fmt.Errorf("failed to decode image base64: %w", err)
	}

	// Verify it's a valid PNG
	if len(imageData) < 8 || !bytes.Equal(imageData[:8], []byte{137, 80, 78, 71, 13, 10, 26, 10}) {
		return errors.New("image data is not a valid PNG file")
	}

	// Find the IEND chunk position
	iendMarker := []byte{0x49, 0x45, 0x4E, 0x44}
	iendPos := bytes.Index(imageData, iendMarker)
	if iendPos == -1 {
		return errors.New("IEND chunk not found in PNG")
	}

	// Position right before IEND chunk (subtract 4 bytes for length field)
	insertPos := iendPos - 4

	// Create the tEXt chunk with keyword "chara"
	keyword := "chara"
	chunkData := []byte(keyword)
	chunkData = append(chunkData, 0) // null separator
	chunkData = append(chunkData, []byte(metadataBase64)...)

	// Calculate chunk length
	chunkLen := uint32(len(chunkData))

	// Build the complete chunk
	var chunk bytes.Buffer

	// Write chunk length (4 bytes, big-endian)
	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, chunkLen)
	chunk.Write(lenBytes)

	// Write chunk type "tEXt" (4 bytes)
	chunkType := []byte("tEXt")
	chunk.Write(chunkType)

	// Write chunk data
	chunk.Write(chunkData)

	// Calculate CRC32 of chunk type + chunk data
	crc := crc32.NewIEEE()
	crc.Write(chunkType)
	crc.Write(chunkData)
	crcBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(crcBytes, crc.Sum32())
	chunk.Write(crcBytes)

	// Construct the final PNG
	var finalPNG bytes.Buffer
	finalPNG.Write(imageData[:insertPos])
	finalPNG.Write(chunk.Bytes())
	finalPNG.Write(imageData[insertPos:])

	// Write to file
	if err := os.WriteFile(outputFile, finalPNG.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

// WritePNGFromCard creates a PNG file with embedded character card metadata from a CharacterCard struct
// outputFile: path to save the output PNG file
// imageBase64: base64-encoded PNG image
// card: CharacterCard struct to embed
func WritePNGFromCard(outputFile, imageBase64 string, card *CharacterCard) error {
	// Convert card to JSON
	jsonData, err := json.Marshal(card)
	if err != nil {
		return fmt.Errorf("failed to marshal card to JSON: %w", err)
	}

	// Encode JSON to base64
	metadataBase64 := base64.StdEncoding.EncodeToString(jsonData)

	// Use WritePNG to create the file
	return WritePNG(outputFile, imageBase64, metadataBase64)
}
