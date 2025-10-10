package readerCharV3_test

import (
	"testing"

	"github.com/jonathanhecl/readerCharV3"
)

func TestReadPNG(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		want    string
		wantErr bool
	}{
		{
			name:    "test",
			file:    "test.png",
			want:    "test",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := readerCharV3.ReadPNG(tt.file)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ReadPNG() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ReadPNG() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("ReadPNG() = %v, want %v", got, tt.want)
			}
		})
	}
}
