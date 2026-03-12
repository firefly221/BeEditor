package main

import (
	"errors"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type Mode int

const (
	ModeInsert Mode = iota
	ModeCommand
)

type Editor struct {
	quit             bool
	mode             Mode
	cursorX, cursorY int
	buffer           *Buffer
	view             *View
	screen           tcell.Screen
	file             *File
	cmd              *CommandLine
	status           string
}

type Buffer struct {
	lines [][]rune
}

type View struct {
	topLine int
	leftCol int
}

func NewEditor(screen tcell.Screen) *Editor {

	buffer := &Buffer{
		lines: [][]rune{{}},
	}

	view := &View{
		topLine: 0,
		leftCol: 0,
	}

	newFile := &File{
		Path:       "",
		IsModified: false,
	}

	return &Editor{
		file:    newFile,
		cursorX: 0,
		cursorY: 0,
		buffer:  buffer,
		view:    view,
		screen:  screen,
		quit:    false,
	}
}

func (e *Editor) Scroll() {
	w, h := e.screen.Size()
	textHeight := h - 1

	if e.cursorY >= e.view.topLine+textHeight {
		e.view.topLine = e.cursorY - textHeight + 1
	}

	if e.cursorY < e.view.topLine {
		e.view.topLine = e.cursorY
	}

	if e.cursorX >= e.view.leftCol+w {
		e.view.leftCol = e.cursorX - w + 1
	}

	if e.cursorX < e.view.leftCol {
		e.view.leftCol = e.cursorX
	}
}

func (e *Editor) InsertRune(r rune) {
	line := e.buffer.lines[e.cursorY]

	newLine := append(line[:e.cursorX],
		append([]rune{r}, line[e.cursorX:]...)...)

	e.buffer.lines[e.cursorY] = newLine
	e.cursorX++
}

func (e *Editor) InsertNewLine() {
	x := e.cursorX
	y := e.cursorY

	line := e.buffer.lines[y]

	left := append([]rune(nil), line[:x]...)
	right := append([]rune(nil), line[x:]...)

	e.buffer.lines[y] = left

	e.buffer.lines = append(
		e.buffer.lines[:y+1],
		append([][]rune{right}, e.buffer.lines[y+1:]...)...,
	)

	e.cursorY++
	e.cursorX = 0
}

func (e *Editor) Backspace() {

	if e.cursorY == 0 && e.cursorX == 0 {
		return
	}

	y := e.cursorY
	x := e.cursorX

	if x > 0 {
		line := e.buffer.lines[y]

		e.buffer.lines[y] = append(line[:x-1], line[x:]...)
		e.cursorX--
		return
	}

	prevY := y - 1
	prev := e.buffer.lines[prevY]
	cur := e.buffer.lines[y]

	joined := make([]rune, 0, len(prev)+len(cur))
	joined = append(joined, prev...)
	joined = append(joined, cur...)

	e.buffer.lines[prevY] = joined

	e.buffer.lines = append(e.buffer.lines[:y], e.buffer.lines[y+1:]...)

	e.cursorY = prevY
	e.cursorX = len(prev)
}

func (e *Editor) Open(path string) error {
	e.view.topLine, e.view.leftCol = 0, 0
	e.cmd = nil
	e.mode = ModeInsert

	e.file.Path = path

	lines, err := e.file.LoadFile()
	if err != nil {
		return err
	}

	if len(lines) == 0 {
		lines = [][]rune{{}}
	}
	e.buffer.lines = lines

	e.cursorX, e.cursorY = 0, 0
	return nil
}

func (e *Editor) Save() error {
	if strings.TrimSpace(e.file.Path) == "" {
		return errors.New("no file path (use SaveAs)")
	}

	if err := e.file.SaveFile(e.buffer.lines); err != nil {
		return err
	}

	return nil
}

func (e *Editor) SaveAs(path string) error {
	e.file.Path = path
	if err := e.file.SaveFile(e.buffer.lines); err != nil {
		return err
	}

	return nil
}

func (e *Editor) New() {
	e.file.Path = ""
	e.file.IsModified = false
	e.buffer.lines = [][]rune{{}}
	e.cursorX, e.cursorY = 0, 0
}

func (e *Editor) HandleKeys(ev *tcell.EventKey) {

	if e.mode != ModeCommand && e.status != "" {
		e.ClearStatus()
	}

	switch e.mode {

	case ModeInsert:

		switch ev.Key() {

		case tcell.KeyEscape:
			e.mode = ModeCommand
			e.cmd = &CommandLine{
				line:    []rune{},
				cursorX: 0,
			}
			return

		case tcell.KeyUp:
			if e.cursorY > 0 {
				e.cursorY--
				for e.cursorX > len(e.buffer.lines[e.cursorY]) {
					e.cursorX--
				}
			}

		case tcell.KeyDown:
			if e.cursorY < len(e.buffer.lines)-1 {
				e.cursorY++
				for e.cursorX > len(e.buffer.lines[e.cursorY]) {
					e.cursorX--
				}
			}

		case tcell.KeyLeft:
			if e.cursorX > 0 {
				e.cursorX--
			}

		case tcell.KeyRight:
			if e.cursorX < len(e.buffer.lines[e.cursorY]) {
				e.cursorX++
			}

		case tcell.KeyRune:
			e.InsertRune(ev.Rune())

		case tcell.KeyEnter:
			e.InsertNewLine()

		case tcell.KeyBackspace:
			e.Backspace()

		case tcell.KeyCtrlS:
			e.SaveAs(".txt")

		case tcell.KeyCtrlN:
			e.New()

		}

		e.Scroll()

	case ModeCommand:
		switch ev.Key() {

		case tcell.KeyEscape:
			e.mode = ModeInsert
			e.cmd = nil
		case tcell.KeyEnter:
			e.executeCommand(string(e.cmd.line))
			e.mode = ModeInsert
			e.cmd = nil
		case tcell.KeyBackspace:
			if e.cmd.cursorX > 0 {
				e.cmd.line = append(e.cmd.line[:e.cmd.cursorX-1],
					e.cmd.line[e.cmd.cursorX:]...)
				e.cmd.cursorX--
			}
		default:
			if ev.Rune() != 0 {
				e.cmd.line = append(e.cmd.line[:e.cmd.cursorX],
					append([]rune{ev.Rune()}, e.cmd.line[e.cmd.cursorX:]...)...)
				e.cmd.cursorX++
			}

		}
	}

}

func (e *Editor) SetStatus(msg string) {
	e.status = msg
}

func (e *Editor) ClearStatus() {
	e.status = ""
}
func (e *Editor) HandleResize(ev *tcell.EventResize) {
	//Handling Resize
}
