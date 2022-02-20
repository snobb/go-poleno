package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/snobb/go-poleno/pkg/processor"
)

var version string

func main() {
	var name string
	var ver bool
	flag.StringVar(&name, "n", "name", "field to show in the header")
	flag.BoolVar(&ver, "v", false, "show version")
	flag.Parse()

	if ver {
		fmt.Println(version)
		return
	}

	in := bufio.NewScanner(os.Stdin)
	out := processor.New(os.Stdout, name)
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

			if _, err := out.Write(bytes); err != nil {
				// not json - just ignore the error and print the original line
				fmt.Println(string(bytes))
			}
		}
	}
}
