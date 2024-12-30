package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
)

type key struct {
	basename string
	hash     [sha256.Size]byte
}

func scan(r io.Reader) func(yield func(t string) bool) {
	s := bufio.NewScanner(r)
	return func(yield func(t string) bool) {
		for {
			if !s.Scan() {
				return
			}
			if !yield(s.Text()) {
				return
			}
		}
	}
}

func warnf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, fmt.Sprintf("%s: %s", format, os.Args[0]), a...)
}

func dief(format string, a ...any) {
	warnf(format, a...)
	os.Exit(1)
}

func main() {
	kong.Parse(&struct{}{})

	uni := map[key]string{}
	sum := sha256.New()
	var hash [sha256.Size]byte
	for path := range scan(os.Stdin) {
		f, err := os.Open(path)
		if err != nil {
			dief("open failed: %s\n", err)
		}

		if _, err := io.Copy(sum, f); err != nil {
			if err := f.Close(); err != nil {
				warnf("close failed: %s", err)
			}
			dief("copy failed: %s\n", err)
		}

		copy(hash[:], sum.Sum(nil))
		if err := f.Close(); err != nil {
			warnf("close failed: %s", err)
		}
		sum.Reset()

		k := key{basename: filepath.Base(path), hash: hash}
		uni[k] = path
	}

	for _, path := range uni {
		fmt.Println(path)
	}
}
