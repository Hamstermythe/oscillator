package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/sdlcanvas"
	"github.com/veandco/go-sdl2/sdl"
)

var osc = &Oscillator{}
var clientInterface = &ClientInterface{}
var echelle float64
var wndWidth, wndHeight int
var souris = Souris{}
var running = true

func setWindow(screenX, screenY int) {
	echelle = float64(screenX) / 1920
	wndWidth = screenX
	wndHeight = screenY
}

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
		if !running {
			wnd.Destroy()
			os.Exit(1)
			return
		}
		souris.Action()
		// Effacer l'écran avec une couleur (noir)
		if !clientInterface.ReloadingWave && !clientInterface.ReloadingSelector {
			cv.SetFillStyle(0, 0, 0, 255)
			cv.FillRect(0, 0, float64(wndWidth), float64(wndHeight))
			clientInterface.drawUncurrentWave(cv)
			clientInterface.drawOscillator(cv)
		}
		// Présenter le rendu à l'écran
		if clientInterface.AddOscillator {
			clientInterface.AddOscillator = false
			newOsc := clientInterface.newOscillator(cv)
			clientInterface.Oscillator = append(clientInterface.Oscillator, newOsc)
			clientInterface.CurrentOscillator = len(clientInterface.Oscillator) - 1
			osc = clientInterface.Oscillator[clientInterface.CurrentOscillator]
			clientInterface.ReloadSelector = true
			clientInterface.ReloadWave = true
		}
		if clientInterface.DeleteOscillator {
			clientInterface.DeleteOscillator = false
			var removerOsc []*Oscillator
			for i, o := range clientInterface.Oscillator {
				if i != clientInterface.CurrentOscillator {
					removerOsc = append(removerOsc, o)
				}
			}
			clientInterface.Oscillator = removerOsc
			var removerWave [][]float32
			for i, w := range clientInterface.Wave {
				if i != clientInterface.CurrentOscillator {
					removerWave = append(removerWave, w)
				}
			}
			clientInterface.Wave = removerWave
			clientInterface.CurrentOscillator = len(removerOsc) - 1
			osc = clientInterface.Oscillator[clientInterface.CurrentOscillator]
			clientInterface.ReloadSelector = true
			clientInterface.ReloadWave = true
		}
		if clientInterface.ReloadWave && !clientInterface.ReloadingWave {
			clientInterface.ReloadWave = false
			clientInterface.ReloadingWave = true
			go clientInterface.reloadWave()
		} else {
			clientInterface.drawCurrentWave(cv)
		}
		if clientInterface.ReloadSelector && !clientInterface.ReloadingSelector {
			clientInterface.ReloadSelector = false
			clientInterface.ReloadingSelector = true
			go clientInterface.reloadSelectorOscillator(cv)
		} else {
			clientInterface.drawSelector(cv)
		}
		// Introduire un léger délai pour limiter la boucle à environ 60 itérations par seconde
		sdl.Delay(16)
	})
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

func (ci *ClientInterface) drawSelector(cv *canvas.Canvas) {
	cv.SetFont(clientInterface.Style.Font, clientInterface.Style.FontSize*2)
	width := cv.MeasureText(">").Width + 10
	if !ci.DeroulantSelector.Value {
		cv.SetFillStyle(0, 255, 0, 255)
		cv.FillText(
			">",
			float64(ci.DeroulantSelector.PositionUp.X)+float64(ci.DeroulantSelector.PositionUp.W/2)-width/2,
			float64(ci.DeroulantSelector.PositionUp.Y)+float64(ci.DeroulantSelector.PositionUp.H/2)+clientInterface.Style.FontSize/2,
		)
		cv.SetFont(clientInterface.Style.Font, clientInterface.Style.FontSize)
		return
	}
	cv.SetFillStyle(255, 0, 0, 255)
	cv.FillText(
		"<",
		float64(ci.DeroulantSelector.PositionUp.X)+float64(ci.DeroulantSelector.PositionUp.W/2)-width/2,
		float64(ci.DeroulantSelector.PositionUp.Y)+float64(ci.DeroulantSelector.PositionUp.H/2)+clientInterface.Style.FontSize/2,
	)
	cv.SetFont(clientInterface.Style.Font, clientInterface.Style.FontSize)
	for i, sel := range ci.Selector {
		cv.SetFillStyle(sel.Color.R, sel.Color.G, sel.Color.B, sel.Color.A)
		cv.FillRect(
			float64(sel.Position.X),
			float64(sel.Position.Y),
			float64(sel.Position.W),
			float64(sel.Position.H),
		)
		cv.SetFillStyle(0, 0, 0, 255)
		cv.FillText(
			sel.Value,
			float64(sel.Position.X)+float64(sel.TextMargin),
			float64(sel.Position.Y)+float64(sel.Position.H)-(clientInterface.Style.FontSize/2),
		)
		if i-2 == ci.CurrentOscillator {
			cv.SetFillStyle(0, 0, 0, 255)
			cv.FillRect(
				float64(sel.Position.X),
				float64(sel.Position.Y),
				float64(sel.Position.W),
				float64(sel.Position.H)*0.05,
			)
			cv.FillRect(
				float64(sel.Position.X),
				float64(sel.Position.Y)+float64(sel.Position.H)-float64(sel.Position.H)*0.05,
				float64(sel.Position.W),
				float64(sel.Position.H)*0.05,
			)
			cv.FillRect(
				float64(sel.Position.X),
				float64(sel.Position.Y),
				float64(sel.Position.W)*0.05,
				float64(sel.Position.H),
			)
			cv.FillRect(
				float64(sel.Position.X)+float64(sel.Position.W)-float64(sel.Position.W)*0.05,
				float64(sel.Position.Y),
				float64(sel.Position.W)*0.05,
				float64(sel.Position.H),
			)
		}
	}
}

func (ci *ClientInterface) drawCurrentWave(cv *canvas.Canvas) {
	if ci.CurrentOscillator >= len(ci.Wave) {
		return
	}
	// onde de l'oscillateur -------------------------------------------------------------------------------------------------------
	indexOfGreaterWave := 0
	for index, wave := range ci.Wave {
		if len(ci.Wave[indexOfGreaterWave]) < len(wave) {
			indexOfGreaterWave = index
		}
	}
	//cv.SetFillStyle(255, 255, 255, 255)
	cv.SetFillStyle(osc.Color.Value.R, osc.Color.Value.G, osc.Color.Value.B, osc.Color.Value.A)
	wave := ci.Wave[ci.CurrentOscillator]
	samplesThis := len(wave)
	samplesGreater := len(ci.Wave[indexOfGreaterWave])
	//pointMultiplier := 4
	marge := 100.0
	pointNumber := (float64(wndWidth) - marge) //* pointMultiplier
	chunk := int(float64(samplesGreater) / pointNumber)
	//cv.SetFillStyle(0, 200, 200, 255)
	for x := 0; x < int(pointNumber) && x*chunk < samplesThis; x++ {
		t := float64(x * chunk)
		y := wave[int(t)]
		zeroY := float64(wndHeight) / 4
		cv.FillRect(
			float64(x+(int(marge)/2)), // / pointMultiplier),
			zeroY-float64(y)*float64(wndHeight)/5,
			1,
			1,
		)
	}
}

func (ci *ClientInterface) drawUncurrentWave(cv *canvas.Canvas) {
	indexOfGreaterWave := 0
	for index, wave := range ci.Wave {
		if len(ci.Wave[indexOfGreaterWave]) < len(wave) {
			indexOfGreaterWave = index
		}
	}
	for index, oscil := range ci.Oscillator {
		if index > len(ci.Wave)-1 {
			break
		}
		if index == ci.CurrentOscillator {
			continue
		}
		// onde de l'oscillateur -------------------------------------------------------------------------------------------------------
		cv.SetFillStyle(oscil.Color.Value.R, oscil.Color.Value.G, oscil.Color.Value.B, oscil.Color.Value.A)
		wave := ci.Wave[index]
		samplesThis := len(wave)
		samplesGreater := len(ci.Wave[indexOfGreaterWave])
		marge := 100.0
		pointNumber := (float64(wndWidth) - marge)
		chunk := int(float64(samplesGreater) / pointNumber)
		for x := 0; x < int(pointNumber) && x*chunk < samplesThis; x++ {
			t := float64(x * chunk)
			y := wave[int(t)]
			zeroY := float64(wndHeight) / 4
			cv.FillRect(
				float64(x+(int(marge)/2)), // / pointMultiplier),
				zeroY-float64(y)*float64(wndHeight)/5,
				1,
				1,
			)
		}
	}
}

func (ci *ClientInterface) drawOscillator(cv *canvas.Canvas) {

	// configurables bool fields -----------------------------------------------------------------------------------------------------------------------

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

	text = "Only Negative"
	marginHorizontal = cv.MeasureText(text).Width + 10
	cv.SetFillStyle(255, 255, 255, 255)
	cv.FillText(
		text,
		float64(osc.OnlyNegative.PositionUp.X)-marginHorizontal,
		float64(osc.OnlyNegative.PositionUp.Y)+float64(osc.OnlyNegative.PositionUp.H)+(clientInterface.Style.FontSize/2),
	) //-clientInterface.Style.FontSize)
	if osc.OnlyNegative.Value {
		cv.SetFillStyle(0, 255, 0, 255)
	} else {
		cv.SetFillStyle(0, 50, 0, 255)
	}
	cv.FillRect(
		float64(osc.OnlyNegative.PositionUp.X),
		float64(osc.OnlyNegative.PositionUp.Y),
		float64(osc.OnlyNegative.PositionUp.W),
		float64(osc.OnlyNegative.PositionUp.H),
	)
	if osc.OnlyNegative.Value {
		cv.SetFillStyle(255, 255, 255, 255)
	} else {
		cv.SetFillStyle(100, 130, 100, 255)
	}
	cv.FillText(
		"on",
		float64(osc.OnlyNegative.PositionUp.X)+((float64(osc.OnlyNegative.PositionUp.W)-cv.MeasureText("on").Width)/2),
		float64(osc.OnlyNegative.PositionUp.Y)+clientInterface.Style.FontSize,
	)
	if !osc.OnlyNegative.Value {
		cv.SetFillStyle(255, 0, 0, 255)
	} else {
		cv.SetFillStyle(50, 0, 0, 255)
	}
	cv.FillRect(
		float64(osc.OnlyNegative.PositionDown.X),
		float64(osc.OnlyNegative.PositionDown.Y),
		float64(osc.OnlyNegative.PositionDown.W),
		float64(osc.OnlyNegative.PositionDown.H),
	)
	if !osc.OnlyNegative.Value {
		cv.SetFillStyle(255, 255, 255, 255)
	} else {
		cv.SetFillStyle(130, 100, 100, 100)
	}
	cv.FillText(
		"off",
		float64(osc.OnlyNegative.PositionDown.X)+((float64(osc.OnlyNegative.PositionUp.W)-cv.MeasureText("off").Width)/2),
		float64(osc.OnlyNegative.PositionDown.Y)+clientInterface.Style.FontSize,
	)

	text = "Only Increase"
	marginHorizontal = cv.MeasureText(text).Width + 10
	cv.SetFillStyle(255, 255, 255, 255)
	cv.FillText(
		text,
		float64(osc.OnlyIncreased.PositionUp.X)-marginHorizontal,
		float64(osc.OnlyIncreased.PositionUp.Y)+float64(osc.OnlyIncreased.PositionUp.H)+(clientInterface.Style.FontSize/2),
	) //-clientInterface.Style.FontSize)
	if osc.OnlyIncreased.Value {
		cv.SetFillStyle(0, 255, 0, 255)
	} else {
		cv.SetFillStyle(0, 50, 0, 255)
	}
	cv.FillRect(
		float64(osc.OnlyIncreased.PositionUp.X),
		float64(osc.OnlyIncreased.PositionUp.Y),
		float64(osc.OnlyIncreased.PositionUp.W),
		float64(osc.OnlyIncreased.PositionUp.H),
	)
	if osc.OnlyIncreased.Value {
		cv.SetFillStyle(255, 255, 255, 255)
	} else {
		cv.SetFillStyle(100, 130, 100, 255)
	}
	cv.FillText(
		"on",
		float64(osc.OnlyIncreased.PositionUp.X)+((float64(osc.OnlyIncreased.PositionUp.W)-cv.MeasureText("on").Width)/2),
		float64(osc.OnlyIncreased.PositionUp.Y)+clientInterface.Style.FontSize,
	)
	if !osc.OnlyIncreased.Value {
		cv.SetFillStyle(255, 0, 0, 255)
	} else {
		cv.SetFillStyle(50, 0, 0, 255)
	}
	cv.FillRect(
		float64(osc.OnlyIncreased.PositionDown.X),
		float64(osc.OnlyIncreased.PositionDown.Y),
		float64(osc.OnlyIncreased.PositionDown.W),
		float64(osc.OnlyIncreased.PositionDown.H),
	)
	if !osc.OnlyIncreased.Value {
		cv.SetFillStyle(255, 255, 255, 255)
	} else {
		cv.SetFillStyle(130, 100, 100, 100)
	}
	cv.FillText(
		"off",
		float64(osc.OnlyIncreased.PositionDown.X)+((float64(osc.OnlyIncreased.PositionUp.W)-cv.MeasureText("off").Width)/2),
		float64(osc.OnlyIncreased.PositionDown.Y)+clientInterface.Style.FontSize,
	)

	text = "Only Decrease"
	marginHorizontal = cv.MeasureText(text).Width + 10
	cv.SetFillStyle(255, 255, 255, 255)
	cv.FillText(
		text,
		float64(osc.OnlyDecreased.PositionUp.X)-marginHorizontal,
		float64(osc.OnlyDecreased.PositionUp.Y)+float64(osc.OnlyDecreased.PositionUp.H)+(clientInterface.Style.FontSize/2),
	) //-clientInterface.Style.FontSize)
	if osc.OnlyDecreased.Value {
		cv.SetFillStyle(0, 255, 0, 255)
	} else {
		cv.SetFillStyle(0, 50, 0, 255)
	}
	cv.FillRect(
		float64(osc.OnlyDecreased.PositionUp.X),
		float64(osc.OnlyDecreased.PositionUp.Y),
		float64(osc.OnlyDecreased.PositionUp.W),
		float64(osc.OnlyDecreased.PositionUp.H),
	)
	if osc.OnlyDecreased.Value {
		cv.SetFillStyle(255, 255, 255, 255)
	} else {
		cv.SetFillStyle(100, 130, 100, 255)
	}
	cv.FillText(
		"on",
		float64(osc.OnlyDecreased.PositionUp.X)+((float64(osc.OnlyDecreased.PositionUp.W)-cv.MeasureText("on").Width)/2),
		float64(osc.OnlyDecreased.PositionUp.Y)+clientInterface.Style.FontSize,
	)
	if !osc.OnlyDecreased.Value {
		cv.SetFillStyle(255, 0, 0, 255)
	} else {
		cv.SetFillStyle(50, 0, 0, 255)
	}
	cv.FillRect(
		float64(osc.OnlyDecreased.PositionDown.X),
		float64(osc.OnlyDecreased.PositionDown.Y),
		float64(osc.OnlyDecreased.PositionDown.W),
		float64(osc.OnlyDecreased.PositionDown.H),
	)
	if !osc.OnlyDecreased.Value {
		cv.SetFillStyle(255, 255, 255, 255)
	} else {
		cv.SetFillStyle(130, 100, 100, 100)
	}
	cv.FillText(
		"off",
		float64(osc.OnlyDecreased.PositionDown.X)+((float64(osc.OnlyDecreased.PositionUp.W)-cv.MeasureText("off").Width)/2),
		float64(osc.OnlyDecreased.PositionDown.Y)+clientInterface.Style.FontSize,
	)

	// configurables fields column 1-----------------------------------------------------------------------------------------------------------------------
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

	text = "Sample Rate   " + strconv.FormatFloat(osc.SampleRate.Value, 'f', 2, 64)
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

	text = "Kick         " + strconv.FormatFloat(osc.Kick.Value, 'f', 2, 64)
	cv.SetFillStyle(255, 255, 255, 255)
	cv.FillText(
		text,
		float64(osc.Kick.PositionDown.X),
		float64(osc.Kick.PositionDown.Y)-clientInterface.Style.FontSize,
	)
	cv.SetFillStyle(0, 255, 0, 255)
	cv.FillRect(
		float64(osc.Kick.PositionUp.X),
		float64(osc.Kick.PositionUp.Y),
		float64(osc.Kick.PositionUp.W),
		float64(osc.Kick.PositionUp.H),
	)
	cv.SetFillStyle(255, 255, 255, 255)
	cv.FillText(
		"+",
		float64(osc.Kick.PositionUp.X)+margeHorizontalePlus,
		float64(osc.Kick.PositionUp.Y)+margeVertical,
	)
	cv.SetFillStyle(255, 0, 0, 255)
	cv.FillRect(
		float64(osc.Kick.PositionDown.X),
		float64(osc.Kick.PositionDown.Y),
		float64(osc.Kick.PositionDown.W),
		float64(osc.Kick.PositionDown.H),
	)
	cv.SetFillStyle(255, 255, 255, 255)
	cv.FillText(
		"-",
		float64(osc.Kick.PositionDown.X)+margeHorizontaleMoins,
		float64(osc.Kick.PositionDown.Y)+margeVertical,
	)

	text = "Asymétrie horizontale " + strconv.FormatFloat(osc.Hauteur.Value, 'f', 2, 64)
	cv.SetFillStyle(255, 255, 255, 255)
	cv.FillText(
		text,
		float64(osc.Hauteur.PositionDown.X),
		float64(osc.Hauteur.PositionDown.Y)-clientInterface.Style.FontSize,
	)
	cv.SetFillStyle(0, 255, 0, 255)
	cv.FillRect(
		float64(osc.Hauteur.PositionUp.X),
		float64(osc.Hauteur.PositionUp.Y),
		float64(osc.Hauteur.PositionUp.W),
		float64(osc.Hauteur.PositionUp.H),
	)
	cv.SetFillStyle(255, 255, 255, 255)
	cv.FillText(
		"+",
		float64(osc.Hauteur.PositionUp.X)+margeHorizontalePlus,
		float64(osc.Hauteur.PositionUp.Y)+margeVertical,
	)
	cv.SetFillStyle(255, 0, 0, 255)
	cv.FillRect(
		float64(osc.Hauteur.PositionDown.X),
		float64(osc.Hauteur.PositionDown.Y),
		float64(osc.Hauteur.PositionDown.W),
		float64(osc.Hauteur.PositionDown.H),
	)
	cv.SetFillStyle(255, 255, 255, 255)
	cv.FillText(
		"-",
		float64(osc.Hauteur.PositionDown.X)+margeHorizontaleMoins,
		float64(osc.Hauteur.PositionDown.Y)+margeVertical,
	)

	text = "Asymétrie verticale " + strconv.FormatFloat(osc.AsymetrieY.Value, 'f', 2, 64)
	cv.SetFillStyle(255, 255, 255, 255)
	cv.FillText(
		text,
		float64(osc.AsymetrieY.PositionDown.X),
		float64(osc.AsymetrieY.PositionDown.Y)-clientInterface.Style.FontSize,
	)
	cv.SetFillStyle(0, 255, 0, 255)
	cv.FillRect(
		float64(osc.AsymetrieY.PositionUp.X),
		float64(osc.AsymetrieY.PositionUp.Y),
		float64(osc.AsymetrieY.PositionUp.W),
		float64(osc.AsymetrieY.PositionUp.H),
	)
	cv.SetFillStyle(255, 255, 255, 255)
	cv.FillText(
		"+",
		float64(osc.AsymetrieY.PositionUp.X)+margeHorizontalePlus,
		float64(osc.AsymetrieY.PositionUp.Y)+margeVertical,
	)
	cv.SetFillStyle(255, 0, 0, 255)
	cv.FillRect(
		float64(osc.AsymetrieY.PositionDown.X),
		float64(osc.AsymetrieY.PositionDown.Y),
		float64(osc.AsymetrieY.PositionDown.W),
		float64(osc.AsymetrieY.PositionDown.H),
	)
	cv.SetFillStyle(255, 255, 255, 255)
	cv.FillText(
		"-",
		float64(osc.AsymetrieY.PositionDown.X)+margeHorizontaleMoins,
		float64(osc.AsymetrieY.PositionDown.Y)+margeVertical,
	)

	//configurables fields column 3-------------------------------------------------------------------------------------------------------------------
	text = "Read Duration " + strconv.FormatFloat(osc.Lecture.Value, 'f', 2, 64)
	cv.SetFillStyle(255, 255, 255, 255)
	cv.FillText(
		text,
		float64(osc.Lecture.PositionDown.X),
		float64(osc.Lecture.PositionDown.Y)-clientInterface.Style.FontSize,
	)
	cv.SetFillStyle(0, 255, 0, 255)
	cv.FillRect(
		float64(osc.Lecture.PositionUp.X),
		float64(osc.Lecture.PositionUp.Y),
		float64(osc.Lecture.PositionUp.W),
		float64(osc.Lecture.PositionUp.H),
	)
	cv.SetFillStyle(255, 255, 255, 255)
	cv.FillText(
		"+",
		float64(osc.Lecture.PositionUp.X)+margeHorizontalePlus,
		float64(osc.Lecture.PositionUp.Y)+margeVertical,
	)
	cv.SetFillStyle(255, 0, 0, 255)
	cv.FillRect(
		float64(osc.Lecture.PositionDown.X),
		float64(osc.Lecture.PositionDown.Y),
		float64(osc.Lecture.PositionDown.W),
		float64(osc.Lecture.PositionDown.H),
	)
	cv.SetFillStyle(255, 255, 255, 255)
	cv.FillText(
		"-",
		float64(osc.Lecture.PositionDown.X)+margeHorizontaleMoins,
		float64(osc.Lecture.PositionDown.Y)+margeVertical,
	)

	text = "Pause Duration " + strconv.FormatFloat(osc.Pause.Value, 'f', 2, 64)
	cv.SetFillStyle(255, 255, 255, 255)
	cv.FillText(
		text,
		float64(osc.Pause.PositionDown.X),
		float64(osc.Pause.PositionDown.Y)-clientInterface.Style.FontSize,
	)
	cv.SetFillStyle(0, 255, 0, 255)
	cv.FillRect(
		float64(osc.Pause.PositionUp.X),
		float64(osc.Pause.PositionUp.Y),
		float64(osc.Pause.PositionUp.W),
		float64(osc.Pause.PositionUp.H),
	)
	cv.SetFillStyle(255, 255, 255, 255)
	cv.FillText(
		"+",
		float64(osc.Pause.PositionUp.X)+margeHorizontalePlus,
		float64(osc.Pause.PositionUp.Y)+margeVertical,
	)
	cv.SetFillStyle(255, 0, 0, 255)
	cv.FillRect(
		float64(osc.Pause.PositionDown.X),
		float64(osc.Pause.PositionDown.Y),
		float64(osc.Pause.PositionDown.W),
		float64(osc.Pause.PositionDown.H),
	)
	cv.SetFillStyle(255, 255, 255, 255)
	cv.FillText(
		"-",
		float64(osc.Pause.PositionDown.X)+margeHorizontaleMoins,
		float64(osc.Pause.PositionDown.Y)+margeVertical,
	)

	text = "Start " + strconv.FormatFloat(osc.Start.Value, 'f', 2, 64)
	cv.SetFillStyle(255, 255, 255, 255)
	cv.FillText(
		text,
		float64(osc.Start.PositionDown.X),
		float64(osc.Start.PositionDown.Y)-clientInterface.Style.FontSize,
	)
	cv.SetFillStyle(0, 255, 0, 255)
	cv.FillRect(
		float64(osc.Start.PositionUp.X),
		float64(osc.Start.PositionUp.Y),
		float64(osc.Start.PositionUp.W),
		float64(osc.Start.PositionUp.H),
	)
	cv.SetFillStyle(255, 255, 255, 255)
	cv.FillText(
		"+",
		float64(osc.Start.PositionUp.X)+margeHorizontalePlus,
		float64(osc.Start.PositionUp.Y)+margeVertical,
	)
	cv.SetFillStyle(255, 0, 0, 255)
	cv.FillRect(
		float64(osc.Start.PositionDown.X),
		float64(osc.Start.PositionDown.Y),
		float64(osc.Start.PositionDown.W),
		float64(osc.Start.PositionDown.H),
	)
	cv.SetFillStyle(255, 255, 255, 255)
	cv.FillText(
		"-",
		float64(osc.Start.PositionDown.X)+margeHorizontaleMoins,
		float64(osc.Start.PositionDown.Y)+margeVertical,
	)

	text = "End " + strconv.FormatFloat(osc.End.Value, 'f', 2, 64)
	cv.SetFillStyle(255, 255, 255, 255)
	cv.FillText(
		text,
		float64(osc.End.PositionDown.X),
		float64(osc.End.PositionDown.Y)-clientInterface.Style.FontSize,
	)
	cv.SetFillStyle(0, 255, 0, 255)
	cv.FillRect(
		float64(osc.End.PositionUp.X),
		float64(osc.End.PositionUp.Y),
		float64(osc.End.PositionUp.W),
		float64(osc.End.PositionUp.H),
	)
	cv.SetFillStyle(255, 255, 255, 255)
	cv.FillText(
		"+",
		float64(osc.End.PositionUp.X)+margeHorizontalePlus,
		float64(osc.End.PositionUp.Y)+margeVertical,
	)
	cv.SetFillStyle(255, 0, 0, 255)
	cv.FillRect(
		float64(osc.End.PositionDown.X),
		float64(osc.End.PositionDown.Y),
		float64(osc.End.PositionDown.W),
		float64(osc.End.PositionDown.H),
	)
	cv.SetFillStyle(255, 255, 255, 255)
	cv.FillText(
		"-",
		float64(osc.End.PositionDown.X)+margeHorizontaleMoins,
		float64(osc.End.PositionDown.Y)+margeVertical,
	)

	// configurables fields column 4---------------------------------------------------------------------------------------------------------------------

	text = "Color R: " + strconv.Itoa(int(osc.Color.Value.R)) + " G: " + strconv.Itoa(int(osc.Color.Value.G)) + " B: " + strconv.Itoa(int(osc.Color.Value.B))
	cv.SetFillStyle(255, 255, 255, 255)
	cv.FillText(
		text,
		float64(osc.Color.PositionDownR.X),
		float64(osc.Color.PositionDownR.Y)-clientInterface.Style.FontSize,
	)
	cv.SetFillStyle(0, 255, 0, 255)
	cv.FillRect(
		float64(osc.Color.PositionUpR.X),
		float64(osc.Color.PositionUpR.Y),
		float64(osc.Color.PositionUpR.W),
		float64(osc.Color.PositionUpR.H),
	)
	cv.SetFillStyle(255, 255, 255, 255)
	cv.FillText(
		"+",
		float64(osc.Color.PositionUpR.X)+margeHorizontalePlus,
		float64(osc.Color.PositionUpR.Y)+margeVertical,
	)
	cv.SetFillStyle(255, 0, 0, 255)
	cv.FillRect(
		float64(osc.Color.PositionDownR.X),
		float64(osc.Color.PositionDownR.Y),
		float64(osc.Color.PositionDownR.W),
		float64(osc.Color.PositionDownR.H),
	)
	cv.SetFillStyle(255, 255, 255, 255)
	cv.FillText(
		"-",
		float64(osc.Color.PositionDownR.X)+margeHorizontaleMoins,
		float64(osc.Color.PositionDownR.Y)+margeVertical,
	)
	cv.SetFillStyle(0, 255, 0, 255)
	cv.FillRect(
		float64(osc.Color.PositionUpG.X),
		float64(osc.Color.PositionUpG.Y),
		float64(osc.Color.PositionUpG.W),
		float64(osc.Color.PositionUpG.H),
	)
	cv.SetFillStyle(255, 255, 255, 255)
	cv.FillText(
		"+",
		float64(osc.Color.PositionUpG.X)+margeHorizontalePlus,
		float64(osc.Color.PositionUpG.Y)+margeVertical,
	)
	cv.SetFillStyle(255, 0, 0, 255)
	cv.FillRect(
		float64(osc.Color.PositionDownG.X),
		float64(osc.Color.PositionDownG.Y),
		float64(osc.Color.PositionDownG.W),
		float64(osc.Color.PositionDownG.H),
	)
	cv.SetFillStyle(255, 255, 255, 255)
	cv.FillText(
		"-",
		float64(osc.Color.PositionDownG.X)+margeHorizontaleMoins,
		float64(osc.Color.PositionDownG.Y)+margeVertical,
	)
	cv.SetFillStyle(0, 255, 0, 255)
	cv.FillRect(
		float64(osc.Color.PositionUpB.X),
		float64(osc.Color.PositionUpB.Y),
		float64(osc.Color.PositionUpB.W),
		float64(osc.Color.PositionUpB.H),
	)
	cv.SetFillStyle(255, 255, 255, 255)
	cv.FillText(
		"+",
		float64(osc.Color.PositionUpB.X)+margeHorizontalePlus,
		float64(osc.Color.PositionUpB.Y)+margeVertical,
	)
	cv.SetFillStyle(255, 0, 0, 255)
	cv.FillRect(
		float64(osc.Color.PositionDownB.X),
		float64(osc.Color.PositionDownB.Y),
		float64(osc.Color.PositionDownB.W),
		float64(osc.Color.PositionDownB.H),
	)
	cv.SetFillStyle(255, 255, 255, 255)
	cv.FillText(
		"-",
		float64(osc.Color.PositionDownB.X)+margeHorizontaleMoins,
		float64(osc.Color.PositionDownB.Y)+margeVertical,
	)

	// utils buttons------------------------------------------------------------------------------------------------------------------------------------

	cv.SetFillStyle(255, 100, 100, 255)
	cv.FillRect(
		float64(clientInterface.ExportWave.X),
		float64(clientInterface.ExportWave.Y),
		float64(clientInterface.ExportWave.W),
		float64(clientInterface.ExportWave.H),
	)
	cv.SetFillStyle(255, 255, 255, 255)
	margeVertical = (float64(clientInterface.ExportWave.H) - clientInterface.Style.FontSize) / 2
	margeHorizontale := (float64(clientInterface.ExportWave.W) - cv.MeasureText("Export").Width) / 2
	cv.FillText(
		"Export",
		float64(clientInterface.ExportWave.X)+margeVertical,
		float64(clientInterface.ExportWave.Y)+margeHorizontale,
	)

	cv.SetFillStyle(255, 100, 100, 255)
	cv.FillRect(
		float64(clientInterface.Enregistrer.X),
		float64(clientInterface.Enregistrer.Y),
		float64(clientInterface.Enregistrer.W),
		float64(clientInterface.Enregistrer.H),
	)
	cv.SetFillStyle(255, 255, 255, 255)
	margeVertical = (float64(clientInterface.Enregistrer.H) - clientInterface.Style.FontSize) / 2
	margeHorizontale = (float64(clientInterface.Enregistrer.W) - cv.MeasureText("Save").Width) / 2
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

}
