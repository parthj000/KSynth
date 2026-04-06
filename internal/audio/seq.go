package audio

import "errors"

func (e *Engine) SetSequence(slot int, events []SequenceEvent, loopLengthSamples int) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if slot < 0 || slot >= len(e.sequences) {
		return errors.New("invalid sequence slot")
	}

	if len(events) == 0 || loopLengthSamples <= 0 {
		e.sequences[slot] = SequenceLoop{}
		return nil
	}

	e.sequences[slot] = SequenceLoop{
		events:   append([]SequenceEvent(nil), events...),
		length:   loopLengthSamples,
		position: 0,
		enabled:  true,
	}

	return nil
}

func (e *Engine) ClearSequence(slot int) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if slot < 0 || slot >= len(e.sequences) {
		return errors.New("invalid sequence slot")
	}

	e.sequences[slot] = SequenceLoop{}
	return nil
}

func (e *Engine) ClearAllSequences() {
	e.mu.Lock()
	defer e.mu.Unlock()

	for i := range e.sequences {
		e.sequences[i] = SequenceLoop{}
	}
}
