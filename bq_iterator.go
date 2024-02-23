package iterfuncs

import (
	"iter"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

// Range return a function that can be used as a range over function.
// It iteratrs through the bigquery iterator and yields results to the caller.
func Range[E any](iter *bigquery.RowIterator) iter.Seq2[*E, error] {
	return func(yield func(*E, error) bool) {
		i := 1
		for {
			var row E
			err := iter.Next(&row)
			if err != nil {
				if err != iterator.Done {
					// call error handler
					if !yield(nil, err) {
						return
					}
				}
				return
			}

			if !yield(&row, nil) {
				return
			}
			i++
		}
	}
}
