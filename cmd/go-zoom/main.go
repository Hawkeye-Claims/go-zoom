package main

import (
	"fmt"
	"io"
	"os"
)

const (
	envZoomAccountID    = "ZOOM_ACCOUNT_ID"
	envZoomClientID     = "ZOOM_CLIENT_ID"
	envZoomClientSecret = "ZOOM_CLIENT_SECRET"
)

// cli wraps CLI output streams and command handlers.
type cli struct {
	stdout io.Writer
	stderr io.Writer
}

// main is the CLI entrypoint.
func main() {
	c := &cli{stdout: os.Stdout, stderr: os.Stderr}
	if err := c.run(os.Args[1:]); err != nil {
		_, _ = fmt.Fprintf(c.stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
