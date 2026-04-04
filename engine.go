package main

import (
	"github.com/ebitengine/oto/v3"
)

type Voice struct {
	buffer []float32
	pos    int
	active bool
}

type Engine struct {
	voices [32]Voice
}

type AudioStream struct {
	engine *Engine
	buffer []float32
}

func (e *Engine) Trigger(buf []float32) {

	for i := 0; i < 32; i++ {
		if !e.voices[i].active {
			e.voices[i].buffer = buf
			e.voices[i].pos = 0
			e.voices[i].active = true
			return
		}
	}
}

func (e *Engine) Mix(out []float32) {
	for i := range out {

		sum := float32(0)

		for v := 0; v < 32; v++ {

			voice := &e.voices[v]

			if !voice.active {
				continue
			}

			if voice.pos >= len(voice.buffer) {
				voice.active = false
				continue
			}

			sum += voice.buffer[voice.pos]
			voice.pos++
		}

		out[i] = sum
	}
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
