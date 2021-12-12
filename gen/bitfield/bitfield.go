// Copyright 2018 The Go Authors. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
// * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
// * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
// * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// Package bitfield converts annotated structs into integer values.
//
// Any field that is marked with a bitfield tag is compacted. The tag value has
// two parts. The part before the comma determines the method name for a
// generated type. If left blank the name of the field is used.
// The part after the comma determines the number of bits to use for the
// representation.
package bitfield

import (
	"bytes"
	"fmt"
	"github.com/exp625/gones/internal/plz"
	"io"
	"reflect"
	"strconv"
	"strings"
)

// Config determines settings for packing and generation. If a Config is used,
// the same Config should be used for packing and generation.
type Config struct {
	// NumBits fixes the maximum allowed bits for the integer representation.
	// If NumBits is not 8, 16, 32, or 64, the actual underlying integer size
	// will be the next largest available.
	NumBits uint

	// If Package is set, code generation will write a package clause.
	Package string

	// TypeName is the name for the generated type. By default it is the name
	// of the type of the value passed to Gen.
	TypeName string
}

var nullConfig = &Config{}

func pack(x interface{}, c *Config) (packed uint64, nBit uint, err error) {
	if c == nil {
		c = nullConfig
	}
	nBits := c.NumBits
	v := reflect.ValueOf(x)
	v = reflect.Indirect(v)
	t := v.Type()
	pos := 64 - nBits
	if nBits == 0 {
		pos = 0
	}
	for i := 0; i < v.NumField(); i++ {
		v := v.Field(i)
		field := t.Field(i)
		f, err := parseField(field)

		if err != nil {
			return 0, 0, err
		}
		if f.nBits == 0 {
			continue
		}
		value := uint64(0)
		switch v.Kind() {
		case reflect.Bool:
			if v.Bool() {
				value = 1
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			value = v.Uint()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			x := v.Int()
			if x < 0 {
				return 0, 0, fmt.Errorf("bitfield: negative value for field %q not allowed", field.Name)
			}
			value = uint64(x)
		}
		if value > (1<<f.nBits)-1 {
			return 0, 0, fmt.Errorf("bitfield: value %#x of field %q does not fit in %d bits", value, field.Name, f.nBits)
		}
		shift := 64 - pos - f.nBits
		if pos += f.nBits; pos > 64 {
			return 0, 0, fmt.Errorf("bitfield: no more bits left for field %q", field.Name)
		}
		packed |= value << shift
	}
	if nBits == 0 {
		nBits = posToBits(pos)
		packed >>= 64 - nBits
	}
	return packed, nBits, nil
}

type field struct {
	name  string
	value uint64
	nBits uint
}

// parseField parses a tag of the form [<name>][:<nBits>][,<pos>[..<end>]]
func parseField(field reflect.StructField) (f field, err error) {
	s, ok := field.Tag.Lookup("bitfield")
	if !ok {
		return f, nil
	}
	switch field.Type.Kind() {
	case reflect.Bool:
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	default:
		return f, fmt.Errorf("bitfield: field %q is not an integer or bool type", field.Name)
	}
	bits := s
	f.name = ""

	if i := strings.IndexByte(s, ','); i >= 0 {
		bits = s[:i]
		f.name = s[i+1:]
	}
	if bits != "" {
		nBits, err := strconv.ParseUint(bits, 10, 8)
		if err != nil {
			return f, fmt.Errorf("bitfield: invalid bit size for field %q: %v", field.Name, err)
		}
		f.nBits = uint(nBits)
	}
	if f.nBits == 0 {
		if field.Type.Kind() == reflect.Bool {
			f.nBits = 1
		} else {
			f.nBits = uint(field.Type.Bits())
		}
	}
	if f.name == "" {
		f.name = field.Name
	}
	return f, err
}

func posToBits(pos uint) (bits uint) {
	switch {
	case pos <= 8:
		bits = 8
	case pos <= 16:
		bits = 16
	case pos <= 32:
		bits = 32
	case pos <= 64:
		bits = 64
	default:
		panic("unreachable")
	}
	return bits
}

// Gen generates code for unpacking integers created with Pack.
func Gen(w io.Writer, x interface{}, structName string, c *Config) error {
	if c == nil {
		c = nullConfig
	}
	_, nBits, err := pack(x, c)
	if err != nil {
		return err
	}

	t := reflect.TypeOf(x)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	firstChar := []rune(structName)[0]

	buf := &bytes.Buffer{}

	pos := uint(0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		f, _ := parseField(field)
		if f.nBits == 0 {
			continue
		}
		shift := nBits - pos - f.nBits
		pos += f.nBits

		retType := field.Type.Name()
		if f.name == "_" {
			continue
		}
		plz.Just(fmt.Fprintf(buf, "\nfunc (%c *%s) %s() %s {\n", firstChar, structName, f.name, retType))
		if field.Type.Kind() == reflect.Bool {
			plz.Just(fmt.Fprintf(buf, "\tconst bit = 1 << %d\n", shift))
			plz.Just(fmt.Fprintf(buf, "\treturn *%c&bit == bit\n", firstChar))
		} else {
			plz.Just(fmt.Fprintf(buf, "\treturn %s((*%c >> %d) & %#x)\n", retType, firstChar, shift, (1<<f.nBits)-1))
		}
		plz.Just(fmt.Fprintf(buf, "}\n"))

		if field.Type.Kind() == reflect.Bool {
			plz.Just(fmt.Fprintf(buf, "\nfunc (%c *%s) Set%s(value bool) {\n", firstChar, structName, f.name))
			plz.Just(fmt.Fprintf(buf, "\tconst bit = uint%d(1) << %d\n", nBits, shift))
			plz.Just(fmt.Fprintf(buf, "\tvalueInt := uint%d(0)\n", nBits))
			plz.Just(fmt.Fprintf(buf, "\tif value {\n"))
			plz.Just(fmt.Fprintf(buf, "\t\tvalueInt = 1\n"))
			plz.Just(fmt.Fprintf(buf, "\t}\n"))
			plz.Just(fmt.Fprintf(buf, "\n*%c = %s((uint%d(*%c) & ^bit) | valueInt << %d)\n", firstChar, structName, nBits, firstChar, shift))
		} else {
			plz.Just(fmt.Fprintf(buf, "\nfunc (%c *%s) Set%s(value %s) {\n", firstChar, structName, f.name, field.Type.String()))
			plz.Just(fmt.Fprintf(buf, "\tconst mask = uint%d(((1 << %d) - 1) << %d)\n", nBits, f.nBits, shift))
			plz.Just(fmt.Fprintf(buf, "\t*%c = %s((uint%d(*%c) & ^mask) | uint%d(value) << %d)\n", firstChar, structName, nBits, firstChar, nBits, shift))
		}
		plz.Just(fmt.Fprintf(buf, "}\n"))

	}

	if c.Package != "" {
		plz.Just(fmt.Fprintf(w, "// Code generated by golang.org/x/text/internal/gen/bitfield. DO NOT EDIT.\n"))
		plz.Just(fmt.Fprintf(w, "package %s\n", c.Package))
	}

	bits := posToBits(pos)
	plz.Just(fmt.Fprintf(w, "type %s uint%d", structName, bits))

	if _, err := io.Copy(w, buf); err != nil {
		return fmt.Errorf("bitfield: write failed: %v", err)
	}
	return nil
}
