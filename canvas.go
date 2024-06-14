package main

import (
	"strconv"

	"github.com/tfriedel6/canvas"
)

type openning struct {
	running  bool
	x        int
	font     *canvas.Font
	fontSize float64
}

var open = openning{
	running: true,
	x:       0,
}

func (o *openning) loadFont(cv *canvas.Canvas) {
	font, err := cv.LoadFont("font/Ubuntu-MI.ttf")
	if err != nil {
		panic(err)
	}
	o.font = font
}

func (o *openning) drawOpenning(cv *canvas.Canvas) {
	cv.SetFillStyle(255, 255, 255, 255)
	cv.SetFont(open.font, open.fontSize)
	margeHorizontale := (float64(wndWidth) - cv.MeasureText("FullFreeGameplay").Width) / 2
	margeVerticale := ((float64(wndHeight) - clientInterface.Style.FontSize*10) / 2) + open.fontSize
	cv.FillText(
		"FullFreeGameplay",
		margeHorizontale,
		margeVerticale,
	)
	cv.SetFont(clientInterface.Style.Font, clientInterface.Style.FontSize)
	chunkSize := float64(wndWidth) / 255.0
	for i := 0; i < 255*2; i++ {
		if i > 255 {
			cv.SetFillStyle(0, 0, 0, 255)
		} else {
			cv.SetFillStyle(0, 0, 0, i)
		}
		x := (float64(i) * chunkSize) + float64(o.x)
		cv.FillRect(
			x,
			0,
			chunkSize,
			float64(wndHeight),
		)
	}
	o.x += int(chunkSize)
	if o.x > int(float64(wndWidth)*1.2) {
		o.running = false
		cv.SetFont(clientInterface.Style.Font, clientInterface.Style.FontSize)
	}
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

func (ci *ClientInterface) drawInterface(cv *canvas.Canvas) {

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

	text = "Max   " + strconv.FormatFloat(osc.MaxAmplitude.Value, 'f', 2, 64)
	cv.SetFillStyle(255, 255, 255, 255)
	cv.FillText(
		text,
		float64(osc.MaxAmplitude.PositionDown.X),
		float64(osc.MaxAmplitude.PositionDown.Y)-clientInterface.Style.FontSize,
	)
	cv.SetFillStyle(0, 255, 0, 255)
	cv.FillRect(
		float64(osc.MaxAmplitude.PositionUp.X),
		float64(osc.MaxAmplitude.PositionUp.Y),
		float64(osc.MaxAmplitude.PositionUp.W),
		float64(osc.MaxAmplitude.PositionUp.H),
	)
	cv.SetFillStyle(255, 255, 255, 255)
	cv.FillText(
		"+",
		float64(osc.MaxAmplitude.PositionUp.X)+margeHorizontalePlus,
		float64(osc.MaxAmplitude.PositionUp.Y)+margeVertical,
	)
	cv.SetFillStyle(255, 0, 0, 255)
	cv.FillRect(
		float64(osc.MaxAmplitude.PositionDown.X),
		float64(osc.MaxAmplitude.PositionDown.Y),
		float64(osc.MaxAmplitude.PositionDown.W),
		float64(osc.MaxAmplitude.PositionDown.H),
	)
	cv.SetFillStyle(255, 255, 255, 255)
	cv.FillText(
		"-",
		float64(osc.MaxAmplitude.PositionDown.X)+margeHorizontaleMoins,
		float64(osc.MaxAmplitude.PositionDown.Y)+margeVertical,
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

	if clientInterface.FileName.Value {
		cv.SetFillStyle(255, 0, 0, 255)
		cv.FillRect(
			float64(clientInterface.CloseFileName.Position.X),
			float64(clientInterface.CloseFileName.Position.Y),
			float64(clientInterface.CloseFileName.Position.W),
			float64(clientInterface.CloseFileName.Position.H),
		)
		margeVertical = ((float64(clientInterface.CloseFileName.Position.H) - clientInterface.Style.FontSize) / 2) + clientInterface.Style.FontSize
		margeHorizontale := (float64(clientInterface.CloseFileName.Position.W) - cv.MeasureText(clientInterface.CloseFileName.Value).Width) / 2
		cv.SetFillStyle(255, 255, 255, 255)
		cv.FillText(
			clientInterface.CloseFileName.Value,
			float64(clientInterface.CloseFileName.Position.X)+margeHorizontale,
			float64(clientInterface.CloseFileName.Position.Y)+margeVertical,
		)

		cv.SetFillStyle(100, 100, 100, 255)
		cv.FillRect(
			float64(clientInterface.FileName.PositionUp.X),
			float64(clientInterface.FileName.PositionUp.Y),
			float64(clientInterface.FileName.PositionUp.W),
			float64(clientInterface.FileName.PositionUp.H),
		)
		text = clientInterface.FileName.Champ
		for cv.MeasureText(text).Width > float64(clientInterface.FileName.PositionUp.W) {
			text = text[1:]
		}
		margeHorizontale = (float64(clientInterface.FileName.PositionUp.W) - cv.MeasureText(text).Width) / 2
		cv.SetFillStyle(255, 255, 255, 255)
		cv.FillText(
			text,
			float64(clientInterface.FileName.PositionUp.X)+margeHorizontale,
			float64(clientInterface.FileName.PositionUp.Y)+margeVertical,
		)
		text = "Choose filename and press enter"
		margeHorizontale = (float64(clientInterface.FileName.PositionUp.W) - cv.MeasureText(text).Width) / 2
		cv.SetFillStyle(255, 255, 255, 255)
		cv.FillText(
			text,
			float64(clientInterface.FileName.PositionUp.X)+margeHorizontale,
			float64(clientInterface.FileName.PositionUp.Y)+margeVertical-float64(clientInterface.FileName.PositionUp.H),
		)
	}

	cv.SetFillStyle(255, 100, 100, 255)
	cv.FillRect(
		float64(clientInterface.OpenMix.X),
		float64(clientInterface.OpenMix.Y),
		float64(clientInterface.OpenMix.W),
		float64(clientInterface.OpenMix.H),
	)
	cv.SetFillStyle(255, 255, 255, 255)
	margeVertical = ((float64(clientInterface.OpenMix.H) - clientInterface.Style.FontSize) / 2) + clientInterface.Style.FontSize
	margeHorizontale := (float64(clientInterface.OpenMix.W) - cv.MeasureText("Open Mix").Width) / 2
	cv.FillText(
		"Open mix",
		float64(clientInterface.OpenMix.X)+margeHorizontale, //+margeVertical,
		float64(clientInterface.OpenMix.Y)+margeVertical,    //+margeHorizontale+clientInterface.Style.FontSize,
	)

	cv.SetFillStyle(255, 100, 100, 255)
	cv.FillRect(
		float64(clientInterface.Enregistrer.X),
		float64(clientInterface.Enregistrer.Y),
		float64(clientInterface.Enregistrer.W),
		float64(clientInterface.Enregistrer.H),
	)
	cv.SetFillStyle(255, 255, 255, 255)
	margeVertical = ((float64(clientInterface.Enregistrer.H) - clientInterface.Style.FontSize) / 2) + clientInterface.Style.FontSize
	margeHorizontale = (float64(clientInterface.Enregistrer.W) - cv.MeasureText("Save osc").Width) / 2
	cv.FillText(
		"Save osc",
		float64(clientInterface.Enregistrer.X)+margeHorizontale, //+margeVertical,
		float64(clientInterface.Enregistrer.Y)+margeVertical,    //+margeHorizontale,
	)

	cv.SetFillStyle(255, 100, 100, 255)
	cv.FillRect(
		float64(clientInterface.SaveAll.X),
		float64(clientInterface.SaveAll.Y),
		float64(clientInterface.SaveAll.W),
		float64(clientInterface.SaveAll.H),
	)
	cv.SetFillStyle(255, 255, 255, 255)
	margeVertical = ((float64(clientInterface.SaveAll.H) - clientInterface.Style.FontSize) / 2) + clientInterface.Style.FontSize
	margeHorizontale = (float64(clientInterface.SaveAll.W) - cv.MeasureText("Save all").Width) / 2
	cv.FillText(
		"Save mix",
		float64(clientInterface.SaveAll.X)+margeHorizontale, //+margeVertical,
		float64(clientInterface.SaveAll.Y)+margeVertical,    //+margeHorizontale+clientInterface.Style.FontSize,
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

	/*
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
	*/
}
