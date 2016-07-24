package main

import "fmt"

// FIXME buffer should may be hold runes and not bytes? UTF files needs testing

// BufferT represents the buffer as shown on the screen, so the aggregation of files sorted for time
// and filtered
type BufferT struct {
	lineps []*byte  // array of pointers to start of lines of aggregation of files
	files  []*FileT // list of files added to the buffer
	// filters []FilterT // list of filters added to the buffer
}

// NewBuffer allocates an empty new buffer
func NewBuffer() *BufferT {
	var buffer BufferT
	buffer.lineps = make([]*byte, 0, 1024)
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

	// FIXME dummy code for testing to handle ONE file
	b.lineps = make([]*byte, f.linecount, f.linecount)

	for lnr := range f.lines[:len(f.lines)-1] {
		fmt.Println(lnr, f.lines[lnr])
		b.lineps[lnr] = &f.contents[f.lines[lnr]]
	}
	// END oF dummy code
}
