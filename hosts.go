package main

import (
	"bufio"
	"os"
	"strings"
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
	reader := bufio.NewReader(osfile)
	for {
		nextline, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}
		if nextline[0] == '#' {
			continue
		}

		f := strings.Split(string(nextline), " ")
		if len(f) <= 1 {
			continue
		}

		ip := f[0]
		h.ip2name[ip] = strings.Trim(f[1], " \t\n")

	}
}
