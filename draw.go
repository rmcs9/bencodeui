package main


import (
    "fmt"
    "strings"
    "sort"

    "github.com/jroimartin/gocui"
    "github.com/rmcs9/benparser"
)


// ----------- DIRECTORY DRAW FUNCTIONS ---------------------
func drawMapDir(v *gocui.View, m benparser.Benval, prefix string) {
    benmap := m.(benparser.Benmap)

    keys := benmap.Keys(); sort.Strings(keys)
    for  i, key := range keys {
        fstring := ""
        newprefix := ""
        if i == len(keys) - 1 {
            fstring = fmt.Sprintf("%s└─%s\n", prefix, key)
            newprefix = prefix +  "  "
        } else {
            fstring = fmt.Sprintf("%s├─%s\n", prefix, key)
            newprefix = prefix + "│ "
        }

        if sub, _ := benmap.Query(key); (*sub).Kind() == benparser.Map {
            fmt.Fprint(v, fstring)
            drawMapDir(v, *sub, newprefix)
        } else if (*sub).Kind() == benparser.List {
            fmt.Fprint(v, fstring)
            drawListDir(v, *sub, newprefix)
        } else {
            fmt.Fprint(v, fstring)
        }

    }
}

func drawListDir(v *gocui.View, l benparser.Benval, prefix string) {
    benlist := l.(benparser.Benlist) 

    for i := range benlist.Len() {
        fstring := ""
        newprefix := ""
        if i == benlist.Len() - 1 {
            fstring = fmt.Sprintf("%s└─[%d]\n", prefix, i)
            newprefix = prefix + "  "
        } else {
            fstring = fmt.Sprintf("%s├─[%d]\n", prefix, i)
            newprefix = prefix + "│ "
        }

        if sub := benlist.Get(i); (*sub).Kind() == benparser.Map {
            fmt.Fprint(v, fstring)
            drawMapDir(v, *sub, newprefix)
        } else if (*sub).Kind() == benparser.List {
            fmt.Fprint(v, fstring)
            drawListDir(v, *sub, newprefix)
        } else {
            fmt.Fprint(v, fstring)
        }
    }
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


func drawContent(v *gocui.View) {
    v.Clear()
    switch benval.Kind() {
    case benparser.Map: drawMapContent(v, &benval, 1);
    case benparser.List: drawListContent(v, &benval, 1)
    case benparser.Int: drawIntContent(v, &benval)
    case benparser.String: drawStringContent(v, &benval)
    }
}

func drawMapContent(v *gocui.View, m *benparser.Benval, lvl int) {
    benmap := (*m).(benparser.Benmap)

    if m == target{
        fmt.Fprint(v, "\x1b[34;7m")
    }

    fmt.Fprintf(v, strings.Repeat("\t\t", lvl - 1) + "%s\n", "d")
    keys := benmap.Keys(); sort.Strings(keys)
    for _, key := range keys {
        fmt.Fprint(v, strings.Repeat("\t\t", lvl))
        fmt.Fprintf(v, "%d:%s", len(key), key)
        sub, _ := benmap.Query(key)
        switch (*sub).Kind() {
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
    if m == target {
        fmt.Fprint(v, "\x1b[0m")
    }
}

func drawListContent(v *gocui.View, l *benparser.Benval, lvl int) {
    benlist := (*l).(benparser.Benlist)

    if l == target {
        fmt.Fprint(v, "\x1b[34;7m")
    }

    fmt.Fprintf(v, strings.Repeat("\t\t", lvl - 1) + "%s\n", "l")
    for i := range benlist.Len() {
        switch sub := benlist.Get(i); (*sub).Kind() {
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
    if l == target {
        fmt.Fprint(v, "\x1b[0m")
    }
}

func drawIntContent(v *gocui.View, i *benparser.Benval) {
    benint := (*i).(benparser.Benint)
    if i == target {
        fmt.Fprintf(v, "\x1b[34;7m%s\x1b[0m", benint.Raw())
    } else {
        fmt.Fprintf(v, "%s", benint.Raw())
    }
}

func drawStringContent(v *gocui.View, s *benparser.Benval) {
    benstring := (*s).(benparser.Benstring)

    if s == target {
        fmt.Fprint(v, "\x1b[34;7m")
    }

    //if the string is 100 bytes or less, draw it
    if len(benstring.Get()) <= 100 {
        fmt.Fprintf(v, "%d:%s", len(benstring.Raw()), benstring.Raw())
        //else, hide it... TODO: maybe later on add a show function to draw the bytes anyway
    } else {
        fmt.Fprintf(v, "%d:[***LARGE BYTESTRING HIDDEN***]", len(benstring.Get()))
    }

    if s == target {
        fmt.Fprint(v, "\x1b[0m")
    }
}
