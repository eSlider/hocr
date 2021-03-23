// Copyright 2021 Nick White.
// Use of this source code is governed by the GPLv3
// license that can be found in the LICENSE file.

// extracthocrlines copies the text and corresponding image section
// for each line of a HOCR file into separate files, which is
// useful for OCR training
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"rescribe.xyz/utils/pkg/hocr"
	"rescribe.xyz/utils/pkg/line"
)

const usage = `Usage: extracthocrlines [-d] [-e] file.hocr [file.hocr]

Copies the text and corresponding image section for each line
of a HOCR file into separate files, which is useful for OCR
training.
`

// saveline saves the text and image for a line in a directory
func saveline(l line.Detail, dir string) error {
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		return err
	}

	base := filepath.Join(dir, l.OcrName+"_"+l.Name)

	f, err := os.Create(base + ".png")
	if err != nil {
		return fmt.Errorf("Error creating file %s: %v", base+".png", err)
	}

	err = l.Img.CopyLineTo(f)
	if err != nil {
		return fmt.Errorf("Error writing line image for %s: %v", base+".png", err)
	}

	f, err = os.Create(base + ".txt")
	if err != nil {
		return fmt.Errorf("Error creating file %s: %v", base+".txt", err)
	}

	_, err = io.WriteString(f, l.Text)
	if err != nil {
		return fmt.Errorf("Error writing line text for %s: %v", base+".txt", err)
	}

	return nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usage)
		flag.PrintDefaults()
	}
	dir := flag.String("d", ".", "Directory to save lines in")
	embeddedimgpath := flag.Bool("e", false, "Use image path embedded in hOCR (rather than the path of the .hocr file with a .png suffix)")
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	for _, f := range flag.Args() {
		var err error
		var newlines line.Details
		if *embeddedimgpath {
			newlines, err = hocr.GetLineDetails(f)
		} else {
			imgName := strings.TrimSuffix(f, ".hocr") + ".png"
			newlines, err = hocr.GetLineDetailsCustomImg(f, imgName)
		}
		if err != nil {
			log.Fatal(err)
		}

		for _, l := range newlines {
			if l.Img == nil {
				continue
			}
			err = saveline(l, *dir)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
