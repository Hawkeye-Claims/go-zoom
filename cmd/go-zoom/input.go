package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

// readJSONInput parses JSON from either an inline flag value or a JSON file path.
func readJSONInput[T any](jsonInput string, jsonFile string) (T, error) {
	var payload T
	if jsonInput == "" && jsonFile == "" {
		return payload, errors.New("one of --json or --json-file is required")
	}
	if jsonInput != "" && jsonFile != "" {
		return payload, errors.New("--json and --json-file cannot be used together")
	}

	var raw []byte
	if jsonInput != "" {
		raw = []byte(jsonInput)
	} else {
		fileBytes, err := os.ReadFile(jsonFile)
		if err != nil {
			return payload, fmt.Errorf("failed to read json file: %w", err)
		}
		raw = fileBytes
	}

	if err := unmarshalJSONStrict(raw, &payload); err != nil {
		return payload, err
	}
	return payload, nil
}

// unmarshalJSONStrict decodes exactly one JSON value and rejects unknown fields.
func unmarshalJSONStrict(raw []byte, target any) error {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("invalid json input: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return errors.New("invalid json input: multiple JSON values provided")
	}
	return nil
}
