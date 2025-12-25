package processor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/snobb/go-poleno/internal/termcolour"
)

var levelColours = map[string]string{
	"error": "red",
	"info":  "white",
	"warn":  "yellow",
	"debug": "grey",
	"trace": "grey",
}

// Processor process inbound data and outputs processed result to the provided writer.
type Processor struct {
	Out          io.Writer
	HeaderFields []string
	LevelField   string
	MsgField     string
}

// Write processes a line of input data and return n of bytes written
// to the out writer or error.
// IN buffer must contain a valid and full json document.
func (p *Processor) Write(in []byte) (n int, err error) {
	var data map[string]interface{}

	if err := json.Unmarshal(in, &data); err != nil {
		return 0, err
	}

	payload, err := p.compile(data)
	if err != nil {
		return 0, err
	}

	if _, err := p.Out.Write(payload); err != nil {
		return len(payload), err
	}

	return len(payload), nil
}

func levelToColour(level string) string {
	colour, ok := levelColours[strings.ToLower(level)]
	if !ok {
		return termcolour.Reset()
	}

	if cHash := termcolour.Lookup(colour); cHash != "" {
		return cHash
	}

	return termcolour.Reset()
}

func (p *Processor) compile(data map[string]interface{}) ([]byte, error) {
	var out bytes.Buffer
	var level string

	level, ok := data[p.LevelField].(string)
	if ok {
		out.WriteString(levelToColour(level))
		delete(data, p.LevelField)
	}

	for _, field := range p.HeaderFields {
		if value, ok := data[field]; ok {
			out.WriteString(fmt.Sprintf("%v ", value))
			delete(data, field)
		}
	}

	if level != "" {
		out.WriteString(fmt.Sprintf("%s ", strings.ToUpper(level)))
		delete(data, p.LevelField)
	}

	if msg, ok := data[p.MsgField]; ok {
		out.WriteString(fmt.Sprintf(":: %s ", msg))
		delete(data, p.MsgField)
	}

	rest, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, err
	}

	out.Write(rest)
	out.WriteString(termcolour.Reset())
	out.WriteString("\n")

	return out.Bytes(), nil
}
