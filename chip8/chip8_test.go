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
}
