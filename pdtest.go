package main

import (
  "fmt";
  "pdfread";
)

func main() {
  pd := pdfread.Load("tvlvolltext.pdf");
  if pd != nil {
    fmt.Printf("%d pages\n", len(pd.Pages()))
  }
}
