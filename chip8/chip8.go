package chip8

import "fmt"

type Emulator struct {
	// 4 kilobytes of RAM
	// Note: 000-1FF were reserved for interpreter in early versions, so start accessible RAM from #200(512) to support older ROMs
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

	e.loadFontData()
}

func (e *Emulator) Print() {
	fmt.Println(e.Memory)
}
