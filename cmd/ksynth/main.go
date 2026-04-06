package main

import (
	"flag"
	"fmt"
	"math"
	"path/filepath"
	"time"

	"KSynth/internal/audio"
	"KSynth/internal/sequencer"
	"github.com/eiannone/keyboard"
)

func main() {
	help := flag.Bool("help", false, "show KSynth usage")
	flag.BoolVar(help, "h", false, "show KSynth usage")
	sound := flag.String("sound", string(audio.SoundModeOrgan), "initial sound mode: sine, piano, organ")
	flag.Parse()

	if *help {
		printBanner()
		printHelp()
		return
	}

	engine := audio.NewEngine(0.3)
	if err := engine.SetSoundMode(audio.SoundMode(*sound)); err != nil {
		panic(err)
	}
	seq := sequencer.NewBank()
	recorder := audio.StartAudio(engine)

	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	printBanner()
	printControls()

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
		case '7':
			if err := engine.SetSoundMode(audio.SoundModeSine); err != nil {
				fmt.Println("Sound:", err)
				continue
			}
			fmt.Println("Sound mode: sine")
			continue
		case '8':
			if err := engine.SetSoundMode(audio.SoundModePiano); err != nil {
				fmt.Println("Sound:", err)
				continue
			}
			fmt.Println("Sound mode: piano")
			continue
		case '9':
			if err := engine.SetSoundMode(audio.SoundModeOrgan); err != nil {
				fmt.Println("Sound:", err)
				continue
			}
			fmt.Println("Sound mode: organ")
			continue
		case '+':
			octaveShift++
			fmt.Printf("Octave shift: %+d\n", octaveShift)
			continue
		case '-':
			octaveShift--
			fmt.Printf("Octave shift: %+d\n", octaveShift)
			continue
		case 'r':
			var seconds float64
			fmt.Print("Export duration in seconds: ")
			fmt.Scan(&seconds)

			filename := fmt.Sprintf("recording-%s.wav", time.Now().Format("20060102-150405"))
			path := filepath.Join(".", filename)
			if err := recorder.Start(path, seconds); err != nil {
				fmt.Println("Export:", err)
				continue
			}

			fmt.Printf("Exporting live output to %s for %.2f seconds\n", filename, seconds)
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

func printBanner() {
	fmt.Println(` _  __ ____              _   _     
| |/ // ___| _   _ _ __ | |_| |__  
| ' / \___ \| | | | '_ \| __| '_ \ 
| . \  ___) | |_| | | | | |_| | | |
|_|\_\|____/ \__, |_| |_|\__|_| |_|
             |___/                 `)
}

func printControls() {
	fmt.Println("Terminal synth and sequencer")
	fmt.Println("Note keys: a s d f g h with sharps on w e t y u")
	fmt.Println("Press 1-4 to select a sequencer slot")
	fmt.Println("Press 7 for sine, 8 for piano, 9 for organ")
	fmt.Println("Press Space to start/stop recording on the selected slot")
	fmt.Println("Press + or - to shift octave up or down")
	fmt.Println("Press r to export live output to a wav file")
	fmt.Println("Press j to toggle sustain on the last triggered note")
	fmt.Println("Press v to inspect voices, k to stop a voice, Esc to quit")
}

func printHelp() {
	fmt.Println("Usage: ksynth [flags]")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("  -h, -help    Show this help message")
	fmt.Println("  -sound       Initial sound mode: sine, piano, organ")
	fmt.Println("")
	printControls()
	fmt.Println("")
	fmt.Println("Export:")
	fmt.Println("  Press r while running, then enter a duration in seconds.")
	fmt.Println("  KSynth writes a wav file into the current directory.")
}
