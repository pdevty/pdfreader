package svg

import (
  "fmt";
  "graf";
  "util";
)

type SvgT struct {
  Parent       graf.PdfDrawerT;
  currentPoint [][]byte;
  firstPoint   [][]byte;
  path         []string;
  p            int;
}

func (s *SvgT) append(p string) {
  if s.path == nil {
    s.path = make([]string, 1024);
    s.p = 0;
  } else if s.p >= len(s.path) {
    t := make([]string, len(s.path)+1024);
    s.p = 0;
    for k := range s.path {
      t[k] = s.path[k];
      s.p++;
    }
    s.path = t;
  }
  s.path[s.p] = p;
  s.p++;
}

func (s *SvgT) SvgPath() []byte {
  if s.path == nil {
    return []byte{}
  }
  return util.JoinStrings(s.path[0:s.p], ' ');
}

func (s *SvgT) DropPath() { s.path = nil }

func (s *SvgT) CurrentPoint() [][]byte { return s.currentPoint }

func (s *SvgT) MoveTo(coord [][]byte) {
  s.currentPoint = coord;
  s.firstPoint = coord;
  s.append(fmt.Sprintf("M %s,%s", coord[0], coord[1]));
}

func (s *SvgT) LineTo(coord [][]byte) {
  s.currentPoint = coord;
  s.append(fmt.Sprintf("L %s,%s", coord[0], coord[1]));
}

func (s *SvgT) CurveTo(coords [][]byte) {
  s.currentPoint = coords[4:5];
  s.append(fmt.Sprintf("C %s,%s %s,%s %s,%s",
    coords[0], coords[1],
    coords[2], coords[3],
    coords[4], coords[5]));
}

func (s *SvgT) Rectangle(coords [][]byte) {}

func (s *SvgT) ClosePath() {}

func (s *SvgT) Stroke() {}

func (s *SvgT) Fill() {}

func (s *SvgT) EOFill() {}

func (s *SvgT) FillAndStroke() {}

func (s *SvgT) EOFillAndStroke() {}

func (s *SvgT) Clip() {}

func (s *SvgT) EOClip() {}

func (s *SvgT) Concat(matrix [][]byte) {}
