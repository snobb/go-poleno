package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/snobb/go-poleno/pkg/processor"
)

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

			if _, err := out.Write(bytes); err != nil {
				// not json - just ignore the error and print the original line
				fmt.Println(string(bytes))
			}
		}
	}
}
