package main

import (
	"time"

	termbox "github.com/nsf/termbox-go"
)

// Screen class to keep state of termbox
type Screen struct {
	w, h   int      // screensize
	files  []*FileT // list of files
	buffer *BufferT // buffer to be shown
	offset int      // offset of forst line into buffer = 0 top of file 1 = second line on top of screen
}

// NewScreen inits termbox
func NewScreen(files []*FileT, buffer *BufferT) *Screen {
	newscreen := new(Screen)
	newscreen.files = files
	newscreen.buffer = buffer
	newscreen.offset = 0

	// init termbox
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	// defer termbox.Close()

	newscreen.w, newscreen.h = termbox.Size()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	return newscreen
}

// eventLoop will not return unless program is ended
func (s *Screen) eventLoop() {
	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	s.draw()

loop:
	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
				break loop
			}
			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyArrowDown {
				if s.offset < s.buffer.linecount-s.h {
					s.offset++
				}
			}
			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyArrowUp {
				if s.offset > 0 {
					s.offset--
				}
			}
			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyPgdn {
				if s.offset < s.buffer.linecount-s.h {
					s.offset += s.h
				} else {
					s.offset = s.buffer.linecount - s.h
				}
			}
			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyPgup {
				if s.offset-s.h > 0 {
					s.offset -= s.h
				} else {
					s.offset = 0
				}
			}
		default:
			s.draw()
			time.Sleep(10 * time.Millisecond)
		}
	}

}

// draw paints whatever is needed
func (s *Screen) draw() {

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	s.w, s.h = termbox.Size()

	// FIXME how to deal with long lines > w ?
	for y := 0; y < s.h; y++ {
		if y+s.offset >= s.buffer.linecount {
			break
		}
		for x := 0; x < s.w; x++ {
			linep := s.buffer.lineps[y+s.offset]
			cont := s.buffer.contps[y+s.offset]
			if (*cont)[linep+x] == 10 {
				break
			}
			rune := rune((*cont)[linep+x])
			termbox.SetCell(x, y, rune, termbox.ColorBlack, termbox.ColorDefault)
		}
	}

	// full redraw
	termbox.Flush()

}
