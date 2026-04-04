package main

import (
	"fmt"
	"math"

	"github.com/eiannone/keyboard"
)

type Calculator struct {
	Name string
}

func (calc *Calculator) add(a int, b int) int {
	return a + b
}

func generateSine(freq float64, duration float64) []float32 {

	sampleRate := 44100
	totalSamples := int(duration * float64(sampleRate))

	buf := make([]float32, totalSamples)

	for i := 0; i < totalSamples; i++ {

		// time is 1/frequency
		//
		t := float64(i) / float64(sampleRate)

		//sin wave coordinate
		val := math.Sin(2 * math.Pi * freq * t)

		// amplitude
		buf[i] = float32(val * 0.3) // volume
	}

	return buf
}
func main() {

	engine := &Engine{}
	startAudio(engine)

	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	fmt.Println("Hi the stimulator is started")
	// sine1 := generateSine(440, 0.1)
	sine2 := generateSine(660, 0.05)

	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		if key == keyboard.KeyEsc {
			break
		}

		engine.Trigger(sine2)

		// fmt.Printf("You pressed: rune %q, key %X\n", char, key)
	}
}
