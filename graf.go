package graf

import (
  "fancy";
  "pdfread";
  "svg";
)

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

func NewStack(n int) *StackT {
  r := new(StackT);
  r.st = make([][]byte, n);
  return r;
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
  CloseDrawing();
  SetIdentity();
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
  Ops    map[string]func(pd *PdfDrawerT);
}

var PdfOps = map[string]func(pd *PdfDrawerT){
  "cm": func(pd *PdfDrawerT) { pd.Draw.Concat(pd.Stack.Drop(6)) },

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

  "rg": func(pd *PdfDrawerT) {
    pd.Color.SetRGBFill(pd.Stack.Drop(3));
    pd.Ops["sc"] = pd.Ops["rg"];
  },
  "RG": func(pd *PdfDrawerT) {
    pd.Color.SetRGBStroke(pd.Stack.Drop(3));
    pd.Ops["SC"] = pd.Ops["RG"];

  },
  "g": func(pd *PdfDrawerT) {
    pd.Color.SetGrayFill(pd.Stack.Pop());
    pd.Ops["sc"] = pd.Ops["g"];

  },
  "G": func(pd *PdfDrawerT) {
    pd.Color.SetGrayStroke(pd.Stack.Pop());
    pd.Ops["SC"] = pd.Ops["G"];

  },
  "k": func(pd *PdfDrawerT) {
    pd.Color.SetCMYKFill(pd.Stack.Drop(4));
    pd.Ops["sc"] = pd.Ops["k"];

  },
  "K": func(pd *PdfDrawerT) {
    pd.Color.SetCMYKStroke(pd.Stack.Drop(4));
    pd.Ops["SC"] = pd.Ops["K"];

  },

  "w": func(pd *PdfDrawerT) { pd.Config.SetLineWidth(pd.Stack.Pop()) },
  "J": func(pd *PdfDrawerT) { pd.Config.SetLineCap(pd.Stack.Pop()) },
  "j": func(pd *PdfDrawerT) { pd.Config.SetLineJoin(pd.Stack.Pop()) },
  "M": func(pd *PdfDrawerT) { pd.Config.SetMiterLimit(pd.Stack.Pop()) },

  "i": func(pd *PdfDrawerT) {
    pd.Stack.Pop() // FIXME
  },

  "gs": func(pd *PdfDrawerT) {
    pd.Stack.Pop();
    pd.Draw.SetIdentity(); // FIXME ;)
  },
}

func (pd *PdfDrawerT) Interpret(rdr fancy.Reader) {
  for {
    t, _ := pdfread.SimpleToken(rdr);
    if len(t) == 0 {
      break
    }
    if f, ok := pd.Ops[string(t)]; ok {
      f(pd)
    } else {
      pd.Stack.Push(t)
    }
  }
}

func NewTestSvg() *PdfDrawerT {
  r := new(PdfDrawerT);
  r.Stack = NewStack(10240);
  t := svg.NewDrawer();
  r.Draw = t;
  r.Color = t;
  r.Config = t;
  r.Ops = make(map[string]func(pd *PdfDrawerT));
  for k := range PdfOps {
    r.Ops[k] = PdfOps[k]
  }
  return r;
}
