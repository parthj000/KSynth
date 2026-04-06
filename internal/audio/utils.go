package audio

import "math"

func envelopeRamp(position int, width int) float64 {
	if width <= 0 {
		return 1
	}
	if position <= 0 {
		return 0
	}
	if position >= width {
		return 1
	}
	return float64(position) / float64(width)
}

func sequenceEnvelope(mode SoundMode, elapsed int, total int) float64 {
	if total <= 0 {
		return 0
	}

	attackSamples, releaseSamples, _ := envelopeForMode(mode)
	attack := envelopeRamp(elapsed, attackSamples)
	remaining := total - elapsed
	release := envelopeRamp(remaining, releaseSamples)
	return attack * release * bodyForMode(mode, elapsed)
}

func sampleForMode(mode SoundMode, phase float64, elapsed int) float64 {
	switch mode {
	case SoundModeSine:
		return sineWave(phase)
	case SoundModePiano:
		return pianoWave(phase, elapsed)
	case SoundModeOrgan:
		return organWave(phase)
	default:
		return organWave(phase)
	}
}

func envelopeForMode(mode SoundMode) (attack int, release int, sustainedLevel float64) {
	switch mode {
	case SoundModeSine:
		return 48, 240, 1
	case SoundModePiano:
		return 24, 2200, 0.9
	case SoundModeOrgan:
		return 96, 900, 1
	default:
		return 96, 900, 1
	}
}

func bodyForMode(mode SoundMode, elapsed int) float64 {
	switch mode {
	case SoundModePiano:
		return pianoBodyDecay(elapsed)
	default:
		return 1
	}
}

func sineWave(phase float64) float64 {
	return math.Sin(2 * math.Pi * phase)
}

func organWave(phase float64) float64 {
	fundamental := math.Sin(2 * math.Pi * phase)
	second := math.Sin(2 * math.Pi * phase * 2)
	third := math.Sin(2 * math.Pi * phase * 3)
	fourth := math.Sin(2 * math.Pi * phase * 4)
	fifth := math.Sin(2 * math.Pi * phase * 5)
	sixth := math.Sin(2 * math.Pi * phase * 6)

	tone := 0.58*fundamental +
		0.24*second +
		0.12*third +
		0.08*fourth +
		0.05*fifth +
		0.03*sixth

	return tone
}

func pianoWave(phase float64, elapsed int) float64 {
	brightness := 0.45 + 0.55*math.Exp(-float64(elapsed)/float64(int(SampleRate)*2))

	fundamental := math.Sin(2 * math.Pi * phase)
	second := math.Sin(2 * math.Pi * phase * 2)
	third := math.Sin(2 * math.Pi * phase * 3)
	fourth := math.Sin(2 * math.Pi * phase * 4)
	inHarmonic := math.Sin(2 * math.Pi * phase * 6.8)

	tone := 0.72*fundamental +
		0.22*brightness*second +
		0.13*brightness*third +
		0.07*(0.7+0.3*brightness)*fourth +
		0.03*brightness*inHarmonic

	return tone
}

func pianoBodyDecay(elapsed int) float64 {
	return 0.55 + 0.45*math.Exp(-float64(elapsed)/float64(int(SampleRate)*3/2))
}
