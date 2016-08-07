package main

// (c) Holger Berger 2016, under GPL

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
	if file.lines[0].time.String() != "2016-07-24 06:29:28 +0200 CEST" {
		t.Fail()
	}

	file, _ = ReadFlexFile("testdata/log2")

	if file.lines[0].time.String() != "2016-07-26 00:36:17 +0200 CEST" ||
		file.lines[1].time.String() != "2016-07-26 01:24:00 +0200 CEST" ||
		file.lines[0].host != "al2oss1" || file.lines[0].hoststart != 33 ||
		file.lines[0].hostend != 40 {
		t.Fail()
	}

	file, _ = ReadFlexFile("testdata/log3")

	if file.lines[0].time.String() != "2016-08-02 19:54:09 +0200 CEST" ||
		file.lines[1].time.String() != "2016-08-02 19:55:47 +0200 CEST" ||
		file.lines[0].host != "derwat" || file.lines[0].hoststart != 16 ||
		file.lines[0].hostend != 22 {
		t.Fail()
	}

}
