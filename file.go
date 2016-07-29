package main

// old file class, no longer used, coudl not read compressed files
// can be removed when everything is transfered to line based memory format

import (
	"io/ioutil"
	"regexp"
	"time"
)

// FileT contains metadata and data of a read logfile
type FileT struct {
	filename  string         // name of the file
	linecount int            // number of line sin the file
	lines     []int          // offsets to linestarts
	contents  []byte         // loaded contents of file
	times     []time.Time    // timestamp for each line
	location  *time.Location // cache for location
}

// ReadFile reads all the file into a buffer
func ReadFile(filename string) (*FileT, error) {
	var (
		newfile FileT
		err     error
	)
	newfile.filename = filename
	newfile.contents, err = ioutil.ReadFile(filename)

	newfile.location, err = time.LoadLocation("Europe/Berlin")
	if err != nil {
		panic("could not load timezone")
	}

	// fmt.Print("reading ", filename, " with ")
	newfile.indexLines()
	newfile.parseLines()

	return &newfile, err
}

// indexLines searches all newlines to store line starts
func (f *FileT) indexLines() {

	// bail out if file is empty
	if len(f.contents) < 1 {
		f.linecount = 0
		f.lines = make([]int, 1, 1)
		f.lines[0] = 0
		return
	}

	// we go 2x over data, first count
	count := 1 // first line extra
	for _, b := range f.contents {
		if b == 10 {
			count++
		}
	}
	f.lines = make([]int, count, count)
	f.linecount = count - 1
	f.times = make([]time.Time, count, count)

	// fmt.Println("size:", len(f.lines))

	// second index
	count = 0
	f.lines[count] = 0
	count++
	for pos, b := range f.contents {
		if b == 10 {
			f.lines[count] = pos + 1
			count++
		}
	}
	// FIXME is the last pointer pointing behind contents?

	// fmt.Println(f.linecount, "lines")
}

// parseLines parses timestamps of all lines and stores them in FileT
// in two phases, first it matches with a regex to strip of rest of line, second
// it using the time.Parse functions
func (f *FileT) parseLines() {

	// Jul 24 06:29:28
	fmt1, _ := regexp.Compile(`[a-zA-Z]{3} [0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}`)
	// 2016-07-26T00:36:17.903571+02:00
	fmt2, _ := regexp.Compile(`([0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2})\.[0-9]+(\+[0-9]{2}):([0-9]{2})`)

	for line := range f.lines[:f.linecount] {
		pos1 := f.lines[line]
		pos2 := f.lines[line+1]
		datestr := f.contents[pos1:pos2]

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
