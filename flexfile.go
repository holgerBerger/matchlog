package main

import (
	"bufio"
	"compress/gzip"
	"os"
	"regexp"
	"strings"
	"time"
)

// rewrite of FileT allowing compressed files to be read
// and changing to line based storage in memory

// FlexFileT contains metadata and data of a read logfile
type FlexFileT struct {
	filename  string         // name of the file
	linecount int            // number of line sin the file
	lines     [][]byte       // offsets to linestarts
	times     []time.Time    // timestamp for each line
	location  *time.Location // cache for location
}

// ReadFlexFile reads compressed or uncompressed files
func ReadFlexFile(filename string) (*FlexFileT, error) {
	var (
		newfile FlexFileT
		err     error
		reader  *bufio.Reader
	)
	newfile.filename = filename

	osfile, err := os.Open(filename)
	if err != nil {
		return &newfile, err
	}

	newfile.location, err = time.LoadLocation("Europe/Berlin")
	if err != nil {
		panic("could not load timezone")
	}

	if strings.HasSuffix(filename, ".gz") {
		newreader, err := gzip.NewReader(osfile)
		if err != nil {
			reader = bufio.NewReader(osfile)
		} else {
			reader = bufio.NewReader(newreader)
		}
	} else {
		reader = bufio.NewReader(osfile)
	}

	newfile.lines = make([][]byte, 0, 1024)

	linecount := 0
	for {
		nextline, err := reader.ReadBytes('\n')
		newfile.lines = append(newfile.lines, nextline)
		if err != nil {
			break
		}
		linecount++
	}
	newfile.linecount = linecount

	newfile.times = make([]time.Time, linecount, linecount)
	newfile.parseLines()

	return &newfile, nil
}

// parseLines parses timestamps of all lines and stores them in FileT
// in two phases, first it matches with a regex to strip of rest of line, second
// it using the time.Parse functions
func (f *FlexFileT) parseLines() {

	// Jul 24 06:29:28
	fmt1, _ := regexp.Compile(`[a-zA-Z]{3} [0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}`)
	// 2016-07-26T00:36:17.903571+02:00
	fmt2, _ := regexp.Compile(`([0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2})\.[0-9]+(\+[0-9]{2}):([0-9]{2})`)

	for line, datestr := range f.lines[:f.linecount] {

		// this is the reference time from the time modul, templates show this time
		// Mon Jan 2 15:04:05 MST 2006

		match1 := fmt1.Find(datestr)
		if match1 != nil {
			t, err := time.ParseInLocation("2006 Jan 02 15:04:05", "2016 "+string(match1), f.location)
			if err == nil {
				f.times[line] = t
				continue
			}
		} else {
			match2 := fmt2.FindSubmatch(datestr)
			if match2 != nil {
				t, err := time.ParseInLocation("2006-01-02T15:04:05", string(match2[1]), f.location)
				if err == nil {
					f.times[line] = t
					continue
				}
			} else {
				// FIXME wtf do we do here?
				panic("unknown date format in " + string(datestr))
			}
		}
	}
}
