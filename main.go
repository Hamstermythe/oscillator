package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/tfriedel6/canvas/sdlcanvas"
	"github.com/veandco/go-sdl2/sdl"
)

var osc = &Oscillator{}
var clientInterface = &ClientInterface{}
var echelle float64
var wndWidth, wndHeight int
var souris = Souris{}
var running = true

func main() {
	// Initialiser SDL2 pour la vi
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}
	defer sdl.Quit()
	// Obtention de la taille de l'écran
	displayIndex := 0 // 0 pour l'écran principal, 1 pour le secondaire, etc.
	displayBounds, err := sdl.GetDisplayBounds(displayIndex)
	if err != nil {
		panic(err)
	}
	screenWidth := int(displayBounds.W)
	scfreenHeight := int(displayBounds.H)
	setWindow(screenWidth, scfreenHeight)
	print(wndWidth, "\n", wndHeight, "\n")
	wnd, cv, err := sdlcanvas.CreateWindow(wndWidth, wndHeight, "Versus")
	if err != nil {
		panic(err)
	}
	wnd.Window.SetBordered(false)
	wnd.Window.SetResizable(false)
	width, height := wnd.Window.GetSize()
	fmt.Println("width: ", width, "height: ", height)
	wnd = addEvent(wnd)

	clientInterface.InitInterface(cv)

	wnd.MainLoop(func() {
		if clientInterface.Style.FontSize != 0 {
			fmt.Println("clientInterface.Style.FontSize: ", clientInterface.Style.FontSize)
		}
		souris.Action()
		// Effacer l'écran avec une couleur (noir)
		cv.SetFillStyle(0, 0, 0, 255)
		cv.FillRect(0, 0, float64(wndWidth), float64(wndHeight))

		// configurables fields column 1-----------------------------------------------------------------------------------------------------------------------
		margeVertical := ((float64(osc.Amplitude.PositionDown.H) - clientInterface.Style.FontSize) / 2) + clientInterface.Style.FontSize
		margeHorizontalePlus := (float64(osc.Amplitude.PositionUp.W) - cv.MeasureText("+").Width) / 2
		margeHorizontaleMoins := (float64(osc.Amplitude.PositionDown.W) - cv.MeasureText("-").Width) / 2
		text := "Only Positive"
		marginHorizontal := cv.MeasureText(text).Width + 10
		cv.SetFillStyle(255, 255, 255, 255)
		cv.FillText(
			text,
			float64(osc.OnlyPositive.PositionUp.X)-marginHorizontal,
			float64(osc.OnlyPositive.PositionUp.Y)+float64(osc.OnlyPositive.PositionUp.H)+(clientInterface.Style.FontSize/2),
		) //-clientInterface.Style.FontSize)
		if osc.OnlyPositive.Value {
			cv.SetFillStyle(0, 255, 0, 255)
		} else {
			cv.SetFillStyle(0, 50, 0, 255)
		}
		cv.FillRect(
			float64(osc.OnlyPositive.PositionUp.X),
			float64(osc.OnlyPositive.PositionUp.Y),
			float64(osc.OnlyPositive.PositionUp.W),
			float64(osc.OnlyPositive.PositionUp.H),
		)
		if osc.OnlyPositive.Value {
			cv.SetFillStyle(255, 255, 255, 255)
		} else {
			cv.SetFillStyle(100, 130, 100, 255)
		}
		cv.FillText(
			"on",
			float64(osc.OnlyPositive.PositionUp.X)+((float64(osc.OnlyPositive.PositionUp.W)-cv.MeasureText("on").Width)/2),
			float64(osc.OnlyPositive.PositionUp.Y)+clientInterface.Style.FontSize,
		)
		if !osc.OnlyPositive.Value {
			cv.SetFillStyle(255, 0, 0, 255)
		} else {
			cv.SetFillStyle(50, 0, 0, 255)
		}
		cv.FillRect(
			float64(osc.OnlyPositive.PositionDown.X),
			float64(osc.OnlyPositive.PositionDown.Y),
			float64(osc.OnlyPositive.PositionDown.W),
			float64(osc.OnlyPositive.PositionDown.H),
		)
		if !osc.OnlyPositive.Value {
			cv.SetFillStyle(255, 255, 255, 255)
		} else {
			cv.SetFillStyle(130, 100, 100, 100)
		}
		cv.FillText(
			"off",
			float64(osc.OnlyPositive.PositionDown.X)+((float64(osc.OnlyPositive.PositionUp.W)-cv.MeasureText("off").Width)/2),
			float64(osc.OnlyPositive.PositionDown.Y)+clientInterface.Style.FontSize,
		)

		text = "Sound Duration   " + strconv.FormatFloat(osc.SoundDuration.Value, 'f', 2, 64)
		cv.SetFillStyle(255, 255, 255, 255)
		cv.FillText(
			text,
			float64(osc.SoundDuration.PositionDown.X),
			float64(osc.SoundDuration.PositionDown.Y)-clientInterface.Style.FontSize,
		)
		cv.SetFillStyle(0, 255, 0, 255)
		cv.FillRect(
			float64(osc.SoundDuration.PositionUp.X),
			float64(osc.SoundDuration.PositionUp.Y),
			float64(osc.SoundDuration.PositionUp.W),
			float64(osc.SoundDuration.PositionUp.H),
		)
		cv.SetFillStyle(255, 255, 255, 255)
		cv.FillText(
			"+",
			float64(osc.SoundDuration.PositionUp.X)+margeHorizontalePlus,
			float64(osc.SoundDuration.PositionUp.Y)+margeVertical,
		)
		cv.SetFillStyle(255, 0, 0, 255)
		cv.FillRect(
			float64(osc.SoundDuration.PositionDown.X),
			float64(osc.SoundDuration.PositionDown.Y),
			float64(osc.SoundDuration.PositionDown.W),
			float64(osc.SoundDuration.PositionDown.H),
		)
		cv.SetFillStyle(255, 255, 255, 255)
		cv.FillText(
			"-",
			float64(osc.SoundDuration.PositionDown.X)+margeHorizontaleMoins,
			float64(osc.SoundDuration.PositionDown.Y)+margeVertical,
		)

		text = "Amplitude     " + strconv.FormatFloat(osc.Amplitude.Value, 'f', 2, 64)
		cv.SetFillStyle(255, 255, 255, 255)
		cv.FillText(
			text,
			float64(osc.Amplitude.PositionDown.X),
			float64(osc.Amplitude.PositionDown.Y)-clientInterface.Style.FontSize,
		)
		cv.SetFillStyle(0, 255, 0, 255)
		cv.FillRect(
			float64(osc.Amplitude.PositionUp.X),
			float64(osc.Amplitude.PositionUp.Y),
			float64(osc.Amplitude.PositionUp.W),
			float64(osc.Amplitude.PositionUp.H),
		)
		cv.SetFillStyle(255, 255, 255, 255)
		cv.FillText(
			"+",
			float64(osc.Amplitude.PositionUp.X)+margeHorizontalePlus,
			float64(osc.Amplitude.PositionUp.Y)+margeVertical,
		)
		cv.SetFillStyle(255, 0, 0, 255)
		cv.FillRect(
			float64(osc.Amplitude.PositionDown.X),
			float64(osc.Amplitude.PositionDown.Y),
			float64(osc.Amplitude.PositionDown.W),
			float64(osc.Amplitude.PositionDown.H),
		)
		cv.SetFillStyle(255, 255, 255, 255)
		cv.FillText(
			"-",
			float64(osc.Amplitude.PositionDown.X)+margeHorizontaleMoins,
			float64(osc.Amplitude.PositionDown.Y)+margeVertical,
		)

		text = "Sample Rate   " + strconv.Itoa(osc.SampleRate.Value)
		cv.SetFillStyle(255, 255, 255, 255)
		cv.FillText(
			text,
			float64(osc.SampleRate.PositionDown.X),
			float64(osc.SampleRate.PositionDown.Y)-clientInterface.Style.FontSize,
		)
		cv.SetFillStyle(0, 255, 0, 255)
		cv.FillRect(
			float64(osc.SampleRate.PositionUp.X),
			float64(osc.SampleRate.PositionUp.Y),
			float64(osc.SampleRate.PositionUp.W),
			float64(osc.SampleRate.PositionUp.H),
		)
		cv.SetFillStyle(255, 255, 255, 255)
		cv.FillText(
			"+",
			float64(osc.SampleRate.PositionUp.X)+margeHorizontalePlus,
			float64(osc.SampleRate.PositionUp.Y)+margeVertical,
		)
		cv.SetFillStyle(255, 0, 0, 255)
		cv.FillRect(
			float64(osc.SampleRate.PositionDown.X),
			float64(osc.SampleRate.PositionDown.Y),
			float64(osc.SampleRate.PositionDown.W),
			float64(osc.SampleRate.PositionDown.H),
		)
		cv.SetFillStyle(255, 255, 255, 255)
		cv.FillText(
			"-",
			float64(osc.SampleRate.PositionDown.X)+margeHorizontaleMoins,
			float64(osc.SampleRate.PositionDown.Y)+margeVertical,
		)

		text = "Frequency     " + strconv.FormatFloat(osc.Frequency.Value, 'f', 2, 64)
		cv.SetFillStyle(255, 255, 255, 255)
		cv.FillText(
			text,
			float64(osc.Frequency.PositionDown.X),
			float64(osc.Frequency.PositionDown.Y)-clientInterface.Style.FontSize,
		)
		cv.SetFillStyle(0, 255, 0, 255)
		cv.FillRect(
			float64(osc.Frequency.PositionUp.X),
			float64(osc.Frequency.PositionUp.Y),
			float64(osc.Frequency.PositionUp.W),
			float64(osc.Frequency.PositionUp.H),
		)
		cv.SetFillStyle(255, 255, 255, 255)
		cv.FillText(
			"+",
			float64(osc.Frequency.PositionUp.X)+margeHorizontalePlus,
			float64(osc.Frequency.PositionUp.Y)+margeVertical,
		)
		cv.SetFillStyle(255, 0, 0, 255)
		cv.FillRect(
			float64(osc.Frequency.PositionDown.X),
			float64(osc.Frequency.PositionDown.Y),
			float64(osc.Frequency.PositionDown.W),
			float64(osc.Frequency.PositionDown.H),
		)
		cv.SetFillStyle(255, 255, 255, 255)
		cv.FillText(
			"-",
			float64(osc.Frequency.PositionDown.X)+margeHorizontaleMoins,
			float64(osc.Frequency.PositionDown.Y)+margeVertical,
		)

		text = "Phase         " + strconv.FormatFloat(osc.Phase.Value, 'f', 2, 64)
		cv.SetFillStyle(255, 255, 255, 255)
		cv.FillText(
			text,
			float64(osc.Phase.PositionDown.X),
			float64(osc.Phase.PositionDown.Y)-clientInterface.Style.FontSize,
		)
		cv.SetFillStyle(0, 255, 0, 255)
		cv.FillRect(
			float64(osc.Phase.PositionUp.X),
			float64(osc.Phase.PositionUp.Y),
			float64(osc.Phase.PositionUp.W),
			float64(osc.Phase.PositionUp.H),
		)
		cv.SetFillStyle(255, 255, 255, 255)
		cv.FillText(
			"+",
			float64(osc.Phase.PositionUp.X)+margeHorizontalePlus,
			float64(osc.Phase.PositionUp.Y)+margeVertical,
		)
		cv.SetFillStyle(255, 0, 0, 255)
		cv.FillRect(
			float64(osc.Phase.PositionDown.X),
			float64(osc.Phase.PositionDown.Y),
			float64(osc.Phase.PositionDown.W),
			float64(osc.Phase.PositionDown.H),
		)
		cv.SetFillStyle(255, 255, 255, 255)
		cv.FillText(
			"-",
			float64(osc.Phase.PositionDown.X)+margeHorizontaleMoins,
			float64(osc.Phase.PositionDown.Y)+margeVertical,
		)

		// configurables fields column 2------------------------------------------------------------------------------------------------------------------------
		text = "Waveform"
		cv.SetFillStyle(255, 255, 255, 255)
		cv.FillText(
			text,
			float64(osc.Waveform.PositionSine.X)+((float64(osc.Waveform.PositionSine.W*2)-cv.MeasureText(text).Width)/2),
			float64(osc.Waveform.PositionSine.Y)-clientInterface.Style.FontSize,
		)
		if osc.Waveform.Value == "sine" {
			cv.SetFillStyle(0, 255, 0, 255)
		} else {
			cv.SetFillStyle(255, 0, 0, 255)
		}
		cv.FillRect(
			float64(osc.Waveform.PositionSine.X),
			float64(osc.Waveform.PositionSine.Y),
			float64(osc.Waveform.PositionSine.W),
			float64(osc.Waveform.PositionSine.H),
		)
		cv.SetFillStyle(255, 255, 255, 255)
		cv.FillText(
			"sine",
			float64(osc.Waveform.PositionSine.X)+((float64(osc.Waveform.PositionSine.W)-cv.MeasureText("sine").Width)/2),
			float64(osc.Waveform.PositionSine.Y)+clientInterface.Style.FontSize,
		)
		if osc.Waveform.Value == "triangle" {
			cv.SetFillStyle(0, 255, 0, 255)
		} else {
			cv.SetFillStyle(255, 0, 0, 255)
		}
		cv.FillRect(
			float64(osc.Waveform.PositionTriangle.X),
			float64(osc.Waveform.PositionTriangle.Y),
			float64(osc.Waveform.PositionTriangle.W),
			float64(osc.Waveform.PositionTriangle.H),
		)
		cv.SetFillStyle(255, 255, 255, 255)
		cv.FillText(
			"triangle",
			float64(osc.Waveform.PositionTriangle.X)+((float64(osc.Waveform.PositionTriangle.W)-cv.MeasureText("triangle").Width)/2),
			float64(osc.Waveform.PositionTriangle.Y)+clientInterface.Style.FontSize,
		)
		if osc.Waveform.Value == "square" {
			cv.SetFillStyle(0, 255, 0, 255)
		} else {
			cv.SetFillStyle(255, 0, 0, 255)
		}
		cv.FillRect(
			float64(osc.Waveform.PositionSquare.X),
			float64(osc.Waveform.PositionSquare.Y),
			float64(osc.Waveform.PositionSquare.W),
			float64(osc.Waveform.PositionSquare.H),
		)
		cv.SetFillStyle(255, 255, 255, 255)
		cv.FillText(
			"square",
			float64(osc.Waveform.PositionSquare.X)+((float64(osc.Waveform.PositionSquare.W)-cv.MeasureText("square").Width)/2),
			float64(osc.Waveform.PositionSquare.Y)+clientInterface.Style.FontSize,
		)
		if osc.Waveform.Value == "flat" {
			cv.SetFillStyle(0, 255, 0, 255)
		} else {
			cv.SetFillStyle(255, 0, 0, 255)
		}
		cv.FillRect(
			float64(osc.Waveform.PositionFlat.X),
			float64(osc.Waveform.PositionFlat.Y),
			float64(osc.Waveform.PositionFlat.W),
			float64(osc.Waveform.PositionFlat.H),
		)
		cv.SetFillStyle(255, 255, 255, 255)
		cv.FillText(
			"flat",
			float64(osc.Waveform.PositionFlat.X)+((float64(osc.Waveform.PositionFlat.W)-cv.MeasureText("flat").Width)/2),
			float64(osc.Waveform.PositionFlat.Y)+clientInterface.Style.FontSize,
		)

		// utils buttons------------------------------------------------------------------------------------------------------------------------------------
		cv.SetFillStyle(255, 100, 100, 255)
		cv.FillRect(
			float64(clientInterface.Enregistrer.X),
			float64(clientInterface.Enregistrer.Y),
			float64(clientInterface.Enregistrer.W),
			float64(clientInterface.Enregistrer.H),
		)
		cv.SetFillStyle(255, 255, 255, 255)
		margeVertical = (float64(clientInterface.Enregistrer.H) - clientInterface.Style.FontSize) / 2
		margeHorizontale := (float64(clientInterface.Enregistrer.W) - cv.MeasureText("Save").Width) / 2
		cv.FillText(
			"Save",
			float64(clientInterface.Enregistrer.X)+margeVertical,
			float64(clientInterface.Enregistrer.Y)+margeHorizontale,
		)

		cv.SetFillStyle(0, 0, 255, 255)
		cv.FillRect(
			float64(clientInterface.Lire.X),
			float64(clientInterface.Lire.Y),
			float64(clientInterface.Lire.W),
			float64(clientInterface.Lire.H),
		)
		cv.SetFillStyle(255, 255, 255, 255)
		margeVertical = (float64(clientInterface.Lire.H) - clientInterface.Style.FontSize) / 2
		margeHorizontale = (float64(clientInterface.Lire.W) - cv.MeasureText("Play").Width) / 2
		cv.FillText(
			"Play",
			float64(clientInterface.Lire.X)+margeVertical,
			float64(clientInterface.Lire.Y)+margeHorizontale,
		)

		cv.SetFillStyle(255, 0, 0, 255)
		cv.FillRect(
			float64(clientInterface.Stopper.X),
			float64(clientInterface.Stopper.Y),
			float64(clientInterface.Stopper.W),
			float64(clientInterface.Stopper.H),
		)
		cv.SetFillStyle(255, 255, 255, 255)
		margeVertical = (float64(clientInterface.Stopper.H) - clientInterface.Style.FontSize) / 2
		margeHorizontale = (float64(clientInterface.Stopper.W) - cv.MeasureText("Stop").Width) / 2
		cv.FillText(
			"Stop",
			float64(clientInterface.Stopper.X)+margeVertical,
			float64(clientInterface.Stopper.Y)+margeHorizontale,
		)

		// onde de l'oscillateur -------------------------------------------------------------------------------------------------------
		cv.SetFillStyle(255, 255, 255, 255)
		for x := 0; x < wndWidth; x++ {
			t := float64(x) / float64(wndWidth)
			y := osc.Value(t)
			if osc.OnlyPositive.Value {
				if y < 0 {
					y = 0
				}
			}
			cv.FillRect(
				float64(x),
				float64(300-int(y*300)),
				1,
				1,
			)
		}

		// Présenter le rendu à l'écran
		//cv.Present()
		if !running {
			wnd.Destroy()
			os.Exit(1)
			return
		}
		// Introduire un léger délai pour limiter la boucle à environ 60 itérations par seconde
		sdl.Delay(16)
	})
}
