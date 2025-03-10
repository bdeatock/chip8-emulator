package chip8

import (
	"testing"
)

func TestClearDisplay(t *testing.T) {
	e := New()

	e.Display[0] = true
	e.Display[10] = true

	e.clearDisplay()

	// Check all pixels are cleared
	for i, pixel := range e.Display {
		if pixel {
			t.Errorf("Pixel at position %d is still set after clearDisplay()", i)
		}
	}
}

func TestFlipPixel(t *testing.T) {
	e := New()

	initialState := e.Display[10*DisplayWidth+5]
	if initialState != false {
		t.Errorf("Expected initial pixel state to be false")
	}

	result := e.flipPixel(5, 10)
	if !result {
		t.Errorf("flipPixel should return true after flipping from false")
	}

	if !e.Display[10*DisplayWidth+5] {
		t.Errorf("Pixel should be true after flipping from false")
	}

	result = e.flipPixel(5, 10)
	if result {
		t.Errorf("flipPixel should return false after flipping from true")
	}

	if e.Display[10*DisplayWidth+5] {
		t.Errorf("Pixel should be false after flipping from true")
	}
}

func TestDrawSprite(t *testing.T) {
	e := New()

	t.Run("Basic sprite drawing", func(t *testing.T) {
		e.clearDisplay()

		// Put a simple line sprite into memory
		e.I = 0x300
		e.Memory[0x300] = 0x80 // 10000000
		e.Memory[0x301] = 0x80 // 10000000
		e.Memory[0x302] = 0x80 // 10000000

		e.drawSprite(5, 10, 3)

		if !e.Display[10*DisplayWidth+5] {
			t.Errorf("Pixel at (5,10) should be set")
		}
		if !e.Display[11*DisplayWidth+5] {
			t.Errorf("Pixel at (5,11) should be set")
		}
		if !e.Display[12*DisplayWidth+5] {
			t.Errorf("Pixel at (5,12) should be set")
		}

		// Check that all other pixels are still clear
		for y := range DisplayHeight {
			for x := range DisplayWidth {
				if (y == 10 && x == 5) || (y == 11 && x == 5) || (y == 12 && x == 5) {
					continue
				}
				if e.Display[y*DisplayWidth+x] {
					t.Errorf("Pixel at (%d,%d) should be clear", x, y)
				}
			}
		}

		if e.Registers[0xF] != 0 {
			t.Errorf("Collision flag should not be set")
		}

		e.drawSprite(5, 10, 3)

		if e.Display[10*DisplayWidth+5] {
			t.Errorf("Pixel at (5,10) should be unset after second draw")
		}

		if e.Registers[0xF] != 1 {
			t.Errorf("Collision flag should be set after collision")
		}
	})

	t.Run("Coordinate wrapping", func(t *testing.T) {
		e.clearDisplay()
		e.I = 0x300
		e.Memory[0x300] = 0x80 // 10000000

		// Draw at X=68 (which should wrap to X=4)
		e.drawSprite(68, 5, 1)

		if !e.Display[5*DisplayWidth+4] {
			t.Errorf("Pixel at (4,5) should be set (wrapped from X=68)")
		}

		// Draw at Y=34 (which should wrap to Y=2)
		e.clearDisplay()
		e.drawSprite(3, 34, 1)

		if !e.Display[2*DisplayWidth+3] {
			t.Errorf("Pixel at (3,2) should be set (wrapped from Y=34)")
		}
	})

	t.Run("Sprite clipping at edge", func(t *testing.T) {
		e.clearDisplay()
		e.I = 0x300
		e.Memory[0x300] = 0xFF
		e.drawSprite(DisplayWidth-2, 5, 1)

		if !e.Display[5*DisplayWidth+(DisplayWidth-2)] {
			t.Errorf("Pixel at (%d,5) should be set", DisplayWidth-2)
		}
		if !e.Display[5*DisplayWidth+(DisplayWidth-1)] {
			t.Errorf("Pixel at (%d,5) should be set", DisplayWidth-1)
		}

		for x := range 6 {
			if e.Display[5*DisplayWidth+x] {
				t.Errorf("Pixel at (%d,5) should NOT be set (sprite should not wrap)", x)
			}
		}
	})
}
