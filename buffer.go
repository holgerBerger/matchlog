package main

/*
	assenmble the consolidated buffer of all loaded files,
	this is the sorted data being displayed

	(c) Holger Berger 2016, under GPL
*/

import (
	"bytes"
	"regexp"
	"time"

	termbox "github.com/nsf/termbox-go"
)

// FIXME buffer should may be hold runes and not bytes? UTF files needs testing

// BufferT represents the buffer as shown on the screen, so the aggregation of files sorted for time
// and filtered
type BufferT struct {
	linecount  int                          // total number of lines
	lines      []LineT                      // line data
	files      []*FlexFileT                 // list of files added to the buffer
	rules      RulesT                       // color rules to apply
	hostcolors map[string]termbox.Attribute // slice with hostnames
	maxcolor   termbox.Attribute            // color for hosts
	hosts      *HostsT                      // hostfiles for IP to hostname mapping
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

// addFile adds a already loaded file to the buffer, host coloring is here
// but the data is not joined into the buffer, call sortFile to do this
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
		b.lines = make([]LineT, b.linecount+f.linecount, b.linecount+f.linecount)
		b.linecount += f.linecount

		for _, line := range f.lines {
			hostname := line.host
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

	regex, _ := regexp.Compile(`([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)`)

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
				if b.files[file].lines[filelinecounter[file]].time.Before(smallest) {
					smallest = b.files[file].lines[filelinecounter[file]].time
					smallestfile = file
				}
			}
		}

		// copy line to buffer, unifying time format
		origline := b.files[smallestfile].lines[filelinecounter[smallestfile]]
		origline.line = append([]byte(origline.time.Format("Mon Jan 02 15:04:05 ")), origline.line[origline.hoststart:]...)
		// correct the location of the hostname in the line
		origline.hostend -= origline.hoststart - 20 // 20 is the length of time format
		origline.hoststart = 20
		b.lines[lnr] = origline

		// if we have a hosts file loaded, check for IP addresses
		if b.hosts != nil {
			match := regex.FindSubmatch(b.lines[lnr].line)
			if match != nil {
				ipname, ok := b.hosts.ip2name[string(match[1])]
				if ok {
					b.lines[lnr].line = bytes.Replace(b.lines[lnr].line, match[1], []byte(ipname), -1)
				}
			}
		}

		filelinecounter[smallestfile]++
	}

}

// AddHosts adds already loaded hostfiles to buffer, will be used for IP replacement
func (b *BufferT) AddHosts(hosts *HostsT) {
	b.hosts = hosts
}
