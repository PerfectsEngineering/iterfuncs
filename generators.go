package iterfuncs

import "iter"

// Range returns a sequence of integers from start to end with step increment.
func Range(start, end, step int) iter.Seq[int] {
	if step > 0 {
		// Forward iteration from start
		return func(yield func(int) bool) {
			for i := start; i < end; i += step {
				if !yield(i) {
					return
				}
			}
		}
	}

	if step < 0 {
		// Backward iteration from start
		return func(yield func(int) bool) {
			for i := start; i > end; i += step {
				if !yield(i) {
					return
				}
			}
		}
	}

	return func(yield func(int) bool) {
		// yield only start value when step 0
		yield(start)
	}
}

// RangeInfinite returns a sequence of integers from start to infinity with step increment.
func RangeInfinite(start, step int) iter.Seq[int] {
	if step > 0 {
		// Forward iteration from start
		return func(yield func(int) bool) {
			for i := start; ; i += step {
				if !yield(i) {
					return
				}
			}
		}
	}

	if step < 0 {
		// Backward iteration from start
		return func(yield func(int) bool) {
			for i := start; ; i += step {
				if !yield(i) {
					return
				}
			}
		}
	}

	return func(yield func(int) bool) {
		// yield only start value when step 0
		yield(start)
	}
}

// Range2 performs the exact same function as Range but its implementation
// uses the RangeInfinite function and the iter.Pull utility function.
func Range2(start, end, step int) iter.Seq[int] {
	next, stop := iter.Pull(RangeInfinite(start, step))
	return func(yield func(int) bool) {
		defer stop()
		for {
			value, ok := next()
			if !ok {
				return
			}

			if step > 0 && value > end {
				return
			}

			if step < 0 && value < end {
				return
			}

			if !yield(value) {
				return
			}
		}
	}
}
