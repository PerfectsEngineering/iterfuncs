package iterfuncs

import (
	"bufio"
	"iter"
	"os"

	"github.com/rs/zerolog/log"
)

func ReadFsLines(filepath string) iter.Seq2[[]byte, error] {
	return func(yield func([]byte, error) bool) {
		file, err := os.Open(filepath)
		if err != nil {
			yield(nil, err)
			return
		}
		defer func() {
			log.Debug().Msg("closing file")
			file.Close()
		}()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if !yield(scanner.Bytes(), nil) {
				return
			}
		}

		if err := scanner.Err(); err != nil {
			if !yield(nil, err) {
				return
			}
		}
	}
}

// WithFile returns a sequence that yeilds a file
// (and an error if there was an error opening the file).
func WithFile(filepath string) iter.Seq2[*os.File, error] {
	return func(yield func(*os.File, error) bool) {
		file, err := os.Open(filepath)
		if err != nil {
			yield(nil, err)
			return
		}
		// ensure the file is always closed
		defer file.Close()

		yield(file, nil)
	}
}
