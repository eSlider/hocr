// Copyright 2019 Nick White.
// Use of this source code is governed by the GPLv3
// license that can be found in the LICENSE file.

// boxtotxt converts a Tesseract .box file to plain Text
package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"rescribe.xyz/utils/pkg/hocr"
)

// Word of a line
type Word struct {
	Text string
	//Meta *map[string]string
}

// Line of a page
type Line struct {
	//Words []string //[]Word
	Text string
	Meta *map[string]string

	//Meta  *map[string]string
}

// Page of a document
type Page struct {
	Lines []Line
	//Meta  *map[string]string
}

// Document is a collection of Pages
type Document struct {
	Pages []Page
}

// NewDocument creates a new document
func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: boxtotxt in.box\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	var filePath string
	if flag.NArg() < 1 {
		//flag.Usage()
		filePath = "data/test.hocr"
		//os.Exit(1)
	} else {
		filePath = flag.Arg(0)

	}
	// read the file
	file, err := os.ReadFile(filePath)
	h, err := hocr.Parse(file)
	if err != nil {
		log.Fatal(err)
	}

	// Prepare Pages
	var pages []Page
	for _, p := range h.Pages {

		// Prepare Lines
		var lines []Line
		for _, l := range p.Lines {

			// Prepare Words
			//var words []Word
			//for _, w := range l.Words {
			//	words = append(words, Word{
			//		Text: w.GetCleanText(),
			//		//Meta: w.GetMeta(),
			//	})
			//}
			//var words []string
			//for _, w := range l.Words {
			//	words = append(words, w.GetCleanText())
			//}
			//
			//// Fill Lines
			//lines = append(lines, Line{
			//	Text: hocr.LineText(l),
			//	//Meta:  l.GetMeta(),
			//	Words: words,
			//})

			// Fill Lines
			lines = append(lines, Line{
				Text: hocr.LineText(l),
				Meta: l.GetMeta(),
				//Words: words,
			})
		}

		pages = append(pages, Page{
			//Meta:  p.GetMeta(),
			Lines: lines,
		})

	}

	r := Document{Pages: pages}

	d, err := yaml.Marshal(r)
	//d, err := json.Marshal(r)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Print(string(d))

}
