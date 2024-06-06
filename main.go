package main

import (
	"fmt"

	"github.com/tfriedel6/canvas/sdlcanvas"
	"github.com/veandco/go-sdl2/sdl"
)

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
			}
			/*
				else if event.State == sdl.RELEASED {
					mouseUpUpdate(int(event.Button))
				}
			*/
		}
	}
	return wnd
}

func mouseMoveUpdate(x, y float64) {
	souris.X = int(x)
	souris.Y = int(y)
}

func mouseDownUpdate(button, x, y int) {
	if x > int(osc.Amplitude.PositionUp.X) && x < int(osc.Amplitude.PositionUp.X+osc.Amplitude.PositionUp.W) && y > int(osc.Amplitude.PositionUp.Y) && y < int(osc.Amplitude.PositionUp.Y+osc.Amplitude.PositionUp.H) {
		osc.Amplitude.Value += 0.1
		fmt.Println("osc.Amplitude.Value up : ", osc.Amplitude.Value)
	}
	if x > int(osc.Amplitude.PositionDown.X) && x < int(osc.Amplitude.PositionDown.X+osc.Amplitude.PositionDown.W) && y > int(osc.Amplitude.PositionDown.Y) && y < int(osc.Amplitude.PositionDown.Y+osc.Amplitude.PositionDown.H) {
		osc.Amplitude.Value -= 0.1
		fmt.Println("osc.Amplitude.Value down : ", osc.Amplitude.Value)
	}
	if x > int(osc.Frequency.PositionUp.X) && x < int(osc.Frequency.PositionUp.X+osc.Frequency.PositionUp.W) && y > int(osc.Frequency.PositionUp.Y) && y < int(osc.Frequency.PositionUp.Y+osc.Frequency.PositionUp.H) {
		osc.Frequency.Value += 0.1
		fmt.Println("osc.Frequency.Value up : ", osc.Frequency.Value)
	}
	if x > int(osc.Frequency.PositionDown.X) && x < int(osc.Frequency.PositionDown.X+osc.Frequency.PositionDown.W) && y > int(osc.Frequency.PositionDown.Y) && y < int(osc.Frequency.PositionDown.Y+osc.Frequency.PositionDown.H) {
		osc.Frequency.Value -= 0.1
		fmt.Println("osc.Frequency.Value down : ", osc.Frequency.Value)
	}
	if x > int(osc.Phase.PositionUp.X) && x < int(osc.Phase.PositionUp.X+osc.Phase.PositionUp.W) && y > int(osc.Phase.PositionUp.Y) && y < int(osc.Phase.PositionUp.Y+osc.Phase.PositionUp.H) {
		osc.Phase.Value += 0.1
		fmt.Println("osc.Phase.Value up : ", osc.Phase.Value)
	}
	if x > int(osc.Phase.PositionDown.X) && x < int(osc.Phase.PositionDown.X+osc.Phase.PositionDown.W) && y > int(osc.Phase.PositionDown.Y) && y < int(osc.Phase.PositionDown.Y+osc.Phase.PositionDown.H) {
		osc.Phase.Value -= 0.1
		fmt.Println("osc.Phase.Value down : ", osc.Phase.Value)
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
func InitInterface() {
	osc = &Oscillator{
		Amplitude: Amplitude{
			Value:        1,
			PositionUp:   sdl.Rect{X: 50, Y: int32(wndHeight - 400), W: 70, H: 50},
			PositionDown: sdl.Rect{X: 150, Y: int32(wndHeight - 400), W: 70, H: 50},
		},
		Frequency: Frequency{
			Value:        1,
			PositionUp:   sdl.Rect{X: 50, Y: int32(wndHeight - 300), W: 70, H: 50},
			PositionDown: sdl.Rect{X: 150, Y: int32(wndHeight - 300), W: 70, H: 50},
		},
		Phase: Phase{
			Value:        0,
			PositionUp:   sdl.Rect{X: 50, Y: int32(wndHeight - 200), W: 70, H: 50},
			PositionDown: sdl.Rect{X: 150, Y: int32(wndHeight - 200), W: 70, H: 50},
		},
	}
	clientInterface = ClientInterface{
		Oscillator:  osc,
		Enregistrer: sdl.Rect{X: int32(wndWidth - 150), Y: int32(wndHeight - 250), W: 100, H: 50},
		Lire:        sdl.Rect{X: int32(wndWidth - 250), Y: int32(wndHeight - 100), W: 100, H: 50},
		Stopper:     sdl.Rect{X: int32(wndWidth - 150), Y: int32(wndHeight - 100), W: 100, H: 50},
	}

}

func setWindow(screenX, screenY int) {
	echelle = float64(screenX) / 1920
	wndWidth = screenX
	wndHeight = screenY
}

type Souris struct {
	X, Y int
}

var osc = &Oscillator{}
var clientInterface = ClientInterface{}
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
	//wnd.Window.SetPosition(0, 0)
	//wnd.Window.SetFullscreen(sdl.WINDOW_FULLSCREEN_DESKTOP)
	wnd.Window.SetResizable(false)
	//wnd.Window.SetPosition(sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED)
	//wnd.Window.SetSize(int32(wndWidth), int32(wndHeight))
	width, height := wnd.Window.GetSize()
	fmt.Println("width: ", width, "height: ", height)
	wnd = addEvent(wnd)

	InitInterface()

	//for running {
	wnd.MainLoop(func() {
		// Effacer l'écran avec une couleur (noir)
		cv.SetFillStyle(0, 0, 0, 255)
		cv.FillRect(0, 0, float64(wndWidth), float64(wndHeight))

		// Dessiner les boutons de l'oscillateur
		cv.SetFillStyle(255, 0, 0, 255)
		cv.FillRect(float64(osc.Amplitude.PositionUp.X), float64(osc.Amplitude.PositionUp.Y), float64(osc.Amplitude.PositionUp.W), float64(osc.Amplitude.PositionUp.H))
		cv.SetFillStyle(0, 255, 0, 255)
		cv.FillRect(float64(osc.Amplitude.PositionDown.X), float64(osc.Amplitude.PositionDown.Y), float64(osc.Amplitude.PositionDown.W), float64(osc.Amplitude.PositionDown.H))
		cv.SetFillStyle(255, 0, 0, 255)
		cv.FillRect(float64(osc.Frequency.PositionUp.X), float64(osc.Frequency.PositionUp.Y), float64(osc.Frequency.PositionUp.W), float64(osc.Frequency.PositionUp.H))
		cv.SetFillStyle(0, 255, 0, 255)
		cv.FillRect(float64(osc.Frequency.PositionDown.X), float64(osc.Frequency.PositionDown.Y), float64(osc.Frequency.PositionDown.W), float64(osc.Frequency.PositionDown.H))
		cv.SetFillStyle(255, 0, 0, 255)
		cv.FillRect(float64(osc.Phase.PositionUp.X), float64(osc.Phase.PositionUp.Y), float64(osc.Phase.PositionUp.W), float64(osc.Phase.PositionUp.H))
		cv.SetFillStyle(0, 255, 0, 255)
		cv.FillRect(float64(osc.Phase.PositionDown.X), float64(osc.Phase.PositionDown.Y), float64(osc.Phase.PositionDown.W), float64(osc.Phase.PositionDown.H))

		cv.SetFillStyle(255, 100, 100, 255)
		cv.FillRect(float64(clientInterface.Enregistrer.X), float64(clientInterface.Enregistrer.Y), float64(clientInterface.Enregistrer.W), float64(clientInterface.Enregistrer.H))
		cv.SetFillStyle(0, 0, 255, 255)
		cv.FillRect(float64(clientInterface.Lire.X), float64(clientInterface.Lire.Y), float64(clientInterface.Lire.W), float64(clientInterface.Lire.H))
		cv.SetFillStyle(255, 0, 0, 255)
		cv.FillRect(float64(clientInterface.Stopper.X), float64(clientInterface.Stopper.Y), float64(clientInterface.Stopper.W), float64(clientInterface.Stopper.H))

		// Dessiner l'onde de l'oscillateur
		cv.SetFillStyle(255, 255, 255, 255)
		for x := 0; x < 800; x++ {
			t := float64(x) / 800
			y := osc.Value(t)
			cv.FillRect(float64(x), float64(300-int(y*300)), 1, 1)
		}

		// Présenter le rendu à l'écran
		//cv.Present()

		// Introduire un léger délai pour limiter la boucle à environ 60 itérations par seconde
		sdl.Delay(16)
	})
}
