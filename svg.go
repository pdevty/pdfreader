package svg

import (
  "fmt";
  "util";
)

type SvgT struct {
  currentPoint [][]byte;
  firstPoint   [][]byte;
  path         []string;
  p            int;
  groups       int;
  strokeColor  string;
  fillColor    string;
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
  s.append(fmt.Sprintf("M%s %s", coord[0], coord[1]));
}

func (s *SvgT) LineTo(coord [][]byte) {
  s.currentPoint = coord;
  s.append(fmt.Sprintf("L%s %s", coord[0], coord[1]));
}

func (s *SvgT) CurveTo(coords [][]byte) {
  s.currentPoint = coords[4:5];
  s.append(fmt.Sprintf("C%s %s %s %s %s %s",
    coords[0], coords[1],
    coords[2], coords[3],
    coords[4], coords[5]));
}

func (s *SvgT) Rectangle(coords [][]byte) {}

func (s *SvgT) ClosePath() { s.append("Z") }

func (s *SvgT) Stroke() {
  fmt.Printf("<path d=\"%s\" fill=\"none\" stroke-width=\"1\" stroke=\"%s\" />\n\n", s.SvgPath(), s.strokeColor);
  s.path = nil;
}

func (s *SvgT) Fill() {
  fmt.Printf("<path d=\"%s\" fill=\"%s\" stroke=\"none\" />\n\n", s.SvgPath(), s.fillColor);
  s.path = nil;
}

func (s *SvgT) EOFill() {
  fmt.Printf("<path d=\"%s\" />\n\n", s.SvgPath());
  s.path = nil;
}

func (s *SvgT) FillAndStroke() {
  fmt.Printf("<path d=\"%s\" />\n\n", s.SvgPath());
  s.path = nil;
}

func (s *SvgT) EOFillAndStroke() {
  fmt.Printf("<path d=\"%s\" />\n\n", s.SvgPath());
  s.path = nil;
}

func (s *SvgT) Clip() {}

func (s *SvgT) EOClip() {}

func (s *SvgT) Concat(m [][]byte) {
  fmt.Printf("<g transform=\"matrix(%s,%s,%s,%s,%s,%s)\">\n\n",
    m[0], m[1], m[2], m[3], m[4], m[5]);
  s.groups++;
}

func (s *SvgT) CloseDrawing() {
  for s.groups > 0 {
    fmt.Printf("</g>\n");
    s.groups--;
  }
}

func NewDrawer() *SvgT { return new(SvgT) }


func percent(c []byte) []byte { // convert 0..1 color lossless to percent
  r := make([]byte, len(c)+2);
  p := 0;
  d := -111;
  q := 0;
  for p < len(c) {
    if d == p-3 {
      r[q] = '.';
      q++;
    }
    if c[p] == '.' {
      d = p
    } else {
      r[q] = c[p];
      q++;
    }
    p++;
  }
  if d == -111 || d == p-1 {
    r[q] = '0';
    q++;
    r[q] = '0';
    q++;
  }
  if d == p-2 {
    r[q] = '0';
    q++;
  }
  for p = 0; p < q-1 && r[p] == '0'; p++ {
  }
  return r[p:q];
}

func (s *SvgT) SetGrayStroke(a []byte)      {}
func (s *SvgT) SetGrayFill(a []byte)        {}
func (s *SvgT) SetCMYKStroke(cmyk [][]byte) {}
func (s *SvgT) SetCMYKFill(cmyk [][]byte)   {}
func (s *SvgT) SetRGBStroke(rgb [][]byte) {
  s.strokeColor = fmt.Sprintf("rgb(%s%%,%s%%,%s%%)",
    percent(rgb[0]),
    percent(rgb[1]),
    percent(rgb[2]))
}
func (s *SvgT) SetRGBFill(rgb [][]byte) {
  s.fillColor = fmt.Sprintf("rgb(%s%%,%s%%,%s%%)",
    percent(rgb[0]),
    percent(rgb[1]),
    percent(rgb[2]))
}
