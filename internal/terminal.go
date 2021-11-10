package internal

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// TerminalDimensions represents the terminal dimensions
type TerminalDimensions struct {
	width  uint
	height uint
}

// ReadTerminalDimensions detects and returns the terminal dimensions
func ReadTerminalDimensions() (TerminalDimensions, error) {
	var terminal TerminalDimensions
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return terminal, fmt.Errorf("Failed to read terminal dimensions: %w", err)
	}
	dimensions := strings.Split(strings.TrimSuffix(string(out), "\n"), " ")
	if len(dimensions) != 2 {
		return terminal, fmt.Errorf("stty returned an unexpected output: %s", out)
	}

	height, err := strconv.Atoi(dimensions[0])
	if err != nil {
		return terminal, fmt.Errorf("stty returned an unexpected value for height: %s", dimensions[0])
	}
	width, err := strconv.Atoi(dimensions[1])
	if err != nil {
		return terminal, fmt.Errorf("stty returned an unexpected value for width: %s", dimensions[0])
	}

	terminal.width = uint(width)
	terminal.height = uint(height)

	return terminal, nil
}
