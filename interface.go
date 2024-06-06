package main

import (
	"fmt"
	"log"

	"github.com/tfriedel6/canvas"
	"github.com/veandco/go-sdl2/sdl"
)

type Style struct {
	Font     *canvas.Font
	FontSize float64
}

type ClientInterface struct {
	Style       Style
	Oscillator  *Oscillator
	Enregistrer sdl.Rect
	Lire        sdl.Rect
	Stopper     sdl.Rect
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
	measure := cv.MeasureText("Only Positive").Width
	fmt.Println("measuer:", measure)
	osc = &Oscillator{
		BitsPerSample: 16, //bits.UintSize,
		// configurable fields column 1
		OnlyPositive: OnlyPositive{
			Value:        false,
			PositionUp:   sdl.Rect{X: 50 + int32(measure) + 10, Y: int32(wndHeight - 625), W: 70, H: 25},
			PositionDown: sdl.Rect{X: 50 + int32(measure) + 10, Y: int32(wndHeight - 600), W: 70, H: 25},
		},
		SoundDuration: SoundDuration{
			Value:        5,
			PositionUp:   sdl.Rect{X: 150, Y: int32(wndHeight - 500), W: 70, H: 50},
			PositionDown: sdl.Rect{X: 50, Y: int32(wndHeight - 500), W: 70, H: 50},
		},
		Amplitude: Amplitude{
			Value:        1,
			PositionUp:   sdl.Rect{X: 150, Y: int32(wndHeight - 400), W: 70, H: 50},
			PositionDown: sdl.Rect{X: 50, Y: int32(wndHeight - 400), W: 70, H: 50},
		},
		SampleRate: SampleRate{
			Value:        44100,
			PositionUp:   sdl.Rect{X: 150, Y: int32(wndHeight - 300), W: 70, H: 50},
			PositionDown: sdl.Rect{X: 50, Y: int32(wndHeight - 300), W: 70, H: 50},
		},
		Frequency: Frequency{
			Value:        1,
			PositionUp:   sdl.Rect{X: 150, Y: int32(wndHeight - 200), W: 70, H: 50},
			PositionDown: sdl.Rect{X: 50, Y: int32(wndHeight - 200), W: 70, H: 50},
		},
		Phase: Phase{
			Value:        0,
			PositionUp:   sdl.Rect{X: 150, Y: int32(wndHeight - 100), W: 70, H: 50},
			PositionDown: sdl.Rect{X: 50, Y: int32(wndHeight - 100), W: 70, H: 50},
		},
		// configurable fields column 2
		Waveform: Waveform{
			Value:            "sine",
			PositionSine:     sdl.Rect{X: 300, Y: int32(wndHeight - 525), W: 80, H: 25},
			PositionTriangle: sdl.Rect{X: 380, Y: int32(wndHeight - 525), W: 80, H: 25},
			PositionSquare:   sdl.Rect{X: 300, Y: int32(wndHeight - 500), W: 80, H: 25},
			PositionFlat:     sdl.Rect{X: 380, Y: int32(wndHeight - 500), W: 80, H: 25},
		},
		Kick: Kick{
			Value:        0,
			PositionUp:   sdl.Rect{X: 400, Y: int32(wndHeight - 400), W: 70, H: 50},
			PositionDown: sdl.Rect{X: 300, Y: int32(wndHeight - 400), W: 70, H: 50},
		},
		AsymetrieX: Asymetrie{
			Value:        0,
			PositionUp:   sdl.Rect{X: 400, Y: int32(wndHeight - 300), W: 70, H: 50},
			PositionDown: sdl.Rect{X: 300, Y: int32(wndHeight - 300), W: 70, H: 50},
		},
		AsymetrieY: Asymetrie{
			Value:        0,
			PositionUp:   sdl.Rect{X: 400, Y: int32(wndHeight - 200), W: 70, H: 50},
			PositionDown: sdl.Rect{X: 300, Y: int32(wndHeight - 200), W: 70, H: 50},
		},
	}
	ci.Oscillator = osc
	ci.Enregistrer = sdl.Rect{X: int32(wndWidth - 150), Y: int32(wndHeight - 250), W: 100, H: 50}
	ci.Lire = sdl.Rect{X: int32(wndWidth - 250), Y: int32(wndHeight - 100), W: 100, H: 50}
	ci.Stopper = sdl.Rect{X: int32(wndWidth - 150), Y: int32(wndHeight - 100), W: 100, H: 50}
	/*
		ci = &ClientInterface{
			Oscillator:  osc,
			Enregistrer: sdl.Rect{X: int32(wndWidth - 150), Y: int32(wndHeight - 250), W: 100, H: 50},
			Lire:        sdl.Rect{X: int32(wndWidth - 250), Y: int32(wndHeight - 100), W: 100, H: 50},
			Stopper:     sdl.Rect{X: int32(wndWidth - 150), Y: int32(wndHeight - 100), W: 100, H: 50},
		}
	*/
	fmt.Println(ci.Style)
	//cv.SetFont(ci.Style.Font, ci.Style.FontSize)

}

func setWindow(screenX, screenY int) {
	echelle = float64(screenX) / 1920
	wndWidth = screenX
	wndHeight = screenY
}
