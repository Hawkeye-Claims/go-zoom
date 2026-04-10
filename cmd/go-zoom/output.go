package main

import (
	"encoding/json"
	"fmt"
)

// writeJSON writes pretty-printed JSON to stdout.
func (c *cli) writeJSON(v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal output: %w", err)
	}
	_, _ = fmt.Fprintln(c.stdout, string(b))
	return nil
}
