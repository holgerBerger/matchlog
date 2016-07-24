package main

import "io/ioutil"

// File contains metadata and data of a read logfile
type FileT struct {
	filename  string
	linecount int
	lines     []int
	contents  []byte
}

// ReadFile reads all the file into a buffer
func ReadFile(filename string) (*FileT, error) {
	var (
		newfile FileT
		err     error
	)
	newfile.filename = filename
	newfile.contents, err = ioutil.ReadFile(filename)

	// fmt.Print("reading ", filename, " with ")
	newfile.indexLines()

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
