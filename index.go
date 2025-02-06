package main


import "github.com/rmcs9/benparser"
import "sort"

var index map[int]*benparser.Benval = make(map[int]*benparser.Benval)
var target *benparser.Benval
var id int = -1

func indexInit(bval *benparser.Benval) {
    if (*bval).Kind() == benparser.Int || (*bval).Kind() == benparser.String {
        index[0] = bval
        return
    }

    indexInitCollection(bval)
}

func indexInitCollection(bval *benparser.Benval) {
    switch (*bval).Kind() {
    case benparser.Map: 

        index[id] = bval
        id++
        bmap := (*bval).(benparser.Benmap)
        keys := bmap.Keys(); sort.Strings(keys)

        for _, key := range keys {
            val, _ := bmap.Query(key)
            indexInitCollection(val)
        }

    case benparser.List:

        index[id] = bval
        id++
        blist := (*bval).(benparser.Benlist)
        for i := range blist.Len() {
            indexInitCollection(blist.Get(i))
        }

    case benparser.Int: 
        index[id] = bval
        id++
    case benparser.String:
        index[id] = bval
        id++
    }
}

