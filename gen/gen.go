// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gen contains common code for the various code generation tools in the
// text repository. Its usage ensures consistency between tools.
//
// This package defines command line flags that are common to most generation
// tools. The flags allow for specifying specific Unicode and CLDR versions
// in the public Unicode data repository (https://www.unicode.org/Public).
//
// A local Unicode data mirror can be set through the flag -local or the
// environment variable UNICODE_DIR. The former takes precedence. The local
// directory should follow the same structure as the public repository.
//
// IANA data can also optionally be mirrored by putting it in the iana directory
// rooted at the top of the local mirror. Beware, though, that IANA data is not
// versioned. So it is up to the developer to use the right version.
package gen // import "golang.org/x/text/internal/gen"

import (
	"fmt"
	"go/format"
	"io"
	"log"
	"os"
)

const header = `// Code generated by running "go generate" in github.com/exp625/gones. DO NOT EDIT.

`

// WriteGoFile prepends a standard file comment and package statement to the
// given bytes, applies gofmt, and writes them to a file with the given name.
// It will call log.Fatal if there are any errors.
func WriteGoFile(filename, pkg string, b []byte) {
	w, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Could not create file %s: %v", filename, err)
	}
	defer w.Close()
	if _, err = WriteGo(w, pkg, "", b); err != nil {
		log.Fatalf("Error writing file %s: %v", filename, err)
	}
}

// WriteGo prepends a standard file comment and package statement to the given
// bytes, applies gofmt, and writes them to w.
func WriteGo(w io.Writer, pkg, tags string, b []byte) (n int, err error) {
	src := []byte(header)
	if tags != "" {
		src = append(src, fmt.Sprintf("// +build %s\n\n", tags)...)
	}
	src = append(src, fmt.Sprintf("package %s\n\n", pkg)...)
	src = append(src, b...)
	formatted, err := format.Source(src)
	if err != nil {
		// Print the generated code even in case of an error so that the
		// returned error can be meaningfully interpreted.
		n, _ = w.Write(src)
		return n, err
	}
	return w.Write(formatted)
}
