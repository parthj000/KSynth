package main

import (
	"fmt"
	"math"

	"KSynth/internal/audio"
	"KSynth/internal/sequencer"
	"github.com/eiannone/keyboard"
)

func main() {
	engine := audio.NewEngine(0.3)
	seq := sequencer.NewBank()
	audio.StartAudio(engine)

	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	fmt.Println("Hi the stimulator is started")
	fmt.Println("Note keys: a s d f g h with sharps on w e t y u")
	fmt.Println("Press 1-4 to select a sequencer slot")
	fmt.Println("Press Space to start/stop recording on the selected slot")
	fmt.Println("Press + or - to shift octave up or down")
	fmt.Println("Press j to toggle sustain on the last triggered note")
	fmt.Println("Press v to inspect voices, k to stop a voice, Esc to quit")

	freqMap := map[rune]float64{
		'a': 261.63,
		's': 293.66,
		'd': 329.63,
		'f': 349.23,
		'g': 392.00,
		'h': 440.00,
		'w': 277.18,
		'e': 311.13,
		't': 369.99,
		'y': 415.30,
		'u': 466.16,
	}
	octaveShift := 0

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		if key == keyboard.KeyEsc {
			engine.ClearAllSequences()
			break
		}

		switch char {
		case '1', '2', '3', '4':
			seq.SelectSlot(int(char - '1'))
			continue
		case '+':
			octaveShift++
			fmt.Printf("Octave shift: %+d\n", octaveShift)
			continue
		case '-':
			octaveShift--
			fmt.Printf("Octave shift: %+d\n", octaveShift)
			continue
		}

		if key == keyboard.KeySpace || char == ' ' {
			seq.ToggleRecording(engine)
			continue
		}

		if freq, ok := freqMap[char]; ok {
			freq *= math.Pow(2, float64(octaveShift))
			engine.Trigger(freq)
			seq.RecordNote(freq)
			continue
		}

		switch char {
		case 'v':
			fmt.Println("Voices:", engine.SnapshotVoices())
		case 'j':
			index, sustained, err := engine.ToggleSustainLastVoice()
			if err != nil {
				fmt.Println("Sustain:", err)
				continue
			}
			fmt.Printf("Voice %d sustain=%t\n", index, sustained)
		case 'k':
			var input int
			fmt.Print("Enter voice index: ")
			fmt.Scan(&input)
			if err := engine.Vanish(input); err != nil {
				fmt.Println("Vanish:", err)
			}
		}
	}
}
