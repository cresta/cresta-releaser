package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type outputFormatter interface {
	WriteStringSlice(into io.Writer, data []string) error
}

type NewlineFormatter struct{}

func (n *NewlineFormatter) WriteStringSlice(into io.Writer, data []string) error {
	total := strings.Join(data, "\n")
	total += "\n"
	_, err := into.Write([]byte(total))
	if err != nil {
		return fmt.Errorf("failed to write to output: %s", err)
	}
	return nil
}

type JSONFormatter struct{}

func (J *JSONFormatter) WriteStringSlice(into io.Writer, data []string) error {
	return json.NewEncoder(into).Encode(data)
}

var _ outputFormatter = &NewlineFormatter{}
var _ outputFormatter = &JSONFormatter{}
