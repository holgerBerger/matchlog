package main

import (
	"bufio"
	"os"
	"regexp"
)

// HostsT represesents a hostfile
type HostsT struct {
	ip2name map[string]string
}

// NewHosts initializes a Host object (no file loaded)
func NewHosts() *HostsT {
	h := HostsT{}
	h.ip2name = make(map[string]string)
	return &h
}

// AddFile adds a file to a Hosts object
func (h *HostsT) AddFile(filename string) {
	osfile, err := os.Open(filename)
	if err != nil {
		panic("can not read hostfile.")
	}

	regexp, _ := regexp.Compile(`\s*([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)\s+(\S*).*`)

	reader := bufio.NewReader(osfile)
	for {
		nextline, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}
		if nextline[0] == '#' {
			continue
		}

		m := regexp.FindSubmatch(nextline)

		if m != nil && len(m) > 2 {
			ip := string(m[1])
			name := string(m[2])
			h.ip2name[ip] = name
		}

	}
}
