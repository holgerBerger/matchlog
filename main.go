package main

import (
	"os"

	termbox "github.com/nsf/termbox-go"
)

var files []*File

// main programm, reading arguments
func main() {

	// read files from commandline
	for _, filename := range os.Args[1:] {
		f, err := ReadFile(filename)
		if err == nil {
			files = append(files, f)
		}
	}

	// create a screen, this is the event loop
	screen := NewScreen()

	// run the event loop
	screen.eventLoop()
	termbox.Close()

}
