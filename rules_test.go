package main

import (
	"testing"

	termbox "github.com/nsf/termbox-go"
)

func TestRules(t *testing.T) {
	rules := DefaultRules()
	if rules.Match([]byte("failure on deck 13")) != termbox.ColorRed {
		t.Fail()
	}
	if rules.Match([]byte("heavy weather warning")) != termbox.ColorYellow {
		t.Fail()
	}
	if rules.Match([]byte("Failure after WARNING")) != termbox.ColorRed {
		t.Fail()
	}
	if rules.Match([]byte("boring text, no matches, ok?")) != termbox.ColorGreen {
		t.Fail()
	}
	if rules.Match([]byte("this time really nothing.")) != termbox.ColorDefault {
		t.Fail()
	}
}
