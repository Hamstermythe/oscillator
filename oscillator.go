package main

import (
	"fmt"
	"math"
	"os"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/youpy/go-wav"
)

type OnlyPositive struct {
	Value        bool
	PositionUp   sdl.Rect
	PositionDown sdl.Rect
}

type SoundDuration struct {
	Value        float64
	PositionUp   sdl.Rect
	PositionDown sdl.Rect
}

// Structures pour l'oscillateur
type Amplitude struct {
	Value        float64
	PositionUp   sdl.Rect
	PositionDown sdl.Rect
}

type SampleRate struct {
	Value        int
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

type Waveform struct {
	// sine || triangle || square || Flat
	Value            string
	PositionSine     sdl.Rect
	PositionTriangle sdl.Rect
	PositionSquare   sdl.Rect
	PositionFlat     sdl.Rect
}

type Kick struct {
	Value        float64
	PositionUp   sdl.Rect
	PositionDown sdl.Rect
}

type Oscillator struct {
	// unsettable fields
	BitsPerSample uint16
	// configurable fields column 1
	OnlyPositive  OnlyPositive
	SoundDuration SoundDuration
	Amplitude     Amplitude
	SampleRate    SampleRate
	Frequency     Frequency
	Phase         Phase
	// configurable fields column 2
	Waveform Waveform
	Kick     Kick
}

// Fonction skewedSine
func (o *Oscillator) kick(phase float64) float64 {
	sine := math.Sin(phase)
	if sine > 0 {
		return math.Pow(sine, 1/o.Kick.Value)
	}
	return -math.Pow(-sine, o.Kick.Value)
}

func (o *Oscillator) Value(t float64) float64 {
	phase := o.Frequency.Value*t + o.Phase.Value
	switch o.Waveform.Value {
	case "sine":
		return o.Amplitude.Value * o.kick(phase)
	case "triangle":
		triangle := (2 / math.Pi) * math.Asin(math.Sin(phase))
		return o.Amplitude.Value * o.kick(triangle)
	case "square":
		if math.Sin(phase) >= 0 {
			return o.Amplitude.Value
		}
		return -o.Amplitude.Value
	case "Flat":
		return o.Amplitude.Value
	}
	return 0
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
		val := o.Value(t)
		if o.OnlyPositive.Value {
			if val < 0 {
				val = 0
			}
		}
		data[i] = float32(val)
	}
	return data
}

func (o *Oscillator) SaveToWav(filename string) error {
	ratioTimeBits := float64(o.BitsPerSample / 8)
	fmt.Println("ratioTimeBits: ", ratioTimeBits, ", o.BitsPerSample: ", o.BitsPerSample)
	samples := int(float64(o.SampleRate.Value) * o.SoundDuration.Value * ratioTimeBits)
	//samples := int((44100.0 / float64(o.SampleRate.Value)) * (44100.0 * 5.0))
	//samples := int(float64(o.SampleRate.Value) * (baseNumber / float64(o.SampleRate.Value)))
	//samples := o.SampleRate.Value * int(20.0/osc.Duration)
	data := o.GenerateWave(samples, o.SampleRate.Value)

	// Create WAV file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Erreur lors de la création du fichier WAV : %w", err)
	}
	defer file.Close()

	writer := wav.NewWriter(file, uint32(samples), 1, uint32(o.SampleRate.Value), o.BitsPerSample)
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
