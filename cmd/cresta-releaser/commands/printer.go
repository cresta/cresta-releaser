package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/k0kubun/pp/v3"
)

type OutputFormatter interface {
	WriteStringSlice(into io.Writer, data []string) error
	WriteObject(into io.Writer, obj interface{}) error
	WriteString(stdout io.Writer, text string) error
}

type NewlineFormatter struct{}

func (n *NewlineFormatter) WriteString(stdout io.Writer, text string) error {
	_, err := io.WriteString(stdout, text)
	if err != nil {
		return err
	}
	if strings.HasPrefix(text, "\n") {
		return nil
	}
	_, err = io.WriteString(stdout, "\n")
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

var _ OutputFormatter = &NewlineFormatter{}
var _ OutputFormatter = &JSONFormatter{}
