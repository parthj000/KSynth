package audio

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
