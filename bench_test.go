package main

import (
	"os"
	"path/filepath"
	"testing"
)

var fixtures = []struct {
	path string
	want int
}{
	{"canada.json", 33},
	{"code.json", 3},
	{"citm.json", 1227563},
	{"warandpeace.txt", 649134},
}

func BenchmarkWhitespace(b *testing.B) {
	b.Run("array", withFixtures(b, benchmarkWhitespaceArray))
	b.Run("array (inlined)", withFixtures(b, benchmarkWhitespaceArrayInlined))
	b.Run("check", withFixtures(b, benchmarkWhitespaceCheck))
	b.Run("shift", withFixtures(b, benchmarkWhitespaceShift))
	b.Run("shift (inlined)", withFixtures(b, benchmarkWhitespaceShiftInlined))
	b.Run("switch", withFixtures(b, benchmarkWhitespaceSwitch))
	b.Run("if", withFixtures(b, benchmarkWhitespaceIf))
	b.Run("if (inlined)", withFixtures(b, benchmarkWhitespaceIfInlined))
}

func withFixtures(b *testing.B, fn func(b *testing.B, input []byte, want int)) func(b *testing.B) {
	return func(b *testing.B) {
		for _, fix := range fixtures {
			data, err := os.ReadFile(filepath.Join("testdata", fix.path))
			if err != nil {
				b.Fatal(err)
			}
			b.Run(fix.path, func(b *testing.B) {
				b.SetBytes(int64(len(data)))
				fn(b, data, fix.want)
			})
		}
	}
}

var whitespace = [256]bool{
	' ':  true,
	'\t': true,
	'\n': true,
	'\r': true,
}

func isWhitespaceArray(c byte) bool {
	return whitespace[c]
}

func benchmarkWhitespaceArray(b *testing.B, input []byte, want int) {
	for i := 0; i < b.N; i++ {
		n := 0
		for _, c := range input {
			if isWhitespaceArray(c) {
				n++
			}
		}
		if n != want {
			b.Fatalf("expected: %v, got: %v", want, n)
		}
	}
}

func benchmarkWhitespaceArrayInlined(b *testing.B, input []byte, want int) {
	for i := 0; i < b.N; i++ {
		n := 0
		for _, c := range input {
			if whitespace[c] {
				n++
			}
		}
		if n != want {
			b.Fatalf("expected: %v, got: %v", want, n)
		}
	}
}

func isWhitespaceShift(c byte) bool {
	const whitespace uint64 = 1<<' ' | 1<<'\t' | 1<<'\r' | 1<<'\n'
	return whitespace&(1<<c) > 0
}

func benchmarkWhitespaceShift(b *testing.B, input []byte, want int) {
	for i := 0; i < b.N; i++ {
		n := 0
		for _, c := range input {
			if isWhitespaceShift(c) {
				n++
			}
		}
		if n != want {
			b.Fatalf("expected: %v, got: %v", want, n)
		}
	}
}

func benchmarkWhitespaceShiftInlined(b *testing.B, input []byte, want int) {
	const whitespace uint64 = 1<<' ' | 1<<'\t' | 1<<'\r' | 1<<'\n'
	for i := 0; i < b.N; i++ {
		n := 0
		for _, c := range input {
			if whitespace&(1<<c) > 0 {
				n++
			}
		}
		if n != want {
			b.Fatalf("expected: %v, got: %v", want, n)
		}
	}
}

func isWhitespaceSwitch(c byte) bool {
	switch c {
	case ' ', '\n', '\t', '\r':
		return true
	default:
		return false
	}
}

func benchmarkWhitespaceSwitch(b *testing.B, input []byte, want int) {
	for i := 0; i < b.N; i++ {
		n := 0
		for _, c := range input {
			if isWhitespaceSwitch(c) {
				n++
			}
		}
		if n != want {
			b.Fatalf("expected: %v, got: %v", want, n)
		}
	}
}

func isWhitespaceIf(c byte) bool {
	return c <= ' ' && (c == ' ' || c == '\n' || c == '\t' || c == '\r')
}

func benchmarkWhitespaceIf(b *testing.B, input []byte, want int) {
	for i := 0; i < b.N; i++ {
		n := 0
		for _, c := range input {
			if isWhitespaceIf(c) {
				n++
			}
		}
		if n != want {
			b.Fatalf("expected: %v, got: %v", want, n)
		}
	}
}

func benchmarkWhitespaceIfInlined(b *testing.B, input []byte, want int) {
	for i := 0; i < b.N; i++ {
		n := 0
		for _, c := range input {
			if c <= ' ' && (c == ' ' || c == '\n' || c == '\t' || c == '\r') {
				n++
			}
		}
		if n != want {
			b.Fatalf("expected: %v, got: %v", want, n)
		}
	}
}

func isWhitespaceCheck(c byte) bool {
	return c <= ' ' && whitespace[c]
}

func benchmarkWhitespaceCheck(b *testing.B, input []byte, want int) {
	for i := 0; i < b.N; i++ {
		n := 0
		for _, c := range input {
			if isWhitespaceCheck(c) {
				n++
			}
		}
		if n != want {
			b.Fatalf("expected: %v, got: %v", want, n)
		}
	}
}
