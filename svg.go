package svg

import (
  "pdfread";
  "svgdraw";
  "svgtext";
  "strm";
  "util";
  "fancy";
  "fmt";
  "os";
)

func complain(err string) {
  fmt.Printf("%s", err);
  os.Exit(1);
}

func Page(pd *pdfread.PdfReaderT, page int) []byte {
  pg := pd.Pages();
  if page >= len(pg) {
    complain("Page does not exist!\n")
  }
  mbox := util.StringArray(pd.Arr(pd.Att("/MediaBox", pg[page])));
  draw := svgdraw.NewTestSvg();
  svgtext.New(pd, draw).Page = page;
  w := strm.Mul(strm.Sub(mbox[2], mbox[0]), "1.25");
  h := strm.Mul(strm.Sub(mbox[3], mbox[1]), "1.25");
  draw.Write.Out(
    "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"no\"?>\n"+
      "<svg\n"+
      "   xmlns:svg=\"http://www.w3.org/2000/svg\"\n"+
      "   xmlns=\"http://www.w3.org/2000/svg\"\n"+
      "   version=\"1.0\"\n"+
      "   width=\"%s\"\n"+
      "   height=\"%s\">\n"+
      "<g transform=\"matrix(1.25,0,0,-1.25,%s,%s)\">\n",
    w, h,
    strm.Mul(mbox[0], "-1.25"),
    strm.Mul(mbox[3], "1.25"));
  cont := pd.ForcedArray(pd.Dic(pg[page])["/Contents"]);
  _, ps := pd.DecodedStream(cont[0]);
  draw.Interpret(fancy.SliceReader(ps));
  draw.Draw.CloseDrawing();
  draw.Write.Out("</g>\n</svg>\n");
  return draw.Write.Content;
}
