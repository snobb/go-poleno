package processor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

const (
	fTime  = "time"
	fMsg   = "msg"
	fLevel = "level"
	fHost  = "hostname"
)

var colourMap = map[string]string{
	"reset": "\x1b[0m",
	"white": "\x1b[1m",
	"grey":  "\x1b[2m",

	"black":   "\x1b[30m",
	"red":     "\x1b[31m",
	"green":   "\x1b[32m",
	"yellow":  "\x1b[33m",
	"blue":    "\x1b[34m",
	"magenta": "\x1b[35m",
	"cyan":    "\x1b[36m",
}

var levelColours = map[string]string{
	"error": "red",
	"info":  "white",
	"warn":  "yellow",
	"debug": "grey",
	"trace": "grey",
}

// Processor process inbound data and outputs processed result to the provided writer
type Processor struct {
	out       io.Writer
	fields    map[string]string
	nameField string
}

// New creates a new processor function
func New(out io.Writer, nameField string) *Processor {
	return &Processor{
		out: out,
		fields: map[string]string{
			fTime:     "time",
			fMsg:      "msg",
			fLevel:    "level",
			fHost:     "hostname",
			nameField: nameField,
		},
		nameField: nameField,
	}
}

// Write processes a line of input data and return n of bytes written
// to the out writer or error.
// IN buffer must contain a valid and full json document.
func (p *Processor) Write(in []byte) (n int, err error) {
	var data map[string]interface{}

	if err := json.Unmarshal(in, &data); err != nil {
		return 0, err
	}

	bytes := p.compile(data)
	_, _ = p.out.Write(bytes)

	return len(bytes), nil
}

func levelToColour(level string) string {
	colour, ok := levelColours[level]
	if !ok {
		colour = "reset"
	}

	if cHash, ok := colourMap[colour]; ok {
		return cHash
	}

	return colourMap["reset"]
}

func (p *Processor) compile(data map[string]interface{}) []byte {
	var out bytes.Buffer
	var level string

	level, ok := data[p.fields[fLevel]].(string)
	if ok {
		out.WriteString(levelToColour(level))
		delete(data, p.fields[fLevel])
	}

	for _, field := range []string{fTime, fHost, p.nameField} {
		if value, ok := data[p.fields[field]]; ok {
			out.WriteString(fmt.Sprintf("%s ", value))
			delete(data, p.fields[field])
		}
	}

	if level != "" {
		out.WriteString(fmt.Sprintf("%s ", strings.ToUpper(level)))
	}

	if msg, ok := data[p.fields[fMsg]]; ok {
		out.WriteString(fmt.Sprintf(":: %s ", msg))
		delete(data, p.fields[fMsg])
	}

	rest, _ := json.MarshalIndent(data, "", "  ")

	out.Write(rest)
	out.WriteString(colourMap["reset"])
	out.WriteString("\n")

	return out.Bytes()
}
