package managerCharAI

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
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

func ReadPNGFromReader(r *bytes.Reader) (*CharacterCard, error) {
	base64, err := extractBase64FromReader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to extract base64: %w", err)
	}

	var card CharacterCard
	if err := json.Unmarshal(base64, &card); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return &card, nil
}

// ReadPNG reads a PNG file and extracts the Character Card base64 metadata
func ReadPNG(file string) (string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	// Extract base64 data from PNG
	base64Data, err := extractBase64FromReader(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("failed to extract base64: %w", err)
	}

	return base64.StdEncoding.EncodeToString(base64Data), nil
}

// ReadPNGAsCard reads a PNG file and extracts the Character Card as a struct
func ReadPNGAsCard(file string) (*CharacterCard, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return ReadPNGFromReader(bytes.NewReader(data))
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

type pngChunkHeader struct {
	Length uint32
	Type   [4]byte
}

// extractBase64FromPNG extracts the base64-encoded character data from a PNG file
func extractBase64FromReader(imgReader *bytes.Reader) ([]byte, error) {
	// Verify PNG signature
	signature := make([]byte, 8)

	if _, err := io.ReadFull(imgReader, signature); err != nil ||
		!bytes.Equal(signature, []byte{137, 80, 78, 71, 13, 10, 26, 10}) {
		return nil, errors.New("not a valid PNG file")
	}

	var chunk pngChunkHeader
	var rawData []byte
	var found bool
	for {
		err := binary.Read(imgReader, binary.BigEndian, &chunk)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if string(chunk.Type[:]) == "tEXt" {
			rawData = make([]byte, chunk.Length)
			if _, err := imgReader.Read(rawData); err != nil {
				return nil, err
			}

			if string(rawData[:6]) == "chara\x00" {
				found = true
				break
			}

			// tEXt but not chara: skip only CRC (4 bytes)
			_, err = io.CopyN(io.Discard, imgReader, 4)
			if err != nil {
				return nil, err
			}
		} else {
			// Non-tEXt: skip data + CRC
			_, err = io.CopyN(io.Discard, imgReader, int64(chunk.Length)+4)
			if err != nil {
				return nil, err
			}
		}
	}

	if !found {
		return nil, errors.New("no character card metadata found in PNG chunks")
	}

	decoded := make([]byte, base64.RawStdEncoding.DecodedLen(len(rawData[6:])))
	n, err := base64.StdEncoding.Decode(decoded, rawData[6:])
	if err != nil {
		return nil, err
	}
	return decoded[:n], nil
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
