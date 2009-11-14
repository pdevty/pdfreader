package main

import (
  "fmt";
  "pdfread";
)

func main() {
  pd := pdfread.Load("tvlvolltext.pdf"); // change to your file name here
  if pd != nil {
    fmt.Printf("--%s--\n", pd.Obj(pd.Trailer["/Root"]));
    fmt.Printf("--%s--\n", pd.Obj(pd.Trailer["/Root"]));
  }
}
