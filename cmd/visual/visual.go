package main

import (
	"log"

	"github.com/gdamore/tcell"
	"github.com/hinshun/astar/cmd/visual/terrain"
	"github.com/hinshun/astar/cmd/visual/view"
	"github.com/pkg/errors"
)

func loop(screen tcell.Screen, v *view.View) error {
	for {
		screenEvent := screen.PollEvent()
		switch event := screenEvent.(type) {
		case *tcell.EventKey:
			switch event.Key() {
			case tcell.KeyCtrlC:
				// Exit
				return nil
			default:
				switch event.Rune() {
				case 'q':
					// Exit
					return nil
				case 'c':
					v.Clear()
				case 'p':
					v.SetMode(view.ModePlain)
				case 's':
					v.SetMode(view.ModeSwamp)
				case 'w':
					v.SetMode(view.ModeWall)
				case 'a':
					v.SetMode(view.ModeAgent)
				case 'g':
					v.SetMode(view.ModeGoal)
				}
			}
		case *tcell.EventMouse:
			if event.Buttons()&tcell.Button1 != 0 {
				x, y := event.Position()
				v.Click(x, y)
			}
		case *tcell.EventResize:
			width, height := event.Size()
			v.Resize(width, height)
		}
		v.Update()
	}
}

func run() error {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	screen, err := tcell.NewScreen()
	if err != nil {
		return errors.Wrap(err, "failed to create new screen")
	}

	err = screen.Init()
	if err != nil {
		return errors.Wrap(err, "failed to initialize screen")
	}
	defer screen.Fini()

	screen.EnableMouse()
	screen.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlack))
	screen.Clear()

	grid, err := terrain.NewGrid("shard1", "W28N48")
	if err != nil {
		return errors.Wrap(err, "failed to create new grid")
	}

	v := view.NewView(screen, grid)
	v.Update()

	err = loop(screen, v)
	if err != nil {
		return errors.Wrap(err, "exited loop due to error")
	}

	return nil
}

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}
