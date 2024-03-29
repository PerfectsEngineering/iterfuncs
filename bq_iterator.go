package iterfuncs

import (
	"context"
	"iter"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

// BqRange return a function that can be used as a range over function.
// It iteratrs through the bigquery iterator and yields results to the caller.
func BqRange[E any](iter *bigquery.RowIterator) iter.Seq2[*E, error] {
	return func(yield func(*E, error) bool) {
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
		}
	}
}

func BqQueryRange[E any](ctx context.Context, query *bigquery.Query) iter.Seq2[*E, error] {
	return func(yield func(*E, error) bool) {
		// Run the query and get the results
		iter, err := query.Read(ctx)
		if err != nil {
			yield(nil, err)
			return
		}

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
		}
	}
}
