package main

import (
	"fmt"
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

var curs int = 0

var cursorDown = func(g *gocui.Gui, v *gocui.View) error {
    if curs == len(v.BufferLines()) - 2 {
        return nil
    }
    v.MoveCursor(0, 1, false)
    curs++
    cview, err := g.View("content")
    if err != nil {
        log.Fatal(err)
    }
    iview, err := g.View("info")
    if err != nil {
        log.Fatal(err)
    }
    cview.Clear()
    target = index[curs]
    drawContent(cview, false)
    iview.Clear() 
    fmt.Fprintf(iview, "cursor: %d", curs)
    return nil
}

var cursorUp = func(g *gocui.Gui, v *gocui.View) error {
    if curs == 0 {
        return nil
    }
    v.MoveCursor(0, -1, false)
    curs--
    cview, err := g.View("content")
    if err != nil {
        log.Fatal(err)
    }
    iview, err := g.View("info")
    if err != nil {
        log.Fatal(err)
    }
    cview.Clear()
    target = index[curs]
    drawContent(cview, false)
    iview.Clear() 
    fmt.Fprintf(iview, "cursor: %d", curs)
    return nil
}


func layout(g *gocui.Gui) error {
    maxX, maxY := g.Size()
    if v, err := g.SetView("dir", 0, 0, 30, (maxY / 2) - 2); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.Title = "dir"
        v.SetCursor(0,0)
        v.SelFgColor = gocui.ColorBlack
        v.SelBgColor = gocui.ColorBlue
        v.Highlight = true

        switch benval.Kind() {
            case benparser.Map : 
                drawMapDir(v, benval, 1)
            case benparser.List:
                drawListDir(v, benval, 1)
            case benparser.Int :
                drawIntDir(v)
            case benparser.String: 
                drawStringDir(v)
        }
    }

    target = &benval

    if v, err := g.SetView("info", 0, (maxY / 2) - 1, 30, maxY - 1); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        fmt.Fprintln(v, target)
        v.Title = "info"
    }

    if v, err := g.SetView("content", 31, 0, maxX -1, maxY -1); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.Title = "content"
        v.Wrap = true
        drawContent(v, true)
    }

    // target = index[0]
    g.SetCurrentView("dir")
    return nil
}
