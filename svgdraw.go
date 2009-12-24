package svgdraw

/* SVG driver for graf.go.

Copyright (c) 2009 Helmar Wodtke

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

import (
	"fmt"
	"util"
	"strm"
	"graf"
)

type SvgT struct {
	Drw    *graf.PdfDrawerT
	path   []string
	p      int
	groups int
}

func (s *SvgT) append(p string) {
	if s.path == nil {
		s.path = make([]string, 1024)
		s.p = 0
	} else if s.p >= len(s.path) {
		t := make([]string, len(s.path)+1024)
		s.p = 0
		for k := range s.path {
			t[k] = s.path[k]
			s.p++
		}
		s.path = t
	}
	s.path[s.p] = p
	s.p++
}

func (s *SvgT) SvgPath() string {
	if s.path == nil {
		return "path d=\"\""
	}
	return fmt.Sprintf("path d=\"%s\"",
		util.JoinStrings(s.path[0:s.p], ' '))
}

func (s *SvgT) DropPath()             { s.path = nil }
func (s *SvgT) MoveTo(coord [][]byte) { s.append(fmt.Sprintf("M%s %s", coord[0], coord[1])) }
func (s *SvgT) LineTo(coord [][]byte) { s.append(fmt.Sprintf("L%s %s", coord[0], coord[1])) }

func (s *SvgT) CurveTo(coords [][]byte) {
	s.append(fmt.Sprintf("C%s %s %s %s %s %s",
		coords[0], coords[1],
		coords[2], coords[3],
		coords[4], coords[5]))
}

func (s *SvgT) Rectangle(coords [][]byte) {
	s.append(fmt.Sprintf("M%s %s V%s H%s V%s H%s Z",
		coords[0], coords[1],
		strm.Add(string(coords[1]), string(coords[3])),
		strm.Add(string(coords[0]), string(coords[2])),
		coords[1], coords[0]))
}

func (s *SvgT) ClosePath() { s.append("Z") }

func (s *SvgT) Stroke() {
	s.Drw.Write.Out("<%s fill=\"none\" stroke-width=\"%s\" stroke=\"%s\" />\n\n",
		s.SvgPath(), s.Drw.ConfigD.LineWidth, s.Drw.ConfigD.StrokeColor)
}

func (s *SvgT) Fill() {
	fill := s.Drw.ConfigD.FillColor
	if fill == "" {
		fill = "none"
	}
	s.Drw.Write.Out("<%s fill=\"%s\" stroke=\"none\" />\n\n",
		s.SvgPath(), fill)
}

func (s *SvgT) EOFill() { s.Fill() }

func (s *SvgT) FillAndStroke() {
	fill := s.Drw.ConfigD.FillColor
	if fill == "" {
		fill = "none"
	}
	s.Drw.Write.Out("<%s fill=\"%s\" stroke-width=\"%s\" stroke=\"%s\" />\n\n",
		s.SvgPath(), fill, s.Drw.ConfigD.LineWidth, s.Drw.ConfigD.StrokeColor)
}

func (s *SvgT) EOFillAndStroke() { s.FillAndStroke() }
func (s *SvgT) Clip()            {}
func (s *SvgT) EOClip()          {}

func (s *SvgT) Concat(m [][]byte) {
	s.Drw.Write.Out("<g transform=\"matrix(%s,%s,%s,%s,%s,%s)\">\n\n",
		m[0], m[1], m[2], m[3], m[4], m[5])
	s.groups++
}

func (s *SvgT) SetIdentity() {
	for s.groups > 0 {
		s.Drw.Write.Out("</g>\n")
		s.groups--
	}
}

func (s *SvgT) CloseDrawing() { s.SetIdentity() }

func percent(c []byte) []byte { // convert 0..1 color lossless to percent
	r := make([]byte, len(c)+2)
	p := 0
	d := -111
	q := 0
	for p < len(c) {
		if d == p-3 {
			r[q] = '.'
			q++
		}
		if c[p] == '.' {
			d = p
		} else {
			r[q] = c[p]
			q++
		}
		p++
	}
	if d == -111 || d == p-1 {
		r[q] = '0'
		q++
		r[q] = '0'
		q++
	}
	if d == p-2 {
		r[q] = '0'
		q++
	}
	for p = 0; p < q-1 && r[p] == '0'; p++ {
	}
	return r[p:q]
}

func (s *SvgT) Gray(a []byte) string {
	c := percent(a)
	return fmt.Sprintf("rgb(%s%%,%s%%,%s%%)", c, c, c)
}
func (s *SvgT) CMYK(cmyk [][]byte) string {
	return fmt.Sprintf("cmyk(%s%%,%s%%,%s%%,%s%%)",
		percent(cmyk[0]),
		percent(cmyk[1]),
		percent(cmyk[2]),
		percent(cmyk[3]))
}
func (s *SvgT) RGB(rgb [][]byte) string {
	return fmt.Sprintf("rgb(%s%%,%s%%,%s%%)",
		percent(rgb[0]),
		percent(rgb[1]),
		percent(rgb[2]))
}

func NewTestSvg() *graf.PdfDrawerT {
	t := new(SvgT)
	t.Drw = graf.NewPdfDrawer()
	t.Drw.ConfigD.SetColors(t)
	t.Drw.Draw = t
	return t.Drw
}
