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
		fmt.Println(string(buffer.lines[0].line), len(buffer.lines[0].line))
		fmt.Println(string(buffer.lines[1].line), len(buffer.lines[1].line))
		fmt.Println(string(buffer.lines[2].line), len(buffer.lines[2].line))
	*/

	// check length of lines
	if len(buffer.lines[0].line) != 49 || len(buffer.lines[1].line) != 122 {
		t.Fail()
	}

	// fmt.Println(bytes.Index(buffer.lines[2].line, []byte("node@o2ib4")))

	// check IP replacement
	if bytes.Index(buffer.lines[2].line, []byte("node@o2ib4")) != 118 {
		t.Fail()
	}

	// fmt.Println(buffer.lines[0].hoststart, buffer.lines[0].hostend)

	// check location of colored hostname
	if buffer.lines[0].hoststart != 20 || buffer.lines[0].hostend != 27 {
		t.Fail()
	}

}
