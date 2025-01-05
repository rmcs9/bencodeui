package bencode


import (
	"log"
	"strconv"
)


func get_bytes() []byte {
	if fm.Peek(0) == 'd' {
		return get_dict()	
	} else if fm.Peek(0) == 'l' {
		return get_list()
	} else if fm.Peek(0) == 'i' {
		return get_int()
	} else {
		return get_other()
	}
}

func get_dict() []byte {
	bytes := make([]byte, 0)
	bytes = append(bytes, fm.Pop(1)...)

	for fm.Peek(0) != 'e' {
		bytes = append(bytes, get_other()...)

		bytes = append(bytes, get_bytes()...)
	}

	bytes = append(bytes, fm.Pop(1)...)

	return bytes
}

func get_list() []byte {
	bytes := make([]byte, 0)
	bytes = append(bytes, fm.Pop(1)...)

	for fm.Peek(0) != 'e' {
		bytes = append(bytes, get_bytes()...)
	}
	bytes = append(bytes, fm.Pop(1)...)
	return bytes
}

func get_int() []byte {
	return fm.Pop(fm.Find('e') + 1)
}

func get_other() []byte {
	bytes := make([]byte, 0)

	i := fm.Find(':')
	b := fm.Pop(i)

	bytes = append(bytes, b...)
	bytes = append(bytes, fm.Pop(1)...)
	
	numOfBytes, err := strconv.Atoi(string(b))
	
	if err != nil {
		log.Fatal(err)
	}

	b = fm.Pop(numOfBytes)
	bytes = append(bytes, b...)
	return bytes
}
