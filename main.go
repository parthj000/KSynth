package main

import (
	"fmt"

	"github.com/eiannone/keyboard"
)

func main() {

	engine := &Engine{}
	fmt.Println((*engine).voices)
	startAudio(engine)

	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	fmt.Println("Hi the stimulator is started")

	var freqMap = map[rune]float64{
		'a': 261.63, // C4
		's': 293.66, // D4
		'd': 329.63, // E4
		'f': 349.23, // F4
		'g': 392.00, // G4
		'h': 440.00, // A4
		'w': 277.18, // C#4
		'e': 311.13, // D#4
		't': 369.99, // F#4
		'y': 415.30, // G#4
		'u': 466.16, // A#4
	}

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		if key == keyboard.KeyEsc {
			break
		}

		if freq, ok := freqMap[char]; ok {
			engine.Trigger(freq)
			continue
		}

		switch char {

		case 'v':
			fmt.Println("Voices:", engine.voices)

		case 'k':
			var input int
			fmt.Print("Enter voice index: ")
			fmt.Scan(&input)
			engine.Vanish(input)
		}
	}
}
