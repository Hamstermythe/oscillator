package main

import (
	"os"

	"github.com/tfriedel6/canvas/sdlcanvas"
	"github.com/veandco/go-sdl2/sdl"
)

type Son struct {
	arr      []float32
	duration float64
}

var osc = &Oscillator{}
var clientInterface = &ClientInterface{}
var echelle float64
var wndWidth, wndHeight int
var souris = Souris{}
var running = true
var chanSon = make(chan *Son)

func setWindow(screenX, screenY int) {
	echelle = float64(screenX) / 1920
	wndWidth = screenX
	wndHeight = screenY
	open.x = -wndWidth
	open.fontSize = float64(wndWidth) / 12
}

func main() {
	// Initialiser SDL2
	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO); err != nil {
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
	wnd = addEvent(wnd)

	open.loadFont(cv)
	clientInterface.InitInterface(cv)

	go playAudio()

	wnd.MainLoop(func() {
		if !running {
			wnd.Destroy()
			os.Exit(1)
			return
		}
		// Présenter l'openning
		if open.running {
			open.drawOpenning(cv)
		} else {
			souris.Action()
			// Effacer l'écran
			if !clientInterface.ReloadingWave && !clientInterface.ReloadingSelector {
				cv.SetFillStyle(0, 0, 0, 255)
				cv.FillRect(0, 0, float64(wndWidth), float64(wndHeight))
				clientInterface.drawUncurrentWave(cv)
				clientInterface.drawInterface(cv)
			}
			// Ajouter puis recharger les oscillateurs
			if clientInterface.AddOscillator {
				clientInterface.AddOscillator = false
				newOsc := clientInterface.newOscillator(cv)
				clientInterface.Oscillator = append(clientInterface.Oscillator, newOsc)
				clientInterface.CurrentOscillator = len(clientInterface.Oscillator) - 1
				osc = clientInterface.Oscillator[clientInterface.CurrentOscillator]
				clientInterface.ReloadSelector = true
				clientInterface.ReloadWave = true
			}
			// Supprimer puis recharger un/les oscillateur
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
			// Recharger l'onde
			if clientInterface.ReloadWave && !clientInterface.ReloadingWave {
				clientInterface.ReloadWave = false
				clientInterface.ReloadingWave = true
				go clientInterface.reloadWave()
			} else {
				// Dessiner l'onde
				clientInterface.drawCurrentWave(cv)
			}
			// Recharger le sélecteur
			if clientInterface.ReloadSelector && !clientInterface.ReloadingSelector {
				clientInterface.ReloadSelector = false
				clientInterface.ReloadingSelector = true
				go clientInterface.reloadSelectorOscillator(cv)
			} else {
				// Dessiner l'oscillateur courant
				clientInterface.drawSelector(cv)
			}

			sdl.Delay(16)
		}
	})
}
