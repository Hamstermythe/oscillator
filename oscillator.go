package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"

	"github.com/youpy/go-wav"
)

type Oscillator struct {
	// unsettable fields
	LastGreater   float64
	Increase      bool
	BitsPerSample uint16
	// configurable bool fields
	OnlyPositive  ButtonBool
	OnlyNegative  ButtonBool
	OnlyIncreased ButtonBool
	OnlyDecreased ButtonBool
	// configurable fields column 1
	SoundDuration ButtonPlusMoins
	Amplitude     ButtonPlusMoins
	MaxAmplitude  ButtonPlusMoins
	Frequency     ButtonPlusMoins
	Phase         ButtonPlusMoins
	// configurable fields column 2
	Waveform   ButtonWaveform
	Kick       ButtonPlusMoins
	Hauteur    ButtonPlusMoins
	AsymetrieY ButtonPlusMoins
	// configurable fields column 3
	Lecture ButtonPlusMoins
	Pause   ButtonPlusMoins
	Start   ButtonPlusMoins
	End     ButtonPlusMoins
	// configurable fields column 4
	Color ButtonColor
}

func (o *Oscillator) kick(phase float64) float64 {
	sine := math.Sin(phase)
	if sine > 0 {
		return math.Pow(sine, 1/o.Kick.Value)
	}
	return -math.Pow(-sine, 1/o.Kick.Value)
}

func (o *Oscillator) hauteur(phase, p, value float64) float64 {
	if value == 0 {
		return value
	}
	mulitple := 1 - math.Abs(math.Sin(phase))
	oscil := (o.Hauteur.Value) * mulitple
	if int(p)%2 == 0 {
		return value + oscil
	}
	return value - oscil
}

func (o *Oscillator) InverseurY(phase float64) float64 {
	sine := math.Sin(phase)
	if sine > 0 {
		return math.Pow(sine, o.AsymetrieY.Value)
	}
	return -math.Pow(-sine, o.AsymetrieY.Value)
}

func (o *Oscillator) cutAmplitude(value float64) float64 {
	if value > o.MaxAmplitude.Value {
		return o.MaxAmplitude.Value
	} else if value < -o.MaxAmplitude.Value {
		return -o.MaxAmplitude.Value
	}
	return value
}

// retourne true si la valeur est contenu dans l'interval de lecture
func (o *Oscillator) boucleController(p float64) bool {
	// position = p / sampleRate
	// si position < démarrage => val = 0
	// si position > arret => val = 0
	// temps de rotation de lecture = lecture + pause
	// position en temps = position / sampleRate
	// si position en temps > lecture => val = 0
	if p < o.Start.Value*clientInterface.SampleRate*2 || p > o.End.Value*clientInterface.SampleRate*2 {
		return false
	}
	if o.Lecture.Value < 0.09 && o.Lecture.Value > -0.09 {
		return true
	}
	totDuration := o.Lecture.Value + o.Pause.Value
	boucleNumber := o.SoundDuration.Value / totDuration
	if boucleNumber == 0 {
		boucleNumber = 1
	}
	pos := (o.SoundDuration.Value * clientInterface.SampleRate * 2) / boucleNumber
	ratioPos := float64(int(p-(o.Start.Value*clientInterface.SampleRate*2))%int(pos)) / pos
	return ratioPos < (o.Lecture.Value / totDuration)
}

func (o *Oscillator) Value(p float64) float64 {
	if !o.boucleController(p) {
		return 0
	}
	t := float64(p) / float64(clientInterface.SampleRate) //(o.SoundDuration.Value*o.Frequency.Value)) // / float64(o.Frequency.Value))
	phase := o.Frequency.Value*t + o.Phase.Value
	//phase := t + o.Phase.Value/o.Frequency.Value
	val := 0.0
	switch o.Waveform.Value {
	case "sine":
		if o.Kick.Value != 0 {
			val = o.Amplitude.Value * o.kick(phase)
		} else {
			val = o.Amplitude.Value * math.Cos(phase)
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
	if o.AsymetrieY.Value > 0 || o.AsymetrieY.Value < 0 {
		asymVal = asymVal * o.InverseurY(phase)
	}
	if o.Hauteur.Value > 0 || o.Hauteur.Value < 0 {
		asymVal = o.hauteur(phase, p, val)
	}
	asymVal = o.cutAmplitude(asymVal)

	asymVal = o.Amplitude.Value * asymVal
	//if math.Sin(phase) == 0.0 {
	if math.Sin(phase) < 0.01 && math.Sin(phase) > -0.01 {
		o.Increase = true
	}
	//if math.Sin(phase) == 1 || math.Sin(phase) == -1 {
	if math.Sin(phase) > 0.99 || math.Sin(phase) < -0.99 {
		o.Increase = false
		o.LastGreater = asymVal
	}
	if o.OnlyPositive.Value && asymVal < 0 {
		asymVal = 0.0
	}
	if o.OnlyNegative.Value && asymVal > 0 {
		asymVal = 0.0
	}
	if o.OnlyIncreased.Value && o.Increase {
		asymVal = 0.0
	}
	if o.OnlyDecreased.Value && !o.Increase {
		asymVal = 0.0
	}
	return asymVal
}

func (o *Oscillator) GenerateWave() ([]float32, int) {
	ratioTimeBits := float64(o.BitsPerSample / 8)
	samples := int(o.SoundDuration.Value * clientInterface.SampleRate * ratioTimeBits)
	//samples := int(o.Frequency.Value * o.SoundDuration.Value * o.SampleRate.Value)
	data := make([]float32, samples)
	for p := 0; p < samples; p++ {
		val := o.Value(float64(p))
		data[p] = float32(val)
	}
	return data, samples
}

// enregistre l'oscillateur dans un fichier JSON dans le dossier res/user
func (o *Oscillator) save(path string) {
	arrByte, err := json.Marshal(o)
	if err != nil {
		fmt.Println("Erreur lors de la conversion de l'oscillateur en JSON:", err.Error())
	}
	err = os.WriteFile(path, arrByte, 0644)
	if err != nil {
		fmt.Println("Erreur lors de l'écriture du fichier JSON:", err.Error())
	}
}

// charge un oscillateur depuis un fichier JSON dans le dossier res/user
func (o *Oscillator) load(filename string) {
	arrByte, err := os.ReadFile("res/user" + filename)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier JSON:", err.Error())
	}
	var loadedOscillator Oscillator
	err = json.Unmarshal(arrByte, &loadedOscillator)
	if err != nil {
		fmt.Println("Erreur lors de la conversion du JSON en oscillateur:", err.Error())
		return
	}
	clientInterface.Oscillator = append(clientInterface.Oscillator, &loadedOscillator)
	clientInterface.CurrentOscillator = len(clientInterface.Oscillator) - 1
	clientInterface.ReloadSelector = true
	clientInterface.ReloadWave = true
}

// retourne le noms des fichiers enregistrés dans le dossier /res/user
func (o *Oscillator) parseNameOfFileRegistred() []string {
	files, err := os.ReadDir("res/user")
	if err != nil {
		fmt.Println("Erreur lors de la lecture du dossier /res/user :", err.Error())
	}
	var filesName []string
	for _, file := range files {
		filesName = append(filesName, file.Name())
	}
	return filesName
}

func (o *Oscillator) readWaveFile(filename string) ([]float32, float64) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, 0 //, fmt.Errorf("Erreur lors de l'ouverture du fichier WAV : %w", err)
	}
	defer file.Close()

	reader := wav.NewReader(file)
	samples, err := reader.ReadSamples()
	if err != nil {
		return nil, 0 //, fmt.Errorf("Erreur lors de la lecture des échantillons du fichier WAV : %w", err)
	}
	var fileData []float32
	for _, sample := range samples {
		fileData = append(fileData, float32(sample.Values[0])/(32767*8))
	}
	soundDuration := float64(len(samples)) / float64(clientInterface.SampleRate)

	return fileData, soundDuration //len(samples)
}
