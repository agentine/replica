package replica

import (
	"testing"
)

type simpleBench struct {
	Name  string
	Value int
	Flag  bool
}

type nestedBench struct {
	Name    string
	Inner   *simpleBench
	Tags    []string
	Meta    map[string]int
	Numbers [10]int
}

func BenchmarkCopySimpleStruct(b *testing.B) {
	orig := simpleBench{Name: "bench", Value: 42, Flag: true}
	b.ResetTimer()
	for range b.N {
		_, _ = Copy(orig)
	}
}

func BenchmarkCopyNestedStruct(b *testing.B) {
	orig := nestedBench{
		Name:    "nested",
		Inner:   &simpleBench{Name: "inner", Value: 99, Flag: false},
		Tags:    []string{"a", "b", "c", "d", "e"},
		Meta:    map[string]int{"x": 1, "y": 2, "z": 3},
		Numbers: [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
	}
	b.ResetTimer()
	for range b.N {
		_, _ = Copy(orig)
	}
}

func BenchmarkCopyLargeSlice(b *testing.B) {
	orig := make([]int, 10000)
	for i := range orig {
		orig[i] = i
	}
	b.ResetTimer()
	for range b.N {
		_, _ = Copy(orig)
	}
}

func BenchmarkCopyMap(b *testing.B) {
	orig := make(map[string]int, 100)
	for i := range 100 {
		orig[string(rune('A'+i%26))+string(rune('0'+i/26))] = i
	}
	b.ResetTimer()
	for range b.N {
		_, _ = Copy(orig)
	}
}
