package scenario

import (
	"encoding/json"
	"fmt"
	"os"
)

// SaveFerexFile marshals the FerexFile struct to JSON and writes it to the specified filePath.
func SaveFerexFile(filePath string, data *FerexFile) error {
	if data == nil {
		return fmt.Errorf("cannot save nil FerexFile data")
	}

	jsonData, err := json.MarshalIndent(data, "", "  ") // Using MarshalIndent for pretty printing
	if err != nil {
		return fmt.Errorf("failed to marshal FerexFile data to JSON: %w", err)
	}

	err = os.WriteFile(filePath, jsonData, 0644) // 0644 provides read/write for owner, read for group/others
	if err != nil {
		return fmt.Errorf("failed to write FerexFile to '%s': %w", filePath, err)
	}

	return nil
}

// LoadFerexFile reads a JSON file from filePath and unmarshals it into a FerexFile struct.
func LoadFerexFile(filePath string) (*FerexFile, error) {
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read FerexFile from '%s': %w", filePath, err)
	}

	var data FerexFile
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal FerexFile JSON from '%s': %w", filePath, err)
	}

	return &data, nil
}
