package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/tfriedel6/canvas"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/youpy/go-wav"
)

type ButtonString struct {
	Name     string
	Value    string
	Position sdl.Rect
}
type ButtonBool struct {
	Value        bool
	PositionUp   sdl.Rect
	PositionDown sdl.Rect
}
type ButtonPlusMoins struct {
	Value        float64
	PositionUp   sdl.Rect
	PositionDown sdl.Rect
}
type ButtonWaveform struct {
	// sine || triangle || square || Flat
	Value            string
	PositionSine     sdl.Rect
	PositionTriangle sdl.Rect
	PositionSquare   sdl.Rect
	PositionFlat     sdl.Rect
}
type ButtonColor struct {
	Value         sdl.Color
	PositionUpR   sdl.Rect
	PositionDownR sdl.Rect
	PositionUpG   sdl.Rect
	PositionDownG sdl.Rect
	PositionUpB   sdl.Rect
	PositionDownB sdl.Rect
}
type Style struct {
	Font     *canvas.Font
	FontSize float64
}
type SelectorOscillator struct {
	TextMargin int
	ButtonString
	sdl.Color
}
type Champ struct {
	Champ string
	Cible string
	ButtonBool
}

type ClientInterface struct {
	SampleRate        float64
	BitsPerSample     uint16
	SeeCondensedWave  bool
	AddOscillator     bool
	DeleteOscillator  bool
	ReloadSelector    bool
	ReloadingSelector bool
	ReloadWave        bool
	ReloadingWave     bool
	CurrentOscillator int
	Style             Style
	Oscillator        []*Oscillator
	DeroulantSelector ButtonBool
	Selector          []SelectorOscillator
	Wave              [][]float32
	CondensedWave     []float32
	CloseFileName     ButtonString
	FileName          Champ
	OpenMix           sdl.Rect
	Enregistrer       sdl.Rect
	SaveAll           sdl.Rect
	Lire              sdl.Rect
	//Stopper           sdl.Rect
}

func (ci *ClientInterface) InitInterface(cv *canvas.Canvas) {
	var err error
	ci.SampleRate = 44100
	ci.BitsPerSample = 16
	ci.Style.FontSize = 20.0
	ci.Style.Font, err = cv.LoadFont("font/Gaulois.ttf")
	if err != nil {
		log.Fatal(err)
	}
	cv.SetFont(ci.Style.Font, ci.Style.FontSize)
	osc = clientInterface.newOscillator(cv)
	ci.Oscillator = append(ci.Oscillator, osc)
	longueurChamp := cv.MeasureText("Choose filename and press enter").Width + 10
	ci.CloseFileName = ButtonString{
		Name:  "Close filename",
		Value: "X",
		Position: sdl.Rect{
			X: int32(wndWidth - int(longueurChamp) - 50),
			Y: int32(wndHeight - 350),
			W: 50,
			H: 50,
		},
	}
	ci.FileName = Champ{
		Champ: "",
		ButtonBool: ButtonBool{
			Value:        false,
			PositionUp:   sdl.Rect{X: int32(wndWidth - int(longueurChamp)), Y: int32(wndHeight - 350), W: int32(longueurChamp), H: 50},
			PositionDown: sdl.Rect{X: int32(wndWidth - int(longueurChamp)), Y: int32(wndHeight - 350), W: int32(longueurChamp), H: 50},
		},
	}
	ci.OpenMix = sdl.Rect{X: int32(wndWidth - 300), Y: int32(wndHeight - 250), W: 100, H: 50}
	ci.Enregistrer = sdl.Rect{X: int32(wndWidth - 150), Y: int32(wndHeight - 250), W: 100, H: 50}
	ci.SaveAll = sdl.Rect{X: int32(wndWidth - 150), Y: int32(wndHeight - 190), W: 100, H: 50}
	ci.Lire = sdl.Rect{X: int32(wndWidth - 250), Y: int32(wndHeight - 100), W: 100, H: 50}
	//ci.Stopper = sdl.Rect{X: int32(wndWidth - 150), Y: int32(wndHeight - 100), W: 100, H: 50}
	width := cv.MeasureText(">"+strconv.Itoa(len(ci.Oscillator))).Width + 10
	ci.DeroulantSelector = ButtonBool{
		Value: true,
		PositionUp: sdl.Rect{
			X: int32(wndWidth - int(width)),
			Y: 0,
			W: int32(width),
			H: int32(ci.Style.FontSize*2) + 10,
		},
		PositionDown: sdl.Rect{},
	}
	ci.reloadSelectorOscillator(cv)
	ci.AddOscillator = false
	ci.ReloadWave = true
	ci.ReloadingWave = false
	ci.ReloadingSelector = false
	ci.CurrentOscillator = 0
}

func (ci *ClientInterface) reloadWave() {
	wave, _ := osc.GenerateWave()
	if len(ci.Wave)-1 < ci.CurrentOscillator || len(ci.Wave) == 0 {
		ci.Wave = append(ci.Wave, wave)
	} else {
		ci.Wave[ci.CurrentOscillator] = wave
	}
	if strings.Contains(souris.Click, "Only") {
		souris.Click = ""
	}
	ci.ReloadingWave = false
}

func (ci *ClientInterface) reloadSelectorOscillator(cv *canvas.Canvas) {
	textMargin := 10
	height := ci.Style.FontSize + float64(textMargin) //cv.MeasureText("Osc "+strconv.Itoa(len(ci.Oscillator))).ActualBoundingBoxDescent + float64(textMargin)
	width := cv.MeasureText("Osc "+strconv.Itoa(len(ci.Oscillator))).Width + float64(textMargin)
	widthAdd := cv.MeasureText("Add Osc").Width + float64(textMargin)
	startY := (ci.Style.FontSize * 2) + float64(textMargin)
	startX := wndWidth - int(width)
	ci.Selector = []SelectorOscillator{
		{
			textMargin / 2,
			ButtonString{
				"Add Osc",
				"Add Osc",
				sdl.Rect{
					X: int32(wndWidth) - int32(widthAdd),
					Y: int32(startY),
					W: int32(widthAdd),
					H: int32(height),
				},
			},
			sdl.Color{
				R: 0,
				G: 255,
				B: 0,
				A: 255,
			},
		},
		{
			textMargin / 2,
			ButtonString{
				"Del Osc",
				"Del Osc",
				sdl.Rect{
					X: int32(wndWidth) - (int32(widthAdd) * 2) - 14,
					Y: int32(startY),
					W: int32(widthAdd),
					H: int32(height),
				},
			},
			sdl.Color{
				R: 255,
				G: 0,
				B: 0,
				A: 255,
			},
		},
	}
	startY += height
	for i := 0; i < len(ci.Oscillator); i++ {
		selector := SelectorOscillator{
			textMargin / 2,
			ButtonString{
				"Osc " + strconv.Itoa(i+1),
				"Osc " + strconv.Itoa(i+1),
				sdl.Rect{
					X: int32(startX), //int32(wndWidth) - int32(width),
					Y: int32(startY),
					W: int32(width),
					H: int32(height),
				},
			},
			ci.Oscillator[i].Color.Value,
		}
		ci.Selector = append(ci.Selector, selector)
		startY += height
		if startY > float64(wndHeight/2) {
			startY = (ci.Style.FontSize * 2) + float64(textMargin) + height
			startX -= int(width)
		}
	}
	ci.ReloadingSelector = false
}

func (ci *ClientInterface) newOscillator(cv *canvas.Canvas) *Oscillator {
	measureOnlyPos := cv.MeasureText("Only Positive").Width
	measureOnlyNeg := cv.MeasureText("Only Negative").Width
	measureOnlyIncreas := cv.MeasureText("Cut Up").Width
	measureOnlyDecreas := cv.MeasureText("Cut Down").Width
	fmt.Println("measuer:", measureOnlyPos)
	return &Oscillator{
		BitsPerSample: 16, //bits.UintSize,
		//configurable bool fields
		OnlyPositive: ButtonBool{
			Value:        false,
			PositionUp:   sdl.Rect{X: 50 + int32(measureOnlyPos) + 10, Y: int32(wndHeight - 625), W: 70, H: 25},
			PositionDown: sdl.Rect{X: 50 + int32(measureOnlyPos) + 10, Y: int32(wndHeight - 600), W: 70, H: 25},
		},
		OnlyNegative: ButtonBool{
			Value:        false,
			PositionUp:   sdl.Rect{X: 300 + int32(measureOnlyNeg) + 10, Y: int32(wndHeight - 625), W: 70, H: 25},
			PositionDown: sdl.Rect{X: 300 + int32(measureOnlyNeg) + 10, Y: int32(wndHeight - 600), W: 70, H: 25},
		},
		OnlyIncreased: ButtonBool{
			Value:        false,
			PositionUp:   sdl.Rect{X: 600 + int32(measureOnlyIncreas) + 10, Y: int32(wndHeight - 625), W: 70, H: 25},
			PositionDown: sdl.Rect{X: 600 + int32(measureOnlyIncreas) + 10, Y: int32(wndHeight - 600), W: 70, H: 25},
		},
		OnlyDecreased: ButtonBool{
			Value:        false,
			PositionUp:   sdl.Rect{X: 900 + int32(measureOnlyDecreas) + 10, Y: int32(wndHeight - 625), W: 70, H: 25},
			PositionDown: sdl.Rect{X: 900 + int32(measureOnlyDecreas) + 10, Y: int32(wndHeight - 600), W: 70, H: 25},
		},

		// configurable fields column 1
		SoundDuration: ButtonPlusMoins{
			Value:        5,
			PositionUp:   sdl.Rect{X: 150, Y: int32(wndHeight - 500), W: 70, H: 50},
			PositionDown: sdl.Rect{X: 50, Y: int32(wndHeight - 500), W: 70, H: 50},
		},
		Amplitude: ButtonPlusMoins{
			Value:        1,
			PositionUp:   sdl.Rect{X: 150, Y: int32(wndHeight - 400), W: 70, H: 50},
			PositionDown: sdl.Rect{X: 50, Y: int32(wndHeight - 400), W: 70, H: 50},
		},
		MaxAmplitude: ButtonPlusMoins{
			Value:        1,
			PositionUp:   sdl.Rect{X: 150, Y: int32(wndHeight - 300), W: 70, H: 50},
			PositionDown: sdl.Rect{X: 50, Y: int32(wndHeight - 300), W: 70, H: 50},
		},
		Frequency: ButtonPlusMoins{
			Value:        1,
			PositionUp:   sdl.Rect{X: 150, Y: int32(wndHeight - 200), W: 70, H: 50},
			PositionDown: sdl.Rect{X: 50, Y: int32(wndHeight - 200), W: 70, H: 50},
		},
		Phase: ButtonPlusMoins{
			Value:        0,
			PositionUp:   sdl.Rect{X: 150, Y: int32(wndHeight - 100), W: 70, H: 50},
			PositionDown: sdl.Rect{X: 50, Y: int32(wndHeight - 100), W: 70, H: 50},
		},
		// configurable fields column 2
		Waveform: ButtonWaveform{
			Value:            "sine",
			PositionSine:     sdl.Rect{X: 300, Y: int32(wndHeight - 500), W: 80, H: 25},
			PositionTriangle: sdl.Rect{X: 380, Y: int32(wndHeight - 500), W: 80, H: 25},
			PositionSquare:   sdl.Rect{X: 300, Y: int32(wndHeight - 475), W: 80, H: 25},
			PositionFlat:     sdl.Rect{X: 380, Y: int32(wndHeight - 475), W: 80, H: 25},
		},
		Kick: ButtonPlusMoins{
			Value:        0,
			PositionUp:   sdl.Rect{X: 400, Y: int32(wndHeight - 400), W: 70, H: 50},
			PositionDown: sdl.Rect{X: 300, Y: int32(wndHeight - 400), W: 70, H: 50},
		},
		Hauteur: ButtonPlusMoins{
			Value:        0,
			PositionUp:   sdl.Rect{X: 400, Y: int32(wndHeight - 300), W: 70, H: 50},
			PositionDown: sdl.Rect{X: 300, Y: int32(wndHeight - 300), W: 70, H: 50},
		},
		AsymetrieY: ButtonPlusMoins{
			Value:        0,
			PositionUp:   sdl.Rect{X: 400, Y: int32(wndHeight - 200), W: 70, H: 50},
			PositionDown: sdl.Rect{X: 300, Y: int32(wndHeight - 200), W: 70, H: 50},
		},
		// configurable fields column 3
		Lecture: ButtonPlusMoins{
			Value:        5,
			PositionUp:   sdl.Rect{X: 650, Y: int32(wndHeight - 500), W: 70, H: 50},
			PositionDown: sdl.Rect{X: 550, Y: int32(wndHeight - 500), W: 70, H: 50},
		},
		Pause: ButtonPlusMoins{
			Value:        0,
			PositionUp:   sdl.Rect{X: 650, Y: int32(wndHeight - 400), W: 70, H: 50},
			PositionDown: sdl.Rect{X: 550, Y: int32(wndHeight - 400), W: 70, H: 50},
		},
		Start: ButtonPlusMoins{
			Value:        0,
			PositionUp:   sdl.Rect{X: 650, Y: int32(wndHeight - 300), W: 70, H: 50},
			PositionDown: sdl.Rect{X: 550, Y: int32(wndHeight - 300), W: 70, H: 50},
		},
		End: ButtonPlusMoins{
			Value:        5,
			PositionUp:   sdl.Rect{X: 650, Y: int32(wndHeight - 200), W: 70, H: 50},
			PositionDown: sdl.Rect{X: 550, Y: int32(wndHeight - 200), W: 70, H: 50},
		},
		// configurable fields column 4
		Color: ButtonColor{
			Value:         sdl.Color{R: 0, G: 255, B: 255, A: 255},
			PositionUpR:   sdl.Rect{X: 900, Y: int32(wndHeight - 500), W: 70, H: 50},
			PositionDownR: sdl.Rect{X: 800, Y: int32(wndHeight - 500), W: 70, H: 50},
			PositionUpG:   sdl.Rect{X: 900, Y: int32(wndHeight - 400), W: 70, H: 50},
			PositionDownG: sdl.Rect{X: 800, Y: int32(wndHeight - 400), W: 70, H: 50},
			PositionUpB:   sdl.Rect{X: 900, Y: int32(wndHeight - 300), W: 70, H: 50},
			PositionDownB: sdl.Rect{X: 800, Y: int32(wndHeight - 300), W: 70, H: 50},
		},
	}
}

func (ci *ClientInterface) condenseWave() {
	indexOfGreaterWave := 0
	for i := 0; i < len(ci.Wave); i++ {
		if len(ci.Wave[indexOfGreaterWave]) < len(ci.Wave[i]) {
			indexOfGreaterWave = i
		}
	}
	samplesGreater := len(ci.Wave[indexOfGreaterWave])
	condesedWave := make([]float32, samplesGreater)
	diviser := make([]float32, samplesGreater)
	for i := 0; i < len(ci.Wave); i++ {
		wave := ci.Wave[i]
		samplesThis := len(ci.Wave[i])
		for j := 0; j < samplesThis; j++ {
			if wave[j] != 0 {
				condesedWave[j] += wave[j]
				diviser[j]++
			}
		}
	}
	for i := 0; i < samplesGreater; i++ {
		condesedWave[i] /= diviser[i]
	}
	ci.CondensedWave = condesedWave
}

func (ci *ClientInterface) amalgameWave() {
	indexOfGreaterWave := 0
	for i := 0; i < len(ci.Wave); i++ {
		if len(ci.Wave[indexOfGreaterWave]) < len(ci.Wave[i]) {
			indexOfGreaterWave = i
		}
	}
	samplesGreater := len(ci.Wave[indexOfGreaterWave])
	totalWave := make([]float32, samplesGreater*len(ci.Wave))
	for i := 0; i < len(ci.Wave); i++ {
		wave := ci.Wave[i]
		samplesThis := len(wave)
		for j := 0; j < samplesThis; j++ {
			ref := (j * len(ci.Wave)) + i
			totalWave[ref] = wave[j]
		}
	}
	amalgamedWave := make([]float32, samplesGreater)
	ref := 0
	for i := 0; i < samplesGreater; i++ {
		indexInTot := (i * len(ci.Wave)) + ref
		amalgamedWave[i] = totalWave[indexInTot]
		ref++
		if ref == len(ci.Wave) {
			ref = 0
			if len(ci.Wave)%2 == 0 {
				ref -= 1
			}
		}
	}
	ci.CondensedWave = amalgamedWave
}

func (ci *ClientInterface) SaveToWav(filename, exportOrSave string) error {
	if exportOrSave == "export" {
		ci.condenseWave()
	} else if exportOrSave == "save" {
		ci.amalgameWave()
	}
	data := ci.CondensedWave
	samples := len(data)

	// Create WAV file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("erreur lors de la création du fichier WAV : %w", err)
	}
	defer file.Close()

	writer := wav.NewWriter(file, uint32(samples), 1, uint32(ci.SampleRate), ci.BitsPerSample)
	var arrBytes []byte
	for _, sample := range data {
		arrBytes = append(arrBytes, byte(sample*(32767*8)))
	}
	writer.Write(arrBytes)

	/*
		go func(filename string) {
			var cmd *exec.Cmd

			switch runtime.GOOS {
			case "windows":
				cmd = exec.Command("cmd", "/C", "start", filename)
			case "linux":
				cmd = exec.Command("bash", "-c", fmt.Sprintf("mplayer %s", filename))
			default:
				fmt.Println("Unsupported OS")
				return
			}

			err := cmd.Run()
			if err != nil {
				fmt.Printf("Erreur lors de la lecture du fichier WAV : %v\n", err)
			}
		}(filename)
	*/
	return nil
}

func (ci *ClientInterface) saveAll(dirName string) {
	// Create directory
	path := "res/user/all/" + dirName
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		fmt.Printf("Erreur lors de la création du dossier : %v\n", err)
		return
	}
	for i, o := range ci.Oscillator {
		fileName := "osc_" + strconv.Itoa(i) + ".json"
		o.save(path + "/" + fileName)
	}
	ci.SaveToWav(path+"/wave.wav", "save")
}

// charge un oscillateur depuis un fichier JSON dans le dossier res/user
func (ci *ClientInterface) loadOsc(path string) *Oscillator {
	arrByte, err := os.ReadFile("res/user" + path)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier JSON:", err.Error())
	}
	var loadedOscillator Oscillator
	err = json.Unmarshal(arrByte, &loadedOscillator)
	if err != nil {
		fmt.Println("Erreur lors de la conversion du JSON en oscillateur:", err.Error())
		return nil
	}
	return &loadedOscillator
}

// retourne le noms des fichiers enregistrés dans le dossier /res/user
func (ci *ClientInterface) parseNameOfFileRegistred() []string {
	files, err := os.ReadDir("res/user/all")
	if err != nil {
		fmt.Println("Erreur lors de la lecture du dossier /res/user/all :", err.Error())
	}
	var filesName []string
	for _, file := range files {
		filesName = append(filesName, file.Name())
	}
	return filesName
}

func (ci *ClientInterface) loadMix(dirName string) {
	path := "res/user/all/" + dirName
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du dossier /res/user/all/"+dirName, err.Error())
	}
	var oscs []*Oscillator
	for _, file := range files {
		if strings.Contains(file.Name(), "osc_") {
			o := ci.loadOsc("/all/" + dirName + "/" + file.Name())
			oscs = append(oscs, o)
		}
	}
	ci.Oscillator = oscs
	ci.CurrentOscillator = 0
	ci.ReloadSelector = true
	ci.ReloadWave = true
}
