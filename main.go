package main

/*
	main programm
	command line handling, init and passing control to
	screen event loop

	(c) Holger Berger 2016, under GPL
*/

import (
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
	termbox "github.com/nsf/termbox-go"
)

var opts struct {
	Hosts string `long:"hosts" description:"hosts files to load, comma seperated"`
}

// main programm, reading arguments
func main() {
	var hosts *HostsT

	// parse commandlines for --hosts option
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(0)
	}

	if opts.Hosts != "" {
		hosts = NewHosts()
		for _, hf := range strings.Split(opts.Hosts, ",") {
			hosts.AddFile(hf)
		}
	} else {
		hosts = nil
	}

	// files holds list of files
	var files []*FlexFileT

	// read files from commandline
	for _, filename := range os.Args[1:] {
		f, err := ReadFlexFile(filename)
		if err == nil {
			files = append(files, f)
		}
	}

	// create the buffer aggregating the files
	buffer := NewBuffer()
	buffer.AddHosts(hosts)

	// add all files to buffer
	for f := range files {
		buffer.addFile(files[f])
	}
	buffer.sortFile()

	// create a screen
	screen := NewScreen(files, buffer)

	// run the event loop
	screen.eventLoop()

	// cleanup
	termbox.Close()

}
