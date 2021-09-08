package main

import (
	"bufio"
	"os"

	"github.com/snobb/go-poleno/pkg/processor"
)

func main() {
	in := bufio.NewScanner(os.Stdin)
	out := processor.New(os.Stdout)

	for in.Scan() {
		bytes := in.Bytes()
		if len(bytes) == 0 {
			break
		}

		if _, err := out.Process(bytes); err != nil {
			panic(err)
		}
	}
}
