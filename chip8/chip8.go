package chip8

const (
	DisplayWidth  = 64
	DisplayHeight = 32

	// Constants for memory addresses and limits
	ProgramStartAddress = 0x200 // Starting address for most CHIP-8 programs
	StackSize           = 16    // Maximum stack depth
	RegisterCount       = 16    // Number of registers
)

type Emulator struct {
	// 4 kilobytes of RAM
	// Note: 0x000-0x1FF reserved for interpreter in early versions, so start accessible RAM from 0x200 to support older ROMs
	Memory [4096]byte

	// Display
	// 64x32 - pixels can be on/off
	Display [DisplayWidth * DisplayHeight]bool

	// Program Counter
	// Points to current instruction in memory
	PC uint16

	// Index Register
	// Points to locations in memory
	I uint16

	// Stack of 16-bit addresses
	// To call functions and return from them
	Stack [StackSize]uint16
	// Stack pointer
	SP uint8

	// Delay timer
	// Decrements 60 times per second if not 0
	DelayTimer uint8

	// Sound timer
	// Functions like delay timer, but beeps while not 0
	SoundTimer uint8

	// Registers
	// General-purpose variable registers
	Registers [RegisterCount]byte
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
	for i := range e.Registers {
		e.Registers[i] = 0
	}
	for i := range e.Stack {
		e.Stack[i] = 0
	}

	// Reset program counter to start of program memory
	e.PC = ProgramStartAddress

	e.I = 0
	e.SP = 0
	e.DelayTimer = 0
	e.SoundTimer = 0

	e.loadFontData()
}

func (e *Emulator) Print() {
	e.printDisplay()
}
