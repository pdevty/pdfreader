package grafgs

// graphics state - implements the graf.DrawerConfig interface.

import "graf"

type GrafGsT struct {
  Color        graf.DrawerColor;
  CurrentPoint [][]byte;
  FirstPoint   [][]byte;
  StrokeColor  string;
  FillColor    string;
  LineWidth    string;
  MiterLimit   string;
  LineJoin     string;
  LineCap      string;
  Flat         string;
}

func (gs *GrafGsT) SetLineWidth(a []byte)  {
  gs.LineWidth = string(a);
}
func (gs *GrafGsT) SetMiterLimit(a []byte) {
  gs.MiterLimit = string(a);
}
func (gs *GrafGsT) SetLineJoin(a []byte)   {
  gs.LineJoin = string(a);
}
func (gs *GrafGsT) SetLineCap(a []byte)    {
  gs.LineCap = string(a);
}
func (gs *GrafGsT) SetFlat(a []byte)       {
  gs.Flat = string(a);
}
func (gs *GrafGsT) SetGrayStroke(a []byte) {
  gs.StrokeColor = gs.Color.Gray(a);
}
func (gs *GrafGsT) SetGrayFill(a []byte)   {
  gs.FillColor = gs.Color.Gray(a);
}
func (gs *GrafGsT) SetCMYKStroke(cmyk [][]byte) {
  gs.StrokeColor = gs.Color.CMYK(cmyk);
}
func (gs *GrafGsT) SetCMYKFill(cmyk [][]byte) {
  gs.FillColor = gs.Color.CMYK(cmyk);
}
func (gs *GrafGsT) SetRGBStroke(rgb [][]byte) {
  gs.StrokeColor = gs.Color.RGB(rgb);
}
func (gs *GrafGsT) SetRGBFill(rgb [][]byte) {
  gs.FillColor = gs.Color.RGB(rgb);
}


func New(color graf.DrawerColor) *GrafGsT {
  r := new(GrafGsT);
  r.Color = color;
  return r;
}
