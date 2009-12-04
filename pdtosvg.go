package main

// Example program for pdfread.go

// The program takes a PDF file as argument and writes the MediaBoxes and
// defined fonts of the pages.

import (
  "os";
//  "fmt";
  "pdfread";
//  "util";
  "fancy";
  "svg";
  "svgtext";
)

func main() {
  pd := pdfread.Load(os.Args[1]);
  if pd != nil {
    test := svg.NewTestSvg();
    svgtext.New(pd, test);
    pg := pd.Pages();
    cont := pd.ForcedArray(pd.Dic(pg[1])["/Contents"]);
    _, ps := pd.DecodedStream(cont[0]);
    test.Interpret(fancy.SliceReader(ps));
    test.Draw.CloseDrawing();
  }
}
