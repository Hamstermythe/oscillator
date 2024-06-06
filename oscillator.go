package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"

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

type Asymetrie struct {
	Value        float64
	PositionUp   sdl.Rect
	PositionDown sdl.Rect
}

type Oscillator struct {
	// unsettable fields
	Increase      bool
	BitsPerSample uint16
	// configurable fields column 1
	OnlyPositive  OnlyPositive
	SoundDuration SoundDuration
	Amplitude     Amplitude
	SampleRate    SampleRate
	Frequency     Frequency
	Phase         Phase
	// configurable fields column 2
	Waveform   Waveform
	Kick       Kick
	AsymetrieX Asymetrie
	AsymetrieY Asymetrie
}

// Fonction skewedSine
func (o *Oscillator) kick(phase float64) float64 {
	sine := math.Sin(phase)
	if sine > 0 {
		return math.Pow(sine, 1/o.Kick.Value)
	}
	return -math.Pow(-sine, o.Kick.Value)
}

// Fonction pour une montée exponentielle et une descente logarithmique
func (o *Oscillator) asymmetricTransformX(value float64) float64 {
	if value > 0 {
		return math.Pow(value, 1/o.AsymetrieX.Value)
	}
	return -math.Log(-value+1) * o.AsymetrieX.Value
}

// Fonction pour une montée exponentielle et une descente logarithmique
func (o *Oscillator) expAndExp(asymVal, phase float64) float64 {
	if o.Increase {
		return asymVal * math.Pow(math.Abs(math.Sin(phase)), o.AsymetrieY.Value)
	} else {
		return asymVal * (1 - math.Log(math.Abs(math.Sin(phase))+1)/(math.Log(2)*o.AsymetrieY.Value))
	}
}

func (o *Oscillator) Value(t float64) float64 {
	phase := o.Frequency.Value*t + o.Phase.Value
	if math.Sin(phase) < 0.02 && math.Sin(phase) > -0.02 {
		o.Increase = true
	}
	if math.Sin(phase) > 0.98 || math.Sin(phase) < -0.98 {
		o.Increase = false
	}
	val := 0.0
	switch o.Waveform.Value {
	case "sine":
		if o.Kick.Value != 0 {
			val = o.Amplitude.Value * o.kick(phase)
		} else {
			val = o.Amplitude.Value * math.Sin(phase)
		}
	case "triangle":
		triangle := (2 / math.Pi) * math.Asin(math.Sin(phase))
		if o.Kick.Value != 0 {
			val = o.Amplitude.Value * o.kick(triangle)
		} else {
			val = o.Amplitude.Value * triangle //math.Sin(phase)
		}
	case "square":
		if math.Sin(phase) >= 0 {
			if o.Kick.Value != 0 {
				val = o.Amplitude.Value * o.kick(phase)
			} else {
				val = o.Amplitude.Value
			}
		} else {
			if o.Kick.Value != 0 {
				val = o.Amplitude.Value * o.kick(phase)
			} else {
				val = -o.Amplitude.Value
			}
		}
	case "Flat":
		if o.Kick.Value != 0 {
			val = o.Amplitude.Value * o.kick(phase)
		} else {
			val = o.Amplitude.Value
		}
	}
	asymVal := val
	if o.AsymetrieX.Value > 0 || o.AsymetrieX.Value < 0 {
		asymVal = o.asymmetricTransformX(val)
	}
	if o.AsymetrieY.Value > 0 || o.AsymetrieY.Value < 0 {
		asymVal = o.expAndExp(asymVal, phase)
	}

	return o.Amplitude.Value * asymVal
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

func (o *Oscillator) GenerateWave() ([]float32, int) {
	ratioTimeBits := float64(o.BitsPerSample / 8)
	fmt.Println("ratioTimeBits: ", ratioTimeBits, ", o.BitsPerSample: ", o.BitsPerSample)
	samples := int(float64(o.SampleRate.Value) * o.SoundDuration.Value * ratioTimeBits)
	data := make([]float32, samples)
	for i := 0; i < samples; i++ {
		t := float64(i) / float64(o.SampleRate.Value)
		val := o.Value(t)
		if o.OnlyPositive.Value {
			if val < 0 {
				val = 0
			}
		}
		data[i] = float32(val)
	}
	return data, samples
}

func (o *Oscillator) SaveToWav(filename string) error {
	//ratioTimeBits := float64(o.BitsPerSample / 8)
	//fmt.Println("ratioTimeBits: ", ratioTimeBits, ", o.BitsPerSample: ", o.BitsPerSample)
	//samples := int(float64(o.SampleRate.Value) * o.SoundDuration.Value * ratioTimeBits)
	//samples := int((44100.0 / float64(o.SampleRate.Value)) * (44100.0 * 5.0))
	//samples := int(float64(o.SampleRate.Value) * (baseNumber / float64(o.SampleRate.Value)))
	//samples := o.SampleRate.Value * int(20.0/osc.Duration)
	data, samples := o.GenerateWave() //samples, o.SampleRate.Value)

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

	go func(filename string) {
		bashCommand := fmt.Sprintf("mplayer %s ", filename)
		cmd := exec.Command("bash", "-c", bashCommand)
		err = cmd.Run()
		if err != nil {
			fmt.Println("Erreur lors de la lecture du fichier WAV : %w", err)
		}
	}(filename)

	return nil
}

func (o *Oscillator) play(filename string) {
	audio(filename)
}
