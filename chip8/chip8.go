package chip8

import "fmt"

type Emulator struct {
	Memory [4096]byte
}

func New() *Emulator {
	e := &Emulator{}
	e.Reset()
	return e
}

func (e *Emulator) Reset() {
	for i := range e.Memory {
		e.Memory[i] = 0
	}
}

func (e *Emulator) Print() {
	fmt.Println(e.Memory)
}