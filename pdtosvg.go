package main

// Example program for pdfread.go

// The program takes a PDF file and converts a page to SVG.

import (
  "os";
  "fmt";
  "pdfread";
  "strm";
  "fancy";
  "svg";
  "svgtext";
  "util";
)

func complain(err string) {
  fmt.Printf("%susage: pdtosvg foo.pdf [page] >foo.svg\n", err);
  os.Exit(1);
}

func main() {
  if len(os.Args) == 1 || len(os.Args) > 3 {
    complain("")
  }
  page := 0;
  if len(os.Args) > 2 {
    page = strm.Int(os.Args[2], 1) - 1;
    if page < 0 {
      complain("Bad page!\n\n")
    }
  }
  pd := pdfread.Load(os.Args[1]);
  if pd == nil {
    complain("Could not load pdf file!\n\n")
  }
  pg := pd.Pages();
  if page >= len(pg) {
    complain("Page does not exist!\n\n")
  }
  mbox := util.StringArray(pd.Arr(pd.Att("/MediaBox", pg[page])));

  w := strm.Mul(strm.Sub(mbox[2], mbox[0]), "1.25");
  h := strm.Mul(strm.Sub(mbox[3], mbox[1]), "1.25");

  fmt.Printf(
    "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"no\"?>\n"
      "<svg\n"
      "   xmlns:svg=\"http://www.w3.org/2000/svg\"\n"
      "   xmlns=\"http://www.w3.org/2000/svg\"\n"
      "   version=\"1.0\"\n"
      "   width=\"%s\"\n"
      "   height=\"%s\">\n"
      "<g transform=\"matrix(1.25,0,0,-1.25,%s,%s)\">\n",
    w, h,
    strm.Mul(mbox[0], "-1.25"),
    strm.Mul(mbox[3], "1.25"));

  test := svg.NewTestSvg();
  svgtext.New(pd, test).Page = page;

  cont := pd.ForcedArray(pd.Dic(pg[page])["/Contents"]);
  _, ps := pd.DecodedStream(cont[0]);

  test.Interpret(fancy.SliceReader(ps));
  test.Draw.CloseDrawing();

  fmt.Printf("</g>\n</svg>\n");

}
