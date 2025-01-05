package bencode


type fileManager struct {
	pointer int
	file []byte
}


func (p *fileManager) Absorb(i int){
	p.pointer += i
}


func (p *fileManager) Peek(i int) (byte) {
	return p.file[p.pointer + i]
}


func (p *fileManager) IsDone() (bool) {
	return p.pointer >= len(p.file)
}


func (p *fileManager) Pop(i int) []byte {
	popped := p.file[p.pointer:p.pointer + i] 
	p.Absorb(i)
	return popped
}

func (p *fileManager) Find(b byte) int {
	for i := 0; !p.IsDone(); i++ {
		if p.Peek(i) == b { return i }
	}
	return -1
}

func (p *fileManager) ResetPointer(i int) {
	p.pointer = i
}

func (p *fileManager) GetPoint() int {
	return p.pointer
}

func BuildFM(f []byte) *fileManager {
	var file *fileManager = new(fileManager) 
	file.file = f
	file.pointer = 0
	return file
}
