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

    fmt.Fprintf(v, strings.Repeat("\t\t", lvl - 1) + "%s\n", "d") 

    for _, key := range benmap.Keys() {
        fstring := strings.Repeat("\t\t", lvl) + "%s\n"

        fmt.Fprintf(v, fstring, key)

        if sub, _ := benmap.Query(key); (*sub).Kind() == benparser.Map {
            drawMapDir(v, *sub, lvl + 1) 
        } else if (*sub).Kind() == benparser.List {
            drawListDir(v, *sub, lvl + 1)
        }

    }

    fmt.Fprintf(v, strings.Repeat("\t\t", lvl - 1) + "%s\n", "e") 
}

func drawListDir(v *gocui.View, l benparser.Benval, lvl int) {
    benlist := l.(benparser.Benlist) 

    fmt.Fprintf(v, strings.Repeat("\t\t", lvl - 1) + "%s\n", "l")

    fstring := strings.Repeat("\t\t", lvl) + "[%d]\n"
    for i := range benlist.Len() {
        fmt.Fprintf(v, fstring, i)

        if sub := benlist.Get(i); (*sub).Kind() == benparser.Map {
            drawMapDir(v, *sub, lvl + 1) 
        } else if (*sub).Kind() == benparser.List {
            drawListDir(v, *sub, lvl + 1)
        }
    }

    fmt.Fprintf(v, strings.Repeat("\t\t", lvl - 1) + "%s\n", "e")
}

func drawIntDir(v *gocui.View) {
    benint := benval.(benparser.Benint) 

    fmt.Fprintf(v, "%d", benint.Get())
}

func drawStringDir(v *gocui.View) {
    benstring := benval.(benparser.Benstring) 

    fmt.Fprintf(v, "byte string SIZE: %d", len(benstring.Get()))
}

// ----------- CONTENT WINDOW DRAW FUNCTIONS ------------------

func drawMapContent(v *gocui.View, m benparser.Benval, lvl int) {
    benmap := m.(benparser.Benmap)

    fmt.Fprintf(v, strings.Repeat("\t\t", lvl - 1) + "%s\n", "d")
    for _, key := range benmap.Keys() {
        fmt.Fprint(v, strings.Repeat("\t\t", lvl))
        fmt.Fprintf(v, "%d:%s", len(key), key)
        subp, _ := benmap.Query(key)
        sub := *subp
        switch sub.Kind() {
            case benparser.Map: 
                fmt.Fprint(v, "\n")
                drawMapContent(v, sub, lvl + 1)
            case benparser.List: 
                fmt.Fprint(v, "\n")
                drawListContent(v, sub, lvl + 1)
            case benparser.Int: 
                drawIntContent(v, sub)
            case benparser.String: 
                drawStringContent(v, sub)
        }
        fmt.Fprint(v, "\n")
    }
    fmt.Fprintf(v, strings.Repeat("\t\t", lvl - 1) + "%s", "e")
}

func drawListContent(v *gocui.View, l benparser.Benval, lvl int) {
    benlist := l.(benparser.Benlist)

    fmt.Fprintf(v, strings.Repeat("\t\t", lvl - 1) + "%s\n", "l")
    for i := range benlist.Len() {
        switch sub := (*benlist.Get(i)); sub.Kind() {
            case benparser.Map: 
                drawMapContent(v, sub, lvl + 1)
            case benparser.List: 
                drawListContent(v, sub, lvl + 1)
            case benparser.Int: 
                fmt.Fprint(v, strings.Repeat("\t\t", lvl))
                drawIntContent(v, sub)
            case benparser.String: 
                fmt.Fprint(v, strings.Repeat("\t\t", lvl))
                drawStringContent(v, sub)
        }
        fmt.Fprint(v, "\n")
    }
    fmt.Fprintf(v, strings.Repeat("\t\t", lvl - 1) + "%s", "e")
}

func drawIntContent(v *gocui.View, i benparser.Benval) {
    benint := i.(benparser.Benint)

    fmt.Fprintf(v, "%s", benint.Raw())
}

func drawStringContent(v *gocui.View, s benparser.Benval) {
    benstring := s.(benparser.Benstring)

    //if the string is 100 bytes or less, draw it
    if len(benstring.Get()) <= 100 {
        fmt.Fprintf(v, "%d:%s", len(benstring.Raw()), benstring.Raw())
        //else, hide it... TODO: maybe later on add a show function to draw the bytes anyway
    } else {
        fmt.Fprintf(v, "%d:LARGE BYTESTRING HIDDEN...", len(benstring.Get()))
    }
}
