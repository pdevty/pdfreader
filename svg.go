package svg

import (
  "fmt";
  "util";
  "strconv";
  "graf";
  "grafgs";
)

type SvgT struct {
  conf *grafgs.GrafGsT;
  currentPoint [][]byte;
  firstPoint   [][]byte;
  path         []string;
  p            int;
  groups       int;
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

func (s *SvgT) SvgPath() string {
  if s.path == nil {
    return "path d=\"\""
  }
  return fmt.Sprintf("path d=\"%s\"",
    util.JoinStrings(s.path[0:s.p], ' '));
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

func (s *SvgT) Rectangle(coords [][]byte) {

  x, _ := strconv.Atof(string(coords[0]));
  y, _ := strconv.Atof(string(coords[1]));
  w, _ := strconv.Atof(string(coords[2]));
  h, _ := strconv.Atof(string(coords[3]));

  s.currentPoint = coords[0:1];
  s.append(fmt.Sprintf("M%s %s V%f H%f V%s H%s Z",
    coords[0], coords[1],
    y+h, x+w,
    coords[1], coords[0]));
}

func (s *SvgT) ClosePath() { s.append("Z") }

func (s *SvgT) Stroke() {
  fmt.Printf("<%s fill=\"none\" stroke-width=\"%s\" stroke=\"%s\" />\n\n", s.SvgPath(), s.conf.LineWidth, s.conf.StrokeColor);
  s.DropPath();
}

func (s *SvgT) Fill() {
  fmt.Printf("<%s fill=\"%s\" stroke=\"none\" />\n\n", s.SvgPath(), s.conf.FillColor);
  s.DropPath();
}

func (s *SvgT) EOFill() {
  fmt.Printf("<%s />\n\n", s.SvgPath());
  s.DropPath();
}

func (s *SvgT) FillAndStroke() {
  fmt.Printf("<%s />\n\n", s.SvgPath());
  s.DropPath();
}

func (s *SvgT) EOFillAndStroke() {
  fmt.Printf("<%s />\n\n", s.SvgPath());
  s.DropPath();
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

func (s *SvgT) SetIdentity() { s.CloseDrawing() }

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

func (s *SvgT) Gray(a []byte) string {
  c := percent(a);
  return fmt.Sprintf("rgb(%s%%,%s%%,%s%%)", c, c, c);
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
  r := new(graf.PdfDrawerT);
  r.Stack = graf.NewStack(10240);
  t := NewDrawer();
  r.Draw = t;
  t.conf = grafgs.New(t);
  r.Config = t.conf;
  r.Ops = make(map[string]func(pd *graf.PdfDrawerT));
  for k := range graf.PdfOps {
    r.Ops[k] = graf.PdfOps[k]
  }
  return r;
}
