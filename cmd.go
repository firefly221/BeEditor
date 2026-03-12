package main

import (
	"strings"

	"github.com/gdamore/tcell/v2"
)

type CommandLine struct {
	line    []rune
	cursorX int
	leftCol int
}

func (e *Editor) DisplayCommandLine() {
	w, h := e.screen.Size()
	y := h - 1
	style := tcell.StyleDefault

	for x := 0; x < w; x++ {
		e.screen.SetContent(x, y, ' ', nil, style)
	}

	if e.mode == ModeCommand {
		e.screen.SetContent(0, y, ':', nil, style)

		if e.cmd == nil {
			return
		}

		for i, r := range e.cmd.line {
			x := 1 + i
			if x >= w {
				break
			}
			e.screen.SetContent(x, y, r, nil, style)
		}
		return
	}

	for i, r := range e.status {
		if i >= w {
			break
		}
		e.screen.SetContent(i, y, r, nil, style)
	}
}

func (e *Editor) executeCommand(cmd string) {
	cmd = strings.TrimSpace(cmd)

	switch {
	case cmd == "quit":
		e.quit = true

	case cmd == "save":
		if err := e.Save(); err != nil {
			e.SetStatus(err.Error())
		} else {
			e.SetStatus("plik zapisany")
		}

	case strings.HasPrefix(cmd, "saveas "):
		path := strings.TrimSpace(strings.TrimPrefix(cmd, "saveas "))
		if err := e.SaveAs(path); err != nil {
			e.SetStatus(err.Error())
		} else {
			e.SetStatus("plik zapisany jako " + path)
		}

	case cmd == "savequit":
		if err := e.Save(); err != nil {
			e.SetStatus(err.Error())
		} else {
			e.quit = true
		}
	case strings.HasPrefix(cmd, "open "):
		path := strings.TrimSpace(strings.TrimPrefix(cmd, "open "))
		if err := e.Open(path); err != nil {
			e.SetStatus(err.Error())
		} else {
			e.SetStatus("otwarto " + path)
		}

	default:
		e.SetStatus("nieznana komenda: " + cmd)
	}
}
