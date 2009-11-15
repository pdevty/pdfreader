package pdfread

import (
  "io";
  "regexp";
  "strconv";
)

// limits
const MAX_PDF_UPDATES = 1024
const MAX_PDF_STRING = 1024
const MAX_PDF_DICT = 1024 * 16
const MAX_PDF_ARRAYSIZE = 1024

// types
type PdfReaderT struct {
  File      string;            // name of the file
  Bin       []byte;            // contents of the file
  Startxref int;               // starting of xref table
  Xref      map[int]int;       // "pointers" of the xref table
  Trailer   map[string][]byte; // trailer dictionary of the file
  rcache    map[string][]byte; // resolver cache
  rncache   map[string]int;    // resolver cache (positions in file)
  pages     [][]byte;          // pages cache
}

var _Bytes = []byte{}

func max(a, b int) int {
  if a < b {
    return b
  }
  return a;
}
func min(a, b int) int {
  if a < b {
    return a
  }
  return b;
}
func end(a []byte, n int) int { return max(0, len(a)-n) }
func num(n string) int {
  if i, e := strconv.Atoi(n); e == nil {
    return i
  }
  return 0;
}

// xrefStart() queries the start of the xref-table in a PDF file.

var xref = regexp.MustCompile(
  "startxref[\t ]*(\r?\n|\r)"
    "[\t ]*([0-9]+)[\t ]*(\r?\n|\r)"
    "[\t ]*%%EOF")

func xrefStart(pdf []byte) int {
  ps := xref.AllMatchesString(string(pdf[end(pdf, 1024):len(pdf)]), 0);
  if ps == nil {
    return -1
  }
  return num(xref.MatchStrings(ps[len(ps)-1])[2]);
}

// xrefSkip() queries the start of the trailer for a (partial) xref-table.

var sxrf = regexp.MustCompile(
  "^[\r\n\t ]*xref[\t ]*(\r?\n|\r)")
var nxrf = regexp.MustCompile(
  "^[\t ]*([0-9]+)[\t ]+([0-9]+)[\t ]*(\r?\n|\r)")

func xrefSkip(pdf []byte, xref int) int {
  ps := sxrf.MatchStrings(string(pdf[xref:min(xref+100, len(pdf))]));
  if ps == nil {
    return -1
  }
  for {
    xref += len(ps[0]);
    ps = nxrf.MatchStrings(string(pdf[xref:min(xref+100, len(pdf))]));
    if ps == nil {
      break
    }
    xref += num(ps[2]) * 20;
  }
  return xref;
}

// skipPdfString() skips over a string in PDF input.
func skipPdfString(pdf *[]byte, p int) int {
  if (*pdf)[p] != '(' {
    return -1
  }
  dc := MAX_PDF_STRING;
  for d := 1; d > 0 && dc > 0; dc-- {
    p++;
    switch (*pdf)[p] {
    case '\\':
      p++
    case '(':
      d++
    case ')':
      d--
    }
  }
  if dc <= 0 {
    return -1
  }
  return p;
}

// skipPdfComposite() skips over composite data. This is very forgiving about
// errors.
func skipPdfComposite(pdf *[]byte, p int) int {
  if (*pdf)[p] == '(' {
    return skipPdfString(pdf, p)
  }
  if (*pdf)[p] != '<' && (*pdf)[p] != '[' {
    return -1
  }
  dc := MAX_PDF_DICT;
  for d := 1; d > 0 && dc > 0; dc-- {
    p++;
    switch (*pdf)[p] {
    case '[', '<':
      d++
    case ']', '>':
      d--
    case '(':
      p = skipPdfString(pdf, p);
      if p == -1 {
        return -1
      }
    }
  }
  if dc <= 0 {
    return -1
  }
  return p;
}

// Token() separates a token from PDF input.

var ptok = regexp.MustCompile(
  "^[\r\n\t ]*("
    "([0-9]+)([\r\n\t ]+([0-9]+)[\r\n\t ]+R)?" // object reference or number
    "|"
    "/[A-Za-z0-9]+" // dictionary key
    "|"
    "<" // hex string or dictionary
    "|"
    "[" // array
    "|"
    "[(]" // string
    "|"
    "%[^\r\n]+[\r\n]" // comment
    "|"
    "[A-Za-z][A-Za-z0-9]+" // keyword
    ")")

func Token(pdf *[]byte, p int) (int, []byte) {
redo:
  ps := ptok.Execute((*pdf)[p:len((*pdf))]);
  if ps == nil {
    return -1, _Bytes
  }
  q := p + ps[1];
  s := (*pdf)[p+ps[2] : q];
  switch s[0] {
  case '%':
    p = q;
    goto redo;
  case '<', '[', '(':
    q = skipPdfComposite(pdf, p+ps[2]) + 1;
    if q == -1 {
      return -1, _Bytes
    }
    s = (*pdf)[p+ps[2] : q];
  }
  return q, s;
}

// Dictionary() makes a map/hash from PDF dictionary data.
func Dictionary(s []byte) map[string][]byte {
  if len(s) < 4 {
    return nil
  }
  e := len(s) - 1;
  if s[0] != s[1] || s[0] != '<' || s[e] != s[e-1] || s[e] != '>' {
    return nil
  }
  r := make(map[string][]byte);
  s = s[2 : e-1];
  p := 0;
  for {
    t := _Bytes;
    p, t = Token(&s, p);
    if p < 0 {
      break
    }
    k := string(t);
    p, t = Token(&s, p);
    if k[0] != '/' || p < 0 {
      return nil
    }
    r[k] = t;
  }
  return r;
}

// Array() extracts an array from PDF data.
func Array(s []byte) [][]byte {
  if len(s) < 2 || s[0] != '[' || s[len(s)-1] != ']' {
    return nil
  }
  s = s[1 : len(s)-1];
  r := make([][]byte, MAX_PDF_ARRAYSIZE);
  b := 0;
  for p := 0; p >= 0; b++ {
    p, r[b] = Token(&s, p)
  }
  if b == 1 {
    return nil
  }
  return r[0 : b-1];
}

// xrefRead() reads the xref table(s) of a PDF file. This is not recursive
// in favour of not to have to keep track of already used starting points
// for xrefs.

var vals = regexp.MustCompile("^[^0-9]*([0-9]+)[\t ]+([0-9]+)[\r\n\t ]+")

func xrefRead(pdf []byte, p int) map[int]int {
  var back [MAX_PDF_UPDATES]int;
  b := 0;
  s := _Bytes;
  for ok := true; ok; {
    back[b] = p;
    b++;
    p = xrefSkip(pdf, p);
    p, s = Token(&pdf, p);
    if string(s) != "trailer" {
      return nil
    }
    p, s = Token(&pdf, p);
    s, ok = Dictionary(s)["/Prev"];
    p = num(string(s));
  }
  r := make(map[int]int);
  for b != 0 {
    b--;
    p = back[b];
    for {
      m := vals.MatchStrings(string(pdf[p : p+32]));
      if m == nil {
        break
      }
      p += len(m[0]);
      o := num(m[1]);
      for c := num(m[2]); c > 0; c-- {
        m = vals.MatchStrings(string(pdf[p : p+20]));
        if m == nil {
          return nil
        }
        if pdf[p+len(m[0])] != 'n' {
          r[o] = 0, false
        } else {
          r[o] = num(m[1])
        }
        p += 20;
        o++;
      }
    }
  }
  return r;
}

// object() extracts the top informations of a PDF "object". For streams
// this would be the dictionary as bytes.  It also returns the position in
// binary data where one has to continue to read for this "object".

var obj = regexp.MustCompile("^[\r\n\t ]*"
  "([0-9]+)"
  "[\r\n\t ]+"
  "[0-9]+"
  "[\r\n\t ]+"
  "obj")

func object(xr map[int]int, pdf []byte, o int) (int, []byte) {
  p, ok := xr[o];
  if !ok {
    return -1, _Bytes
  }
  m := obj.MatchStrings(string(pdf[p : p+64]));
  if m == nil || num(m[1]) != o {
    return -1, _Bytes
  }
  return Token(&pdf, p+len(m[0]));
}

// pd.Resolve() resolves a reference in the PDF file. You'll probably need
// this method for reading streams only.
var res = regexp.MustCompile("^[\r\n\t ]*"
  "([0-9]+)"
  "[\r\n\t ]+"
  "[0-9]+"
  "[\r\n\t ]+"
  "R$")

func (pd *PdfReaderT) Resolve(s []byte) (int, []byte) {
  n := -1;
  if len(s) >= 5 && s[len(s)-1] == 'R' {
    z, ok := pd.rcache[string(s)];
    if ok {
      return pd.rncache[string(s)], z
    }
    done := make(map[int]int);
    orig := s;
  redo:
    m := res.MatchStrings(string(s));
    if m != nil {
      n = num(m[1]);
      if _, wrong := done[n]; wrong {
        return -1, _Bytes
      }
      done[n] = 1;
      n, s = object(pd.Xref, pd.Bin, n);
      if z, ok = pd.rcache[string(s)]; !ok {
        goto redo
      }
      s = z;
      n = pd.rncache[string(s)];
    }
    pd.rcache[string(orig)] = s;
    pd.rncache[string(orig)] = n;
  }
  return n, s;
}

// pd.Obj() is the universal method to access contents of PDF objects or
// data tokens in i.e.  dictionaries.  For reading streams you'll have to
// utilize pd.Resolve().
func (pd *PdfReaderT) Obj(reference []byte) []byte {
  _, r := pd.Resolve(reference);
  return r;
}

func (pd *PdfReaderT) Num(reference []byte) int {
  return num(string(pd.Obj(reference)))
}

func (pd *PdfReaderT) Dic(reference []byte) map[string][]byte {
  return Dictionary(pd.Obj(reference))
}

func (pd *PdfReaderT) Arr(reference []byte) [][]byte {
  return Array(pd.Obj(reference))
}

// pd.Pages() returns an array with references to the pages of the PDF.
func (pd *PdfReaderT) Pages() [][]byte {
  if pd.pages != nil {
    return pd.pages
  }
  pages := pd.Dic(pd.Dic(pd.Trailer["/Root"])["/Pages"]);
  pd.pages = make([][]byte, pd.Num(pages["/Count"]));
  cp := 0;
  done := make(map[string]int);
  var q func(p [][]byte);
  q = func(p [][]byte) {
    for k := range p {
      if _, wrong := done[string(p[k])]; !wrong {
        done[string(p[k])] = 1;
        if kids, ok := pd.Dic(p[k])["/Kids"]; ok {
          q(pd.Arr(kids))
        } else {
          pd.pages[cp] = p[k];
          cp++;
        }
      }
    }
  };
  q(pd.Arr(pages["/Kids"]));
  return pd.pages;
}

func (pd *PdfReaderT) Attribute(a string, src []byte) []byte {
  d := pd.Dic(src);
  done := make(map[string]int);
  r, ok := d[a];
  for !ok {
    r, ok = d["/Parent"];
    if _, wrong := done[string(r)]; wrong || !ok {
      return _Bytes
    }
    done[string(r)] = 1;
    d = pd.Dic(r);
    r, ok = d[a];
  }
  return r;
}

func (pd *PdfReaderT) Att(a string, src []byte) []byte {
  return pd.Obj(pd.Attribute(a, src))
}

// Load() loads a PDF file of a given name.
func Load(fn string) *PdfReaderT {
  a, e := io.ReadFile(fn);
  //  fmt.Printf("Here we do have it!\n");
  if e != nil {
    return nil
  }
  r := new(PdfReaderT);
  r.File = fn;
  r.Bin = a;
  if r.Startxref = xrefStart(a); r.Startxref == -1 {
    return nil
  }
  if r.Xref = xrefRead(a, r.Startxref); r.Xref == nil {
    return nil
  }
  p, s := Token(&a, xrefSkip(a, r.Startxref));
  if string(s) != "trailer" {
    return nil
  }
  p, s = Token(&a, p);
  if r.Trailer = Dictionary(s); r.Trailer == nil {
    return nil
  }
  r.rcache = make(map[string][]byte);
  r.rncache = make(map[string]int);
  return r;
}
