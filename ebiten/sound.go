package main

import (
	"encoding/binary"
	"math"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

// stream implements audio.ReadCloser interface for sine wave gen
type stream struct {
	pos int64
}

// Read generates sine wave data and writes it to buffer, then updates the stream
// position to ensure continuous playback across multiple calls
func (s *stream) Read(buf []byte) (int, error) {
	const bytesPerSample = 8
	const samplesPerCycle = sampleRate / frequency

	sampleCount := len(buf) / bytesPerSample

	for i := 0; i < sampleCount; i++ {
		// Calculate sine wave value
		phase := float64((s.pos/bytesPerSample)+int64(i)) / float64(samplesPerCycle)
		sineValue := float32(math.Sin(2 * math.Pi * phase))

		bits := math.Float32bits(sineValue)

		// Write identical values to left and right channels
		offset := i * bytesPerSample
		binary.LittleEndian.PutUint32(buf[offset:], bits)
		binary.LittleEndian.PutUint32(buf[offset+4:], bits)
	}

	bytesWritten := sampleCount * bytesPerSample
	s.pos = (s.pos + int64(bytesWritten)) % (samplesPerCycle * bytesPerSample)

	return bytesWritten, nil
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
