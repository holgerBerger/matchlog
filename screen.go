package main

import termbox "github.com/nsf/termbox-go"

// Screen class to keep state of termbox
type Screen struct {
	w, h int
}

// NewScreen inits termbox
func NewScreen() *Screen {
	newscreen := new(Screen)

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
			//		default:
			//			s.draw()
			//			time.Sleep(10 * time.Millisecond)
		}
	}

}

// draw paints whatever is needed
func (s *Screen) draw() {

	fileid := 0 // FIXME

	for y := 0; y < s.h; y++ {
		if y >= files[fileid].linecount {
			break
		}
		for x := 0; x < s.w; x++ {
			if x >= files[fileid].lines[y+1]-files[fileid].lines[y] {
				break
			}
			linep := files[fileid].lines[y]
			if linep+x >= len(files[fileid].contents) {
				break
			}
			rune := rune(files[fileid].contents[linep+x])
			termbox.SetCell(x, y, rune, termbox.ColorBlack, termbox.ColorWhite)
		}
	}

	termbox.Flush()
}
