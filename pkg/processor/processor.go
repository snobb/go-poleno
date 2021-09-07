package processor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

const (
	fTime  = "time"
	fMsg   = "msg"
	fName  = "name"
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

// StdoutProcessor process inbound data and outputs processed result to stdout
type StdoutProcessor struct {
	fields map[string]string
}

// New creates a new processor function
func New() *StdoutProcessor {
	return &StdoutProcessor{
		fields: map[string]string{
			fTime:  "time",
			fMsg:   "msg",
			fName:  "name",
			fLevel: "level",
			fHost:  "hostname",
		},
	}
}

// Process processes a line of input data.
func (s *StdoutProcessor) Process(in []byte) error {
	var out map[string]interface{}

	if err := json.Unmarshal(in, &out); err != nil {
		return err
	}

	fmt.Println(s.compile(out))

	return nil
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

func (s *StdoutProcessor) compile(data map[string]interface{}) string {
	var out bytes.Buffer
	var level string

	level, ok := data[s.fields[fLevel]].(string)
	if ok {
		out.WriteString(levelToColour(level))
		delete(data, s.fields[fLevel])
	}

	for _, field := range []string{fTime, fHost, fName} {
		if value, ok := data[s.fields[field]]; ok {
			out.WriteString(fmt.Sprintf("%s ", value))
			delete(data, s.fields[field])
		}
	}

	if level != "" {
		out.WriteString(fmt.Sprintf("%s ", strings.ToUpper(level)))
	}

	if msg, ok := data[s.fields[fMsg]]; ok {
		out.WriteString(fmt.Sprintf(":: %s ", msg))
		delete(data, s.fields[fMsg])
	}

	rest, _ := json.MarshalIndent(data, "", "  ")

	out.Write(rest)
	out.WriteString(colourMap["reset"])

	return out.String()
}
