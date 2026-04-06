package audio

import "github.com/ebitengine/oto/v3"

func StartAudio(e *Engine) *Recorder {
	ctx, ready, err := oto.NewContext(&oto.NewContextOptions{
		SampleRate:   int(SampleRate),
		ChannelCount: 1,
		Format:       oto.FormatSignedInt16LE,
	})
	if err != nil {
		panic(err)
	}

	<-ready

	stream := &AudioStream{
		engine:   e,
		buffer:   make([]float32, streamBufferSize),
		recorder: NewRecorder(),
	}

	player := ctx.NewPlayer(stream)
	player.Play()
	return stream.recorder
}

func (s *AudioStream) Read(p []byte) (int, error) {
	s.engine.Mix(s.buffer)
	for i := range s.buffer {
		sample := int16(s.buffer[i] * 30000)

		p[i*2] = byte(sample)
		p[i*2+1] = byte(sample >> 8)
	}

	s.recorder.WritePCM(p[:len(s.buffer)*2])

	return len(s.buffer) * 2, nil
}
