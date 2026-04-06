package sequencer

import "time"

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
