package main

import "testing"

// TestFile is testing File Reading and indexing
func TestFile(t *testing.T) {
	file, _ := ReadFile("testdata/testfile")
	if file.linecount != 3 || len(file.lines) != 4 {
		t.Fail()
	}
}

// TestParse checks date parsing
func TestParse(t *testing.T) {
	file, _ := ReadFile("testdata/log1")
	if file.times[0].String() != "2016-07-24 06:29:28 +0200 CEST" {
		t.Fail()
	}

	file, _ = ReadFile("testdata/log2")
	if file.times[0].String() != "2016-07-26 00:36:17 +0200 CEST" ||
		file.times[1].String() != "2016-07-26 01:24:00 +0200 CEST" {
		t.Fail()
	}
}
