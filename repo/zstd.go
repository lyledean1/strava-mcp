package repo

import (
	"encoding/json"
	"fmt"
	"github.com/klauspost/compress/zstd"
	"os"
	"path/filepath"
)

func SaveToZstd(data interface{}, filename string) error {
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	//nolint: errcheck // defer is used to clean up
	defer file.Close()

	encoder, err := zstd.NewWriter(file)
	if err != nil {
		return fmt.Errorf("failed to create zstd encoder: %w", err)
	}

	//nolint: errcheck // defer is used to clean up
	defer encoder.Close()

	jsonEncoder := json.NewEncoder(encoder)
	jsonEncoder.SetIndent("", "  ")
	if err := jsonEncoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

func LoadFromZstd(filename string, target interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	//nolint: errcheck // defer is used to clean up
	defer file.Close()

	decoder, err := zstd.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create zstd decoder: %w", err)
	}
	defer decoder.Close()

	jsonDecoder := json.NewDecoder(decoder)

	if err := jsonDecoder.Decode(target); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}
	return nil
}
