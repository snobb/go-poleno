package main

import (
	"bufio"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/snobb/go-poleno/pkg/processor"
)

func wrapError(err error, line []byte) []byte {
	msg, _ := json.Marshal(map[string]interface{}{
		"time":  time.Now().Format(time.RFC3339),
		"level": "error",
		"name":  "poleno:internal",
		"line":  string(line),
		"error": err.Error(),
	})

	return msg
}

func main() {
	in := bufio.NewScanner(os.Stdin)
	out := processor.New(os.Stdout)
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs,
		syscall.SIGINT, syscall.SIGTERM,
		syscall.SIGQUIT, os.Interrupt)

	for in.Scan() {
		select {
		case <-sigs:
			break

		default:
			bytes := in.Bytes()
			if len(bytes) == 0 {
				break
			}

			if _, err := out.Process(bytes); err != nil {
				if _, err = out.Process(wrapError(err, bytes)); err != nil {
					panic(err)
				}
			}
		}
	}
}
