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
}

// NewScreen inits termbox
func NewScreen(files []*FileT, buffer *BufferT) *Screen {
	newscreen := new(Screen)
	newscreen.files = files
	newscreen.buffer = buffer

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

	offset := 0

loop:
	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
				break loop
			}
			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyArrowDown {
				offset++
			}
			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyArrowUp {
				offset--
			}
		default:
			s.draw()
			time.Sleep(10 * time.Millisecond)
		}
	}

}

// draw paints whatever is needed
func (s *Screen) draw() {

	fileid := 0 // FIXME

	for y := 0; y < s.h; y++ {
		if y >= s.files[fileid].linecount {
			break
		}
		for x := 0; x < s.w; x++ {
			if x >= s.files[fileid].lines[y+1]-s.files[fileid].lines[y] {
				break
			}
			linep := s.files[fileid].lines[y]
			if linep+x >= len(s.files[fileid].contents) {
				break
			}
			rune := rune(s.files[fileid].contents[linep+x])
			termbox.SetCell(x, y, rune, termbox.ColorBlack, termbox.ColorWhite)
		}
	}

	// full redraw
	termbox.Flush()
}
