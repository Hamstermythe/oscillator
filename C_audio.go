package main

// #include <stdint.h>
// typedef uint8_t Uint8;
// void Wave(void *userdata, Uint8 *stream, int len);
import "C"
import (
	"fmt"
	"log"
	"reflect"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

var (
	audioSamples []float32
	currentIndex int
)

//export Wave
func Wave(userdata unsafe.Pointer, stream *C.Uint8, length C.int) {
	n := int(length) / 4 // SDL expects byte length, but we're working with float32 (4 bytes each)
	if currentIndex+n > len(audioSamples) {
		n = len(audioSamples) - currentIndex
	}
	if n <= 0 {
		return
	}

	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(stream)),
		Len:  n,
		Cap:  n,
	}

	buf := *(*[]float32)(unsafe.Pointer(&hdr))
	copy(buf, audioSamples[currentIndex:currentIndex+n])
	currentIndex += n
}

func playAudio() {
	for running {
		audio := <-chanSon
		audioSamples = audio.arr
		currentIndex = 0
		soundDuration := audio.duration * 1000
		spec := &sdl.AudioSpec{
			Freq:     int32(clientInterface.SampleRate * 2),
			Format:   sdl.AUDIO_F32LSB, // <-- little endian float 32, since we are working on a little endian cpu
			Channels: 1,                // <-- Reduced channel to 1 for the sake of simplicity
			Samples:  4096,             //len(clientInterface.CondensedWave)),
			Callback: sdl.AudioCallback(C.Wave),
		}
		if err := sdl.OpenAudio(spec, nil); err != nil {
			log.Println(err)
			return
		}
		sdl.PauseAudio(false)
		sdl.Delay(uint32(soundDuration)) // play audio for long enough to understand whether it works
		sdl.CloseAudio()
	}
}

func restoreFloat32(b byte) float32 {
	v := int16(b) | int16(b)<<8 // assuming 16-bit PCM data
	return float32(v) / 32768.0
}

func readWaveFile(filename string) ([]float32, float64) {
	arrByte, n := sdl.LoadFile(filename)
	if n == 0 {
		fmt.Println("Error loading file length: ", n)
		return nil, 0
	}
	var fileData []float32
	for i := 0; i < len(arrByte)-1; i += 1 {
		fileData = append(fileData, restoreFloat32(arrByte[i]))
	}
	soundDuration := float64(len(fileData)) / float64(clientInterface.SampleRate*2)
	fmt.Println("len(samples): ", len(fileData))
	return fileData, soundDuration
}
