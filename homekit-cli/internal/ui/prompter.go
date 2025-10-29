package ui

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Prompter provides minimal user interaction helpers.
type Prompter struct {
	In  io.Reader
	Out io.Writer
}

// Confirm asks the user for yes/no confirmation.
func (p Prompter) Confirm(question string, defaultYes bool) (bool, error) {
	if p.In == nil || p.Out == nil {
		return defaultYes, nil
	}

	reader := bufio.NewReader(p.In)
	def := "y/N"
	if defaultYes {
		def = "Y/n"
	}

	fmt.Fprintf(p.Out, "%s [%s]: ", question, def)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	input = strings.TrimSpace(strings.ToLower(input))
	switch input {
	case "", "y", "yes":
		return true, nil
	case "n", "no":
		return false, nil
	default:
		return defaultYes, nil
	}
}
