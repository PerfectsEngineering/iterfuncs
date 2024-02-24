package iterfuncs

import (
	"reflect"
	"testing"
)

func TestRange_And_Range2(t *testing.T) {
	type args struct {
		start int
		end   int
		step  int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Test Range Increments",
			args: args{start: 0, end: 10, step: 2},
			want: []int{0, 2, 4, 6, 8},
		},
		{
			name: "Test Range Decrements",
			args: args{start: 10, end: 0, step: -2},
			want: []int{10, 8, 6, 4, 2},
		},
		{
			name: "Test Range Zero Step",
			args: args{start: 0, end: 10, step: 0},
			want: []int{0},
		},
		{
			name: "Test Out of Range Start",
			args: args{start: 10, end: 0, step: 2},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := []int{}
			result2 := []int{}

			for i := range Range(tt.args.start, tt.args.end, tt.args.step) {
				result = append(result, i)
			}

			for i := range Range2(tt.args.start, tt.args.end, tt.args.step) {
				result2 = append(result2, i)
			}

			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("Range() = %v, want %v", result, tt.want)
			}

			if !reflect.DeepEqual(result2, tt.want) {
				t.Errorf("Range2() = %v, want %v", result2, tt.want)
			}
		})
	}
}

func TestRangeInfinite(t *testing.T) {
	type args struct {
		start int
		step  int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Test Range Infinite Increments",
			args: args{start: 0, step: 2},
			want: []int{0, 2, 4, 6, 8},
		},
		{
			name: "Test Range Infinite Decrements",
			args: args{start: 10, step: -2},
			want: []int{10, 8, 6, 4, 2},
		},
		{
			name: "Test Range Infinite Zero Step",
			args: args{start: 0, step: 0},
			want: []int{0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := []int{}

			for i := range RangeInfinite(tt.args.start, tt.args.step) {
				result = append(result, i)
				if len(result) > 4 {
					break
				}
			}

			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("RangeInfinite() = %v, want %v", result, tt.want)
			}
		})
	}
}
