package main 


import (
    "fmt"
    "strings"

    "github.com/jroimartin/gocui"
    "github.com/rmcs9/benparser"
)


// ----------- DIRECTORY DRAW FUNCTIONS ---------------------
func drawMapDir(v *gocui.View, m benparser.Benval, lvl int) {
    benmap := m.(benparser.Benmap)

    fmt.Fprintf(v, strings.Repeat("\t", lvl - 1) + "%s\n", "d") 

    for _, key := range benmap.Keys() {
        fstring := strings.Repeat("\t", lvl) + "%s\n"

        fmt.Fprintf(v, fstring, key)

        if sub, _ := benmap.Query(key); (*sub).Kind() == benparser.Map {
            drawMapDir(v, *sub, lvl + 1) 
        } else if (*sub).Kind() == benparser.List {
            drawListDir(v, *sub, lvl + 1)
        }

    }

    fmt.Fprintf(v, strings.Repeat("\t", lvl - 1) + "%s\n", "e") 
}

func drawListDir(v *gocui.View, l benparser.Benval, lvl int) {
    benlist := l.(benparser.Benlist) 

    fmt.Fprintf(v, strings.Repeat("\t", lvl - 1) + "%s\n", "l")

    fstring := strings.Repeat("\t", lvl) + "[%d]\n"
    for i := range benlist.Len() {
        fmt.Fprintf(v, fstring, i)

        if sub := benlist.Get(i); (*sub).Kind() == benparser.Map {
            drawMapDir(v, *sub, lvl + 1) 
        } else if (*sub).Kind() == benparser.List {
            drawListDir(v, *sub, lvl + 1)
        }
    }

    fmt.Fprintf(v, strings.Repeat("\t", lvl - 1) + "%s\n", "e")
}

func drawIntDir(v *gocui.View) {
    benint := benval.(benparser.Benint) 

    fmt.Fprintf(v, "%d", benint.Get())
}

func drawStringDir(v *gocui.View) {
    benstring := benval.(benparser.Benstring) 

    fmt.Fprintf(v, "byte string SIZE: %d", len(benstring.Get()))
}
