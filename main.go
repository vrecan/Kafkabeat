package main

import (
	"os"

	"github.com/vrecan/kafkabeat/cmd"

	_ "github.com/vrecan/kafkabeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
