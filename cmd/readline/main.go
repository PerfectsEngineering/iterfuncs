package main

import (
	"fmt"
	"io"
	"path"

	"github.com/perfectsengineering/iterfuncs"
	"github.com/rs/zerolog/log"
)

func main() {
	filepath := path.Join(".", "testdata", "file.txt")

	for line, err := range iterfuncs.ReadFsLines(filepath) {
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		log.Print(string(line))
	}

	for file, err := range iterfuncs.WithFile(filepath) {
		if err != nil {
			log.Fatal().Err(err).Send()
		}

		result, err := io.ReadAll(file)
		if err != nil {
			log.Fatal().Err(err).Send()
		}

		log.Debug().Msg(fmt.Sprint("All file contents: ", string(result)))
	}
}
