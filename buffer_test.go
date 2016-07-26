package main

import "testing"

// TestAddFile tests if files are sorted correctly
func TestAddFile(t *testing.T) {

	file1, _ := ReadFile("testdata/log1")
	file2, _ := ReadFile("testdata/log1_b")
	file3, _ := ReadFile("testdata/log2")

	buffer := NewBuffer()

	buffer.addFile(file1)
	buffer.addFile(file2)
	buffer.addFile(file3)
	buffer.sortFile()

	if buffer.lineps[2] != 135 || buffer.lineps[3] != 0 {
		t.Fail()
	}

}
