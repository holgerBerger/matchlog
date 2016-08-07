package main

import "testing"

func TestHosts(t *testing.T) {
	hosts := NewHosts()
	hosts.AddFile("testdata/hosts")
	if len(hosts.ip2name) != 3 || hosts.ip2name["2.3.4.5"] != "NAME1" || hosts.ip2name["1.2.3.4"] != "name1" {
		t.Fail()
	}
}
