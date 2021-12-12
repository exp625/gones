// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gen

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"hash"
	"hash/fnv"
	"io"
	"log"
	"os"
	"strings"
)

// This file contains utilities for generating code.

// TODO: other write methods like:
// - slices, maps, types, etc.

// CodeWriter is a utility for writing structured code. It computes the content
// hash and size of written content. It ensures there are newlines between
// written code blocks.
type CodeWriter struct {
	buf  bytes.Buffer
	Size int
	Hash hash.Hash32 // content hash
	gob  *gob.Encoder
	// For comments we skip the usual one-line separator if they are followed by
	// a code block.
	skipSep bool
}

func (w *CodeWriter) Write(p []byte) (n int, err error) {
	return w.buf.Write(p)
}

// NewCodeWriter returns a new CodeWriter.
func NewCodeWriter() *CodeWriter {
	h := fnv.New32()
	return &CodeWriter{Hash: h, gob: gob.NewEncoder(h)}
}

// WriteGoFile appends the buffer with the total size of all created structures
// and writes it as a Go file to the given file with the given package name.
func (w *CodeWriter) WriteGoFile(filename, pkg string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Could not create file %s: %v", filename, err)
	}
	defer f.Close()
	if _, err = w.WriteGo(f, pkg, ""); err != nil {
		log.Fatalf("Error writing file %s: %v", filename, err)
	}
}

// WriteGo appends the buffer with the total size of all created structures and
// writes it as a Go file to the given writer with the given package name.
func (w *CodeWriter) WriteGo(out io.Writer, pkg, tags string) (n int, err error) {
	sz := w.Size
	if sz > 0 {
		w.WriteComment("Total table size %d bytes (%dKiB); checksum: %X\n", sz, sz/1024, w.Hash.Sum32())
	}
	defer w.buf.Reset()
	return WriteGo(out, pkg, tags, w.buf.Bytes())
}

func (w *CodeWriter) printf(f string, x ...interface{}) {
	fmt.Fprintf(w, f, x...)
}

// WriteComment writes a comment block. All line starts are prefixed with "//".
// Initial empty lines are gobbled. The indentation for the first line is
// stripped from consecutive lines.
func (w *CodeWriter) WriteComment(comment string, args ...interface{}) {
	s := fmt.Sprintf(comment, args...)
	s = strings.Trim(s, "\n")

	// Use at least two newlines to ensure a blank space between the previous
	// block. WriteGoFile will remove extraneous newlines.
	w.printf("\n\n// ")
	w.skipSep = true

	// strip first indent level.
	sep := "\n"
	for ; len(s) > 0 && (s[0] == '\t' || s[0] == ' '); s = s[1:] {
		sep += s[:1]
	}

	strings.NewReplacer(sep, "\n// ", "\n", "\n// ").WriteString(w, s)

	w.printf("\n")
}
