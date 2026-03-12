package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

func (e *Editor) Render() {

	e.screen.Clear()

	top := e.view.topLine
	left := e.view.leftCol
	w, h := e.screen.Size()

	realH := h - 1

	lineNumberWidth := 4

	for screenY := 0; screenY < realH; screenY++ {
		bufY := top + screenY
		if bufY >= len(e.buffer.lines) {
			break
		}

		line := e.buffer.lines[bufY]

		num := fmt.Sprintf("%3d", bufY+1)

		for i, r := range num {
			e.screen.SetContent(i, screenY, r, nil, tcell.StyleDefault)
		}

		for screenX := lineNumberWidth; screenX < w; screenX++ {
			bufX := left + screenX - lineNumberWidth
			ch := ' '
			if bufX >= 0 && bufX < len(line) {
				ch = line[bufX]
			}
			e.screen.SetContent(screenX, screenY, ch, nil, tcell.StyleDefault)
		}
	}

	e.DisplayCommandLine()

	if e.mode == ModeCommand && e.cmd != nil {
		e.screen.ShowCursor(1+e.cmd.cursorX, h-1)
	} else {

		cx := e.cursorX - left
		cy := e.cursorY - top
		if cy >= 0 && cy < realH && cx >= 0 && cx < w {
			e.screen.ShowCursor(cx+lineNumberWidth, cy)
		} else {
			e.screen.HideCursor()
		}

	}

	e.screen.Show()

}
