package main

// Example program for pdfread.go

// The program takes a PDF file and an object reference of a stream.
// The output are the decoded stream contents.
//
// Example:
//  ./pdstream.go foo.pdf "9 0 R"

import (
  "os";
  "fmt";
  "pdfread";
  "util";
  "cmap";
  "fancy";
)

func main() {
  pd := pdfread.Load(os.Args[1]);
  _, d := pd.DecodedStream(util.Bytes(os.Args[2]));
  fmt.Printf("%s", d);

  a := cmap.Read(fancy.SliceReader(d));
  fmt.Printf("\n%v\n", a);
}
