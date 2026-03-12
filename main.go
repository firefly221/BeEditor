package main

import (
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
)

func main() {

	defStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorYellow)

	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}

	if err = s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	s.SetStyle(defStyle)
	s.Clear()

	// Safe quitting
	quit := func() {
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	myEditor := NewEditor(s)

	if len(os.Args) > 1 {
		if err := myEditor.Open(os.Args[1]); err != nil {
			myEditor.SetStatus(err.Error())
		}
	}

	// Main loop
	for !myEditor.quit {

		ev := myEditor.screen.PollEvent()

		switch tev := ev.(type) {
		case *tcell.EventKey:
			myEditor.HandleKeys(tev)
		case *tcell.EventResize:
			myEditor.HandleResize(tev)

		}

		myEditor.Render()

	}

}
