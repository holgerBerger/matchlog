package main

import "testing"

// TestFile is testing File Reading and indexing
func TestFile(t *testing.T) {
	file, _ := ReadFile("testdata/testfile")
	if file.linecount != 3 || len(file.lines) != 4 {
		t.Fail()
	}
}
