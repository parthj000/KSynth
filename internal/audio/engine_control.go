package audio

import "errors"

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

func (e *Engine) Trigger(freq float64) {
	e.mu.Lock()
	defer e.mu.Unlock()

	for i := 0; i < len(e.voices); i++ {
		if !e.voices[i].Active {
			e.startVoice(i, freq, "", 0, e.defaultDuration, false)
			e.lastTriggeredIdx = i
			return
		}
	}
}

func (e *Engine) TriggerLabeled(freq float64, key rune, label string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	for i := 0; i < len(e.voices); i++ {
		if !e.voices[i].Active {
			e.startVoice(i, freq, label, key, e.defaultDuration, false)
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

	e.startVoice(index, freq, "", 0, durationSamples, sustained)
	if index != len(e.voices)-1 {
		e.lastTriggeredIdx = index
	}

	return nil
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

func (e *Engine) ActiveVoices() []VoiceInfo {
	e.mu.Lock()
	defer e.mu.Unlock()

	voices := make([]VoiceInfo, 0, len(e.voices))
	for i, voice := range e.voices {
		if !voice.Active {
			continue
		}

		voices = append(voices, VoiceInfo{
			Index:     i,
			Key:       voice.Key,
			Label:     voice.Label,
			Freq:      voice.Freq,
			Sustained: voice.Sustained,
			SoundMode: e.soundMode,
		})
	}

	return voices
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

func (e *Engine) startVoice(index int, freq float64, label string, key rune, durationSamples int, sustained bool) {
	e.voices[index].Freq = freq
	e.voices[index].Key = key
	e.voices[index].Label = label
	e.voices[index].Phase = 0
	e.voices[index].Active = true
	e.voices[index].Sustained = sustained
	e.voices[index].RemainingSamples = durationSamples
	e.voices[index].TotalSamples = durationSamples
}
