package chip8

const (
	DisplayWidth  = 64
	DisplayHeight = 32
)

type Emulator struct {
	// 4 kilobytes of RAM
	// Note: 0x000-0x1FF reserved for interpreter in early versions, so start accessible RAM from 0x200 to support older ROMs
	Memory [4096]byte

	// Display
	// 64x32 - pixels can be on/off
	Display [DisplayWidth * DisplayHeight]bool
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
	for i := range e.Display {
		e.Display[i] = false
	}

	e.loadFontData()
}

func (e *Emulator) Print() {
	e.printDisplay()
}
