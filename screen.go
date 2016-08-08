package main

/*

	screen and keyboard handling
	using termbox, which is redrawing all screen but simple,
	ncurses would probably need less display bandwidth
	cpu usage minimalized, only redrawing after events

	(c) Holger Berger 2016, under GPL
*/

import (
	"fmt"
	"strings"
	"sync"

	termbox "github.com/nsf/termbox-go"
)

// Screen class to keep state of termbox
type Screen struct {
	w, h          int          // screensize
	files         []*FlexFileT // list of files
	buffer        *BufferT     // buffer to be shown
	offsety       int          // offset of forst line into buffer = 0 top of file 1 = second line on top of screen
	offsetx       int          // offset of first character in line to be shown
	searchInput   bool         // flag if search input is ongoing
	searchForward bool         // search direction, false = backward
	searchString  string       // string to be searched
	lock          sync.Mutex
}

// NewScreen inits termbox
func NewScreen(files []*FlexFileT, buffer *BufferT) *Screen {
	newscreen := new(Screen)
	newscreen.files = files
	newscreen.buffer = buffer
	newscreen.offsety = 0
	newscreen.offsetx = 0
	newscreen.searchInput = false

	// init termbox
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	// defer termbox.Close()

	newscreen.w, newscreen.h = termbox.Size()
	termbox.SetOutputMode(termbox.Output256)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	return newscreen
}

// eventHandler catches events and changes state
func (s *Screen) eventHandler(eventQueue chan termbox.Event, exitQueue chan bool) {

	for {
		select {
		case ev := <-eventQueue:

			// handling of search

			if ev.Type == termbox.EventKey && ev.Ch == '/' {
				termbox.SetCursor(1, s.h-1)
				s.searchInput = true
				s.searchForward = true
				s.draw()
				break
			}

			if ev.Type == termbox.EventKey && ev.Ch == '?' {
				termbox.SetCursor(1, s.h-1)
				s.searchInput = true
				s.searchForward = false
				s.draw()
				break
			}

			if s.searchInput && ev.Type == termbox.EventKey && ev.Key == termbox.KeyEnter {
				s.searchInput = false
				termbox.HideCursor()
				// perform search
				if s.searchForward {
					for lnr, currentline := range s.buffer.lines[s.offsety+1:] {
						// FIXME case insensitive search?!?!
						if strings.Index(string(currentline.line), s.searchString) > 0 {
							s.offsety += lnr + 1
							break
						}
					}
				} else {
					// backward
					for lnr := s.offsety - 1; lnr >= 0; lnr-- {
						// FIXME case insensitive search?!?!
						currentline := s.buffer.lines[lnr]
						if strings.Index(string(currentline.line), s.searchString) > 0 {
							s.offsety = lnr
							break
						}
					}
				}
				s.draw()
				break
			}

			if s.searchInput && ev.Type == termbox.EventKey && (ev.Key == termbox.KeyBackspace || ev.Key == termbox.KeyBackspace2 || ev.Key == termbox.KeyDelete) {
				if len(s.searchString) > 0 {
					s.searchString = s.searchString[:len(s.searchString)-1]
				}
				s.draw()
				break
			}

			if s.searchInput && ev.Type == termbox.EventKey {
				if ev.Key != termbox.KeyArrowDown &&
					ev.Key != termbox.KeyArrowUp &&
					ev.Key != termbox.KeyArrowRight &&
					ev.Key != termbox.KeyArrowLeft {
					s.searchString = s.searchString + string(ev.Ch)
				}
				if ev.Key == termbox.KeyEsc {
					s.searchString = ""
					s.searchInput = false
					termbox.HideCursor()
				}
				s.draw()
				break
			}

			// all other keys

			if ev.Type == termbox.EventKey && (ev.Key == termbox.KeyEsc || ev.Ch == 'q') {
				exitQueue <- true
			}

			if ev.Type == termbox.EventKey && (ev.Key == termbox.KeyHome || ev.Ch == '0') {
				s.offsety = 0
				s.draw()
				break
			}

			if ev.Type == termbox.EventKey && (ev.Key == termbox.KeyEnd || ev.Ch == 'G') {
				s.offsety = s.buffer.linecount - s.h
				s.draw()
				break
			}

			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyArrowDown {
				if s.offsety < s.buffer.linecount-s.h {
					s.offsety++
				}
				s.draw()
				break
			}
			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyArrowUp {
				if s.offsety > 0 {
					s.offsety--
				}
				s.draw()
				break
			}
			if ev.Type == termbox.EventKey && (ev.Key == termbox.KeyPgdn || ev.Key == termbox.KeySpace) {
				if s.offsety < s.buffer.linecount-s.h {
					s.offsety += s.h
				} else {
					s.offsety = s.buffer.linecount - s.h
				}
				s.draw()
				break
			}
			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyPgup {
				if s.offsety-s.h > 0 {
					s.offsety -= s.h
				} else {
					s.offsety = 0
				}
				s.draw()
				break
			}

			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyArrowRight {
				s.offsetx++
				s.draw()
				break
			}
			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyArrowLeft {
				if s.offsetx > 0 {
					s.offsetx--
				}
				s.draw()
				break
			}

			if ev.Type == termbox.EventResize {
				s.w, s.h = termbox.Size()
				s.draw()
			}
		}
	}
}

// eventLoop will not return unless program is ended
func (s *Screen) eventLoop() {
	eventQueue := make(chan termbox.Event, 1)
	exitQueue := make(chan bool, 1)
	// handle events like keypress and resize
	go s.eventHandler(eventQueue, exitQueue)

	// catch events and send to event handler
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	// inital draw
	s.draw()

	// endless loop, waiting for end event
loop:
	for {
		select {
		case <-exitQueue:
			break loop
		}
	}
}

// draw paints whatever is needed
func (s *Screen) draw() {
	// as we use go routines, we need a lock here,
	// to avoid a redraw triggered by goroutine to interfere
	// with another, and we protect termbox.Flush at the same time
	s.lock.Lock()

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	for y := 0; y < s.h-1; y++ {
		if y+s.offsety >= s.buffer.linecount {
			break
		}

		linep := s.buffer.lines[y+s.offsety].line
		var color termbox.Attribute

		// if line is not empty, match the line
		if len(s.buffer.lines[y+s.offsety].line) > 0 {
			color = s.buffer.rules.Match(linep)
		} else {
			color = termbox.ColorDefault
		}

		// render the line
		for x := 0; x < s.w; x++ {
			if x+s.offsetx >= len(linep) || linep[x+s.offsetx] == '\n' {
				break
			}
			rune := rune(linep[x+s.offsetx])
			if x+s.offsetx >= s.buffer.lines[y+s.offsety].hoststart &&
				x+s.offsetx < s.buffer.lines[y+s.offsety].hostend {
				hostname := string(linep[s.buffer.lines[y+s.offsety].hoststart:s.buffer.lines[y+s.offsety].hostend])
				termbox.SetCell(x, y, rune, s.buffer.hostcolors[hostname], termbox.ColorDefault)
			} else {
				termbox.SetCell(x, y, rune, color, termbox.ColorDefault)
			}
		}
	}

	// status line
	for x := 0; x < s.w; x++ {
		termbox.SetCell(x, s.h-1, ' ', termbox.ColorBlack|termbox.AttrReverse, termbox.ColorDefault)
	}

	// ruler
	ruler := int(float32(s.offsety) / float32(s.buffer.linecount-s.h) * 100.0)
	rulerstring := fmt.Sprintf("%7d/%7d %3d%%", s.offsety, s.buffer.linecount, ruler)
	for x := 0; x < len(rulerstring); x++ {
		termbox.SetCell(s.w-22+x, s.h-1, rune(rulerstring[x]), termbox.ColorBlack|termbox.AttrReverse, termbox.ColorDefault)
	}

	// input mode for search
	if s.searchInput {
		if s.searchForward {
			termbox.SetCell(0, s.h-1, '/', termbox.ColorBlack|termbox.AttrReverse, termbox.ColorDefault)
		} else {
			termbox.SetCell(0, s.h-1, '?', termbox.ColorBlack|termbox.AttrReverse, termbox.ColorDefault)
		}
		for x := 0; x < len(s.searchString); x++ {
			termbox.SetCell(1+x, s.h-1, rune(s.searchString[x]), termbox.ColorBlack|termbox.AttrReverse, termbox.ColorDefault)
		}
		termbox.SetCursor(len(s.searchString)+1, s.h-1)
	} else {
		// helpstring
		helpstring := "matchlog - ESC,q: quit  /,?: search  Home,End,PgUp,PgDown,Up,Down,Left,Right: navigation"
		for x := 0; x < len(helpstring); x++ {
			termbox.SetCell(1+x, s.h-1, rune(helpstring[x]), termbox.ColorBlack|termbox.AttrReverse, termbox.ColorDefault)
		}
	}

	// full redraw
	termbox.Flush()
	s.lock.Unlock()
}
