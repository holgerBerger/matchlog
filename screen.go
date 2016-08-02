package main

/*

	screen and keyboard handling
	using termbox, which is redrawing all screen, but simple,
	ncurses would probably need less display bandwidth

	(c) Holger Berger 2016, under GPL
*/

import (
	"sync"

	termbox "github.com/nsf/termbox-go"
)

// Screen class to keep state of termbox
type Screen struct {
	w, h    int          // screensize
	files   []*FlexFileT // list of files
	buffer  *BufferT     // buffer to be shown
	offsety int          // offset of forst line into buffer = 0 top of file 1 = second line on top of screen
	offsetx int
	lock    sync.Mutex
}

// NewScreen inits termbox
func NewScreen(files []*FlexFileT, buffer *BufferT) *Screen {
	newscreen := new(Screen)
	newscreen.files = files
	newscreen.buffer = buffer
	newscreen.offsety = 0
	newscreen.offsetx = 0

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
			if ev.Type == termbox.EventKey && (ev.Key == termbox.KeyEsc || ev.Ch == 'q') {
				exitQueue <- true
			}
			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyArrowDown {
				if s.offsety < s.buffer.linecount-s.h {
					s.offsety++
				}
			}
			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyArrowUp {
				if s.offsety > 0 {
					s.offsety--
				}
			}
			if ev.Type == termbox.EventKey && (ev.Key == termbox.KeyPgdn || ev.Key == termbox.KeySpace) {
				if s.offsety < s.buffer.linecount-s.h {
					s.offsety += s.h
				} else {
					s.offsety = s.buffer.linecount - s.h
				}
			}
			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyPgup {
				if s.offsety-s.h > 0 {
					s.offsety -= s.h
				} else {
					s.offsety = 0
				}
			}

			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyArrowRight {
				s.offsetx++
			}
			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyArrowLeft {
				if s.offsetx > 0 {
					s.offsetx--
				}
			}

			if ev.Type == termbox.EventResize {
				s.w, s.h = termbox.Size()
			}

			// something happened, so redraw screen
			s.draw()
		}
	}
}

// eventLoop will not return unless program is ended
func (s *Screen) eventLoop() {
	eventQueue := make(chan termbox.Event, 1)
	exitQueue := make(chan bool, 1)
	// handle events likle keypress and resize
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

	for y := 0; y < s.h; y++ {
		if y+s.offsety >= s.buffer.linecount {
			break
		}

		linep := s.buffer.lineps[y+s.offsety]
		var color termbox.Attribute

		// if line is not empty, match the line
		if len(s.buffer.lineps[y+s.offsety]) > 0 {
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
			termbox.SetCell(x, y, rune, color, termbox.ColorDefault)
		}
	}

	// full redraw
	termbox.Flush()
	s.lock.Unlock()
}
