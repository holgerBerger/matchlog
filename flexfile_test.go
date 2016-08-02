package main

import "testing"

// TestFile is testing File Reading and indexing
func TestFlexFile(t *testing.T) {
	file, _ := ReadFlexFile("testdata/testfile")
	if file.linecount != 3 || len(file.lines) != 4 {
		t.Fail()
	}
	file, _ = ReadFlexFile("testdata/testfile.gz")
	if file.linecount != 3 || len(file.lines) != 4 || file.filename != "testdata/testfile.gz" {
		t.Fail()
	}
	//fmt.Println(string(file.lines[0]))
}

// TestParse checks date parsing
func TestFlexParse(t *testing.T) {
	file, _ := ReadFlexFile("testdata/log1")
	if file.times[0].String() != "2016-07-24 06:29:28 +0200 CEST" {
		t.Fail()
	}

	file, _ = ReadFlexFile("testdata/log2")
	if file.times[0].String() != "2016-07-26 00:36:17 +0200 CEST" ||
		file.times[1].String() != "2016-07-26 01:24:00 +0200 CEST" {
		t.Fail()
	}

	file, _ = ReadFlexFile("testdata/log3")
	if file.times[0].String() != "2016-08-02 19:54:09 +0200 CEST" ||
		file.times[1].String() != "2016-08-02 19:55:47 +0200 CEST" {
		t.Fail()
	}

}
