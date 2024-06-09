package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/tfriedel6/canvas"
	"github.com/veandco/go-sdl2/sdl"
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

type ClientInterface struct {
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
	Enregistrer       sdl.Rect
	Lire              sdl.Rect
	Stopper           sdl.Rect
}

func (ci *ClientInterface) InitInterface(cv *canvas.Canvas) {
	var err error
	//fontSize := 20.0
	ci.Style.FontSize = 20.0
	ci.Style.Font, err = cv.LoadFont("font/Gaulois.ttf")
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(style)
	cv.SetFont(ci.Style.Font, ci.Style.FontSize)
	//fmt.Println(style)
	osc = clientInterface.newOscillator(cv)
	ci.Oscillator = append(ci.Oscillator, osc)
	ci.Enregistrer = sdl.Rect{X: int32(wndWidth - 150), Y: int32(wndHeight - 250), W: 100, H: 50}
	ci.Lire = sdl.Rect{X: int32(wndWidth - 250), Y: int32(wndHeight - 100), W: 100, H: 50}
	ci.Stopper = sdl.Rect{X: int32(wndWidth - 150), Y: int32(wndHeight - 100), W: 100, H: 50}
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
		SampleRate: ButtonPlusMoins{
			Value:        44100,
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

func setWindow(screenX, screenY int) {
	echelle = float64(screenX) / 1920
	wndWidth = screenX
	wndHeight = screenY
}
