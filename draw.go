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

var id int = -1

func drawContent(v *gocui.View, init bool) {
    v.Clear()
    switch benval.Kind() {
    case benparser.Map: drawMapContent(v, &benval, 1, init);
    case benparser.List: drawListContent(v, &benval, 1, init)
    case benparser.Int: drawIntContent(v, &benval, init)
    case benparser.String: drawStringContent(v, &benval, init)
    }
}

func drawMapContent(v *gocui.View, m *benparser.Benval, lvl int, init bool) {
    if init { 
        index[id] = m  
        id++ 
    }
    benmap := (*m).(benparser.Benmap)

    if m == target && !init {
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
    if m == target && !init {
        fmt.Fprint(v, "\x1b[0m")
    }
    // if init { index[id] = m; id++ }
}

func drawListContent(v *gocui.View, l *benparser.Benval, lvl int, init bool) {
    if init { 
        index[id] = l
        id++ 
    }
    benlist := (*l).(benparser.Benlist)

    if l == target && !init {
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
    if l == target && !init {
        fmt.Fprint(v, "\x1b[0m")
    }
    // if init { index[id] = l; id++ }
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


// type Node struct {
//     Str string 
//     Val *benparser.Benval
//     Children []*Node
// }
//
// func walkContent(r benparser.Benval) *Node {
//     root := new(Node)
//     
//     switch r.Kind() {
//         case benparser.Map: walkMap(r.(benparser.Benmap), root) 
//         case benparser.List: walkList(r.(benparser.Benlist), root) 
//         case benparser.Int: walkInt(r.(benparser.Benint), root)
//         case benparser.String: walkString(r.(benparser.Benstring), root)
//     }
//     
//     return root
// }
//
// func walkMap(m benparser.Benmap, n *Node) {
//     str := "d"
//     keys := m.Keys()
//
//     for _, key := range keys {
//         sub, _ := m.Query(key)
//
//
//     }
// }
//
// func walkList(l benparser.Benlist, n *Node) {
//
// }
//
// func walkInt(i benparser.Benint, n *Node) {
//
// }
//
// func walkString(s benparser.Benstring, n *Node) {
//
// }
