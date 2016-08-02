package main

/*
	read and pares plaintext and compressed files (.gz)
	date string parsing happens here

	(c) Holger Berger 2016, under GPL
*/

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
	filename             string         // name of the file
	linecount            int            // number of line sin the file
	lines                [][]byte       // offsets to linestarts
	times                []time.Time    // timestamp for each line
	hosts                []string       // source host of the message
	hostsstart, hostsend []int          // start and index position of hostname
	location             *time.Location // cache for location
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
	newfile.hosts = make([]string, linecount, linecount)
	newfile.hostsend = make([]int, linecount, linecount)
	newfile.hostsstart = make([]int, linecount, linecount)
	newfile.parseLines()

	return &newfile, nil
}

// parseLines parses timestamps of all lines and stores them in FileT
// in two phases, first it matches with a regex to strip of rest of line, second
// it using the time.Parse functions
func (f *FlexFileT) parseLines() {

	// Jul 24 06:29:28
	fmt1, _ := regexp.Compile(`([a-zA-Z]{3}\s+[0-9]+ [0-9]{2}:[0-9]{2}:[0-9]{2}) (\S+)`)
	// 2016-07-26T00:36:17.903571+02:00
	fmt2, _ := regexp.Compile(`([0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2})\.[0-9]+(\+[0-9]{2}):([0-9]{2}) (\S+)`)

	for line, datestr := range f.lines[:f.linecount] {

		// this is the reference time from the time modul, templates show this time
		// Mon Jan 2 15:04:05 MST 2006

		match1 := fmt1.FindSubmatch(datestr)
		index1 := fmt1.FindSubmatchIndex(datestr)
		if match1 != nil {
			// FIXME hard coded year here, needs to replaced
			t, err := time.ParseInLocation("2006 Jan 02 15:04:05", "2016 "+string(match1[1]), f.location)
			if err == nil {
				f.times[line] = t
				f.hosts[line] = string(match1[2])
				f.hostsstart[line] = index1[2*2]
				f.hostsend[line] = index1[2*2+1]
				continue
			} else {
				// FIXME hard coded year here, needs to replaced
				t, err := time.ParseInLocation("2006 Jan  2 15:04:05", "2016 "+string(match1[1]), f.location)
				if err == nil {
					f.times[line] = t
					f.hosts[line] = string(match1[2])
					f.hostsstart[line] = index1[2*2]
					f.hostsend[line] = index1[2*2+1]
					continue
				} else {
					// FIXME wtf do we do here?
					panic("could not parse date" + string(match1[1]))
				}
			}
		} else {
			match2 := fmt2.FindSubmatch(datestr)
			index2 := fmt2.FindSubmatchIndex(datestr)
			if match2 != nil {
				t, err := time.ParseInLocation("2006-01-02T15:04:05", string(match2[1]), f.location)
				if err == nil {
					f.times[line] = t
					f.hosts[line] = string(match2[4])
					f.hostsstart[line] = index2[4*2]
					f.hostsend[line] = index2[4*2+1]
					continue
				}
			} else {
				// FIXME wtf do we do here?
				panic("unknown date format in " + string(datestr))
			}
		}
	}
}
