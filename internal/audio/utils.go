package audio

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

func sequenceEnvelope(elapsed int, total int) float64 {
	if total <= 0 {
		return 0
	}

	attack := envelopeRamp(elapsed, attackSamples)
	remaining := total - elapsed
	release := envelopeRamp(remaining, releaseSamples)
	return attack * release
}
