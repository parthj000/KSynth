package audio

import "sync"

type Voice struct {
	Freq             float64
	Key              rune
	Label            string
	Phase            float64
	Active           bool
	Sustained        bool
	RemainingSamples int
	TotalSamples     int
}

type VoiceInfo struct {
	Index     int
	Key       rune
	Label     string
	Freq      float64
	Sustained bool
	SoundMode SoundMode
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
