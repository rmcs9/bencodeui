package main


import (
    "fmt"
    "strings"
    "sort"

    "github.com/jroimartin/gocui"
    "github.com/rmcs9/benparser"
)


// ----------- DIRECTORY DRAW FUNCTIONS ---------------------

func drawDir(v *gocui.View, benv *benparser.Benval, prefix string) {
    switch (*benv).Kind() {
    case benparser.Map: 
        drawMapDir(v, benv, prefix)
    case benparser.List:
        drawListDir(v, benv, prefix)
    case benparser.Int:
        drawIntDir(v)
    case benparser.String:
        drawStringDir(v)
    }
}
func drawMapDir(v *gocui.View, m *benparser.Benval, prefix string) {
    benmap := (*m).(benparser.Benmap)

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

        fmt.Fprint(v, fstring)
        if sub, _ := benmap.Query(key); (*sub).Kind() == benparser.Map || (*sub).Kind() == benparser.List {
            drawDir(v, sub, newprefix)
        }
    }
}

func drawListDir(v *gocui.View, l *benparser.Benval, prefix string) {
    benlist := (*l).(benparser.Benlist) 

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

        fmt.Fprint(v, fstring)
        if sub := benlist.Get(i); (*sub).Kind() == benparser.Map || (*sub).Kind() == benparser.List {
            drawDir(v, sub, newprefix)
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
func drawContent(v *gocui.View, benv *benparser.Benval, lvl int) {
    switch (*benv).Kind() {
    case benparser.Map: 
        drawMapContent(v, benv, lvl);
    case benparser.List: 
        drawListContent(v, benv, lvl)
    case benparser.Int: 
        drawIntContent(v, benv)
    case benparser.String: 
        drawStringContent(v, benv)
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
        if (*sub).Kind() == benparser.Map || (*sub).Kind() == benparser.List {
            fmt.Fprint(v, "\n")
        }
        drawContent(v, sub, lvl + 1)
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
        sub := benlist.Get(i)
        if (*sub).Kind() == benparser.Int || (*sub).Kind() == benparser.String {
            fmt.Fprint(v, strings.Repeat("\t\t", lvl))
        }
        drawContent(v, sub, lvl + 1)
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
        //else, hide it... (its likely a hash of some kind and has illegal characters that will break the gui)
    } else {
        fmt.Fprintf(v, "%d:[***LARGE BYTESTRING HIDDEN***]", len(benstring.Get()))
    }

    if s == target {
        fmt.Fprint(v, "\x1b[0m")
    }
}

// -----------------------------------------INFO DRAW---------------------------------------------------
func drawInfo(v *gocui.View) {
    v.Clear()

    switch (*target).Kind() {
    case benparser.Map: 
        fmt.Fprintln(v, "TYPE: MAP")
        fmt.Fprintf(v, "BYTE SIZE: %d\n", len((*target).(benparser.Benmap).Raw()))
        fmt.Fprintf(v, "# OF KEYS: %d\n", len((*target).(benparser.Benmap).Keys()))
    case benparser.List:
        fmt.Fprintln(v, "TYPE: LIST")
        fmt.Fprintf(v, "BYTE SIZE: %d\n", len((*target).(benparser.Benlist).Raw()))
        fmt.Fprintf(v, "LIST SIZE: %d\n", (*target).(benparser.Benlist).Len())
    case benparser.Int:
        fmt.Fprintln(v, "TYPE: INT")
        fmt.Fprintf(v, "VALUE: %d", (*target).(benparser.Benint).Get())
    case benparser.String:
        fmt.Fprintln(v, "TYPE: BYTESTRING")
        fmt.Fprintf(v, "BYTE SIZE: %d\n", len((*target).(benparser.Benstring).Get()))
    }
}
