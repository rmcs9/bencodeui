package main


import (
    "fmt"
    "strings"
    "sort"

    "github.com/jroimartin/gocui"
    "github.com/rmcs9/benparser"
)

var index map[int]*benparser.Benval = make(map[int]*benparser.Benval)
var target *benparser.Benval

// ----------- DIRECTORY DRAW FUNCTIONS ---------------------
func drawMapDir(v *gocui.View, m benparser.Benval, lvl int) {
    benmap := m.(benparser.Benmap)

    fmt.Fprintf(v, strings.Repeat("\t\t", lvl - 1) + "%s\n", "d") 

    keys := benmap.Keys(); sort.Strings(keys)
    for  _, key := range keys {
        fstring := strings.Repeat("\t\t", lvl) + "%s"


        if sub, _ := benmap.Query(key); (*sub).Kind() == benparser.Map {
            fmt.Fprintf(v, fstring, key)
            drawMapDir(v, *sub, lvl + 1) 
        } else if (*sub).Kind() == benparser.List {
            fmt.Fprintf(v, fstring, key)
            drawListDir(v, *sub, lvl + 1)
        } else {
            fstring += "\n" 
            fmt.Fprintf(v, fstring, key)
        }

    }

    fmt.Fprintf(v, strings.Repeat("\t\t", lvl - 1) + "%s\n", "e") 
}

func drawListDir(v *gocui.View, l benparser.Benval, lvl int) {
    benlist := l.(benparser.Benlist) 

    fmt.Fprintf(v, strings.Repeat("\t\t", lvl - 1) + "%s\n", "l")

    for i := range benlist.Len() {
        fstring := strings.Repeat("\t\t", lvl) + "[%d]"

        if sub := benlist.Get(i); (*sub).Kind() == benparser.Map {
            fmt.Fprintf(v, fstring, i)
            drawMapDir(v, *sub, lvl + 1) 
        } else if (*sub).Kind() == benparser.List {
            fmt.Fprintf(v, fstring, i)
            drawListDir(v, *sub, lvl + 1)
        } else {
            fstring += "\n"
            fmt.Fprintf(v, fstring, i)
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

var id int = 0

func drawContent(v *gocui.View, init bool) {
    switch benval.Kind() {
        case benparser.Map: drawMapContent(v, &benval, 1, init);
        case benparser.List: drawListContent(v, &benval, 1, init)
        case benparser.Int: drawIntContent(v, &benval, init)
        case benparser.String: drawStringContent(v, &benval, init)
    }
}

func drawMapContent(v *gocui.View, m *benparser.Benval, lvl int, init bool) {
    if init { index[id] = m; id++ }
    benmap := (*m).(benparser.Benmap)

    if m == target {
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
                drawMapContent(v, sub, lvl + 1, init)
            case benparser.List: 
                fmt.Fprint(v, "\n")
                drawListContent(v, sub, lvl + 1, init)
            case benparser.Int: 
                drawIntContent(v, sub, init)
            case benparser.String: 
                drawStringContent(v, sub, init)
        }
        fmt.Fprint(v, "\n")
    }
    fmt.Fprintf(v, strings.Repeat("\t\t", lvl - 1) + "%s", "e")
    if m == target {
        fmt.Fprint(v, "\x1b[0m")
    }
    if init { index[id] = m; id++ }
}

func drawListContent(v *gocui.View, l *benparser.Benval, lvl int, init bool) {
    if init { index[id] = l; id++ }
    benlist := (*l).(benparser.Benlist)

    if l == target {
        fmt.Fprint(v, "\x1b[34;7m")
    }

    fmt.Fprintf(v, strings.Repeat("\t\t", lvl - 1) + "%s\n", "l")
    for i := range benlist.Len() {
        switch sub := benlist.Get(i); (*sub).Kind() {
            case benparser.Map: 
                drawMapContent(v, sub, lvl + 1, init)
            case benparser.List: 
                drawListContent(v, sub, lvl + 1, init)
            case benparser.Int: 
                fmt.Fprint(v, strings.Repeat("\t\t", lvl))
                drawIntContent(v, sub, init)
            case benparser.String: 
                fmt.Fprint(v, strings.Repeat("\t\t", lvl))
                drawStringContent(v, sub, init)
        }
        fmt.Fprint(v, "\n")
    }
    fmt.Fprintf(v, strings.Repeat("\t\t", lvl - 1) + "%s", "e")
    if l == target {
        fmt.Fprint(v, "\x1b[0m")
    }
    if init { index[id] = l; id++ }
}

func drawIntContent(v *gocui.View, i *benparser.Benval, init bool) {
    if init { index[id] = i; id++ }
    benint := (*i).(benparser.Benint)
    if i == target {
        fmt.Fprintf(v, "\x1b[34;7m%s\x1b[0m", benint.Raw())
    } else {
        fmt.Fprintf(v, "%s", benint.Raw())
    }
}

func drawStringContent(v *gocui.View, s *benparser.Benval, init bool) {
    if init { index[id] = s; id++ }
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

// func drawContent(v benparser.Benval, init bool) string {
//     
// }
//
// func mapContent(m benparser.Benmap, lvl int, init, tar bool) string {
//     str := strings.Repeat("\t\t", lvl - 1) + "d\n"
//     keys := m.Keys(); sort.Strings(keys)
//     for _, key := range keys {
//         str += strings.Repeat("\t\t", lvl) + strconv.Itoa(len(key)) + ":" + key
//         subp, _ := m.Query(key) 
//         sub := *subp
//         switch sub.Kind() {
//             case benparser.Map:
//                 str += "\n" 
//                 str += mapContent(sub.(benparser.Benmap), lvl + 1, init, tar)
//             case benparser.List:
//                 str += "\n"
//                 str += listContent(sub.(benparser.Benlist), lvl + 1, init, tar)
//             case benparser.Int:
//                 str += intContent(sub.(benparser.Benint), init, tar)
//             case benparser.String: 
//                 str += stringContent(sub.(benparser.Benstring), init, tar)
//         }
//         str += "\n"
//     }
// }
//
// func listContent(l benparser.Benlist, lvl int, init, tar bool) string {
//
// }
//
// func intContent(i benparser.Benint, init, tar bool) string {
//
// }
//
// func stringContent(s benparser.Benstring, init, tar bool) string {
//
// }
