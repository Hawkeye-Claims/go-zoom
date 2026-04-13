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
	return readJSONInputWithFlags[T](jsonInput, jsonFile, "--json", "--json-file", true)
}

// readOptionalJSONInput parses JSON from inline or file input when provided.
func readOptionalJSONInput[T any](jsonInput string, jsonFile string) (*T, error) {
	return readOptionalJSONInputWithFlags[T](jsonInput, jsonFile, "--json", "--json-file")
}

// readOptionalJSONInputWithFlags parses JSON from inline or file input when provided.
func readOptionalJSONInputWithFlags[T any](jsonInput string, jsonFile string, inlineFlag string, fileFlag string) (*T, error) {
	if jsonInput == "" && jsonFile == "" {
		return nil, nil
	}
	payload, err := readJSONInputWithFlags[T](jsonInput, jsonFile, inlineFlag, fileFlag, false)
	if err != nil {
		return nil, err
	}
	return &payload, nil
}

func readJSONInputWithFlags[T any](jsonInput string, jsonFile string, inlineFlag string, fileFlag string, required bool) (T, error) {
	var payload T
	if jsonInput == "" && jsonFile == "" {
		if required {
			return payload, fmt.Errorf("one of %s or %s is required", inlineFlag, fileFlag)
		}
		return payload, nil
	}
	if jsonInput != "" && jsonFile != "" {
		return payload, fmt.Errorf("%s and %s cannot be used together", inlineFlag, fileFlag)
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
