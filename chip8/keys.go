package chip8

import "fmt"

func (e *Emulator) PressKey(key byte) error {
	if key > 0xF {
		return fmt.Errorf("invalid key: %X", key)
	}
	e.Keypad[key] = true
	return nil
}

func (e *Emulator) ReleaseKey(key byte) error {
	if key > 0xF {
		return fmt.Errorf("invalid key: %X", key)
	}
	e.Keypad[key] = false
	return nil
}
