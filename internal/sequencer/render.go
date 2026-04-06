package sequencer

import (
	"fmt"
	"time"

	"KSynth/internal/audio"
)

func (s *Bank) generateSequence(engine *audio.Engine, slot int) {
	track := s.tracks[slot]
	durationSamples := engine.DefaultDurationSamples()
	events := make([]audio.SequenceEvent, 0, len(track.events))
	for _, event := range track.events {
		events = append(events, audio.SequenceEvent{
			Freq:         event.Freq,
			StartSample:  int(event.Offset.Seconds() * audio.SampleRate),
			LengthSample: durationSamples,
		})
	}

	loopLength := int((track.events[len(track.events)-1].Offset + 350*time.Millisecond).Seconds() * audio.SampleRate)
	if loopLength <= 0 {
		loopLength = durationSamples
	}

	if err := engine.SetSequence(slot, events, loopLength); err != nil {
		fmt.Println("Sequencer:", err)
		return
	}

	fmt.Printf("Sequencer slot %d playback started\n", slot+1)
}
