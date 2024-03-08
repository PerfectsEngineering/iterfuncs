package iterfuncs

import (
	"bufio"
	"iter"
	"os"
)

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

func ReadFsLines(filepath string) iter.Seq2[[]byte, error] {
	pullFile, stop := iter.Pull2(WithFile(filepath))
	return func(yield func([]byte, error) bool) {
		defer stop()
		file, err, ok := pullFile()
		if !ok {
			return
		}
		if err != nil {
			if !yield(nil, err) {
				return
			}
		}

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
