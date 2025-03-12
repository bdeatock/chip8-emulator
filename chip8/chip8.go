package chip8

import (
	"fmt"
	"os"
	"time"
)

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

func (e *Emulator) LoadROM(romPath string) error {
	romData, err := os.ReadFile(romPath)
	if err != nil {
		return fmt.Errorf("failed to read ROM file: %w", err)
	}

	if len(romData) > len(e.Memory)-ProgramStartAddress {
		return fmt.Errorf("ROM too large: %dB (max is %dB)", len(romData), len(e.Memory)-ProgramStartAddress)
	}

	for i := range romData {
		e.Memory[ProgramStartAddress+i] = romData[i]
	}

	return nil
}

func (e *Emulator) Run(cyclesPerSecond int) {
	clock := time.NewTicker(time.Second / time.Duration(cyclesPerSecond))

	go func() {
		for range clock.C {
			e.RunCycle()
		}
	}()
}

func (e *Emulator) RunCycle() {
	// instruction is 16-bits long, combine 2 bytes from memory at program counter
	opcode := uint16(e.Memory[e.PC])<<8 | uint16(e.Memory[e.PC+1])
	e.PC += 2

	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	n := opcode & 0x000F
	nn := opcode & 0x00FF
	nnn := opcode & 0x0FFF

	switch opcode & 0xF000 {
	case 0x0000:
		switch opcode {
		case 0x00E0:
			// 00E0: Clear screen
			e.clearDisplay()
		}
	case 0x1000:
		// 1NNN: Jump
		e.PC = nnn
	case 0x6000:
		// 6XNN: Set
		e.Registers[x] = byte(nn)
	case 0x7000:
		// 7XNN: Add
		e.Registers[x] += byte(nn)
	case 0xA000:
		// ANNN: Set index
		e.I = nnn
	case 0xD000:
		// DXYN: Display
		e.drawSprite(int(e.Registers[x]), int(e.Registers[y]), int(n))
	}
}

func (e *Emulator) Reset() {
	e.clearDisplay()
	for i := range e.Memory {
		e.Memory[i] = 0
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

	fmt.Printf("PC: 0x%04x\n", e.PC)
	fmt.Printf("I : 0x%04x\n", e.I)
	fmt.Println("===Registers===")
	for i := range e.Registers {
		fmt.Printf("Reg %2d: 0x%02x\n", i, e.Registers[i])
	}
}
