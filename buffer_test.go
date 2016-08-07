package main

// (c) Holger Berger 2016, under GPL

import (
	"bytes"
	"testing"
)

// TestAddFile tests if files are sorted correctly
func TestAddFile(t *testing.T) {

	file1, _ := ReadFlexFile("testdata/log1")
	file2, _ := ReadFlexFile("testdata/log1_b")
	file3, _ := ReadFlexFile("testdata/log2")

	hosts := NewHosts()
	hosts.AddFile("testdata/hosts")

	buffer := NewBuffer()
	buffer.AddHosts(hosts)
	buffer.addFile(file1)
	buffer.addFile(file2)
	buffer.addFile(file3)
	buffer.sortFile()

	/*
		fmt.Println(string(buffer.lines[0].line))
		fmt.Println(string(buffer.lines[1].line))
		fmt.Println(string(buffer.lines[2].line))
	*/

	if len(buffer.lines[0].line) != 45 || len(buffer.lines[1].line) != 135 {
		t.Fail()
	}

	if bytes.Index(buffer.lines[2].line, []byte("node@o2ib4")) != 131 {
		t.Fail()
	}

}
