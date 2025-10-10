# managerCharAI

Go library to read and write Character Card AI metadata embedded in PNG files. Supports Character Card V2 and V3 specifications.

## Features

- **Read** Character Cards from PNG files with embedded metadata
- **Write** Character Cards into PNG files as embedded metadata
- **Read** Character Cards from standalone JSON files
- **Write** Character Cards to standalone JSON files  
- **Parse** Character Card structures (V2/V3 compatible)
- **Convert** between PNG and JSON formats
- **Base64 encoding/decoding** support
- **Standard PNG chunk format** (tEXt chunk with "chara" keyword)

## Installation

```bash
go get github.com/jonathanhecl/managerCharAI
```

## Quick Start

### Reading from PNG

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

### Reading from JSON

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/jonathanhecl/managerCharAI"
)

func main() {
    // Read a character card from JSON file
    card, err := managerCharAI.ReadJSON("character.json")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Character: %s\n", card.Name)
    fmt.Printf("Creator: %s\n", card.Data.Creator)
}
```

## Function Summary

### Package Functions

| Function | Description |
|----------|-------------|
| `ReadPNG()` | Extract base64 metadata from PNG |
| `ReadPNGAsCard()` | Read PNG and parse to CharacterCard struct |
| `ReadJSON()` | Read standalone JSON file to CharacterCard struct |
| `WritePNG()` | Create PNG with embedded metadata from base64 strings |
| `WritePNGFromCard()` | Create PNG with embedded metadata from CharacterCard struct |
| `WriteJSON()` | Write CharacterCard struct to standalone JSON file |

### CharacterCard Methods

| Method | Description |
|--------|-------------|
| `card.ToJSON()` | Convert CharacterCard to JSON string |
| `card.Format()` | Get the specification format (e.g., "chara_card_v3") |
| `card.SaveJSON(file)` | Save CharacterCard to JSON file |
| `card.SavePNG(file, imageBase64)` | Save CharacterCard to PNG with embedded metadata |

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

#### `ReadJSON(file string) (*CharacterCard, error)`
Reads a JSON file and parses it into a CharacterCard struct.

**Parameters:**
- `file`: Path to the JSON file

**Returns:**
- `*CharacterCard`: Parsed Character Card struct
- `error`: Error if any

**Example:**
```go
card, err := managerCharAI.ReadJSON("character.json")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Name: %s\n", card.Name)
fmt.Printf("Creator: %s\n", card.Data.Creator)
fmt.Printf("Description: %s\n", card.Description)
```

### Writing Functions

#### `WritePNG(outputFile, imageBase64, metadataBase64 string) error`
Creates a PNG file with embedded Character Card metadata.

**Parameters:**
- `outputFile`: Path where to save the output PNG file
- `imageBase64`: Base64-encoded source PNG image
- `metadataBase64`: Base64-encoded Character Card JSON

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
err := managerCharAI.WritePNG("output.png", imageBase64, metadataBase64)
if err != nil {
    log.Fatal(err)
}
```

#### `WritePNGFromCard(outputFile, imageBase64 string, card *CharacterCard) error`
Creates a PNG file from a CharacterCard struct.

**Parameters:**
- `outputFile`: Path where to save the output PNG file
- `imageBase64`: Base64-encoded source PNG image
- `card`: CharacterCard struct to embed

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
err := managerCharAI.WritePNGFromCard("my_character.png", imageBase64, card)
if err != nil {
    log.Fatal(err)
}
```

#### `WriteJSON(outputFile string, card *CharacterCard) error`
Writes a CharacterCard struct to a JSON file with pretty formatting.

**Parameters:**
- `outputFile`: Path where to save the JSON file
- `card`: CharacterCard struct to write

**Returns:**
- `error`: Error if any

**Example:**
```go
card := &managerCharAI.CharacterCard{
    Name:        "My Character",
    Description: "A brave warrior",
    Spec:        "chara_card_v3",
    SpecVersion: "3.0",
    Tags:        []string{"fantasy", "warrior"},
    Data: managerCharAI.CharacterData{
        Name:        "My Character",
        Creator:     "YourName",
        Description: "Detailed description",
        FirstMes:    "Hello!",
    },
}

err := managerCharAI.WriteJSON("my_character.json", card)
if err != nil {
    log.Fatal(err)
}
```

## CharacterCard Methods

### `card.ToJSON() (string, error)`
Converts the CharacterCard to a JSON string.

**Example:**
```go
card := &managerCharAI.CharacterCard{
    Name: "My Character",
    Spec: "chara_card_v3",
}

jsonString, err := card.ToJSON()
if err != nil {
    log.Fatal(err)
}
fmt.Println(jsonString)
```

### `card.Format() string`
Returns the specification format of the character card.

**Example:**
```go
card, _ := managerCharAI.ReadPNGAsCard("character.png")
fmt.Printf("Format: %s\n", card.Format()) // Output: chara_card_v3
```

### `card.SaveJSON(outputFile string) error`
Saves the CharacterCard directly to a JSON file.

**Example:**
```go
card := &managerCharAI.CharacterCard{
    Name:        "My Character",
    Description: "A brave warrior",
    Spec:        "chara_card_v3",
    SpecVersion: "3.0",
    Data: managerCharAI.CharacterData{
        Name:    "My Character",
        Creator: "YourName",
    },
}

// Save directly using the method
err := card.SaveJSON("my_character.json")
if err != nil {
    log.Fatal(err)
}
```

### `card.SavePNG(outputFile, imageBase64 string) error`
Saves the CharacterCard to a PNG file with embedded metadata.

**Parameters:**
- `outputFile`: Path where to save the PNG file
- `imageBase64`: Base64-encoded PNG image

**Example:**
```go
// Create or load a character card
card := &managerCharAI.CharacterCard{
    Name:        "My Character",
    Description: "A brave warrior",
    Spec:        "chara_card_v3",
    SpecVersion: "3.0",
    Data: managerCharAI.CharacterData{
        Name:    "My Character",
        Creator: "YourName",
    },
}

// Read base image
imageData, _ := os.ReadFile("avatar.png")
imageBase64 := base64.StdEncoding.EncodeToString(imageData)

// Save directly using the method
err := card.SavePNG("my_character.png", imageBase64)
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

### 1. Using Methods - Quick Save

```go
func quickSaveExample() error {
    // Create a character card
    card := &managerCharAI.CharacterCard{
        Name:        "Quick Character",
        Description: "Created and saved quickly",
        Spec:        "chara_card_v3",
        SpecVersion: "3.0",
        Data: managerCharAI.CharacterData{
            Name:    "Quick Character",
            Creator: "QuickCreator",
        },
    }
    
    // Save to JSON using method
    if err := card.SaveJSON("quick_character.json"); err != nil {
        return err
    }
    
    // Save to PNG using method
    imageData, _ := os.ReadFile("avatar.png")
    imageBase64 := base64.StdEncoding.EncodeToString(imageData)
    
    if err := card.SavePNG("quick_character.png", imageBase64); err != nil {
        return err
    }
    return nil
}

### 2. Method vs Function - Two ways to Save

```go
func demonstrateSavingMethods() error {
    card := &managerCharAI.CharacterCard{
        Name:        "Demo Character",
        Spec:        "chara_card_v3",
        SpecVersion: "3.0",
        Data: managerCharAI.CharacterData{
            Name:    "Demo Character",
            Creator: "DemoCreator",
        },
    }
    
    // Method 1: Using package functions
    err := managerCharAI.WriteJSON("output1.json", card)
    if err != nil {
        return err
    }
    
    // Method 2: Using struct methods (more convenient)
    err = card.SaveJSON("output2.json")
    if err != nil {
        return err
    }
    
    // Both produce the same result!
    return nil
}
```

### 3. Extract Character Information from PNG

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

### 4. Modify Character Card and Save

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
    return managerCharAI.WritePNGFromCard(outputFile, imageBase64, card)
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
    return managerCharAI.WritePNGFromCard(outputPath, imageBase64, card)
}
```

### 4. Convert JSON to PNG

```go
func convertJSONToPNG(jsonFile, imageFile, outputFile string) error {
    // Read character card from JSON
    card, err := managerCharAI.ReadJSON(jsonFile)
    if err != nil {
        return err
    }
    
    // Read base image
    imageData, err := os.ReadFile(imageFile)
    if err != nil {
        return err
    }
    imageBase64 := base64.StdEncoding.EncodeToString(imageData)
    
    // Create PNG with embedded character card
    return managerCharAI.WritePNGFromCard(outputFile, imageBase64, card)
}
```

### 5. Extract JSON from PNG

```go
func extractJSONFromPNG(pngFile, outputJSONFile string) error {
    // Read character card from PNG
    card, err := managerCharAI.ReadPNGAsCard(pngFile)
    if err != nil {
        return err
    }
    
    // Write to JSON file
    return managerCharAI.WriteJSON(outputJSONFile, card)
}
```

### 6. Batch Process Multiple Character Cards

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
        
        if err := managerCharAI.WritePNGFromCard(outputPath, imageBase64, card); err != nil {
            fmt.Printf("Error saving %s: %v\n", file.Name(), err)
        }
    }
    
    return nil
}
```

## Testing

Run all tests:

```bash
go test -v
```

Run specific tests:

```bash
# Test PNG reading
go test -v -run TestReadPNG

# Test PNG writing
go test -v -run TestWritePNG

# Test JSON reading
go test -v -run TestReadJSON

# Test JSON writing (creates test_output.json for inspection)
go test -v -run TestWriteJSON

# Test JSON to PNG conversion
go test -v -run TestJSONToPNG

# Test CharacterCard.SaveJSON method
go test -v -run TestCharacterCard_SaveJSON

# Test CharacterCard.SavePNG method
go test -v -run TestCharacterCard_SavePNG
```

### Test Output Files

The tests create output files for inspection:
- `test_output.json` - Example JSON character card (from TestWriteJSON)
- `test_output.png` - Example PNG with embedded metadata (from TestWritePNG)
- `test_method_output.json` - Example using SaveJSON method (from TestCharacterCard_SaveJSON)
- `test_method_output.png` - Example using SavePNG method (from TestCharacterCard_SavePNG)

These files are kept on disk so you can verify the output format and content.

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

## Complete Feature Matrix

| Operation | Package Function | CharacterCard Method |
|-----------|-----------------|---------------------|
| Read PNG to struct | `ReadPNGAsCard(file)` | N/A |
| Read JSON to struct | `ReadJSON(file)` | N/A |
| Write struct to JSON | `WriteJSON(file, card)` | `card.SaveJSON(file)` ✨ |
| Write struct to PNG | `WritePNGFromCard(file, img64, card)` | `card.SavePNG(file, img64)` ✨ |
| Convert to JSON string | N/A | `card.ToJSON()` |
| Get format | N/A | `card.Format()` |

✨ = Convenient method alternative

## Workflow Examples

```
┌─────────────┐
│  JSON File  │──ReadJSON()──────────┐
└─────────────┘                      │
                                     ▼
┌─────────────┐              ┌──────────────┐
│  PNG File   │──ReadPNG()──▶│ CharacterCard│
└─────────────┘              │    Struct    │
                             └──────────────┘
                                     │
                    ┌────────────────┼────────────────┐
                    ▼                ▼                ▼
            card.SaveJSON()  card.SavePNG()   card.ToJSON()
                    │                │                │
                    ▼                ▼                ▼
            ┌─────────────┐  ┌─────────────┐  ┌──────────┐
            │  JSON File  │  │  PNG File   │  │  String  │
            └─────────────┘  └─────────────┘  └──────────┘
```

## License

See LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.