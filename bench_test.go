package main

import (
	"io/ioutil"
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

func BenchmarkWhitespaceArray(b *testing.B) {
	withFixtures(b, benchmarkWhitespaceArray)
}

func BenchmarkWhitespaceArrayInlined(b *testing.B) {
	withFixtures(b, benchmarkWhitespaceArrayInlined)
}

func BenchmarkWhitespaceShift(b *testing.B) {
	withFixtures(b, benchmarkWhitespaceShift)
}

func BenchmarkWhitespaceShiftInlined(b *testing.B) {
	withFixtures(b, benchmarkWhitespaceShiftInlined)
}

func BenchmarkWhitespaceSwitch(b *testing.B) {
	withFixtures(b, benchmarkWhitespaceSwitch)
}

func BenchmarkWhitespaceIf(b *testing.B) {
	withFixtures(b, benchmarkWhitespaceIf)
}

func BenchmarkWhitespaceIfInlined(b *testing.B) {
	withFixtures(b, benchmarkWhitespaceIfInlined)
}

func withFixtures(b *testing.B, fn func(b *testing.B, input []byte, want int)) {
	for _, fix := range fixtures {
		data, err := ioutil.ReadFile(filepath.Join("testdata", fix.path))
		if err != nil {
			b.Fatal(err)
		}
		b.Run(fix.path, func(b *testing.B) {
			b.SetBytes(int64(len(data)))
			fn(b, data, fix.want)
		})
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
