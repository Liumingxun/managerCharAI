package readerCharV3_test

import (
	"testing"

	"github.com/jonathanhecl/readerCharV3"
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
			card, err := readerCharV3.ReadPNG(tt.file)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("ReadPNG() failed: %v", err)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ReadPNG() succeeded unexpectedly")
			}
			if card == nil {
				t.Fatal("ReadPNG() returned nil card")
			}
			
			// Basic validation
			if card.Name == "" {
				t.Error("Character name is empty")
			}
			
			t.Logf("Successfully read character: %s", card.Name)
			t.Logf("Description: %s", card.Description)
		})
	}
}
