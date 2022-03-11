package commands

import (
	"encoding/json"
	"fmt"
	"github.com/k0kubun/pp/v3"
	"io"
	"strings"
)

type outputFormatter interface {
	WriteStringSlice(into io.Writer, data []string) error
	WriteObject(into io.Writer, obj interface{}) error
	WriteString(stdout io.Writer, text string) error
}

type NewlineFormatter struct{}

func (n *NewlineFormatter) WriteString(stdout io.Writer, text string) error {
	_, err := io.WriteString(stdout, text)
	return err
}

func (n *NewlineFormatter) WriteObject(into io.Writer, obj interface{}) error {
	_, err := pp.New().Fprint(into, obj)
	return err
}

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

func (J *JSONFormatter) WriteString(stdout io.Writer, text string) error {
	return J.WriteObject(stdout, text)
}

func (J *JSONFormatter) WriteObject(into io.Writer, obj interface{}) error {
	enc := json.NewEncoder(into)
	enc.SetIndent("", "\t")
	return enc.Encode(obj)
}

func (J *JSONFormatter) WriteStringSlice(into io.Writer, data []string) error {
	return J.WriteObject(into, data)
}

var _ outputFormatter = &NewlineFormatter{}
var _ outputFormatter = &JSONFormatter{}
