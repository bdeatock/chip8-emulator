package chip8

import (
	"testing"
)

func TestOpcodes(t *testing.T) {
	t.Run("00E0 - Clear Screen", func(t *testing.T) {
		e := New()
		// 0x00E0 - clear screen
		e.Memory[0x200] = 0x00
		e.Memory[0x201] = 0xE0

		e.Display[0] = true
		e.Display[10] = true
		e.Display[100] = true

		e.RunCycle()

		for i, pixel := range e.Display {
			if pixel {
				t.Errorf("Pixel at position %d should be cleared", i)
			}
		}
	})

	t.Run("1NNN - Jump", func(t *testing.T) {
		e := New()
		// 0x1350 - jump to 0x0350
		e.Memory[0x200] = 0x13
		e.Memory[0x201] = 0x50

		e.RunCycle()

		if e.PC != 0x350 {
			t.Errorf("PC should be 0x350, got 0x%04X", e.PC)
		}
	})

	t.Run("6XNN - Set Register", func(t *testing.T) {
		e := New()
		// 0x6A42 - set register A to 0x42
		e.Memory[0x200] = 0x6A
		e.Memory[0x201] = 0x42

		e.RunCycle()

		if e.Registers[0xA] != 0x42 {
			t.Errorf("Register A should be 0x42, got 0x%02X", e.Registers[0xA])
		}
	})

	t.Run("7XNN - Add", func(t *testing.T) {
		e := New()
		// 0x7520 - add 0x20 to register 5
		e.Memory[0x200] = 0x75
		e.Memory[0x201] = 0x20
		// register 5 contains 0x10
		e.Registers[0x5] = 0x10

		e.RunCycle()

		if e.Registers[0x5] != 0x30 {
			t.Errorf("Register 5 should be 0x30, got 0x%02X", e.Registers[0x5])
		}
		// 0x7605 - add 0x05 to register 6
		e.Memory[0x202] = 0x76
		e.Memory[0x203] = 0x05
		// register 6 contains 0xFF
		e.Registers[0x6] = 0xFF
		e.RunCycle()

		// Check that the value wrapped around correctly
		if e.Registers[0x6] != 0x04 {
			t.Errorf("Register 6 should be 0x04 after overflow, got 0x%02X", e.Registers[0x6])
		}
	})

	t.Run("ANNN - Set Index", func(t *testing.T) {
		e := New()
		// 0xA123 - set index to 0x123
		e.Memory[0x200] = 0xA1
		e.Memory[0x201] = 0x23

		e.RunCycle()

		if e.I != 0x123 {
			t.Errorf("Index register should be 0x123, got 0x%04X", e.I)
		}
	})

	t.Run("DXYN - Draw Sprite", func(t *testing.T) {
		e := New()
		// 0xD123 - draw 3-high sprite at coords (reg1, reg2)
		e.Memory[0x200] = 0xD1
		e.Memory[0x201] = 0x23
		// sprite address is 0x300
		e.I = 0x300
		// sprite in 0x300 is 10000000
		e.Memory[0x300] = 0x80
		// reg1 is 0x5, reg2 is 0xA
		e.Registers[1] = 0x5
		e.Registers[2] = 0xA

		e.RunCycle()

		if !e.Display[10*DisplayWidth+5] {
			t.Errorf("Sprite should be drawn at (5,10)")
		}
	})

	t.Run("2NNN and 00EE - Call and Return", func(t *testing.T) {
		e := New()
		// 0x2400 - call subroutine at 0x400
		e.Memory[0x200] = 0x24
		e.Memory[0x201] = 0x00

		// 0x00EE - return from subroutine
		e.Memory[0x400] = 0x00
		e.Memory[0x401] = 0xEE

		e.RunCycle() // Execute call

		if e.PC != 0x400 || e.SP != 1 || e.Stack[0] != 0x202 {
			t.Errorf("Call failed: PC=0x%04X, SP=%d, Stack[0]=0x%04X", e.PC, e.SP, e.Stack[0])
		}

		e.RunCycle() // Execute return

		if e.PC != 0x202 || e.SP != 0 {
			t.Errorf("Return failed: PC=0x%04X, SP=%d", e.PC, e.SP)
		}
	})

	t.Run("3XNN - Skip if Equal", func(t *testing.T) {
		e := New()
		// 0x3A42 - skip next instruction if VA == 0x42
		e.Memory[0x200] = 0x3A
		e.Memory[0x201] = 0x42
		e.Registers[0xA] = 0x42

		e.RunCycle()

		if e.PC != 0x204 {
			t.Errorf("Skip if equal (equal case): PC should be 0x204, got 0x%04X", e.PC)
		}

		e.PC = 0x200 // try again with unequal case
		e.Registers[0xA] = 0x43

		e.RunCycle()

		if e.PC != 0x202 {
			t.Errorf("Skip if equal (not equal case): PC should be 0x202, got 0x%04X", e.PC)
		}
	})

	t.Run("4XNN - Skip if Not Equal", func(t *testing.T) {
		e := New()
		// 0x4A42 - skip next instruction if VA != 0x42
		e.Memory[0x200] = 0x4A
		e.Memory[0x201] = 0x42
		e.Registers[0xA] = 0x43

		e.RunCycle()

		if e.PC != 0x204 {
			t.Errorf("Skip if not equal (not equal case): PC should be 0x204, got 0x%04X", e.PC)
		}

		e.PC = 0x200 // try again with equal case
		e.Registers[0xA] = 0x42

		e.RunCycle()

		if e.PC != 0x202 {
			t.Errorf("Skip if not equal (equal case): PC should be 0x202, got 0x%04X", e.PC)
		}
	})

	t.Run("5XY0 - Skip if VX equals VY", func(t *testing.T) {
		e := New()
		// 0x5AB0 - skip next instruction if VA == VB
		e.Memory[0x200] = 0x5A
		e.Memory[0x201] = 0xB0
		e.Registers[0xA] = 0x42
		e.Registers[0xB] = 0x42

		e.RunCycle()

		if e.PC != 0x204 {
			t.Errorf("Skip if VX equals VY (equal case): PC should be 0x204, got 0x%04X", e.PC)
		}

		// try again with unequal case
		e.PC = 0x200
		e.Registers[0xB] = 0x43

		e.RunCycle()

		if e.PC != 0x202 {
			t.Errorf("Skip if VX equals VY (not equal case): PC should be 0x202, got 0x%04X", e.PC)
		}
	})

	t.Run("9XY0 - Skip if VX not equals VY", func(t *testing.T) {
		e := New()
		// 0x9AB0 - skip next instruction if VA != VB
		e.Memory[0x200] = 0x9A
		e.Memory[0x201] = 0xB0
		e.Registers[0xA] = 0x42
		e.Registers[0xB] = 0x43

		e.RunCycle()

		if e.PC != 0x204 {
			t.Errorf("Skip if VX not equals VY (not equal case): PC should be 0x204, got 0x%04X", e.PC)
		}

		// try again with equal case
		e.PC = 0x200
		e.Registers[0xB] = 0x42

		e.RunCycle()

		if e.PC != 0x202 {
			t.Errorf("Skip if VX not equals VY (equal case): PC should be 0x202, got 0x%04X", e.PC)
		}
	})

	t.Run("FX29 - Set I to font address", func(t *testing.T) {
		e := New()

		// 0xFA29 - Set I to location of sprite for digit in VA
		e.Memory[0x200] = 0xFA
		e.Memory[0x201] = 0x29

		testCases := []struct {
			digit    byte
			expected uint16
		}{
			{0x0, FontStartAddress + 0*FontSpriteHeight},
			{0x5, FontStartAddress + 5*FontSpriteHeight},
			{0xA, FontStartAddress + 10*FontSpriteHeight},
			{0xF, FontStartAddress + 15*FontSpriteHeight},
		}

		for _, tc := range testCases {
			e.PC = 0x200 // Reset PC for each test case
			e.Registers[0xA] = tc.digit

			e.RunCycle()

			if e.I != tc.expected {
				t.Errorf("For digit 0x%X, I register should be 0x%04X, got 0x%04X",
					tc.digit, tc.expected, e.I)
			}
		}
	})

	t.Run("8XY0 - Set VX to VY", func(t *testing.T) {
		e := New()
		// 0x8AB0 - set register A to value of register B
		e.Memory[0x200] = 0x8A
		e.Memory[0x201] = 0xB0
		e.Registers[0xA] = 0x00
		e.Registers[0xB] = 0x42

		e.RunCycle()

		if e.Registers[0xA] != 0x42 {
			t.Errorf("Register A should be 0x42, got 0x%02X", e.Registers[0xA])
		}
	})

	t.Run("8XY1 - Bitwise OR", func(t *testing.T) {
		e := New()
		// 0x8AB1 - set register A to A OR B
		e.Memory[0x200] = 0x8A
		e.Memory[0x201] = 0xB1
		e.Registers[0xA] = 0x0F
		e.Registers[0xB] = 0xF0

		e.RunCycle()

		if e.Registers[0xA] != 0xFF {
			t.Errorf("Register A should be 0xFF after OR, got 0x%02X", e.Registers[0xA])
		}
	})

	t.Run("8XY2 - Bitwise AND", func(t *testing.T) {
		e := New()
		// 0x8AB2 - set register A to A AND B
		e.Memory[0x200] = 0x8A
		e.Memory[0x201] = 0xB2
		e.Registers[0xA] = 0x0F
		e.Registers[0xB] = 0xFF

		e.RunCycle()

		if e.Registers[0xA] != 0x0F {
			t.Errorf("Register A should be 0x0F after AND, got 0x%02X", e.Registers[0xA])
		}
	})

	t.Run("8XY3 - Bitwise XOR", func(t *testing.T) {
		e := New()
		// 0x8AB3 - set register A to A XOR B
		e.Memory[0x200] = 0x8A
		e.Memory[0x201] = 0xB3
		e.Registers[0xA] = 0x0F
		e.Registers[0xB] = 0xFF

		e.RunCycle()

		if e.Registers[0xA] != 0xF0 {
			t.Errorf("Register A should be 0xF0 after XOR, got 0x%02X", e.Registers[0xA])
		}
	})

	t.Run("8XY4 - Add with Carry", func(t *testing.T) {
		e := New()
		// Test case without overflow
		// 0x8AB4 - add register B to register A
		e.Memory[0x200] = 0x8A
		e.Memory[0x201] = 0xB4
		e.Registers[0xA] = 0x10
		e.Registers[0xB] = 0x20

		e.RunCycle()

		if e.Registers[0xA] != 0x30 {
			t.Errorf("Register A should be 0x30 after addition, got 0x%02X", e.Registers[0xA])
		}
		if e.Registers[0xF] != 0 {
			t.Errorf("Carry flag should be 0 when no overflow, got %d", e.Registers[0xF])
		}

		// Test case with overflow
		e.PC = 0x200
		e.Registers[0xA] = 0xFF
		e.Registers[0xB] = 0x03

		e.RunCycle()

		if e.Registers[0xA] != 0x02 {
			t.Errorf("Register A should be 0x02 after overflow, got 0x%02X", e.Registers[0xA])
		}
		if e.Registers[0xF] != 1 {
			t.Errorf("Carry flag should be 1 when overflow occurs, got %d", e.Registers[0xF])
		}
	})

	t.Run("8XY5 - Subtract VY from VX with Borrow", func(t *testing.T) {
		e := New()
		// Test case without borrow
		// 0x8AB5 - subtract register B from register A
		e.Memory[0x200] = 0x8A
		e.Memory[0x201] = 0xB5
		e.Registers[0xA] = 0x30
		e.Registers[0xB] = 0x10

		e.RunCycle()

		if e.Registers[0xA] != 0x20 {
			t.Errorf("Register A should be 0x20 after subtraction, got 0x%02X", e.Registers[0xA])
		}
		if e.Registers[0xF] != 1 {
			t.Errorf("Borrow flag should be 1 when no borrow needed, got %d", e.Registers[0xF])
		}

		// Test case with borrow
		e.PC = 0x200
		e.Registers[0xA] = 0x10
		e.Registers[0xB] = 0x20

		e.RunCycle()

		if e.Registers[0xA] != 0xF0 {
			t.Errorf("Register A should be 0xF0 after borrow, got 0x%02X", e.Registers[0xA])
		}
		if e.Registers[0xF] != 0 {
			t.Errorf("Borrow flag should be 0 when borrow needed, got %d", e.Registers[0xF])
		}
	})

	t.Run("8XY6 - Shift Right", func(t *testing.T) {
		// Test modern behavior (shift VX right)
		e := New()
		// 0x8A06 - shift register A right by 1
		e.Memory[0x200] = 0x8A
		e.Memory[0x201] = 0x06
		e.Registers[0xA] = 0x03

		e.RunCycle()

		if e.Registers[0xA] != 0x01 {
			t.Errorf("Register A should be 0x01 after shift right, got 0x%02X", e.Registers[0xA])
		}
		if e.Registers[0xF] != 1 {
			t.Errorf("Flag register should be 1 (least significant bit was 1), got %d", e.Registers[0xF])
		}

		// Test with LSB = 0
		e.PC = 0x200
		e.Registers[0xA] = 0x04

		e.RunCycle()

		if e.Registers[0xA] != 0x02 {
			t.Errorf("Register A should be 0x02 after shift right, got 0x%02X", e.Registers[0xA])
		}
		if e.Registers[0xF] != 0 {
			t.Errorf("Flag register should be 0 (least significant bit was 0), got %d", e.Registers[0xF])
		}

		// Test legacy behavior (set VX to VY then shift)
		config := &EmulatorConfig{legacyShift: true}
		e = New(config)
		e.Memory[0x200] = 0x8A
		e.Memory[0x201] = 0xB6
		e.Registers[0xA] = 0x00
		e.Registers[0xB] = 0x03

		e.RunCycle()

		if e.Registers[0xA] != 0x01 {
			t.Errorf("Register A should be 0x01 after legacy shift right, got 0x%02X", e.Registers[0xA])
		}
		if e.Registers[0xF] != 1 {
			t.Errorf("Flag register should be 1 (least significant bit was 1), got %d", e.Registers[0xF])
		}
	})

	t.Run("8XY7 - Subtract VX from VY", func(t *testing.T) {
		e := New()
		// Test case without borrow
		// 0x8AB7 - set register A to register B minus register A
		e.Memory[0x200] = 0x8A
		e.Memory[0x201] = 0xB7
		e.Registers[0xA] = 0x10
		e.Registers[0xB] = 0x30

		e.RunCycle()

		if e.Registers[0xA] != 0x20 {
			t.Errorf("Register A should be 0x20 after subtraction, got 0x%02X", e.Registers[0xA])
		}
		if e.Registers[0xF] != 1 {
			t.Errorf("Borrow flag should be 1 when no borrow needed, got %d", e.Registers[0xF])
		}

		// Test case with borrow
		e.PC = 0x200
		e.Registers[0xA] = 0x30
		e.Registers[0xB] = 0x20

		e.RunCycle()

		if e.Registers[0xA] != 0xF0 {
			t.Errorf("Register A should be 0xF0 after borrow, got 0x%02X", e.Registers[0xA])
		}
		if e.Registers[0xF] != 0 {
			t.Errorf("Borrow flag should be 0 when borrow needed, got %d", e.Registers[0xF])
		}
	})

	t.Run("8XYE - Shift Left", func(t *testing.T) {
		// Test modern behavior (shift VX left)
		e := New()
		// 0x8A0E - shift register A left by 1
		e.Memory[0x200] = 0x8A
		e.Memory[0x201] = 0x0E
		e.Registers[0xA] = 0x81

		e.RunCycle()

		if e.Registers[0xA] != 0x02 {
			t.Errorf("Register A should be 0x02 after shift left, got 0x%02X", e.Registers[0xA])
		}
		if e.Registers[0xF] != 1 {
			t.Errorf("Flag register should be 1 (most significant bit was 1), got %d", e.Registers[0xF])
		}

		// Test with MSB = 0
		e.PC = 0x200
		e.Registers[0xA] = 0x01

		e.RunCycle()

		if e.Registers[0xA] != 0x02 {
			t.Errorf("Register A should be 0x02 after shift left, got 0x%02X", e.Registers[0xA])
		}
		if e.Registers[0xF] != 0 {
			t.Errorf("Flag register should be 0 (most significant bit was 0), got %d", e.Registers[0xF])
		}

		// Test legacy behavior (set VX to VY then shift)
		config := &EmulatorConfig{legacyShift: true}
		e = New(config)
		e.Memory[0x200] = 0x8A
		e.Memory[0x201] = 0xBE
		e.Registers[0xA] = 0x00
		e.Registers[0xB] = 0x81

		e.RunCycle()

		if e.Registers[0xA] != 0x02 {
			t.Errorf("Register A should be 0x02 after legacy shift left, got 0x%02X", e.Registers[0xA])
		}
		if e.Registers[0xF] != 1 {
			t.Errorf("Flag register should be 1 (MSB was 1), got %d", e.Registers[0xF])
		}
	})

	t.Run("FX55 - Store registers V0-VX (modern mode)", func(t *testing.T) {
		// Create emulator with modern store/load behavior
		config := &EmulatorConfig{legacyStoreLoad: false}
		e := New(config)

		// 0xF355 - store registers V0-V3 at address I
		e.Memory[0x200] = 0xF3
		e.Memory[0x201] = 0x55

		e.I = 0x300

		// Set registers V0-V3
		e.Registers[0] = 0x10
		e.Registers[1] = 0x20
		e.Registers[2] = 0x30
		e.Registers[3] = 0x40

		e.RunCycle()

		if e.Memory[0x300] != 0x10 {
			t.Errorf("Memory at 0x300 should be 0x10, got 0x%02X", e.Memory[0x300])
		}
		if e.Memory[0x301] != 0x20 {
			t.Errorf("Memory at 0x301 should be 0x20, got 0x%02X", e.Memory[0x301])
		}
		if e.Memory[0x302] != 0x30 {
			t.Errorf("Memory at 0x302 should be 0x30, got 0x%02X", e.Memory[0x302])
		}
		if e.Memory[0x303] != 0x40 {
			t.Errorf("Memory at 0x303 should be 0x40, got 0x%02X", e.Memory[0x303])
		}

		if e.I != 0x300 {
			t.Errorf("I register should remain 0x300 in modern mode, got 0x%04X", e.I)
		}
	})

	t.Run("FX55 - Store registers V0-VX (legacy mode)", func(t *testing.T) {
		// Create emulator with legacy store/load behavior
		config := &EmulatorConfig{legacyStoreLoad: true}
		e := New(config)

		// 0xF355 - store registers V0-V3 at address I
		e.Memory[0x200] = 0xF3
		e.Memory[0x201] = 0x55

		e.I = 0x300

		e.Registers[0] = 0x10
		e.Registers[1] = 0x20
		e.Registers[2] = 0x30
		e.Registers[3] = 0x40

		e.RunCycle()

		if e.Memory[0x300] != 0x10 {
			t.Errorf("Memory at 0x300 should be 0x10, got 0x%02X", e.Memory[0x300])
		}
		if e.Memory[0x301] != 0x20 {
			t.Errorf("Memory at 0x301 should be 0x20, got 0x%02X", e.Memory[0x301])
		}
		if e.Memory[0x302] != 0x30 {
			t.Errorf("Memory at 0x302 should be 0x30, got 0x%02X", e.Memory[0x302])
		}
		if e.Memory[0x303] != 0x40 {
			t.Errorf("Memory at 0x303 should be 0x40, got 0x%02X", e.Memory[0x303])
		}

		if e.I != 0x304 {
			t.Errorf("I register should be 0x304 in legacy mode, got 0x%04X", e.I)
		}
	})

	t.Run("FX65 - Load registers V0-VX (modern mode)", func(t *testing.T) {
		// Create emulator with modern store/load behavior
		config := &EmulatorConfig{legacyStoreLoad: false}
		e := New(config)

		// 0xF365 - load registers V0-V3 from address I
		e.Memory[0x200] = 0xF3
		e.Memory[0x201] = 0x65

		e.I = 0x300

		e.Memory[0x300] = 0x15
		e.Memory[0x301] = 0x25
		e.Memory[0x302] = 0x35
		e.Memory[0x303] = 0x45

		e.RunCycle()

		// Check register values
		if e.Registers[0] != 0x15 {
			t.Errorf("Register V0 should be 0x15, got 0x%02X", e.Registers[0])
		}
		if e.Registers[1] != 0x25 {
			t.Errorf("Register V1 should be 0x25, got 0x%02X", e.Registers[1])
		}
		if e.Registers[2] != 0x35 {
			t.Errorf("Register V2 should be 0x35, got 0x%02X", e.Registers[2])
		}
		if e.Registers[3] != 0x45 {
			t.Errorf("Register V3 should be 0x45, got 0x%02X", e.Registers[3])
		}

		if e.I != 0x300 {
			t.Errorf("I register should remain 0x300 in modern mode, got 0x%04X", e.I)
		}
	})

	t.Run("FX65 - Load registers V0-VX (legacy mode)", func(t *testing.T) {
		// Create emulator with legacy store/load behavior
		config := &EmulatorConfig{legacyStoreLoad: true}
		e := New(config)

		// 0xF365 - load registers V0-V3 from address I
		e.Memory[0x200] = 0xF3
		e.Memory[0x201] = 0x65

		e.I = 0x300

		e.Memory[0x300] = 0x15
		e.Memory[0x301] = 0x25
		e.Memory[0x302] = 0x35
		e.Memory[0x303] = 0x45

		e.RunCycle()

		if e.Registers[0] != 0x15 {
			t.Errorf("Register V0 should be 0x15, got 0x%02X", e.Registers[0])
		}
		if e.Registers[1] != 0x25 {
			t.Errorf("Register V1 should be 0x25, got 0x%02X", e.Registers[1])
		}
		if e.Registers[2] != 0x35 {
			t.Errorf("Register V2 should be 0x35, got 0x%02X", e.Registers[2])
		}
		if e.Registers[3] != 0x45 {
			t.Errorf("Register V3 should be 0x45, got 0x%02X", e.Registers[3])
		}

		if e.I != 0x304 {
			t.Errorf("I register should be 0x304 in legacy mode, got 0x%04X", e.I)
		}
	})
}
