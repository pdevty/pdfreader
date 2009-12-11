package main

// use this program with a pfa-font - it is only here for testing

import (
  "type1";
  "io/ioutil";
  "os";
  "fancy";
  "fmt";
  "util";
)

func main() {
  a, _ := ioutil.ReadFile(os.Args[1]);
  g := type1.Read(fancy.SliceReader(a));
  fmt.Printf("%v\n", util.StringArray(g.St.Dump()));
}
