package textutil

import (
	"image"
	"image/color"
	"math"
	"unicode/utf8"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/hajimehoshi/ebiten/v2"
)

var clock int64

func Update() {
	clock++
}

type Text struct {
	orig       fixed.Point26_6
	dot        fixed.Point26_6
	LineHeight fixed.Int26_6
	TabWidth   fixed.Int26_6
	font       font.Face
	buffer     []byte
	image      *ebiten.Image
	prevRune   rune
	scale      float64
	op         *ebiten.DrawImageOptions
}

func New(font font.Face, w int, h int, x int, y int, scale float64) *Text {
	txt := &Text{
		font:       font,
		prevRune:   rune(-1),
		LineHeight: font.Metrics().Height,
		image:      ebiten.NewImage(w, h),
		scale:      scale,
		op:         &ebiten.DrawImageOptions{},
	}
	tab, _ := font.GlyphAdvance(' ')
	txt.TabWidth = tab * 4
	txt.orig = fixed.Point26_6{X: fixed.I(int(float64(x) / txt.scale)), Y: fixed.I(int(float64(y)/txt.scale)) + txt.font.Metrics().CapHeight}
	cr, cg, cb, ca := color.White.RGBA()
	txt.op.ColorM.Scale(float64(cr)/float64(ca), float64(cg)/float64(ca), float64(cb)/float64(ca), float64(ca)/0xffff)
	txt.op.GeoM.Scale(scale, scale)
	txt.Clear()
	return txt
}

func (txt *Text) Color(color color.Color) {
	txt.op.ColorM.Reset()
	cr, cg, cb, ca := color.RGBA()
	txt.op.ColorM.Scale(float64(cr)/float64(ca), float64(cg)/float64(ca), float64(cb)/float64(ca), float64(ca)/0xffff)
}

func (txt *Text) SetPosition(x int, y int) {
	txt.orig = fixed.Point26_6{X: fixed.I(int(float64(x) / txt.scale)), Y: fixed.I(int(float64(y)/txt.scale)) + txt.font.Metrics().CapHeight}
}

func (txt *Text) SetDot(x int, y int) {
	txt.dot = fixed.Point26_6{X: fixed.I(int(float64(x) / txt.scale)), Y: fixed.I(int(float64(y)/txt.scale)) + txt.font.Metrics().CapHeight}
}

func (txt *Text) Clear() {
	txt.prevRune = -1
	txt.dot = txt.orig
	txt.image.Clear()
}

func (txt *Text) Write(p []byte) (n int, err error) {
	txt.buffer = append(txt.buffer, p...)
	txt.drawBuffer()
	return len(p), nil
}

func (txt *Text) WriteString(s string) (n int, err error) {
	txt.buffer = append(txt.buffer, s...)
	txt.drawBuffer()
	return len(s), nil
}

func (txt *Text) WriteByte(c byte) error {
	txt.buffer = append(txt.buffer, c)
	txt.drawBuffer()
	return nil
}

func (txt *Text) Draw(dst *ebiten.Image) {
	op := ebiten.DrawImageOptions{}
	dst.DrawImage(txt.image, &op)
}

func (txt *Text) drawBuffer() {
	if !utf8.FullRune(txt.buffer) {
		return
	}
	for utf8.FullRune(txt.buffer) {
		r, size := utf8.DecodeRune(txt.buffer)
		txt.buffer = txt.buffer[size:]
		if txt.prevRune >= 0 {
			txt.dot.X += txt.font.Kern(txt.prevRune, r)
		}
		if r == '\n' {
			txt.dot.X = txt.orig.X
			txt.dot.Y += txt.LineHeight
			txt.prevRune = rune(-1)
			continue
		}
		if r == '\t' {
			rem := math.Mod(float64(txt.dot.X-txt.orig.X), float64(txt.TabWidth))
			rem = math.Mod(rem, rem+float64(txt.TabWidth))
			if rem == 0 {
				rem = float64(txt.TabWidth)
			}
			txt.dot.X += fixed.Int26_6(rem)
			continue
		}
		img := getGlyphImage(txt.font, r)
		drawGlyph(txt.image, txt.font, r, img, txt.dot.X, txt.dot.Y, txt.op)
		txt.dot.X += glyphAdvance(txt.font, r)
		txt.prevRune = r
	}

	const cacheSoftLimit = 512

	if len(glyphImageCache[txt.font]) > cacheSoftLimit {
		for r, e := range glyphImageCache[txt.font] {
			if e.atime < clock-60 {
				delete(glyphImageCache[txt.font], r)
			}
		}
	}
}

// Copyright 2017 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

func drawGlyph(dst *ebiten.Image, face font.Face, r rune, img *ebiten.Image, dx, dy fixed.Int26_6, op *ebiten.DrawImageOptions) {
	if img == nil {
		return
	}

	b := getGlyphBounds(face, r)
	op2 := &ebiten.DrawImageOptions{}
	if op != nil {
		*op2 = *op
	}
	op2.GeoM.Reset()
	op2.GeoM.Translate(float64((dx+b.Min.X)>>6), float64((dy+b.Min.Y)>>6))
	if op != nil {
		op2.GeoM.Concat(op.GeoM)
	}
	dst.DrawImage(img, op2)
}

var (
	glyphBoundsCache = map[font.Face]map[rune]fixed.Rectangle26_6{}
)

func getGlyphBounds(face font.Face, r rune) fixed.Rectangle26_6 {
	if _, ok := glyphBoundsCache[face]; !ok {
		glyphBoundsCache[face] = map[rune]fixed.Rectangle26_6{}
	}
	if b, ok := glyphBoundsCache[face][r]; ok {
		return b
	}
	b, _, _ := face.GlyphBounds(r)
	glyphBoundsCache[face][r] = b
	return b
}

type glyphImageCacheEntry struct {
	image *ebiten.Image
	atime int64
}

var (
	glyphImageCache = map[font.Face]map[rune]*glyphImageCacheEntry{}
)

func getGlyphImage(face font.Face, r rune) *ebiten.Image {
	if _, ok := glyphImageCache[face]; !ok {
		glyphImageCache[face] = map[rune]*glyphImageCacheEntry{}
	}

	if e, ok := glyphImageCache[face][r]; ok {
		e.atime = clock
		return e.image
	}

	b := getGlyphBounds(face, r)
	w, h := (b.Max.X - b.Min.X).Ceil(), (b.Max.Y - b.Min.Y).Ceil()
	if w == 0 || h == 0 {
		glyphImageCache[face][r] = &glyphImageCacheEntry{
			image: nil,
			atime: clock,
		}
		return nil
	}

	if b.Min.X&((1<<6)-1) != 0 {
		w++
	}
	if b.Min.Y&((1<<6)-1) != 0 {
		h++
	}
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))

	d := font.Drawer{
		Dst:  rgba,
		Src:  image.White,
		Face: face,
	}
	x, y := -b.Min.X, -b.Min.Y
	x, y = fixed.I(x.Ceil()), fixed.I(y.Ceil())
	d.Dot = fixed.Point26_6{X: x, Y: y}
	d.DrawString(string(r))

	img := ebiten.NewImageFromImage(rgba)
	if _, ok := glyphImageCache[face][r]; !ok {
		glyphImageCache[face][r] = &glyphImageCacheEntry{
			image: img,
			atime: clock,
		}
	}

	return img
}
