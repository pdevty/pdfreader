package svgtext

import (
  "pdfread";
  "graf";
//  "fmt";
  "util";
)

func unquot(a []byte) []byte { // stub
  return a[1 : len(a)-1] // removes braces only.
}

type SvgTextT struct {
  Pdf    *pdfread.PdfReaderT;
  Drw    *graf.PdfDrawerT;
  matrix []string;
}

func New() *SvgTextT {
  r := new(SvgTextT);
  r.matrix = []string{"1", "0", "0", "1", "0", "0"};
  return r;
}

func (t *SvgTextT) TMoveTo(s [][]byte) {
  t.matrix[4], t.matrix[5] = string(s[0]), string(s[1])
}

func (t *SvgTextT) TNextLine() {}

func (t *SvgTextT) TSetMatrix(s [][]byte) { t.matrix = util.StringArray(s) }

func (t *SvgTextT) TShow(a []byte) {
  tx := t.Pdf.ForcedArray(a); // FIXME: Should be "ForcedSimpleArray()"
  for k := range tx {
//    fmt.Printf("<!-- %s,%s: %s -->\n", t.matrix[4], t.matrix[5], tx[k])
    _ = k;
  }
}
