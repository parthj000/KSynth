package audio

import (
	"errors"
	"sync"

	"github.com/ebitengine/oto/v3"
)

type Voice struct {
	Freq             float64
	Phase            float64
	Active           bool
	Sustained        bool
	RemainingSamples int
	TotalSamples     int
}

type SequenceEvent struct {
	Freq         float64
	StartSample  int
	LengthSample int
}

type SequenceLoop struct {
	events   []SequenceEvent
	length   int
	position int
	enabled  bool
}

type SoundMode string

const (
	SoundModeSine  SoundMode = "sine"
	SoundModePiano SoundMode = "piano"
	SoundModeOrgan SoundMode = "organ"
)

type Engine struct {
	mu               sync.Mutex
	voices           [32]Voice
	defaultDuration  int
	lastTriggeredIdx int
	sequences        [4]SequenceLoop
	soundMode        SoundMode
}

type AudioStream struct {
	engine   *Engine
	buffer   []float32
	recorder *Recorder
}

const (
	SampleRate        = 44100.0
	streamBufferSize  = 64
	defaultNoteVolume = 0.22
	sequenceVolume    = 0.17
	organAttackSamples  = 96
	organReleaseSamples = 900
)

func NewEngine(noteDurationSeconds float64) *Engine {
	durationSamples := int(noteDurationSeconds * SampleRate)
	if durationSamples <= 0 {
		durationSamples = int(SampleRate / 2)
	}

	return &Engine{
		defaultDuration:  durationSamples,
		lastTriggeredIdx: -1,
		soundMode:        SoundModeOrgan,
	}
}

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

func (e *Engine) Trigger(freq float64) {
	e.mu.Lock()
	defer e.mu.Unlock()

	for i := 0; i < len(e.voices); i++ {
		if !e.voices[i].Active {
			e.startVoice(i, freq, e.defaultDuration, false)
			e.lastTriggeredIdx = i
			return
		}
	}
}

func (e *Engine) TriggerOnVoice(index int, freq float64, durationSamples int, sustained bool) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if index < 0 || index >= len(e.voices) {
		return errors.New("invalid voice index")
	}

	if durationSamples <= 0 {
		durationSamples = e.defaultDuration
	}

	e.startVoice(index, freq, durationSamples, sustained)
	if index != len(e.voices)-1 {
		e.lastTriggeredIdx = index
	}

	return nil
}

func (e *Engine) Mix(out []float32) {
	e.mu.Lock()
	defer e.mu.Unlock()

	for i := range out {
		sum := float32(0)

		for v := 0; v < 32; v++ {
			voice := &e.voices[v]

			if !voice.Active {
				continue
			}

			elapsed := voice.TotalSamples - voice.RemainingSamples
			if elapsed < 0 {
				elapsed = 0
			}

			sample := sampleForMode(e.soundMode, voice.Phase, elapsed)

			voice.Phase += voice.Freq / SampleRate

			if voice.Phase >= 1 {
				voice.Phase -= 1
			}

			sum += float32(sample * defaultNoteVolume * e.voiceEnvelope(voice))

			if voice.Sustained {
				continue
			}

			voice.RemainingSamples--
			if voice.RemainingSamples <= 0 {
				voice.Active = false
			}
		}

		sequenceSample, _ := e.mixSequenceSample()
		sum += sequenceSample

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
	e.mu.Lock()
	defer e.mu.Unlock()

	if index < 0 || index >= len(e.voices) {
		return errors.New("invalid voice index")
	}

	e.voices[index] = Voice{}
	return nil
}

func (e *Engine) DefaultDurationSamples() int {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.defaultDuration
}

func (e *Engine) ToggleSustainLastVoice() (int, bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.lastTriggeredIdx < 0 || e.lastTriggeredIdx >= len(e.voices) {
		return -1, false, errors.New("no note has been triggered yet")
	}

	voice := &e.voices[e.lastTriggeredIdx]
	if !voice.Active {
		return -1, false, errors.New("last triggered voice is no longer active")
	}

	voice.Sustained = !voice.Sustained
	return e.lastTriggeredIdx, voice.Sustained, nil
}

func (e *Engine) SnapshotVoices() [32]Voice {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.voices
}

func (e *Engine) SetSoundMode(mode SoundMode) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	switch mode {
	case SoundModeSine, SoundModePiano, SoundModeOrgan:
		e.soundMode = mode
		return nil
	default:
		return errors.New("invalid sound mode")
	}
}

func (e *Engine) SoundMode() SoundMode {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.soundMode
}

func (e *Engine) startVoice(index int, freq float64, durationSamples int, sustained bool) {
	e.voices[index].Freq = freq
	e.voices[index].Phase = 0
	e.voices[index].Active = true
	e.voices[index].Sustained = sustained
	e.voices[index].RemainingSamples = durationSamples
	e.voices[index].TotalSamples = durationSamples
}

func (e *Engine) mixSequenceSample() (float32, int) {
	sum := 0.0
	activeEvents := 0

	for idx := range e.sequences {
		loop := &e.sequences[idx]
		if !loop.enabled || loop.length <= 0 || len(loop.events) == 0 {
			continue
		}

		currentSample := loop.position

		for _, event := range loop.events {
			if currentSample < event.StartSample || currentSample >= event.StartSample+event.LengthSample {
				continue
			}

			elapsed := currentSample - event.StartSample
			phase := float64(elapsed) * event.Freq / SampleRate
			sum += sampleForMode(e.soundMode, phase, elapsed) * sequenceVolume * sequenceEnvelope(e.soundMode, elapsed, event.LengthSample)
			activeEvents++
		}

		loop.position++
		if loop.position >= loop.length {
			loop.position = 0
		}
	}

	return float32(sum), activeEvents
}

func (e *Engine) voiceEnvelope(voice *Voice) float64 {
	elapsed := voice.TotalSamples - voice.RemainingSamples
	if elapsed < 0 {
		elapsed = 0
	}

	attackSamples, releaseSamples, sustainedLevel := envelopeForMode(e.soundMode)
	attack := envelopeRamp(elapsed, attackSamples)
	if voice.Sustained {
		return attack * sustainedLevel
	}

	release := envelopeRamp(voice.RemainingSamples, releaseSamples)
	return attack * release * bodyForMode(e.soundMode, elapsed)
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
