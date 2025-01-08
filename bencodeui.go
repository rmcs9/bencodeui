package main

import (
	"fmt"
	"log"
    "os"

	"github.com/jroimartin/gocui"
    "github.com/rmcs9/benparser"
)

var benmap benparser.Benmap

func main() {

    if len(os.Args) != 2 {
        log.Fatal("a single bencode file argument must be provided")
    }
    ben := benparser.ParseFile(os.Args[1])
    benmap = ben.(benparser.Benmap)
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

    if err := g.SetKeybinding("", 'j', gocui.ModNone, cursorDown); err != nil {
        log.Panicln(err)
    }

    if err := g.SetKeybinding("", 'k', gocui.ModNone, cursorUp); err != nil {
        log.Panicln(err)
    }

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

var cursorDown = func(g *gocui.Gui, v *gocui.View) error {
    x, y := v.Cursor()
    v.SetCursor(x, y-1)       
    return nil
}

var cursorUp = func(g *gocui.Gui, v *gocui.View) error {
    x,y := v.Cursor()
    v.SetCursor(x, y+1)
    return nil
}

func layout(g *gocui.Gui) error {
    maxX, maxY := g.Size()
    if v, err := g.SetView("dir", 0, 0, 30, (maxY / 2) - 2); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.Title = "dir"
        v.SetCursor(0,1)
        v.SelFgColor = gocui.ColorBlack
        v.SelBgColor = gocui.ColorBlue
        v.Highlight = true

        for _, key := range benmap.Keys() {
            fmt.Fprintf(v, "%s\n", key)
        }
    }

    if v, err := g.SetView("info", 0, (maxY / 2) - 1, 30, maxY - 1); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.Title = "info"
    }

    if v, err := g.SetView("content", 31, 0, maxX -1, maxY -1); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.Title = "content"
        fmt.Fprintf(v, "right panel. maxX is: %d MaxY is: %d", maxX, maxY)
    }

    g.SetCurrentView("dir")
    return nil
}
