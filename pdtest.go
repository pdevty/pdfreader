package main

// Example program for pdfread.go

// The program takes a PDF file as argument and writes the MediaBoxes and
// defined fonts of the pages.

import (
  "os";
  "fmt";
  "pdfread";
  "fancy";
  "graf";
  "svg";
  "svgtext";
)

func main() {
  pd := pdfread.Load(os.Args[1]);
  if pd != nil {
    pg := pd.Pages();
    for k := range pg {
      fmt.Printf("Page %d - MediaBox: %s\n",
        k+1, pd.Att("/MediaBox", pg[k]));
      fonts := pd.PageFonts(pg[k]);
      for l := range fonts {
        fontname := pd.Dic(fonts[l])["/BaseFont"];
        fmt.Printf("  %s = \"%s\"\n",
          l, fontname[1:len(fontname)]);
      }
    }

    //    /* To test PDF streams:
    cont := pd.ForcedArray(pd.Dic(pg[1])["/Contents"]);
    _, ps := pd.DecodedStream(cont[0]);
    fmt.Printf("Length of stream: %d\n%v", len(ps),
      string(ps));
    //    */
    test := svg.NewTestSvg();
    st := svgtext.New();
    st.Conf = new(graf.TextConfigT);
    test.TConf = st.Conf;
    st.Pdf = pd;
    st.Page = 1;
    test.Text = st;
    test.Interpret(fancy.SliceReader(ps));
    test.Draw.CloseDrawing();
    fmt.Printf("%v\n", test.Stack.Dump());
  }
}
