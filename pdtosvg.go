package main

// Example program for pdfread.go

// The program takes a PDF file and converts a page to SVG.

import (
	"os"
	"fmt"
	"strm"
	"pdfread"
	"svg"
)

func complain(err string) {
	fmt.Printf("%susage: pdtosvg foo.pdf [page] >foo.svg\n", err)
	os.Exit(1)
}

func main() {
	if len(os.Args) == 1 || len(os.Args) > 3 {
		complain("")
	}
	page := 0
	if len(os.Args) > 2 {
		page = strm.Int(os.Args[2], 1) - 1
		if page < 0 {
			complain("Bad page!\n\n")
		}
	}
	pd := pdfread.Load(os.Args[1])
	if pd == nil {
		complain("Could not load pdf file!\n\n")
	}
	fmt.Printf("%s", svg.Page(pd, page))
}
