package main

import (
	"io"
	"math"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

type Sound struct {
	audioContext *audio.Context
	player       *audio.Player
	isPlaying    bool
}

// Custom stream that generates a sine wave
type SineWave struct {
	position  int64
	frequency float64
}

func (s *SineWave) Read(buf []byte) (int, error) {
	const amplitude = 0.3
	for i := 0; i < len(buf)/2; i++ {
		// Generate sine wave
		position := float64(s.position) / sampleRate
		sine := math.Sin(2*math.Pi*s.frequency*position) * amplitude

		// Convert to 16-bit PCM
		sample := int16(sine * 32767)

		// Write to buffer (little endian)
		buf[i*2] = byte(sample)
		buf[i*2+1] = byte(sample >> 8)

		s.position++
		// Prevent overflow by wrapping around
		if s.position == sampleRate*10 {
			s.position = 0
		}
	}
	return len(buf), nil
}

// Implement Seek to satisfy io.ReadSeeker interface
func (s *SineWave) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		s.position = offset
	case io.SeekCurrent:
		s.position += offset
	case io.SeekEnd:
		// For a sine wave, "end" doesn't really make sense, but we'll define it as a full cycle
		s.position = sampleRate - offset
	}
	return s.position, nil
}

func (s *SineWave) Close() error {
	return nil
}

func initSound() *Sound {
	audioContext := audio.NewContext(sampleRate)

	// Create a sine wave stream
	sineWave := &SineWave{frequency: frequency}

	// Create an infinite loop stream
	stream := audio.NewInfiniteLoop(sineWave, sampleRate)

	// Create a new player
	player, err := audioContext.NewPlayer(stream)
	if err != nil {
		return &Sound{audioContext: audioContext}
	}

	// Set the volume
	player.SetVolume(0.5)

	return &Sound{
		audioContext: audioContext,
		player:       player,
		isPlaying:    false,
	}
}

func (s *Sound) Play() {
	if s.player != nil && !s.isPlaying {
		s.player.Play()
		s.isPlaying = true
	}
}

func (s *Sound) Stop() {
	if s.player != nil && s.isPlaying {
		s.player.Pause()
		s.isPlaying = false
	}
}

func (g *Game) handleSound() {
	if g.emulator.SoundTimer > 0 {
		if !g.sound.isPlaying {
			g.sound.Play()
		}
	} else {
		if g.sound.isPlaying {
			g.sound.Stop()
		}
	}
}
