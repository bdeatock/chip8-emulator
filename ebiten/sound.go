package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

type stream struct {
	pos int64
}

func (s *stream) Read(buf []byte) (int, error) {
	const bytesPerSample = 8

	n := len(buf) / bytesPerSample * bytesPerSample

	const length = sampleRate / frequency

	for i := range n / bytesPerSample {
		v := math.Float32bits(float32(math.Sin(2 * math.Pi * float64(s.pos/bytesPerSample+int64(i)) / length)))
		buf[8*i] = byte(v)
		buf[8*i+1] = byte(v >> 8)
		buf[8*i+2] = byte(v >> 16)
		buf[8*i+3] = byte(v >> 24)
		buf[8*i+4] = byte(v)
		buf[8*i+5] = byte(v >> 8)
		buf[8*i+6] = byte(v >> 16)
		buf[8*i+7] = byte(v >> 24)
	}

	s.pos += int64(n)
	s.pos %= length * bytesPerSample

	return n, nil
}

func (s *stream) Close() error {
	return nil
}

func (g *Game) initSound() error {
	g.audioContext = audio.NewContext(sampleRate)

	var err error
	if g.audioPlayer, err = g.audioContext.NewPlayerF32(&stream{}); err != nil {
		return err
	}

	g.audioPlayer.SetVolume(0.2)

	return nil
}

func (g *Game) handleSound() {
	if g.emulator.SoundTimer > 0 && !g.audioPlayer.IsPlaying() {
		g.audioPlayer.Play()
	} else if g.emulator.SoundTimer == 0 && g.audioPlayer.IsPlaying() {
		g.audioPlayer.Pause()
	}
}
