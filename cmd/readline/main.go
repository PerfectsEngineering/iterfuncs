package main

import (
	"path"

	"github.com/perfectsengineering/iterfuncs"
	"github.com/rs/zerolog/log"
)

func main() {
	filepath := path.Join(".", "testdata", "file.txt")

	for line, err := range iterfuncs.ReadLines(filepath) {
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		log.Print(string(line))
	}
}
