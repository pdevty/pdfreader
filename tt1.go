package main

// use this program with a pfa-font - it is only here for testing

import (
	"type1"
	"io/ioutil"
	"os"
	"fancy"
	"fmt"
	"util"
	"pfb"
)

func dumpT1(i *type1.TypeOneI) {
	for k := range i.Fonts {
		fmt.Printf("Font: %s %s\n", k, i.Fonts[k])
		df := i.Dic(i.Fonts[k])
		for l := range df {
			fmt.Printf("  %s %s\n", l, df[l])
		}
		fmt.Printf("\nFontInfo:\n")
		d := i.Dic(string(df["/FontInfo"]))
		for l := range d {
			fmt.Printf("  %s %s\n", l, d[l])
		}
		/*
		   fmt.Printf("\n\nCharStrings:");
		   d = i.Dic(string(df["/CharStrings"]));
		   for l := range d {
		     fmt.Printf("  %s %v\n", l, d[l])
		   }
		*/
	}
}

func main() {
	a, _ := ioutil.ReadFile(os.Args[1])
	if a[0] == 128 {
		a = pfb.Decode(a)
	}
	g := type1.Read(fancy.SliceReader(a))
	fmt.Printf("%v\n", util.StringArray(g.St.Dump()))
	dumpT1(g)
}
