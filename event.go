package main

import (
	"fmt"
	"strings"
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
		if osc.End.Value == osc.SoundDuration.Value {
			osc.SoundDuration.Value += 0.1
			osc.End.Value += 0.1
		} else {
			osc.SoundDuration.Value += 0.1
		}
	} else if s.Click == "SoundDuration-" {
		if osc.SoundDuration.Value > 0.1 {
			osc.SoundDuration.Value -= 0.1
			if osc.End.Value > osc.SoundDuration.Value {
				osc.End.Value = osc.SoundDuration.Value
				if osc.Start.Value > osc.End.Value {
					osc.Start.Value = osc.End.Value
				}
			}
		}
	} else if s.Click == "Amplitude+" {
		osc.Amplitude.Value += 0.01
	} else if s.Click == "Amplitude-" {
		osc.Amplitude.Value -= 0.01
	} else if s.Click == "SampleRate+" {
		osc.MaxAmplitude.Value += 0.01
	} else if s.Click == "SampleRate-" {
		if osc.MaxAmplitude.Value > 0.01 {
			osc.MaxAmplitude.Value -= 0.01
		}
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
		osc.Hauteur.Value += 0.1
	} else if s.Click == "AsymetrieX-" {
		osc.Hauteur.Value -= 0.1
	} else if s.Click == "AsymetrieY+" {
		osc.AsymetrieY.Value += 0.1
	} else if s.Click == "AsymetrieY-" {
		osc.AsymetrieY.Value -= 0.1
	}
	if s.Click == "Lecture+" {
		osc.Lecture.Value += 0.1
	} else if s.Click == "Lecture-" && osc.Lecture.Value > 0 {
		osc.Lecture.Value -= 0.01
	} else if s.Click == "Pause+" {
		osc.Pause.Value += 0.01
	} else if s.Click == "Pause-" && osc.Pause.Value > 0 {
		osc.Pause.Value -= 0.01
	} else if s.Click == "Start+" && osc.Start.Value < osc.End.Value {
		osc.Start.Value += 0.1
	} else if s.Click == "Start-" && osc.Start.Value > 0 {
		osc.Start.Value -= 0.1
	} else if s.Click == "End+" && osc.End.Value < osc.SoundDuration.Value {
		osc.End.Value += 0.1
	} else if s.Click == "End-" && osc.End.Value > osc.Start.Value {
		osc.End.Value -= 0.1
	}
	if s.Click == "ColorR+" {
		if osc.Color.Value.R < 255 {
			osc.Color.Value.R += 1
		}
	} else if s.Click == "ColorR-" {
		if osc.Color.Value.R > 0 {
			osc.Color.Value.R -= 1
		}
	} else if s.Click == "ColorG+" {
		if osc.Color.Value.G < 255 {
			osc.Color.Value.G += 1
		}
	} else if s.Click == "ColorG-" {
		if osc.Color.Value.G > 0 {
			osc.Color.Value.G -= 1
		}
	} else if s.Click == "ColorB+" {
		if osc.Color.Value.B < 255 {
			osc.Color.Value.B += 1
		}
	} else if s.Click == "ColorB-" {
		if osc.Color.Value.B > 0 {
			osc.Color.Value.B -= 1
		}
	}

	// correction d'imprécision
	if osc.Amplitude.Value < 0.009 && osc.Amplitude.Value > -0.009 {
		osc.Amplitude.Value = 0.0
	}
	if osc.MaxAmplitude.Value < 0.009 && osc.MaxAmplitude.Value > -0.009 {
		osc.MaxAmplitude.Value = 0.0
	}
	if osc.Frequency.Value < 0.09 && osc.Frequency.Value > -0.09 {
		osc.Frequency.Value = 0.0
	}
	if osc.Phase.Value < 0.09 && osc.Phase.Value > -0.09 {
		osc.Phase.Value = 0.0
	}
	if osc.Kick.Value < 0.09 && osc.Kick.Value > -0.09 {
		osc.Kick.Value = 0.0
	}
	if osc.Hauteur.Value < 0.09 && osc.Hauteur.Value > -0.09 {
		osc.Hauteur.Value = 0.0
	}
	if osc.AsymetrieY.Value < 0.09 && osc.AsymetrieY.Value > -0.09 {
		osc.AsymetrieY.Value = 0.0
	}
	if osc.Lecture.Value < 0.009 && osc.Lecture.Value > -0.009 {
		osc.Lecture.Value = 0.0
	}
	if osc.Pause.Value < 0.009 && osc.Pause.Value > -0.009 {
		osc.Pause.Value = 0.0
	}
	if osc.Start.Value < 0.09 && osc.Start.Value > -0.09 {
		osc.Start.Value = 0.0
	}
	if osc.End.Value < 0.09 && osc.End.Value > -0.09 {
		osc.End.Value = 0.0
	}

	if strings.Contains(s.Click, "Color") {
		clientInterface.ReloadSelector = true
	} else if s.Click != "" {
		clientInterface.ReloadWave = true
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
				mouseDownUpdateOscillatorValues(int(event.Button), int(event.X), int(event.Y))
				mouseDownSelector(int(event.Button), int(event.X), int(event.Y))
			} else if event.State == sdl.RELEASED {
				mouseUpUpdate(int(event.Button))
			}
		case *sdl.TextInputEvent:
			text := event.GetText()
			fmt.Println(text)
			if clientInterface.FileName.Value {
				clientInterface.FileName.Champ += text
			}
		case *sdl.KeyboardEvent:
			// touche effacer
			if clientInterface.FileName.Value && event.Keysym.Scancode == sdl.SCANCODE_BACKSPACE {
				if len(clientInterface.FileName.Champ) > 0 {
					clientInterface.FileName.Champ = clientInterface.FileName.Champ[:len(clientInterface.FileName.Champ)-1]
				}
				// touche echap
			} else if clientInterface.FileName.Value && event.Keysym.Scancode == sdl.SCANCODE_ESCAPE {
				clientInterface.FileName.Champ = ""
				clientInterface.FileName.Value = false
				// touche entrée
			} else if clientInterface.FileName.Value && event.Keysym.Scancode == sdl.SCANCODE_RETURN {
				if len(clientInterface.FileName.Champ) > 0 {
					if clientInterface.FileName.Cible == "All" {
						clientInterface.saveAll(clientInterface.FileName.Champ)
						clientInterface.FileName.Champ = ""
						clientInterface.FileName.Value = false
					} else if clientInterface.FileName.Cible == "Osc" {
						osc.save("res/user/osc/" + clientInterface.FileName.Champ + ".json")
						clientInterface.FileName.Champ = ""
						clientInterface.FileName.Value = false
					} else if clientInterface.FileName.Cible == "Open" {
						clientInterface.loadMix(clientInterface.FileName.Champ)
						clientInterface.FileName.Champ = ""
						clientInterface.FileName.Value = false
					}
				}
			} else if event.Keysym.Scancode == sdl.SCANCODE_ESCAPE || event.Keysym.Scancode == sdl.SCANCODE_F12 {
				running = false
				return
			}
			//stringEvent := sdl.GetKeyName(sdl.Keycode(event.Keysym.Scancode))
			//fmt.Println(stringEvent)
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

func mouseDownSelector(button, x, y int) {
	if button != 1 {
		return
	}
	if x > int(clientInterface.DeroulantSelector.PositionUp.X) && x < int(clientInterface.DeroulantSelector.PositionUp.X+clientInterface.DeroulantSelector.PositionUp.W) && y > int(clientInterface.DeroulantSelector.PositionUp.Y) && y < int(clientInterface.DeroulantSelector.PositionUp.Y+clientInterface.DeroulantSelector.PositionUp.H) {
		clientInterface.DeroulantSelector.Value = !clientInterface.DeroulantSelector.Value
		fmt.Println("DeroulantSelector: ", clientInterface.DeroulantSelector.Value)
	}
	if clientInterface.DeroulantSelector.Value {
		for index, sel := range clientInterface.Selector {
			if x > int(sel.Position.X) && x < int(sel.Position.X+sel.Position.W) && y > int(sel.Position.Y) && y < int(sel.Position.Y+sel.Position.H) {
				if sel.Name == "Add Osc" {
					if len(clientInterface.Oscillator) == len(clientInterface.Wave) && len(clientInterface.Oscillator) == len(clientInterface.Selector)-2 {
						clientInterface.AddOscillator = true
					}
				} else if sel.Name == "Del Osc" {
					if len(clientInterface.Oscillator) <= 1 {
						return
					}
					clientInterface.DeleteOscillator = true
				} else {
					clientInterface.CurrentOscillator = index - 2
					osc = clientInterface.Oscillator[clientInterface.CurrentOscillator]
					fmt.Println("CurrentOscillator : ", clientInterface.CurrentOscillator)
				}
				fmt.Println("Selector: ", sel.Name, "  ", len(clientInterface.Oscillator))
			}
		}
	}
}

func mouseDownUpdateOscillatorValues(button, x, y int) {
	if button != 1 {
		return
	}

	// configurables bool fields -------------------------------------------------------------------------------------------------------------------------------

	if x > int(osc.OnlyPositive.PositionUp.X) && x < int(osc.OnlyPositive.PositionUp.X+osc.OnlyPositive.PositionUp.W) && y > int(osc.OnlyPositive.PositionUp.Y) && y < int(osc.OnlyPositive.PositionUp.Y+osc.OnlyPositive.PositionUp.H) {
		souris.Click = "OnlyPositive"
		osc.OnlyPositive.Value = !osc.OnlyPositive.Value
		fmt.Println("osc.OnlyPositive.Value : ", osc.OnlyPositive.Value)
	}
	if x > int(osc.OnlyPositive.PositionDown.X) && x < int(osc.OnlyPositive.PositionDown.X+osc.OnlyPositive.PositionDown.W) && y > int(osc.OnlyPositive.PositionDown.Y) && y < int(osc.OnlyPositive.PositionDown.Y+osc.OnlyPositive.PositionDown.H) {
		souris.Click = "OnlyPositive"
		osc.OnlyPositive.Value = !osc.OnlyPositive.Value
		fmt.Println("osc.OnlyPositive.Value : ", osc.OnlyPositive.Value)
	}
	if x > int(osc.OnlyNegative.PositionUp.X) && x < int(osc.OnlyNegative.PositionUp.X+osc.OnlyNegative.PositionUp.W) && y > int(osc.OnlyNegative.PositionUp.Y) && y < int(osc.OnlyNegative.PositionUp.Y+osc.OnlyNegative.PositionUp.H) {
		souris.Click = "OnlyNegative"
		osc.OnlyNegative.Value = !osc.OnlyNegative.Value
		fmt.Println("osc.OnlyNegative.Value : ", osc.OnlyNegative.Value)
	}
	if x > int(osc.OnlyNegative.PositionDown.X) && x < int(osc.OnlyNegative.PositionDown.X+osc.OnlyNegative.PositionDown.W) && y > int(osc.OnlyNegative.PositionDown.Y) && y < int(osc.OnlyNegative.PositionDown.Y+osc.OnlyNegative.PositionDown.H) {
		souris.Click = "OnlyNegative"
		osc.OnlyNegative.Value = !osc.OnlyNegative.Value
		fmt.Println("osc.OnlyNegative.Value : ", osc.OnlyNegative.Value)
	}
	if x > int(osc.OnlyIncreased.PositionUp.X) && x < int(osc.OnlyIncreased.PositionUp.X+osc.OnlyIncreased.PositionUp.W) && y > int(osc.OnlyIncreased.PositionUp.Y) && y < int(osc.OnlyIncreased.PositionUp.Y+osc.OnlyIncreased.PositionUp.H) {
		souris.Click = "OnlyIncreased"
		osc.OnlyIncreased.Value = !osc.OnlyIncreased.Value
		fmt.Println("osc.OnlyIncreased.Value : ", osc.OnlyIncreased.Value)
	}
	if x > int(osc.OnlyIncreased.PositionDown.X) && x < int(osc.OnlyIncreased.PositionDown.X+osc.OnlyIncreased.PositionDown.W) && y > int(osc.OnlyIncreased.PositionDown.Y) && y < int(osc.OnlyIncreased.PositionDown.Y+osc.OnlyIncreased.PositionDown.H) {
		souris.Click = "OnlyIncreased"
		osc.OnlyIncreased.Value = !osc.OnlyIncreased.Value
		fmt.Println("osc.OnlyIncreased.Value : ", osc.OnlyIncreased.Value)
	}
	if x > int(osc.OnlyDecreased.PositionUp.X) && x < int(osc.OnlyDecreased.PositionUp.X+osc.OnlyDecreased.PositionUp.W) && y > int(osc.OnlyDecreased.PositionUp.Y) && y < int(osc.OnlyDecreased.PositionUp.Y+osc.OnlyDecreased.PositionUp.H) {
		souris.Click = "OnlyDecreased"
		osc.OnlyDecreased.Value = !osc.OnlyDecreased.Value
		fmt.Println("osc.OnlyDecreased.Value : ", osc.OnlyDecreased.Value)
	}
	if x > int(osc.OnlyDecreased.PositionDown.X) && x < int(osc.OnlyDecreased.PositionDown.X+osc.OnlyDecreased.PositionDown.W) && y > int(osc.OnlyDecreased.PositionDown.Y) && y < int(osc.OnlyDecreased.PositionDown.Y+osc.OnlyDecreased.PositionDown.H) {
		souris.Click = "OnlyDecreased"
		osc.OnlyDecreased.Value = !osc.OnlyDecreased.Value
		fmt.Println("osc.OnlyDecreased.Value : ", osc.OnlyDecreased.Value)
	}

	// configurables fields colunm 1 ---------------------------------------------------------------------------------------------------------------------------

	if x > int(osc.SoundDuration.PositionUp.X) && x < int(osc.SoundDuration.PositionUp.X+osc.SoundDuration.PositionUp.W) && y > int(osc.SoundDuration.PositionUp.Y) && y < int(osc.SoundDuration.PositionUp.Y+osc.SoundDuration.PositionUp.H) {
		souris.Click = "SoundDuration+"
		souris.DateClick = time.Now()
		if osc.End.Value == osc.SoundDuration.Value {
			osc.End.Value += 0.1
			osc.SoundDuration.Value += 0.1
		} else {
			osc.SoundDuration.Value += 0.1
		}
		fmt.Println("osc.SoundDuration.Value up : ", osc.SoundDuration.Value)
	}
	if x > int(osc.SoundDuration.PositionDown.X) && x < int(osc.SoundDuration.PositionDown.X+osc.SoundDuration.PositionDown.W) && y > int(osc.SoundDuration.PositionDown.Y) && y < int(osc.SoundDuration.PositionDown.Y+osc.SoundDuration.PositionDown.H) {
		souris.Click = "SoundDuration-"
		souris.DateClick = time.Now()
		osc.SoundDuration.Value -= 0.1
		if osc.End.Value > osc.SoundDuration.Value {
			osc.End.Value = osc.SoundDuration.Value
			if osc.Start.Value > osc.End.Value {
				osc.Start.Value = osc.End.Value
			}
		}
		fmt.Println("osc.SoundDuration.Value down : ", osc.SoundDuration.Value)
	}
	if x > int(osc.Amplitude.PositionUp.X) && x < int(osc.Amplitude.PositionUp.X+osc.Amplitude.PositionUp.W) && y > int(osc.Amplitude.PositionUp.Y) && y < int(osc.Amplitude.PositionUp.Y+osc.Amplitude.PositionUp.H) {
		souris.Click = "Amplitude+"
		souris.DateClick = time.Now()
		osc.Amplitude.Value += 0.01
		fmt.Println("osc.Amplitude.Value up : ", osc.Amplitude.Value)
	}
	if x > int(osc.Amplitude.PositionDown.X) && x < int(osc.Amplitude.PositionDown.X+osc.Amplitude.PositionDown.W) && y > int(osc.Amplitude.PositionDown.Y) && y < int(osc.Amplitude.PositionDown.Y+osc.Amplitude.PositionDown.H) {
		souris.Click = "Amplitude-"
		souris.DateClick = time.Now()
		osc.Amplitude.Value -= 0.01
		fmt.Println("osc.Amplitude.Value down : ", osc.Amplitude.Value)
	}
	if x > int(osc.MaxAmplitude.PositionUp.X) && x < int(osc.MaxAmplitude.PositionUp.X+osc.MaxAmplitude.PositionUp.W) && y > int(osc.MaxAmplitude.PositionUp.Y) && y < int(osc.MaxAmplitude.PositionUp.Y+osc.MaxAmplitude.PositionUp.H) {
		souris.Click = "SampleRate+"
		souris.DateClick = time.Now()
		osc.MaxAmplitude.Value += 0.01
		fmt.Println("osc.SampleRate.Value up : ", osc.MaxAmplitude.Value)
	}
	if x > int(osc.MaxAmplitude.PositionDown.X) && x < int(osc.MaxAmplitude.PositionDown.X+osc.MaxAmplitude.PositionDown.W) && y > int(osc.MaxAmplitude.PositionDown.Y) && y < int(osc.MaxAmplitude.PositionDown.Y+osc.MaxAmplitude.PositionDown.H) {
		souris.Click = "SampleRate-"
		souris.DateClick = time.Now()
		osc.MaxAmplitude.Value -= 0.01
		fmt.Println("osc.SampleRate.Value down : ", osc.MaxAmplitude.Value)
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

	// configurables fields column 2 ---------------------------------------------------------------------------------------------------------------------------

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
	if x > int(osc.Hauteur.PositionUp.X) && x < int(osc.Hauteur.PositionUp.X+osc.Hauteur.PositionUp.W) && y > int(osc.Hauteur.PositionUp.Y) && y < int(osc.Hauteur.PositionUp.Y+osc.Hauteur.PositionUp.H) {
		souris.Click = "AsymetrieX+"
		souris.DateClick = time.Now()
		osc.Hauteur.Value += 0.1
		fmt.Println("osc.Asymetrie.Value up : ", osc.Hauteur.Value)
	}
	if x > int(osc.Hauteur.PositionDown.X) && x < int(osc.Hauteur.PositionDown.X+osc.Hauteur.PositionDown.W) && y > int(osc.Hauteur.PositionDown.Y) && y < int(osc.Hauteur.PositionDown.Y+osc.Hauteur.PositionDown.H) {
		souris.Click = "AsymetrieX-"
		souris.DateClick = time.Now()
		osc.Hauteur.Value -= 0.1
		fmt.Println("osc.Asymetrie.Value down : ", osc.Hauteur.Value)
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

	// configurables fields column 3--------------------------------------------------------------------------------------------------------------------------

	if x > int(osc.Lecture.PositionUp.X) && x < int(osc.Lecture.PositionUp.X+osc.Lecture.PositionUp.W) && y > int(osc.Lecture.PositionUp.Y) && y < int(osc.Lecture.PositionUp.Y+osc.Lecture.PositionUp.H) {
		souris.Click = "Lecture+"
		souris.DateClick = time.Now()
		osc.Lecture.Value += 0.01
		fmt.Println("osc.Lecture.Value up : ", osc.Lecture.Value)
	}
	if x > int(osc.Lecture.PositionDown.X) && x < int(osc.Lecture.PositionDown.X+osc.Lecture.PositionDown.W) && y > int(osc.Lecture.PositionDown.Y) && y < int(osc.Lecture.PositionDown.Y+osc.Lecture.PositionDown.H) {
		if osc.Lecture.Value > 0 {
			souris.Click = "Lecture-"
			souris.DateClick = time.Now()
			osc.Lecture.Value -= 0.01
			fmt.Println("osc.Lecture.Value down : ", osc.Lecture.Value)
		}
	}
	if x > int(osc.Pause.PositionUp.X) && x < int(osc.Pause.PositionUp.X+osc.Pause.PositionUp.W) && y > int(osc.Pause.PositionUp.Y) && y < int(osc.Pause.PositionUp.Y+osc.Pause.PositionUp.H) {
		souris.Click = "Pause+"
		souris.DateClick = time.Now()
		osc.Pause.Value += 0.01
		fmt.Println("osc.Pause.Value up : ", osc.Pause.Value)
	}
	if x > int(osc.Pause.PositionDown.X) && x < int(osc.Pause.PositionDown.X+osc.Pause.PositionDown.W) && y > int(osc.Pause.PositionDown.Y) && y < int(osc.Pause.PositionDown.Y+osc.Pause.PositionDown.H) {
		if osc.Pause.Value > 0 {
			souris.Click = "Pause-"
			souris.DateClick = time.Now()
			osc.Pause.Value -= 0.01
			fmt.Println("osc.Pause.Value down : ", osc.Pause.Value)
		}
	}
	if x > int(osc.Start.PositionUp.X) && x < int(osc.Start.PositionUp.X+osc.Start.PositionUp.W) && y > int(osc.Start.PositionUp.Y) && y < int(osc.Start.PositionUp.Y+osc.Start.PositionUp.H) {
		if osc.Start.Value < osc.End.Value {
			souris.Click = "Start+"
			souris.DateClick = time.Now()
			osc.Start.Value += 0.01
			fmt.Println("osc.Start.Value up : ", osc.Start.Value)
		}
	}
	if x > int(osc.Start.PositionDown.X) && x < int(osc.Start.PositionDown.X+osc.Start.PositionDown.W) && y > int(osc.Start.PositionDown.Y) && y < int(osc.Start.PositionDown.Y+osc.Start.PositionDown.H) {
		if osc.Start.Value > 0 {
			souris.Click = "Start-"
			souris.DateClick = time.Now()
			osc.Start.Value -= 0.01
			fmt.Println("osc.Start.Value down : ", osc.Start.Value)
		}
	}
	if x > int(osc.End.PositionUp.X) && x < int(osc.End.PositionUp.X+osc.End.PositionUp.W) && y > int(osc.End.PositionUp.Y) && y < int(osc.End.PositionUp.Y+osc.End.PositionUp.H) {
		if osc.End.Value < osc.SoundDuration.Value {
			souris.Click = "End+"
			souris.DateClick = time.Now()
			osc.End.Value += 0.1
			fmt.Println("osc.End.Value up : ", osc.End.Value)
		}
	}
	if x > int(osc.End.PositionDown.X) && x < int(osc.End.PositionDown.X+osc.End.PositionDown.W) && y > int(osc.End.PositionDown.Y) && y < int(osc.End.PositionDown.Y+osc.End.PositionDown.H) {
		if osc.End.Value > osc.Start.Value {
			souris.Click = "End-"
			souris.DateClick = time.Now()
			osc.End.Value -= 0.1
			fmt.Println("osc.End.Value down : ", osc.End.Value)
		}
	}

	// configurables fields column 4--------------------------------------------------------------------------------------------------------------------------

	if x > int(osc.Color.PositionUpR.X) && x < int(osc.Color.PositionUpR.X+osc.Color.PositionUpR.W) && y > int(osc.Color.PositionUpR.Y) && y < int(osc.Color.PositionUpR.Y+osc.Color.PositionUpR.H) {
		if osc.Color.Value.R < 255 {
			souris.Click = "ColorR+"
			souris.DateClick = time.Now()
			osc.Color.Value.R += 1
			fmt.Println("osc.Color.Value.R up : ", osc.Color.Value.R)
		}
	}
	if x > int(osc.Color.PositionDownR.X) && x < int(osc.Color.PositionDownR.X+osc.Color.PositionDownR.W) && y > int(osc.Color.PositionDownR.Y) && y < int(osc.Color.PositionDownR.Y+osc.Color.PositionDownR.H) {
		if osc.Color.Value.R > 0 {
			souris.Click = "ColorR-"
			souris.DateClick = time.Now()
			osc.Color.Value.R -= 1
			fmt.Println("osc.Color.Value.R down : ", osc.Color.Value.R)
		}
	}
	if x > int(osc.Color.PositionUpG.X) && x < int(osc.Color.PositionUpG.X+osc.Color.PositionUpG.W) && y > int(osc.Color.PositionUpG.Y) && y < int(osc.Color.PositionUpG.Y+osc.Color.PositionUpG.H) {
		if osc.Color.Value.G < 255 {
			souris.Click = "ColorG+"
			souris.DateClick = time.Now()
			osc.Color.Value.G += 1
			fmt.Println("osc.Color.Value.G up : ", osc.Color.Value.G)
		}
	}
	if x > int(osc.Color.PositionDownG.X) && x < int(osc.Color.PositionDownG.X+osc.Color.PositionDownG.W) && y > int(osc.Color.PositionDownG.Y) && y < int(osc.Color.PositionDownG.Y+osc.Color.PositionDownG.H) {
		if osc.Color.Value.G > 0 {
			souris.Click = "ColorG-"
			souris.DateClick = time.Now()
			osc.Color.Value.G -= 1
			fmt.Println("osc.Color.Value.G down : ", osc.Color.Value.G)
		}
	}
	if x > int(osc.Color.PositionUpB.X) && x < int(osc.Color.PositionUpB.X+osc.Color.PositionUpB.W) && y > int(osc.Color.PositionUpB.Y) && y < int(osc.Color.PositionUpB.Y+osc.Color.PositionUpB.H) {
		if osc.Color.Value.B < 255 {
			souris.Click = "ColorB+"
			souris.DateClick = time.Now()
			osc.Color.Value.B += 1
			fmt.Println("osc.Color.Value.B up : ", osc.Color.Value.B)
		}
	}
	if x > int(osc.Color.PositionDownB.X) && x < int(osc.Color.PositionDownB.X+osc.Color.PositionDownB.W) && y > int(osc.Color.PositionDownB.Y) && y < int(osc.Color.PositionDownB.Y+osc.Color.PositionDownB.H) {
		if osc.Color.Value.B > 0 {
			souris.Click = "ColorB-"
			souris.DateClick = time.Now()
			osc.Color.Value.B -= 1
			fmt.Println("osc.Color.Value.B down : ", osc.Color.Value.B)
		}
	}

	// client interface--------------------------------------------------------------------------------------------------------------------------

	if x > int(clientInterface.CloseFileName.Position.X) && x < int(clientInterface.CloseFileName.Position.X+clientInterface.CloseFileName.Position.W) && y > int(clientInterface.CloseFileName.Position.Y) && y < int(clientInterface.CloseFileName.Position.Y+clientInterface.CloseFileName.Position.H) {
		clientInterface.FileName.Champ = ""
		clientInterface.FileName.Value = false
	}
	if x > int(clientInterface.OpenMix.X) && x < int(clientInterface.OpenMix.X+clientInterface.OpenMix.W) && y > int(clientInterface.OpenMix.Y) && y < int(clientInterface.OpenMix.Y+clientInterface.OpenMix.H) {
		if clientInterface.FileName.Value {
			if clientInterface.FileName.Cible == "Open" && clientInterface.FileName.Champ != "" {
				clientInterface.loadMix(clientInterface.FileName.Champ) //"osc/" + clientInterface.FileName.Champ + ".json")
				clientInterface.FileName.Value = false
			}
			clientInterface.FileName.Cible = "Open"
			clientInterface.FileName.Champ = ""
			return
		}
		clientInterface.FileName.Value = true
		clientInterface.FileName.Cible = "Open"
		clientInterface.FileName.Champ = ""
	}
	if x > int(clientInterface.Enregistrer.X) && x < int(clientInterface.Enregistrer.X+clientInterface.Enregistrer.W) && y > int(clientInterface.Enregistrer.Y) && y < int(clientInterface.Enregistrer.Y+clientInterface.Enregistrer.H) {
		if clientInterface.FileName.Value {
			if clientInterface.FileName.Cible == "Osc" && clientInterface.FileName.Champ != "" {
				osc.save("osc/" + clientInterface.FileName.Champ + ".json")
				clientInterface.FileName.Value = false
			}
			clientInterface.FileName.Cible = "Osc"
			clientInterface.FileName.Champ = ""
			return
		}
		clientInterface.FileName.Value = true
		clientInterface.FileName.Cible = "Osc"
		clientInterface.FileName.Champ = ""
	}
	if x > int(clientInterface.SaveAll.X) && x < int(clientInterface.SaveAll.X+clientInterface.SaveAll.W) && y > int(clientInterface.SaveAll.Y) && y < int(clientInterface.SaveAll.Y+clientInterface.SaveAll.H) {
		if clientInterface.FileName.Value {
			if clientInterface.FileName.Cible == "All" && clientInterface.FileName.Champ != "" {
				clientInterface.saveAll(clientInterface.FileName.Champ)
				clientInterface.FileName.Value = false
			}
			clientInterface.FileName.Cible = "All"
			clientInterface.FileName.Champ = ""
			return
		}
		clientInterface.FileName.Value = true
		clientInterface.FileName.Cible = "All"
		clientInterface.FileName.Champ = ""
	}
	if x > int(clientInterface.Lire.X) && x < int(clientInterface.Lire.X+clientInterface.Lire.W) && y > int(clientInterface.Lire.Y) && y < int(clientInterface.Lire.Y+clientInterface.Lire.H) {
		clientInterface.SaveToWav("res/temp/output.wav", "save")
		/*
			fileName := "output.wav"
			osc.stop(fileName)
		*/
	}
	/*
		if x > int(clientInterface.Stopper.X) && x < int(clientInterface.Stopper.X+clientInterface.Stopper.W) && y > int(clientInterface.Stopper.Y) && y < int(clientInterface.Stopper.Y+clientInterface.Stopper.H) {
			//osc.stop(fileName)
		}
	*/

	if strings.Contains(souris.Click, "Color") {
		clientInterface.ReloadSelector = true
	} else if souris.Click != "" {
		clientInterface.ReloadWave = true
	}
}
