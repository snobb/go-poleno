package processor

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// interesting fields
const (
	Msg = iota
	Time
	Level
	Name
	Hostname
)

const (
	cGrey   = "\x1b[2m"
	cRed    = "\x1b[31m"
	cWhite  = "\x1b[1m"
	cYellow = "\x1b[33m"
	cReset  = "\x1b[m"
)

// StdoutProcessor process inbound data and outputs processed result to stdout
type StdoutProcessor struct {
	fields map[int]string
}

// New creates a new processor function
func New() *StdoutProcessor {
	return &StdoutProcessor{
		fields: map[int]string{
			Time:     "time",
			Msg:      "msg",
			Name:     "name",
			Level:    "level",
			Hostname: "hostname",
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

func (s *StdoutProcessor) compile(data map[string]interface{}) string {
	var out bytes.Buffer
	var colour string

	if level, ok := data[s.fields[Level]]; ok {
		switch level {
		case "info":
			colour = cWhite

		case "error":
			colour = cRed

		case "warn":
			colour = cYellow

		case "debug":
			colour = cGrey

		case "trace":
			colour = cGrey
		}
		delete(data, s.fields[Level])
		out.WriteString(colour)
	}

	for _, field := range []int{Time, Hostname, Name} {
		if value, ok := data[s.fields[field]]; ok {
			out.WriteString(fmt.Sprintf("%s ", value))
			delete(data, s.fields[field])
		}
	}

	if msg, ok := data[s.fields[Msg]]; ok {
		out.WriteString(fmt.Sprintf(":: %s ", msg))
		delete(data, s.fields[Msg])
	}

	rest, _ := json.MarshalIndent(data, "", "  ")

	out.Write(rest)
	out.WriteString(cReset)

	return out.String()
}
