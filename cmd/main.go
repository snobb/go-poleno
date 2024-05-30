package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/snobb/go-poleno/pkg/processor"
)

var version string

func main() {
	var fields, level, msg string
	var ver bool

	flag.StringVar(&fields, "f", "time,hostname,name", "comma separated list of fields to show in the header")
	flag.StringVar(&level, "l", "level", "name of the log level field")
	flag.StringVar(&msg, "m", "msg", "name of the message field")
	flag.BoolVar(&ver, "v", false, "show version")
	flag.Parse()

	if ver {
		fmt.Println(version)
		return
	}

	in := bufio.NewScanner(os.Stdin)

	out := &processor.Processor{
		Out:          os.Stdout,
		HeaderFields: split(fields),
		LevelField:   level,
		MsgField:     msg,
	}
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs,
		syscall.SIGINT, syscall.SIGTERM,
		syscall.SIGQUIT, os.Interrupt)

	for in.Scan() {
		select {
		case <-sigs:
			return

		default:
			bytes := in.Bytes()
			if len(bytes) == 0 {
				return
			}

			if _, err := out.Write(bytes); err != nil {
				// not json - just ignore the error and print the original line
				fmt.Println(string(bytes))
			}
		}
	}
}

func split(line string) []string {
	re := regexp.MustCompile(`\s*,\s*`)
	return re.Split(strings.TrimSpace(line), -1)
}
