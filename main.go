package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/jarl-tornroos/cloudfrontbeat/beater"
)

func main() {
	err := beat.Run("cloudfrontbeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
