package main

import (
	"log"
    "os"

	"github.com/jroimartin/gocui"
    "github.com/rmcs9/benparser"
)

var benval benparser.Benval

func main() {

    if len(os.Args) != 2 {
        log.Fatal("a single bencode file argument must be provided")
    }
    benval = benparser.ParseFile(os.Args[1])
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()


	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

    if err := g.SetKeybinding("", 'j', gocui.ModNone, moveCursor(1)); err != nil {
        log.Panicln(err)
    }

    if err := g.SetKeybinding("", 'k', gocui.ModNone, moveCursor(-1)); err != nil {
        log.Panicln(err)
    }

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

var curs int = 0

func moveCursor(dy int) func(*gocui.Gui, *gocui.View) error {
    return func(g *gocui.Gui, v *gocui.View) error {
        if (curs == 0 && dy < 0) || (curs == len(v.BufferLines()) - 2 && dy > 0) || len(v.BufferLines()) == 1 {
            return nil
        }

        v.MoveCursor(0, dy, false)
        curs += dy

        cview, err := g.View("content")
        if err != nil {
            return err
        }

        iview, err := g.View("info")
        if err != nil {
            return err
        }
        target = index[curs] 
        cview.Clear()
        drawContent(cview, &benval, 1)
        drawInfo(iview)

        ox, oy := cview.Origin() 
        _, cy := cview.Size()
        if len(cview.BufferLines()) > cy {
            if err := cview.SetOrigin(ox, oy + dy); err != nil {
                return err
            }
        }
        return nil
    }
}

func layout(g *gocui.Gui) error {
    maxX, maxY := g.Size()
    if v, err := g.SetView("dir", 0, 0, 30, maxY - 6); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.SetCursor(0,0)
        v.SelFgColor = gocui.ColorBlack
        v.SelBgColor = gocui.ColorYellow
        v.Highlight = true

        drawDir(v, &benval, "")
    }

    if v, err := g.SetView("content", 31, 0, maxX -1, maxY -1); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.Title = "content"
        v.Wrap = true
        indexInit(&benval, false)
        target = index[0]
        drawContent(v, &benval, 1)
    }

    if v, err := g.SetView("info", 0, maxY - 5, 30, maxY - 1); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.Title = "info"
        drawInfo(v)
    }
    g.SetCurrentView("dir")
    return nil
}
