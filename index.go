package main


import "github.com/rmcs9/benparser"
import "sort"

var index map[int]*benparser.Benval = make(map[int]*benparser.Benval)
var target *benparser.Benval
var id int = 0

func indexInit(bval *benparser.Benval, log bool) {
    switch (*bval).Kind() {
    case benparser.Map: 

        if log {
            index[id] = bval
            id++
        }
        bmap := (*bval).(benparser.Benmap)
        keys := bmap.Keys(); sort.Strings(keys)

        for _, key := range keys {
            val, _ := bmap.Query(key)
            indexInit(val, true)
        }

    case benparser.List:

        if log {
            index[id] = bval
            id++
        }
        blist := (*bval).(benparser.Benlist)
        for i := range blist.Len() {
            indexInit(blist.Get(i), true)
        }

    case benparser.Int: 
        index[id] = bval
        id++
    case benparser.String:
        index[id] = bval
        id++
    }
}

