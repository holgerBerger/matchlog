package main

import "time"

// FIXME buffer should may be hold runes and not bytes? UTF files needs testing

// BufferT represents the buffer as shown on the screen, so the aggregation of files sorted for time
// and filtered
type BufferT struct {
	linecount int       // total number of lines
	lineps    []int     // array of pointers to start of lines of aggregation of files
	contps    []*[]byte // array of pointers to start of contents, for each line of Buffer
	files     []*FileT  // list of files added to the buffer
	// filters []FilterT // list of filters added to the buffer
}

// NewBuffer allocates an empty new buffer
func NewBuffer() *BufferT {
	var buffer BufferT
	buffer.lineps = make([]int, 0, 1024)
	buffer.files = make([]*FileT, 0, 10)
	return &buffer
}

// addFile adds a already loaded file to the buffer, here sorting and filtering takes place
// FIXME this works with pointers in original buffers, this could be changed to a copy,
// which would be nice in case filters like "resolv IP addresses" are changing the buffer,
// so one could have a copy and therefor by able to undo the filtering
func (b *BufferT) addFile(f *FileT) {

	// append file to files if not yet in
	found := false
	for _, fileHelper := range b.files {
		if f.filename == fileHelper.filename {
			found = true
		}
	}
	if !found {
		b.files = append(b.files, f)
	}

	// make space for new file
	b.lineps = make([]int, b.linecount+f.linecount, b.linecount+f.linecount)
	b.contps = make([]*[]byte, b.linecount+f.linecount, b.linecount+f.linecount)
	b.linecount += f.linecount

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
		b.contps[lnr] = &b.files[smallestfile].contents
		filelinecounter[smallestfile]++
	}

}
