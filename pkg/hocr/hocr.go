// Copyright 2019 Nick White.
// Use of this source code is governed by the GPLv3
// license that can be found in the LICENSE file.

// hocr contains structures and functions for parsing and analysing
// hocr files
package hocr

import (
	"encoding/xml"
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Hocr struct {
	Pages []Page `xml:"body>div"`
}

type Page struct {
	Lines []OcrLine `xml:"div>p>span"`
	Title string    `xml:"title,attr"`
}

func (p *Page) GetMeta() *map[string]string {
	return SplitMeta(p.Title)
}

type OcrLine struct {
	Class string    `xml:"class,attr"`
	Id    string    `xml:"id,attr"`
	Title string    `xml:"title,attr"`
	Words []OcrWord `xml:"span"`
	Text  string    `xml:",chardata"`
}

func (l OcrLine) GetShortId() string {
	a := strings.Split(l.Id, "_")
	if len(a) == 3 {
		return a[2]
	}
	return l.Id
}

func (l OcrLine) GetMeta() *map[string]string {
	return SplitMeta(l.Title)
}

type OcrWord struct {
	Class string    `xml:"class,attr"`
	Id    string    `xml:"id,attr"`
	Title string    `xml:"title,attr"`
	Chars []OcrChar `xml:"span"`
	Text  string    `xml:",chardata"`
}

func (w *OcrWord) GetShortId() string {
	a := strings.Split(w.Id, "_")
	if len(a) == 3 {
		return a[2]
	}
	return w.Id
}

// GetCleanText parses a hOCR file and extracts the text from it
func (w *OcrWord) GetCleanText() string {
	return strings.Trim(strings.Trim(w.Text, " "), "\n")
}

// SplitMeta splits a hOCR title string into a map of key/value pairs
func SplitMeta(meta string) *map[string]string {
	a := strings.Split(meta, ";")
	r := make(map[string]string)
	for _, s := range a {
		s = strings.Trim(s, " ")
		kv := strings.Split(s, " ")
		r[kv[0]] = strings.Trim(strings.TrimPrefix(s, kv[0]), " ")
	}
	return &r
}

// GetMeta parses a hOCR file and extracts the metadata from it
func (w *OcrWord) GetMeta() *map[string]string {
	return SplitMeta(w.Title)
}

type OcrChar struct {
	Class string    `xml:"class,attr"`
	Id    string    `xml:"id,attr"`
	Title string    `xml:"title,attr"`
	Chars []OcrChar `xml:"span"`
	Text  string    `xml:",chardata"`
}

// Returns the confidence for a word based on its x_wconf value
func wordConf(s string) (float64, error) {
	re, err := regexp.Compile(`x_wconf ([0-9.]+)`)
	if err != nil {
		return 0.0, err
	}
	conf := re.FindStringSubmatch(s)
	return strconv.ParseFloat(conf[1], 64)
}

// BoxCoords parses bbox coordinate strings
func BoxCoords(s string) ([4]int, error) {
	var coords [4]int
	re, err := regexp.Compile(`bbox ([0-9]+) ([0-9]+) ([0-9]+) ([0-9]+)`)
	if err != nil {
		return coords, err
	}
	coordstr := re.FindStringSubmatch(s)
	for i := range coords {
		c, err := strconv.Atoi(coordstr[i+1])
		if err != nil {
			return coords, err
		}
		coords[i] = c
	}
	return coords, nil
}

func noText(s string) bool {
	t := strings.Trim(s, " \n")
	return len(t) == 0
}

// Parse parses a hOCR file
func Parse(b []byte) (Hocr, error) {
	var hocr Hocr

	err := xml.Unmarshal(b, &hocr)
	if err != nil {
		return hocr, err
	}

	return hocr, nil
}

// GetText parses a hOCR file and extracts the text from it
func GetText(hocrfn string) (string, error) {
	var s string

	file, err := os.ReadFile(hocrfn)
	if err != nil {
		return s, err
	}

	h, err := Parse(file)
	if err != nil {
		return s, err
	}

	for _, p := range h.Pages {
		for _, l := range p.Lines {
			s += LineText(l) + "\n"
		}
	}
	return s, nil
}

// GetAvgConf calculates the average confidence of a hOCR file from
// confidences embedded in each word
func GetAvgConf(hocrfn string) (float64, error) {
	file, err := os.ReadFile(hocrfn)
	if err != nil {
		return 0, err
	}

	h, err := Parse(file)
	if err != nil {
		return 0, err
	}

	var total, num float64
	for _, p := range h.Pages {
		for _, l := range p.Lines {
			for _, w := range l.Words {
				c, err := wordConf(w.Title)
				if err != nil {
					return 0, err
				}
				total += c
				num++
			}
		}
	}
	if num == 0 {
		return 0, errors.New("No words found")
	}
	return total / num, nil
}

// GetWordConfs is a utility function that parses a hocr
// file and returns an array containing the confidences
// of each word therein
func GetWordConfs(hocrfn string) ([]float64, error) {
	var confs []float64

	file, err := os.ReadFile(hocrfn)
	if err != nil {
		return confs, err
	}

	h, err := Parse(file)
	if err != nil {
		return confs, err
	}

	for _, p := range h.Pages {
		for _, l := range p.Lines {
			for _, w := range l.Words {
				c, err := wordConf(w.Title)
				if err != nil {
					return confs, err
				}
				confs = append(confs, c)
			}
		}
	}

	return confs, nil
}
