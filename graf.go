package graf

// Every args are in bytes here to allow lossless formate transformation.
// Using floating point args would cause trouble with rounding - and: []byte
// is not complicated to understand ;)

type StackT struct {
  st [][]byte;
  sp int;
}

func (st *StackT) Push(s []byte) {
  st.st[st.sp] = s;
  st.sp++;
}

func (st *StackT) Drop(n int) [][]byte {
  st.sp -= n;
  return st.st[st.sp : st.sp+n];
}

func (st *StackT) Pop() []byte {
  st.sp--;
  return st.st[st.sp];
}

type Stack interface {
  Push([]byte);
  Pop() []byte;
  Drop(int) (st [][]byte);
}

type Drawer interface {
  DropPath();
  CurrentPoint() [][]byte;
  MoveTo(coord [][]byte);
  LineTo(coord [][]byte);
  CurveTo(coords [][]byte);
  Rectangle(coords [][]byte);
  ClosePath();
  Stroke();
  Fill();
  EOFill();
  FillAndStroke();
  EOFillAndStroke();
  Clip();
  EOClip();
  Concat(matrix [][]byte);
}

type DrawerState interface {
  GSave();
  GRestore();
}

type DrawerConfig interface {
  SetLineWidth(a []byte);
  SetMiterLimit(a []byte);
  SetLineJoin(a []byte);
  SetLineCap(a []byte);
  SetFlat(a []byte);
}

type DrawerColor interface {
  SetGrayStroke(a []byte);
  SetGrayFill(a []byte);
  SetCMYKStroke(cmyk [][]byte);
  SetCMYKFill(cmyk [][]byte);
  SetRGBStroke(rgb [][]byte);
  SetRGBFill(rgb [][]byte);
}

type PdfDrawerT struct {
  Draw   Drawer;
  State  DrawerState;
  Config DrawerConfig;
  Color  DrawerColor;
  Stack  Stack;
}

var PdfOps = map[string]func(pd *PdfDrawerT){
  "m": func(pd *PdfDrawerT) { pd.Draw.MoveTo(pd.Stack.Drop(2)) },
  "l": func(pd *PdfDrawerT) { pd.Draw.LineTo(pd.Stack.Drop(2)) },
  "c": func(pd *PdfDrawerT) { pd.Draw.CurveTo(pd.Stack.Drop(6)) },
  "v": func(pd *PdfDrawerT) {
    c := pd.Draw.CurrentPoint();
    a := pd.Stack.Drop(4);
    pd.Draw.CurveTo([][]byte{c[0], c[1], a[0], a[1], a[2], a[3]});
  },
  "y": func(pd *PdfDrawerT) {
    a := pd.Stack.Drop(4);
    pd.Draw.CurveTo([][]byte{a[0], a[1], a[2], a[3], a[2], a[3]});
  },
  "h": func(pd *PdfDrawerT) { pd.Draw.ClosePath() },
  "re": func(pd *PdfDrawerT) { pd.Draw.Rectangle(pd.Stack.Drop(4)) },

  "S": func(pd *PdfDrawerT) { pd.Draw.Stroke() },
  "s": func(pd *PdfDrawerT) {
    pd.Draw.ClosePath();
    pd.Draw.Stroke();
  },
  "f": func(pd *PdfDrawerT) { pd.Draw.Fill() },
  "F": func(pd *PdfDrawerT) { pd.Draw.Fill() },
  "f*": func(pd *PdfDrawerT) { pd.Draw.EOFill() },
  "B": func(pd *PdfDrawerT) { pd.Draw.FillAndStroke() },
  "B*": func(pd *PdfDrawerT) { pd.Draw.EOFillAndStroke() },
  "b": func(pd *PdfDrawerT) {
    pd.Draw.ClosePath();
    pd.Draw.FillAndStroke();
  },
  "b*": func(pd *PdfDrawerT) {
    pd.Draw.ClosePath();
    pd.Draw.EOFillAndStroke();
  },
  "n": func(pd *PdfDrawerT) { pd.Draw.DropPath() },
}
