package sequencer

import (
	"fmt"
	"time"

	"KSynth/internal/audio"
)

type RecordedNote struct {
	Freq   float64
	Offset time.Duration
}

type Track struct {
	events      []RecordedNote
	recording   bool
	recordStart time.Time
}

type Bank struct {
	tracks        [4]Track
	activeSlot    int
	recordingSlot int
}

func NewBank() *Bank {
	return &Bank{
		activeSlot:    0,
		recordingSlot: -1,
	}
}

func (s *Bank) SelectSlot(slot int) {
	if slot < 0 || slot >= len(s.tracks) {
		return
	}

	s.activeSlot = slot
	fmt.Printf("Selected sequencer slot %d\n", slot+1)
}

func (s *Bank) ToggleRecording(engine *audio.Engine) {
	track := &s.tracks[s.activeSlot]

	if track.recording {
		track.recording = false
		s.recordingSlot = -1
		if len(track.events) == 0 {
			_ = engine.ClearSequence(s.activeSlot)
			fmt.Printf("Slot %d stopped: no notes captured\n", s.activeSlot+1)
			return
		}

		fmt.Printf("Slot %d stopped: %d notes captured\n", s.activeSlot+1, len(track.events))
		s.generateSequence(engine, s.activeSlot)
		return
	}

	if s.recordingSlot >= 0 && s.recordingSlot != s.activeSlot {
		fmt.Printf("Slot %d is recording. Stop it before recording another slot\n", s.recordingSlot+1)
		return
	}

	_ = engine.ClearSequence(s.activeSlot)
	track.events = nil
	track.recording = true
	track.recordStart = time.Now()
	s.recordingSlot = s.activeSlot
	fmt.Printf("Slot %d armed: note presses are being captured\n", s.activeSlot+1)
}

func (s *Bank) RecordNote(freq float64) {
	if s.recordingSlot < 0 {
		return
	}

	track := &s.tracks[s.recordingSlot]
	if !track.recording {
		return
	}

	track.events = append(track.events, RecordedNote{
		Freq:   freq,
		Offset: time.Since(track.recordStart),
	})
}

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
