# managerCharAI

Go library to read and write Character Card AI metadata embedded in PNG files. Supports Character Card V2 and V3 specifications.

## Features

- **Read** Character Cards from PNG files with embedded metadata
- **Write** Character Cards into PNG files as embedded metadata  
- **Parse** Character Card structures (V2/V3 compatible)
- **Base64 encoding/decoding** support
- **Standard PNG chunk format** (tEXt chunk with "chara" keyword)

## Installation

```bash
go get github.com/jonathanhecl/managerCharAI
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/jonathanhecl/managerCharAI"
)

func main() {
    // Read a character card from PNG
    card, err := managerCharAI.ReadPNGAsCard("character.png")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Character: %s\n", card.Name)
    fmt.Printf("Description: %s\n", card.Description)
}
```

## API Reference

### Reading Functions

#### `ReadPNG(file string) (string, error)`
Extracts the base64-encoded Character Card metadata from a PNG file.

**Parameters:**
- `file`: Path to the PNG file

**Returns:**
- `string`: Base64-encoded Character Card JSON
- `error`: Error if any

**Example:**
```go
base64Data, err := managerCharAI.ReadPNG("character.png")
if err != nil {
    log.Fatal(err)
}
fmt.Println("Base64 metadata:", base64Data)
```

#### `ReadPNGAsCard(file string) (*CharacterCard, error)`
Reads a PNG file and parses the Character Card into a struct.

**Parameters:**
- `file`: Path to the PNG file

**Returns:**
- `*CharacterCard`: Parsed Character Card struct
- `error`: Error if any

**Example:**
```go
card, err := managerCharAI.ReadPNGAsCard("character.png")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Name: %s\n", card.Name)
fmt.Printf("Creator: %s\n", card.Data.Creator)
fmt.Printf("Tags: %v\n", card.Tags)
fmt.Printf("First Message: %s\n", card.FirstMes)
```

### Writing Functions

#### `WritePNG(imageBase64, metadataBase64, outputFile string) error`
Creates a PNG file with embedded Character Card metadata.

**Parameters:**
- `imageBase64`: Base64-encoded source PNG image
- `metadataBase64`: Base64-encoded Character Card JSON
- `outputFile`: Path where to save the output PNG file

**Returns:**
- `error`: Error if any

**Example:**
```go
import (
    "encoding/base64"
    "os"
)

// Read source image
imageData, _ := os.ReadFile("avatar.png")
imageBase64 := base64.StdEncoding.EncodeToString(imageData)

// Create metadata JSON
jsonData := `{"name":"My Character","spec":"chara_card_v3","spec_version":"3.0"}`
metadataBase64 := base64.StdEncoding.EncodeToString([]byte(jsonData))

// Write PNG with embedded metadata
err := managerCharAI.WritePNG(imageBase64, metadataBase64, "output.png")
if err != nil {
    log.Fatal(err)
}
```

#### `WritePNGFromCard(imageBase64 string, card *CharacterCard, outputFile string) error`
Creates a PNG file from a CharacterCard struct.

**Parameters:**
- `imageBase64`: Base64-encoded source PNG image
- `card`: CharacterCard struct to embed
- `outputFile`: Path where to save the output PNG file

**Returns:**
- `error`: Error if any

**Example:**
```go
// Create a new character card
card := &managerCharAI.CharacterCard{
    Name:        "My Character",
    Description: "A brave warrior",
    Spec:        "chara_card_v3",
    SpecVersion: "3.0",
    Tags:        []string{"fantasy", "warrior", "male"},
    Data: managerCharAI.CharacterData{
        Name:        "My Character",
        Description: "Detailed description here",
        FirstMes:    "Hello! I'm ready for adventure!",
        Personality: "Brave, loyal, determined",
        Scenario:    "Medieval fantasy setting",
        Tags:        []string{"fantasy", "warrior"},
    },
}

// Read base image
imageData, _ := os.ReadFile("avatar.png")
imageBase64 := base64.StdEncoding.EncodeToString(imageData)

// Create PNG with character card
err := managerCharAI.WritePNGFromCard(imageBase64, card, "my_character.png")
if err != nil {
    log.Fatal(err)
}
```

## Data Structures

### CharacterCard
Main structure representing a complete Character Card.

```go
type CharacterCard struct {
    Avatar         string        // Avatar image reference
    Chat           string        // Chat data
    CreateDate     string        // Creation date
    CreatorComment string        // Creator's comment
    Data           CharacterData // Core character data
    Description    string        // Character description
    Fav            bool          // Favorite flag
    FirstMes       string        // First message
    MesExample     string        // Message examples
    Name           string        // Character name
    Personality    string        // Personality traits
    Scenario       string        // Scenario description
    Spec           string        // Specification version (e.g., "chara_card_v3")
    SpecVersion    string        // Version number (e.g., "3.0")
    Tags           []string      // Character tags
    Talkativeness  string        // Talkativeness level
}
```

### CharacterData
Core character information within the data field.

```go
type CharacterData struct {
    AlternateGreetings      []string   // Alternative greeting messages
    CharacterVersion        string     // Character version
    Creator                 string     // Creator name
    CreatorNotes            string     // Creator's notes
    Description             string     // Detailed description
    Extensions              Extensions // Additional extensions
    FirstMes                string     // First message
    GroupOnlyGreetings      []string   // Group-only greetings
    MesExample              string     // Message examples
    Name                    string     // Character name
    Personality             string     // Personality description
    PostHistoryInstructions string     // Post-history instructions
    Scenario                string     // Scenario details
    SystemPrompt            string     // System prompt
    Tags                    []string   // Tags
}
```

## Common Use Cases

### 1. Extract Character Information from PNG

```go
func extractCharacterInfo(filename string) error {
    card, err := managerCharAI.ReadPNGAsCard(filename)
    if err != nil {
        return err
    }
    
    fmt.Printf("=== Character Card ===\n")
    fmt.Printf("Name: %s\n", card.Name)
    fmt.Printf("Spec: %s v%s\n", card.Spec, card.SpecVersion)
    fmt.Printf("Creator: %s\n", card.Data.Creator)
    fmt.Printf("Tags: %v\n", card.Tags)
    fmt.Printf("Description: %s\n", card.Description)
    
    return nil
}
```

### 2. Modify Character Card and Save

```go
func modifyAndSave(inputFile, outputFile string) error {
    // Read existing card
    card, err := managerCharAI.ReadPNGAsCard(inputFile)
    if err != nil {
        return err
    }
    
    // Modify card
    card.Name = "Modified " + card.Name
    card.Tags = append(card.Tags, "modified")
    
    // Read original image
    imageData, err := os.ReadFile(inputFile)
    if err != nil {
        return err
    }
    imageBase64 := base64.StdEncoding.EncodeToString(imageData)
    
    // Save modified card
    return managerCharAI.WritePNGFromCard(imageBase64, card, outputFile)
}
```

### 3. Create Character Card from Scratch

```go
func createCharacter(imagePath, outputPath string) error {
    // Create character card
    card := &managerCharAI.CharacterCard{
        Name:        "Luna the Mage",
        Description: "A powerful sorceress who protects her village",
        Spec:        "chara_card_v3",
        SpecVersion: "3.0",
        Tags:        []string{"fantasy", "mage", "female", "magic"},
        Fav:         false,
        Data: managerCharAI.CharacterData{
            Name:        "Luna the Mage",
            Creator:     "YourName",
            Description: "Luna is a powerful sorceress...",
            FirstMes:    "Greetings, traveler. What brings you to my tower?",
            Personality: "Wise, mysterious, helpful but cautious",
            Scenario:    "Fantasy medieval setting with magic",
            Tags:        []string{"fantasy", "mage", "female"},
            Extensions: managerCharAI.Extensions{
                Talkativeness: "0.7",
                Fav:          false,
            },
        },
    }
    
    // Read avatar image
    imageData, err := os.ReadFile(imagePath)
    if err != nil {
        return err
    }
    imageBase64 := base64.StdEncoding.EncodeToString(imageData)
    
    // Create PNG with card
    return managerCharAI.WritePNGFromCard(imageBase64, card, outputPath)
}
```

### 4. Batch Process Multiple Character Cards

```go
func batchProcess(inputDir, outputDir string) error {
    files, err := os.ReadDir(inputDir)
    if err != nil {
        return err
    }
    
    for _, file := range files {
        if filepath.Ext(file.Name()) != ".png" {
            continue
        }
        
        inputPath := filepath.Join(inputDir, file.Name())
        
        // Read card
        card, err := managerCharAI.ReadPNGAsCard(inputPath)
        if err != nil {
            fmt.Printf("Skipping %s: %v\n", file.Name(), err)
            continue
        }
        
        // Process card
        fmt.Printf("Processing: %s (%s)\n", card.Name, file.Name())
        
        // Save to output directory
        outputPath := filepath.Join(outputDir, file.Name())
        imageData, _ := os.ReadFile(inputPath)
        imageBase64 := base64.StdEncoding.EncodeToString(imageData)
        
        if err := managerCharAI.WritePNGFromCard(imageBase64, card, outputPath); err != nil {
            fmt.Printf("Error saving %s: %v\n", file.Name(), err)
        }
    }
    
    return nil
}
```

## Testing

Run tests with:

```bash
go test -v
```

Run specific test:

```bash
go test -v -run TestReadPNG
go test -v -run TestWritePNG
```

## Technical Details

- Uses standard PNG tEXt chunk format
- Keyword: `chara`
- Data format: Base64-encoded JSON
- Compatible with Character Card V2 and V3 specifications
- Properly calculates CRC32 for chunk integrity

## Useful Links

- **Testing tool:** [AICharED](https://desune.moe/aichared/)  
- **Spec V2:** [SPEC_V2.md](https://github.com/malfoyslastname/character-card-spec-v2/blob/main/spec_v2.md)  
- **Spec V3:** [SPEC_V3.md](https://github.com/kwaroran/character-card-spec-v3/blob/main/SPEC_V3.md)

## License

See LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.