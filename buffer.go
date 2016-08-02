package main

/*
	assenmble the consolidated buffer of all loaded files,
	this is the sorted data being displayed

	(c) Holger Berger 2016, under GPL
*/

import (
	"time"

	termbox "github.com/nsf/termbox-go"
)

// FIXME buffer should may be hold runes and not bytes? UTF files needs testing

// BufferT represents the buffer as shown on the screen, so the aggregation of files sorted for time
// and filtered
type BufferT struct {
	linecount  int                          // total number of lines
	lineps     [][]byte                     // array if lines
	hostsstart []int                        // start if hostname in line
	hostsend   []int                        // end of hostname in line
	files      []*FlexFileT                 // list of files added to the buffer
	rules      RulesT                       // color rules to apply
	hostcolors map[string]termbox.Attribute // slice with hostnames
	maxcolor   termbox.Attribute
	// filters []FilterT // list of filters added to the buffer
}

// NewBuffer allocates an empty new buffer
func NewBuffer() *BufferT {
	var buffer BufferT
	buffer.files = make([]*FlexFileT, 0, 10)
	buffer.rules = DefaultRules()
	buffer.hostcolors = make(map[string]termbox.Attribute)
	buffer.maxcolor = 17
	return &buffer
}

// addFile adds a already loaded file to the buffer, here sorting and filtering takes place
func (b *BufferT) addFile(f *FlexFileT) {

	// append file to files if not yet in
	found := false
	for _, fileHelper := range b.files {
		if f.filename == fileHelper.filename {
			found = true
		}
	}
	if !found {
		b.files = append(b.files, f)

		// make space for new file
		b.lineps = make([][]byte, b.linecount+f.linecount, b.linecount+f.linecount)
		b.hostsstart = make([]int, b.linecount+f.linecount, b.linecount+f.linecount)
		b.hostsend = make([]int, b.linecount+f.linecount, b.linecount+f.linecount)
		b.linecount += f.linecount

		for _, hostname := range f.hosts {
			_, ok := b.hostcolors[hostname]
			if !ok {
				b.maxcolor += 5
				if b.maxcolor > 232 {
					b.maxcolor = 20 // we wrap around before grey values
				}
				b.hostcolors[hostname] = b.maxcolor
			}
		}
	}

}

// sortFile has to be called after all files have been added with addFile
func (b *BufferT) sortFile() {
	// add file to buffer sorted by time

	// first init some file-local linecounter
	filelinecounter := make([]int, len(b.files), len(b.files))
	for i := range filelinecounter {
		filelinecounter[i] = 0
	}

	for lnr := 0; lnr < b.linecount; lnr++ {
		// find next line in time from all files
		smallest := time.Unix(0xffffffff, 0)
		smallestfile := -1
		for file := range b.files {
			if filelinecounter[file] < b.files[file].linecount {
				if b.files[file].times[filelinecounter[file]].Before(smallest) {
					smallest = b.files[file].times[filelinecounter[file]]
					smallestfile = file
				}
			}
		}
		b.lineps[lnr] = b.files[smallestfile].lines[filelinecounter[smallestfile]]
		b.hostsstart[lnr] = b.files[smallestfile].hostsstart[filelinecounter[smallestfile]]
		b.hostsend[lnr] = b.files[smallestfile].hostsend[filelinecounter[smallestfile]]
		filelinecounter[smallestfile]++
	}

}
