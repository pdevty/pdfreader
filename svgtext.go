package svgtext

import (
  "pdfread";
  "graf";
  "fmt";
  "util";
  "utf8";
  "strm";
  "io";
)

// ------------------------------------------------ Font Substitution

const DEFAULT_FSTYLE = "font-family:Arial;"

func csvtok(d []byte) []byte {
  p := 0;
  for ; p < len(d); p++ {
    if d[p] < 32 {
      break
    }
  }
  return d[0:p];
}

func endcsvl(d []byte) int {
  p := 0;
  for ; p < len(d); p++ {
    if d[p] == 10 {
      break
    }
  }
  return p + 1;
}

var fileNames, // will be needed for SVG font inclusion
  styles map[string]string

func fontnamemap(fn string) int {
  if fileNames == nil {
    fileNames = make(map[string]string)
  }
  if styles == nil {
    styles = make(map[string]string)
  }

  data, _ := io.ReadFile(fn);
  no := 0;
  for p := 0; p < len(data); {
    n := string(csvtok(data[p:len(data)]));
    p += len(n) + 1;
    f := string(csvtok(data[p:len(data)]));
    p += len(f) + 1;
    s := string(csvtok(data[p:len(data)]));
    p += len(s);
    p += endcsvl(data[p:len(data)]);
    fileNames[n] = f;
    styles[n] = s;
    no++;
  }
  return no;
}

var numFonts = fontnamemap("fontnamemap.txt") // initialize fileNames and styles

func FStyle(f string) string {
  if f[0] == '/' { f = f[1:len(f)] }
  if r, ok := styles[f]; ok {
    return r;
  }
  q := 0;
  for ; q < len(f); q++ {
    if f[q] == '+' { break }
  }
  if q < len(f) {
    f = f[q+1: len(f)];
  }
  if r, ok := styles[f]; ok {
    return r;
  }
  return DEFAULT_FSTYLE
}

// ------------------------------------------------


func unquot(a []byte) []byte { // STUB! FIXME!
  if a[0] != '(' {
    return a
  }
  r := make([]byte, len(a)*6); // or so
  p := 0;
  a = a[1 : len(a)-1];
  for k := 0; k < len(a); k++ {
    p += utf8.EncodeRune(int(a[k]), r[p:len(r)])
  }
  return r[0:p]; // removes braces only.
}

const WIDTH_DENSITY = 10000

type SvgTextT struct {
  Pdf    *pdfread.PdfReaderT;
  Drw    *graf.PdfDrawerT;
  Conf   *graf.TextConfigT;
  matrix []string;
  Page   int;
  fonts  pdfread.DictionaryT;
  fontw  map[string][]int64;
  x0, x, y   string;
}

func New() *SvgTextT {
  r := new(SvgTextT);
  r.matrix = []string{"1", "0", "0", "1", "0", "0"};
  return r;
}

func (t *SvgTextT) style(font string) (r string) {
  r = DEFAULT_FSTYLE;
  if t.fonts == nil {
    t.fonts = t.Pdf.PageFonts(t.Pdf.Pages()[t.Page]);
    if t.fonts == nil {
      return
    }
  }
  if dr, ok := t.fonts[font]; ok {
    d := t.Pdf.Dic(dr);
    if fd, ok := d["/FontDescriptor"]; ok { // FIXME: Too simple...
      return FStyle(string(t.Pdf.Dic(fd)["/FontName"]));
    }
  }
  return
}

func (t *SvgTextT) widths(font string) (r []int64) {
  if t.fontw == nil {
    t.fontw = make(map[string][]int64)
  } else if r, ok := t.fontw[font]; ok {
    return r
  }
  r = make([]int64, 256);
  t.fontw[font] = r;
  // initialize like for Courier.
  for k := range r {
    r[k] = 600 * WIDTH_DENSITY / 1000
  }
  if t.fonts == nil {
    t.fonts = t.Pdf.PageFonts(t.Pdf.Pages()[t.Page]);
    if t.fonts == nil {
      return
    }
  }
  if dr, ok := t.fonts[font]; ok {
    d := t.Pdf.Dic(dr);
    fc, ok := d["/FirstChar"];
    if !ok {
      return
    }
    lc, ok := d["/LastChar"];
    if !ok {
      return
    }
    wd, ok := d["/Widths"];
    if !ok {
      return
    }
    p := strm.Int(string(fc), 1);
    q := strm.Int(string(lc), 1);
    a := t.Pdf.Arr(wd);
    for k := p; k < q; k++ {
      r[k] = strm.Int64(string(a[k-p]), WIDTH_DENSITY/1000)
    }
  }
  return;
}

func (t *SvgTextT) Utf8TsAdvance(s []byte) ([]byte, int64) {
  w := t.widths(t.Conf.Font);
  if s[0] != '(' {
    return []byte{}, 0
  }
  z := s[1 : len(s)-1];
  width := int64(0);
  for k := range z {
    width += w[z[k]]
  }
  return unquot(s), width;
}

func (t *SvgTextT) Utf8Advance(s []byte) ([]byte, string) {
  r, a := t.Utf8TsAdvance(s);
  return r, strm.Mul(t.Conf.FontSize, strm.String(a, WIDTH_DENSITY));
}

func (t *SvgTextT) TMoveTo(s [][]byte) {
  t.x0 = strm.Add(t.x0, string(s[0]));
  t.x = t.x0;
  t.y = strm.Add(t.y, string(s[1]));
}

func (t *SvgTextT) TNextLine() {
  t.x = t.x0;
  t.y = strm.Sub(t.y, t.Conf.Leading);
}

func (t *SvgTextT) TSetMatrix(s [][]byte) {
  t.matrix = util.StringArray(s);
  t.x0 = "0";
  t.x = t.x0;
  t.y = "0";
}


func (t *SvgTextT) TShow(a []byte) {
  tx := t.Pdf.ForcedArray(a); // FIXME: Should be "ForcedSimpleArray()"
  for k := range tx {
    if tx[k][0] == '(' {
      tmp, adv := t.Utf8Advance(tx[k]);
      fmt.Printf(
        "<g transform=\"matrix(%s,%s,%s,%s,%s,%s)\">\n"
          "<text x=\"%s\" y=\"%s\""
          " font-size=\"%s\" stroke=\"none\""
          " style=\"%v\""
          " fill=\"black\">%s</text>\n"
          "</g>\n",
        t.matrix[0], t.matrix[1], t.matrix[2], t.matrix[3], t.matrix[4], t.matrix[5],
        t.x, t.y,
        t.Conf.FontSize,
        t.style(t.Conf.Font),
        tmp
);
      t.x = strm.Add(t.x, adv);
    } else {
      t.x = strm.Add(t.x, strm.Mul(string(tx[k]), "0.001"))
    }
  }
}
