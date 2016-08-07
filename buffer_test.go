package main

// (c) Holger Berger 2016, under GPL

import "testing"

// TestAddFile tests if files are sorted correctly
func TestAddFile(t *testing.T) {

	file1, _ := ReadFlexFile("testdata/log1")
	file2, _ := ReadFlexFile("testdata/log1_b")
	file3, _ := ReadFlexFile("testdata/log2")

	buffer := NewBuffer()

	buffer.addFile(file1)
	buffer.addFile(file2)
	buffer.addFile(file3)
	buffer.sortFile()

	/*
		println(len(buffer.lineps[0]))
		println(len(buffer.lineps[1]))
	*/

	if len(buffer.lines[0].line) != 45 || len(buffer.lines[1].line) != 135 {
		t.Fail()
	}

}
