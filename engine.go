package main

import (
	"errors"
	"fmt"
	"math"

	"github.com/ebitengine/oto/v3"
)

type Voice struct {
	freq   float64
	phase  float64
	active bool
}
type Engine struct {
	voices [32]Voice
}

type AudioStream struct {
	engine *Engine
	buffer []float32
}

func (e *Engine) Trigger(freq float64) {
	fmt.Println(e.voices)
	for i := 0; i < 32; i++ {
		if !e.voices[i].active {
			e.voices[i].freq = freq
			e.voices[i].phase = 0
			e.voices[i].active = true
			return
		}
	}
}

func (e *Engine) Mix(out []float32) {
	sampleRate := 44100.0

	for i := range out {
		sum := float32(0)

		for v := 0; v < 32; v++ {
			voice := &e.voices[v]

			if !voice.active {
				continue
			}

			sample := math.Sin(2 * math.Pi * voice.phase)

			voice.phase += voice.freq / sampleRate

			if voice.phase >= 1 {
				voice.phase -= 1
			}

			sum += float32(sample * 0.2)
		}

		if sum > 1 {
			sum = 1
		}
		if sum < -1 {
			sum = -1
		}

		out[i] = sum
	}
}

func (e *Engine) Vanish(index int) error {
	if index < 0 || index >= len(e.voices) {
		return errors.New("invalid voice index")
	}

	e.voices[index].active = false
	return nil
}

func startAudio(e *Engine) {
	ctx, ready, err := oto.NewContext(&oto.NewContextOptions{
		SampleRate:   44100,
		ChannelCount: 1,
		Format:       oto.FormatSignedInt16LE,
	})
	if err != nil {
		panic(err)
	}

	<-ready

	stream := &AudioStream{
		engine: e,
		buffer: make([]float32, 256),
	}

	player := ctx.NewPlayer(stream)
	player.Play()
}

func (s *AudioStream) Read(p []byte) (int, error) {
	s.engine.Mix(s.buffer)
	for i := range s.buffer {

		sample := int16(s.buffer[i] * 30000)

		p[i*2] = byte(sample)
		p[i*2+1] = byte(sample >> 8)
	}

	return len(s.buffer) * 2, nil
}
