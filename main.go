package main

import (
	"os"

	termbox "github.com/nsf/termbox-go"
)

// main programm, reading arguments
func main() {

	// files holds list of files
	var files []*FileT

	// read files from commandline
	for _, filename := range os.Args[1:] {
		f, err := ReadFile(filename)
		if err == nil {
			files = append(files, f)
		}
	}

	// create the buffer aggregating the files
	buffer := NewBuffer()

	// add all files to buffer
	for f := range files {
		buffer.addFile(files[f])
	}
	buffer.sortFile()

	// create a screen, this is the event loop
	screen := NewScreen(files, buffer)

	// run the event loop
	screen.eventLoop()
	termbox.Close()

}
