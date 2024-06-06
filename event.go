package main

import (
	"fmt"
	"time"

	"github.com/tfriedel6/canvas/sdlcanvas"
	"github.com/veandco/go-sdl2/sdl"
)

type Souris struct {
	X, Y      int
	Click     string
	DateClick time.Time
}

func (s *Souris) Action() {
	if time.Since(s.DateClick) < 250*time.Millisecond {
		return
	}
	if s.Click == "SoundDuration+" {
		osc.SoundDuration.Value += 0.1
	} else if s.Click == "SoundDuration-" {
		osc.SoundDuration.Value -= 0.1
	} else if s.Click == "Amplitude+" {
		osc.Amplitude.Value += 0.1
	} else if s.Click == "Amplitude-" {
		osc.Amplitude.Value -= 0.1
	} else if s.Click == "SampleRate+" {
		osc.SampleRate.Value += 100
	} else if s.Click == "SampleRate-" {
		osc.SampleRate.Value -= 100
	} else if s.Click == "Frequency+" {
		osc.Frequency.Value += 1
	} else if s.Click == "Frequency-" {
		osc.Frequency.Value -= 1
	} else if s.Click == "Phase+" {
		osc.Phase.Value += 0.1
	} else if s.Click == "Phase-" {
		osc.Phase.Value -= 0.1
	}
	if s.Click == "Kick+" {
		osc.Kick.Value += 0.1
	} else if s.Click == "Kick-" {
		osc.Kick.Value -= 0.1
	} else if s.Click == "AsymetrieX+" {
		osc.AsymetrieX.Value += 0.1
	} else if s.Click == "AsymetrieX-" {
		osc.AsymetrieX.Value -= 0.1
	} else if s.Click == "AsymetrieY+" {
		osc.AsymetrieY.Value += 0.1
	} else if s.Click == "AsymetrieY-" {
		osc.AsymetrieY.Value -= 0.1
	}
}

func addEvent(wnd *sdlcanvas.Window) *sdlcanvas.Window {
	wnd.Event = func(event sdl.Event) {
		switch event := event.(type) {
		case *sdl.QuitEvent:
			running = false
			return
		case *sdl.DisplayEvent:
		case *sdl.MouseMotionEvent:
			mouseMoveUpdate(float64(event.X), float64(event.Y))
		case *sdl.MouseWheelEvent:
			// mouseWheelEvent(int(event.X), int(event.Y))
		case *sdl.MouseButtonEvent:
			if event.State == sdl.PRESSED {
				mouseDownUpdate(int(event.Button), int(event.X), int(event.Y))
			} else if event.State == sdl.RELEASED {
				mouseUpUpdate(int(event.Button))
			}
		case *sdl.KeyboardEvent:
			if event.Keysym.Scancode == sdl.SCANCODE_ESCAPE || event.Keysym.Scancode == sdl.SCANCODE_F12 {
				running = false
				return
			}
		}
	}
	return wnd
}

func mouseMoveUpdate(x, y float64) {
	souris.X = int(x)
	souris.Y = int(y)
}

func mouseUpUpdate(button int) {
	if button == 1 {
		souris.Click = ""
	}
}

func mouseDownUpdate(button, x, y int) {
	if button != 1 {
		return
	}
	// configurables fields colunm 1
	if x > int(osc.OnlyPositive.PositionUp.X) && x < int(osc.OnlyPositive.PositionUp.X+osc.OnlyPositive.PositionUp.W) && y > int(osc.OnlyPositive.PositionUp.Y) && y < int(osc.OnlyPositive.PositionUp.Y+osc.OnlyPositive.PositionUp.H) {
		osc.OnlyPositive.Value = !osc.OnlyPositive.Value
		fmt.Println("osc.OnlyPositive.Value : ", osc.OnlyPositive.Value)
	}
	if x > int(osc.OnlyPositive.PositionDown.X) && x < int(osc.OnlyPositive.PositionDown.X+osc.OnlyPositive.PositionDown.W) && y > int(osc.OnlyPositive.PositionDown.Y) && y < int(osc.OnlyPositive.PositionDown.Y+osc.OnlyPositive.PositionDown.H) {
		osc.OnlyPositive.Value = !osc.OnlyPositive.Value
		fmt.Println("osc.OnlyPositive.Value : ", osc.OnlyPositive.Value)
	}
	if x > int(osc.SoundDuration.PositionUp.X) && x < int(osc.SoundDuration.PositionUp.X+osc.SoundDuration.PositionUp.W) && y > int(osc.SoundDuration.PositionUp.Y) && y < int(osc.SoundDuration.PositionUp.Y+osc.SoundDuration.PositionUp.H) {
		souris.Click = "SoundDuration+"
		souris.DateClick = time.Now()
		osc.SoundDuration.Value += 0.1
		fmt.Println("osc.SoundDuration.Value up : ", osc.SoundDuration.Value)
	}
	if x > int(osc.SoundDuration.PositionDown.X) && x < int(osc.SoundDuration.PositionDown.X+osc.SoundDuration.PositionDown.W) && y > int(osc.SoundDuration.PositionDown.Y) && y < int(osc.SoundDuration.PositionDown.Y+osc.SoundDuration.PositionDown.H) {
		souris.Click = "SoundDuration-"
		souris.DateClick = time.Now()
		osc.SoundDuration.Value -= 0.1
		fmt.Println("osc.SoundDuration.Value down : ", osc.SoundDuration.Value)
	}
	if x > int(osc.Amplitude.PositionUp.X) && x < int(osc.Amplitude.PositionUp.X+osc.Amplitude.PositionUp.W) && y > int(osc.Amplitude.PositionUp.Y) && y < int(osc.Amplitude.PositionUp.Y+osc.Amplitude.PositionUp.H) {
		souris.Click = "Amplitude+"
		souris.DateClick = time.Now()
		osc.Amplitude.Value += 0.1
		fmt.Println("osc.Amplitude.Value up : ", osc.Amplitude.Value)
	}
	if x > int(osc.Amplitude.PositionDown.X) && x < int(osc.Amplitude.PositionDown.X+osc.Amplitude.PositionDown.W) && y > int(osc.Amplitude.PositionDown.Y) && y < int(osc.Amplitude.PositionDown.Y+osc.Amplitude.PositionDown.H) {
		souris.Click = "Amplitude-"
		souris.DateClick = time.Now()
		osc.Amplitude.Value -= 0.1
		fmt.Println("osc.Amplitude.Value down : ", osc.Amplitude.Value)
	}
	if x > int(osc.SampleRate.PositionUp.X) && x < int(osc.SampleRate.PositionUp.X+osc.SampleRate.PositionUp.W) && y > int(osc.SampleRate.PositionUp.Y) && y < int(osc.SampleRate.PositionUp.Y+osc.SampleRate.PositionUp.H) {
		souris.Click = "SampleRate+"
		souris.DateClick = time.Now()
		osc.SampleRate.Value += 100
		fmt.Println("osc.SampleRate.Value up : ", osc.SampleRate.Value)
	}
	if x > int(osc.SampleRate.PositionDown.X) && x < int(osc.SampleRate.PositionDown.X+osc.SampleRate.PositionDown.W) && y > int(osc.SampleRate.PositionDown.Y) && y < int(osc.SampleRate.PositionDown.Y+osc.SampleRate.PositionDown.H) {
		souris.Click = "SampleRate-"
		souris.DateClick = time.Now()
		osc.SampleRate.Value -= 100
		fmt.Println("osc.SampleRate.Value down : ", osc.SampleRate.Value)
	}
	if x > int(osc.Frequency.PositionUp.X) && x < int(osc.Frequency.PositionUp.X+osc.Frequency.PositionUp.W) && y > int(osc.Frequency.PositionUp.Y) && y < int(osc.Frequency.PositionUp.Y+osc.Frequency.PositionUp.H) {
		souris.Click = "Frequency+"
		souris.DateClick = time.Now()
		osc.Frequency.Value += 0.1
		fmt.Println("osc.Frequency.Value up : ", osc.Frequency.Value)
	}
	if x > int(osc.Frequency.PositionDown.X) && x < int(osc.Frequency.PositionDown.X+osc.Frequency.PositionDown.W) && y > int(osc.Frequency.PositionDown.Y) && y < int(osc.Frequency.PositionDown.Y+osc.Frequency.PositionDown.H) {
		souris.Click = "Frequency-"
		souris.DateClick = time.Now()
		osc.Frequency.Value -= 0.1
		fmt.Println("osc.Frequency.Value down : ", osc.Frequency.Value)
	}
	if x > int(osc.Phase.PositionUp.X) && x < int(osc.Phase.PositionUp.X+osc.Phase.PositionUp.W) && y > int(osc.Phase.PositionUp.Y) && y < int(osc.Phase.PositionUp.Y+osc.Phase.PositionUp.H) {
		souris.Click = "Phase+"
		souris.DateClick = time.Now()
		osc.Phase.Value += 0.1
		fmt.Println("osc.Phase.Value up : ", osc.Phase.Value)
	}
	if x > int(osc.Phase.PositionDown.X) && x < int(osc.Phase.PositionDown.X+osc.Phase.PositionDown.W) && y > int(osc.Phase.PositionDown.Y) && y < int(osc.Phase.PositionDown.Y+osc.Phase.PositionDown.H) {
		souris.Click = "Phase-"
		souris.DateClick = time.Now()
		osc.Phase.Value -= 0.1
		fmt.Println("osc.Phase.Value down : ", osc.Phase.Value)
	}
	// configurables fields column 2
	if x > int(osc.Waveform.PositionSine.X) && x < int(osc.Waveform.PositionSine.X+osc.Waveform.PositionSine.W) && y > int(osc.Waveform.PositionSine.Y) && y < int(osc.Waveform.PositionSine.Y+osc.Waveform.PositionSine.H) {
		osc.Waveform.Value = "sine"
		fmt.Println("osc.Waveform.Value : ", osc.Waveform.Value)
	}
	if x > int(osc.Waveform.PositionTriangle.X) && x < int(osc.Waveform.PositionTriangle.X+osc.Waveform.PositionTriangle.W) && y > int(osc.Waveform.PositionTriangle.Y) && y < int(osc.Waveform.PositionTriangle.Y+osc.Waveform.PositionTriangle.H) {
		osc.Waveform.Value = "triangle"
		fmt.Println("osc.Waveform.Value : ", osc.Waveform.Value)
	}
	if x > int(osc.Waveform.PositionSquare.X) && x < int(osc.Waveform.PositionSquare.X+osc.Waveform.PositionSquare.W) && y > int(osc.Waveform.PositionSquare.Y) && y < int(osc.Waveform.PositionSquare.Y+osc.Waveform.PositionSquare.H) {
		osc.Waveform.Value = "square"
		fmt.Println("osc.Waveform.Value : ", osc.Waveform.Value)
	}
	if x > int(osc.Waveform.PositionFlat.X) && x < int(osc.Waveform.PositionFlat.X+osc.Waveform.PositionFlat.W) && y > int(osc.Waveform.PositionFlat.Y) && y < int(osc.Waveform.PositionFlat.Y+osc.Waveform.PositionFlat.H) {
		osc.Waveform.Value = "flat"
		fmt.Println("osc.Waveform.Value : ", osc.Waveform.Value)
	}
	if x > int(osc.Kick.PositionUp.X) && x < int(osc.Kick.PositionUp.X+osc.Kick.PositionUp.W) && y > int(osc.Kick.PositionUp.Y) && y < int(osc.Kick.PositionUp.Y+osc.Kick.PositionUp.H) {
		souris.Click = "Kick+"
		souris.DateClick = time.Now()
		osc.Kick.Value += 0.1
		fmt.Println("oscs.Kick.Value up : ", osc.Kick.Value)
	}
	if x > int(osc.Kick.PositionDown.X) && x < int(osc.Kick.PositionDown.X+osc.Kick.PositionDown.W) && y > int(osc.Kick.PositionDown.Y) && y < int(osc.Kick.PositionDown.Y+osc.Kick.PositionDown.H) {
		souris.Click = "Kick-"
		souris.DateClick = time.Now()
		osc.Kick.Value -= 0.1
		fmt.Println("osc.Kick.Value down : ", osc.Kick.Value)
	}
	if x > int(osc.AsymetrieX.PositionUp.X) && x < int(osc.AsymetrieX.PositionUp.X+osc.AsymetrieX.PositionUp.W) && y > int(osc.AsymetrieX.PositionUp.Y) && y < int(osc.AsymetrieX.PositionUp.Y+osc.AsymetrieX.PositionUp.H) {
		souris.Click = "AsymetrieX+"
		souris.DateClick = time.Now()
		osc.AsymetrieX.Value += 0.1
		fmt.Println("osc.Asymetrie.Value up : ", osc.AsymetrieX.Value)
	}
	if x > int(osc.AsymetrieX.PositionDown.X) && x < int(osc.AsymetrieX.PositionDown.X+osc.AsymetrieX.PositionDown.W) && y > int(osc.AsymetrieX.PositionDown.Y) && y < int(osc.AsymetrieX.PositionDown.Y+osc.AsymetrieX.PositionDown.H) {
		souris.Click = "AsymetrieX-"
		souris.DateClick = time.Now()
		osc.AsymetrieX.Value -= 0.1
		fmt.Println("osc.Asymetrie.Value down : ", osc.AsymetrieX.Value)
	}
	if x > int(osc.AsymetrieY.PositionUp.X) && x < int(osc.AsymetrieY.PositionUp.X+osc.AsymetrieY.PositionUp.W) && y > int(osc.AsymetrieY.PositionUp.Y) && y < int(osc.AsymetrieY.PositionUp.Y+osc.AsymetrieY.PositionUp.H) {
		souris.Click = "AsymetrieY+"
		souris.DateClick = time.Now()
		osc.AsymetrieY.Value += 0.1
		fmt.Println("osc.Asymetrie.Value up : ", osc.AsymetrieY.Value)
	}
	if x > int(osc.AsymetrieY.PositionDown.X) && x < int(osc.AsymetrieY.PositionDown.X+osc.AsymetrieY.PositionDown.W) && y > int(osc.AsymetrieY.PositionDown.Y) && y < int(osc.AsymetrieY.PositionDown.Y+osc.AsymetrieY.PositionDown.H) {
		souris.Click = "AsymetrieY-"
		souris.DateClick = time.Now()
		osc.AsymetrieY.Value -= 0.1
		fmt.Println("osc.Asymetrie.Value down : ", osc.AsymetrieY.Value)
	}

	if x > int(clientInterface.Enregistrer.X) && x < int(clientInterface.Enregistrer.X+clientInterface.Enregistrer.W) && y > int(clientInterface.Enregistrer.Y) && y < int(clientInterface.Enregistrer.Y+clientInterface.Enregistrer.H) {
		fileName := "output.wav"
		osc.SaveToWav(fileName)
		//osc.play(fileName)
	}
	if x > int(clientInterface.Lire.X) && x < int(clientInterface.Lire.X+clientInterface.Lire.W) && y > int(clientInterface.Lire.Y) && y < int(clientInterface.Lire.Y+clientInterface.Lire.H) {
		fileName := "output.wav"
		osc.play(fileName)
	}
	if x > int(clientInterface.Stopper.X) && x < int(clientInterface.Stopper.X+clientInterface.Stopper.W) && y > int(clientInterface.Stopper.Y) && y < int(clientInterface.Stopper.Y+clientInterface.Stopper.H) {
		//osc.stop(fileName)
	}
}
