package main

import (
	"fmt"
	"math"
	"os"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/youpy/go-wav"
)

// Structures pour l'oscillateur
type Amplitude struct {
	Value        float64
	PositionUp   sdl.Rect
	PositionDown sdl.Rect
}

type Frequency struct {
	Value        float64
	PositionUp   sdl.Rect
	PositionDown sdl.Rect
}

type Phase struct {
	Value        float64
	PositionUp   sdl.Rect
	PositionDown sdl.Rect
}

type Oscillator struct {
	Amplitude Amplitude
	Frequency Frequency
	Phase     Phase
}

type ClientInterface struct {
	Oscillator  *Oscillator
	Enregistrer sdl.Rect
	Lire        sdl.Rect
	Stopper     sdl.Rect
}

func (o *Oscillator) Value(t float64) float64 {
	return o.Amplitude.Value * math.Sin(o.Frequency.Value*t+o.Phase.Value)
}

func (o *Oscillator) Update(dt float64) {
	o.Phase.Value += dt * o.Frequency.Value
}

func (o *Oscillator) Reset() {
	o.Phase.Value = 0
}

func (o *Oscillator) SetAmplitude(a float64) {
	o.Amplitude.Value = a
}

func (o *Oscillator) SetFrequency(f float64) {
	o.Frequency.Value = f
}

func (o *Oscillator) SetPhase(p float64) {
	o.Phase.Value = p
}

func (o *Oscillator) SetParams(a, f, p float64) {
	o.Amplitude.Value = a
	o.Frequency.Value = f
	o.Phase.Value = p
}

func (o *Oscillator) GenerateWave(samples int, sampleRate int) []float32 {
	data := make([]float32, samples)
	for i := 0; i < samples; i++ {
		t := float64(i) / float64(sampleRate)
		data[i] = float32(o.Value(t))
	}
	return data
}

func (o *Oscillator) SaveToWav(filename string) error {
	// Generate 5 seconds of audio
	sampleRate := 44100
	samples := sampleRate * 5
	data := o.GenerateWave(samples, sampleRate)

	// Create WAV file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Erreur lors de la création du fichier WAV : %w", err)
	}
	defer file.Close()

	writer := wav.NewWriter(file, uint32(samples), 1, uint32(sampleRate), 2)
	var arrBytes []byte
	for _, sample := range data {
		arrBytes = append(arrBytes, byte(sample*32767))
		//writer.Write(int(sample * 32767))
	}
	writer.Write(arrBytes)
	return nil
}

func (o *Oscillator) play(filename string) {
	audio(filename)
}
