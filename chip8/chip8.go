package chip8

import (
	"fmt"
	"math/rand"
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

// EmulatorConfig contains configuration options for the CHIP-8 emulator.
type EmulatorConfig struct {
	legacyShift     bool  // chip-48 and super-chip onwards is modern
	legacyJump      bool  // chip-48 and super-chip onwards is modern
	legacyStoreLoad bool  // legacy mode for older games from 1970s and 1980s
	randSeed        int64 // seed for rand
}

// Emulator represents a CHIP-8 emulator with all necessary components
// for executing CHIP-8 programs.
type Emulator struct {
	// 4 kilobytes of RAM
	// Note: 0x000-0x1FF reserved for interpreter in early versions, so
	// start accessible RAM from 0x200 to support older ROMs
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

	timerDelta time.Duration // timer to track time since last timer update

	// Registers
	// General-purpose variable registers
	Registers [RegisterCount]byte

	// Keypad - 16 keys (0-F)
	Keypad [16]bool

	// Config
	config *EmulatorConfig

	// Random number generator
	rng *rand.Rand
}

// NewEmulator creates and initializes a new CHIP-8 emulator with the provided configuration.
func New(config ...*EmulatorConfig) *Emulator {
	e := &Emulator{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	if len(config) > 0 && config[0] != nil {
		e.config = config[0]
		e.rng.Seed(e.config.randSeed)
	} else {
		// Default config
		e.config = &EmulatorConfig{
			legacyShift:     false,
			legacyJump:      true,
			legacyStoreLoad: false,
		}
	}

	e.Reset()
	return e
}

// LoadROM loads a CHIP-8 ROM from the specified file path into the emulator's memory
// starting at address ProgramStartAddress (usually 0x200). Returns an error if the file
// cannot be read or if the ROM is too large to fit in memory.
func (e *Emulator) LoadROMFromPath(romPath string) error {
	romData, err := os.ReadFile(romPath)
	if err != nil {
		return fmt.Errorf("failed to read ROM file: %w", err)
	}

	if len(romData) > len(e.Memory)-ProgramStartAddress {
		return fmt.Errorf("ROM too large: %dB (max is %dB)", len(romData), len(e.Memory)-ProgramStartAddress)
	}

	// Load ROM into memory starting at 0x200
	copy(e.Memory[ProgramStartAddress:], romData)
	return nil
}

func (e *Emulator) LoadROMFromData(romData []byte) error {
	if len(romData) > len(e.Memory)-ProgramStartAddress {
		return fmt.Errorf("ROM too large: %dB (max is %dB)", len(romData), len(e.Memory)-ProgramStartAddress)
	}

	copy(e.Memory[ProgramStartAddress:], romData)
	return nil
}

// Run starts the emulator's main execution loop in a separate goroutine.
// It continuously fetches, decodes, and executes instructions at the rate
// specified by the emulator's clock speed.
func (e *Emulator) Run(cyclesPerSecond int) <-chan error {
	clock := time.NewTicker(time.Second / time.Duration(cyclesPerSecond))

	errCh := make(chan error, 1)

	go func() {
		for range clock.C {
			if err := e.Step(time.Second / time.Duration(cyclesPerSecond)); err != nil {
				errCh <- err
				return
			}
		}
	}()

	return errCh
}

func validateReadAddress(address uint16, offset uint16) error {
	if address+offset > 0xFFF {
		return fmt.Errorf("memory access out of bounds: (0x%04X)", address)
	}
	return nil
}

func validateWriteAddress(address uint16, offset uint16) error {
	if address < ProgramStartAddress {
		return fmt.Errorf("memory access in reserved space: (0x%04X)", address)
	}
	return validateReadAddress(address, offset)
}

// Step executes a single instruction cycle of the emulator.
// This includes fetching the next opcode, decoding it, and executing
// the corresponding operation.
func (e *Emulator) Step(deltaTime time.Duration) error {
	if err := validateReadAddress(e.PC, 1); err != nil {
		return fmt.Errorf("failed to read opcode: %w", err)
	}
	opcode := uint16(e.Memory[e.PC])<<8 | uint16(e.Memory[e.PC+1])
	e.PC += 2

	e.UpdateTimers(deltaTime)

	err := e.executeOpcode(opcode)
	if err != nil {
		return fmt.Errorf("error executing opcode: %w", err)
	}
	return nil
}

func (e *Emulator) executeOpcode(opcode uint16) error {
	x := byte((opcode & 0x0F00) >> 8)
	y := byte((opcode & 0x00F0) >> 4)
	n := byte(opcode & 0x000F)
	nn := byte(opcode & 0x00FF)
	nnn := opcode & 0x0FFF

	switch opcode & 0xF000 {
	case 0x0000:
		switch opcode {
		case 0x00E0:
			// 00E0: Clear screen
			e.clearDisplay()
		case 0x00EE:
			// Return from subroutine
			if e.SP == 0 {
				// stack is empty
				return fmt.Errorf("stack underflow - attempted to return from subroutine with empty stack")
			}
			// Decrement stack pointer first
			e.SP--
			// Set PC to the address from the stack
			e.PC = e.Stack[e.SP]
		}
	case 0x1000:
		// 1NNN: Jump
		e.PC = nnn
	case 0x2000:
		// 2NNN: call subroutine at NNN
		// check stack has room
		if int(e.SP) >= len(e.Stack) {
			return fmt.Errorf("stack overflow - maximum call depth exceeded")
		}
		// push current pc to stack
		e.Stack[e.SP] = e.PC
		e.SP++
		// set pc to new address
		e.PC = nnn
	case 0x3000:
		// 3XNN: Skip next instruction if VX equals NN
		if e.Registers[x] == byte(nn) {
			e.PC += 2
		}
	case 0x4000:
		// 4XNN: Skip next instruction if VX not equal to NN
		if e.Registers[x] != byte(nn) {
			e.PC += 2
		}
	case 0x5000:
		// 5XY0: Skip next instruction if VX equal to VY
		if n == 0 && e.Registers[x] == e.Registers[y] {
			e.PC += 2
		} else if n != 0 {
			return fmt.Errorf("unknown opcode: 0x%X", opcode)
		}
	case 0x6000:
		// 6XNN: Set
		e.Registers[x] = byte(nn)
	case 0x7000:
		// 7XNN: Add
		e.Registers[x] += byte(nn)
	case 0x8000:
		switch n {
		case 0x0:
			// 8XY0: Set VX to value of VY
			e.Registers[x] = e.Registers[y]
		case 0x1:
			// 8XY1: Set VX to bitwise VX OR VY
			e.Registers[x] |= e.Registers[y]
		case 0x2:
			// 8XY2: Set VX to bitwise VX AND VY
			e.Registers[x] &= e.Registers[y]
		case 0x3:
			// 8XY3: Set VX to bitwise VX XOR VY
			e.Registers[x] ^= e.Registers[y]
		case 0x4:
			// 8XY4: Add VY to VX with carry
			sum := uint16(e.Registers[x]) + uint16(e.Registers[y])
			if sum > 0xFF {
				e.Registers[0xF] = 1 // Set carry flag
			} else {
				e.Registers[0xF] = 0
			}
			e.Registers[x] = byte(sum)
		case 0x5:
			// 8XY5: Subtract VY from VX with borrow
			if e.Registers[x] >= e.Registers[y] {
				e.Registers[0xF] = 1 // No borrow needed
			} else {
				e.Registers[0xF] = 0 // Borrow needed
			}
			e.Registers[x] -= e.Registers[y]
		case 0x6:
			// 8XY6: legacy - Set VX to VY shifted 1 bit to right, VF is set to bit shifted out
			//       modern - Shift VX 1 bit to right, VF is set to bit shifted out
			if e.config.legacyShift {
				e.Registers[x] = e.Registers[y]
			}
			// Check rightmost bit before shift
			e.Registers[0xF] = e.Registers[x] & 0x01
			e.Registers[x] = e.Registers[x] >> 1
		case 0x7:
			// 8XY7: Set VX to VY - VX with borrow
			if e.Registers[y] >= e.Registers[x] {
				e.Registers[0xF] = 1 // No borrow needed
			} else {
				e.Registers[0xF] = 0 // Borrow needed
			}
			e.Registers[x] = e.Registers[y] - e.Registers[x]
		case 0xE:
			// 8XYE: legacy - Set VX to VY shifted 1 bit to left, VF is set to bit shifted out
			//       modern - Shift VX 1 bit to left, VF is set to bit shifted out
			if e.config.legacyShift {
				e.Registers[x] = e.Registers[y]
			}
			// Check leftmost bit before shift
			e.Registers[0xF] = (e.Registers[x] & 0x80) >> 7
			e.Registers[x] = e.Registers[x] << 1
		default:
			return fmt.Errorf("unknown opcode: 0x%X", opcode)
		}
	case 0x9000:
		// 9XY0: Skip next instruction if VX not equal to VY
		if n == 0 && e.Registers[x] != e.Registers[y] {
			e.PC += 2
		} else if n != 0 {
			return fmt.Errorf("unknown opcode: 0x%X", opcode)
		}
	case 0xA000:
		// ANNN: Set index
		e.I = nnn
	case 0xB000:
		// BNNN: Jump with offset
		if e.config.legacyJump {
			// jump to address NNN + value in V0
			e.PC = (nnn + uint16(e.Registers[0])) & 0x0FFF
		} else {
			// jump to address NNN + value in X
			e.PC = (nnn + uint16(e.Registers[x])) & 0x0FFF
		}
	case 0xC000:
		// CXNN: sets VX to random number AND NN
		e.Registers[x] = byte(e.rng.Int()) & nn
	case 0xD000:
		// DXYN: Display
		e.drawSprite(int(e.Registers[x]), int(e.Registers[y]), int(n))
	case 0xE000:
		switch nn {
		case 0x9E:
			// EX9E Skip if key X is pressed
			if e.Keypad[e.Registers[x]] {
				e.PC += 2
			}
		case 0xA1:
			// EXA1 Skip if key X is not pressed
			if !e.Keypad[e.Registers[x]] {
				e.PC += 2
			}
		}
	case 0xF000:
		switch nn {
		case 0x07:
			// FX07 Set VX to current value of delay timer
			e.Registers[x] = e.DelayTimer
		case 0x0A:
			// 0xFX0A Wait for a key press and set VX to key
			e.PC -= 2 // If no key pressed we need to repeat instruction
			for key, pressed := range e.Keypad {
				if pressed {
					e.Registers[x] = byte(key)
					e.PC += 2
					break
				}
			}
		case 0x15:
			// 0xFX15 Set delay timer to VX
			e.DelayTimer = e.Registers[x]
		case 0x18:
			// 0xFX18 Set sound timer to VX
			e.SoundTimer = e.Registers[x]
		case 0x1E:
			// 0xFX1E Add VX to I
			// Note: this didn't affect VF on overflow (I > 0x0FFF) in original chip-8, but did in some later versions
			// At least one known game requires it, but unlikely that anything relies on it not happening so
			// implementing it this way
			if (e.I+uint16(e.Registers[x]))&0xF000 != 0 {
				e.Registers[0xF] = 1
			} else {
				e.Registers[0xF] = 0
			}
			e.I += uint16(e.Registers[x])
		case 0x29:
			// 0xFX29 Set I to address of font for hex char in VX
			e.I = FontStartAddress + uint16(e.Registers[x]&0x0F)*FontSpriteHeight
		case 0x33:
			// 0xFX33 Take number in VX, convert to three decimal digits, and store at address in I, I+1, I+2
			if err := validateWriteAddress(e.I, 2); err != nil {
				return fmt.Errorf("decimalise register: %w", err)
			}
			e.Memory[e.I] = e.Registers[x] / 100
			e.Memory[e.I+1] = (e.Registers[x] % 100) / 10
			e.Memory[e.I+2] = e.Registers[x] % 10
		case 0x55:
			// 0xFX55 Store V0-VX at address I
			if err := validateWriteAddress(e.I, uint16(x)); err != nil {
				return fmt.Errorf("failed to store registers: %w", err)
			}
			for i := range uint16(x + 1) {
				e.Memory[e.I+i] = e.Registers[i]
			}
			if e.config.legacyStoreLoad {
				e.I = e.I + uint16(x) + 1
			}
		case 0x65:
			// 0xFX65 Load memory from address I into V0-VX
			if err := validateReadAddress(e.I, uint16(x)); err != nil {
				return fmt.Errorf("failed to load into registers: %w", err)
			}
			for i := range uint16(x + 1) {
				e.Registers[i] = e.Memory[e.I+i]
			}
			if e.config.legacyStoreLoad {
				e.I = e.I + uint16(x) + 1
			}
		default:
			return fmt.Errorf("unknown opcode: 0x%X", opcode)
		}

	default:
		return fmt.Errorf("unknown opcode: 0x%X", opcode)
	}
	return nil
}

// Reset resets the emulator to its initial state, clearing memory, registers,
// and resetting the program counter to the starting address (0x200).
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
	for i := range e.Keypad {
		e.Keypad[i] = false
	}

	// Reset program counter to start of program memory
	e.PC = ProgramStartAddress

	e.I = 0
	e.SP = 0
	e.DelayTimer = 0
	e.SoundTimer = 0
	e.timerDelta = 0

	e.loadFontData()
}

// Prints the emulator display and variables to console
func (e *Emulator) Print() {
	e.printDisplay()

	fmt.Printf("PC: 0x%04x\n", e.PC)
	fmt.Printf("I : 0x%04x\n", e.I)
	fmt.Println("===Registers===")
	for i := range e.Registers {
		fmt.Printf("Reg %2d: 0x%02x\n", i, e.Registers[i])
	}
}

// UpdateTimers decrements the delay and sound timers if they are greater than zero.
// This should be called at a rate of 60Hz according to the CHIP-8 specification.
func (e *Emulator) UpdateTimers(deltaTime time.Duration) {
	e.timerDelta += deltaTime
	if e.timerDelta > time.Second/60.0 {
		e.timerDelta = 0.0
		if e.DelayTimer > 0 {
			e.DelayTimer--
		}
		if e.SoundTimer > 0 {
			e.SoundTimer--
		}
	}
}
