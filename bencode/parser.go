package bencode


import (
	"os"
	"log"
	"strconv"
)

var fm *fileManager

//top level parse function... takes a filepath to a bencode file, returns a map from string -> interface
func ParseFile(path string) (*map[string]interface{}, []byte){
	//read the file 
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	//ensure that the begining of the file is formatted correctly
	if file[0] != 'd'{
		log.Fatal("unexpected token found at the start of bencode file...")
	}
	
	//create a new filemanager object
	fm = BuildFM(file)
	//begin parsing
	ret, bytes := parseTop()
	fm = nil
	return ret, bytes
}

func Parse(b *[]byte) *map[string]interface{} {
	if (*b)[0] != 'd' {
		log.Fatal("unexpected token found at the start of bencode file...")
	}
	fm = BuildFM(*b)
	ret := new(map[string]interface{})
	*ret = parseDict()
	fm = nil
	return ret
}

func parseTop() (*map[string]interface{}, []byte) {
	if fm.Peek(0) != 'd' {
		log.Fatal("invalid start to bencode object...")
	}

	fm.Absorb(1)

	top := new(map[string]interface{})
	*top = make(map[string]interface{})
	var bytes []byte
	for fm.Peek(0) != 'e' {
		key := parseByteString()

		if key == "info" {
			p := fm.GetPoint()
			bytes = get_bytes()
			fm.ResetPointer(p)
		}
		value := parseItem()

		(*top)[key] = value
	}
	fm.Absorb(1)

	return top, bytes
}


// parses a bencode dictionary
// dictionaries in bencode are formatted as follows:
// dkey:valuekey:value....e 
// where key is always a byteString and value is any of bencodes data types (int, string, list, dict)
func parseDict() map[string]interface{} {
	//ensure that the opening 'd' is present
	if fm.Peek(0) != 'd' {
		log.Fatal("invalid start to bencode object...")
	}
	fm.Absorb(1)
	
	info := make(map[string]interface{})
	//while we are not at the end of the dict..
	for fm.Peek(0) != 'e' {
		//take the key
		key := parseByteString()
		//take the value
		value := parseItem()
		//add the pair into the map
		info[key] = value
	}
	//absorb the e
	fm.Absorb(1)
	//return the map
	return info
}

//a midway point that determines what type the upcoming data is
func parseItem() (interface{}){
	if fm.Peek(0) == 'd' {
		return parseDict() 
	} else if fm.Peek(0) == 'i' {
		return parseInt()
	} else if fm.Peek(0) == 'l' {
		return parseList()
	} else {
		return parseBytes()
	}
}

// parses bencode ints... ints are formatted as follows:
// i124782385435e
func parseInt() int {
	//absorb the i
	fm.Absorb(1)

	//find the end of the encoded int
	i := fm.Find('e')
	//obtained the sliced bytes
	num := fm.Pop(i)
	//absorb the e
	fm.Absorb(1)
	//convert the bytes to an int and return
	numint, err := strconv.Atoi(string(num))
	if err != nil {
		log.Fatal(err)
	}
	return numint
}

//parses a bencode list, which is encoded as follows:
// lvaluevaluevaluevalue...e where value can be any valid bencode type
func parseList() []interface{} {
	//absorb the l
	fm.Absorb(1)

	//while not at the end of the list, parse new items
	//and add accumulate them in a slice
	var list []interface{}
	for fm.Peek(0) != 'e' {
		list = append(list, parseItem())
	}
	//absorb the e
	fm.Absorb(1)
	return list
}

//parses a bencode bytestring which is formated as follows
// <bytelength>:bytes 
func parseByteString() string {
	i := fm.Find(':')
	//obtain the length of the bytestring
	numOfBytes, err := strconv.Atoi(string(fm.Pop(i)))
	if err != nil {
		log.Fatal(err)
	}
	//absorb the colon
	fm.Absorb(1)
	//return the byte string
	return string(fm.Pop(numOfBytes))
}

func parseBytes() []byte {
	i := fm.Find(':')
	numOfBytes, err := strconv.Atoi(string(fm.Pop(i)))
	if err != nil {
		log.Fatal(err)
	}

	fm.Absorb(1)
	return fm.Pop(numOfBytes)
}
