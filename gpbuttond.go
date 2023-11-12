package main

import (
	"github.com/bendahl/uinput"
	"github.com/warthog618/gpiod"
	"time"
)

// map GPIO pins to keys - auxiliary pins do not yet have programmed functions
var pinUp = 19
var pinDown = 6
var pinLeft = 26
var pinRight = 5

//var pinAux1 = 16
//var pinAux2 = 20
//var pinAux3 = 21
//var pinAux4 = 13

// create channel for signaling end of routineHold
var exitChannel = make(chan bool)

var kbd, _ = uinput.CreateKeyboard("/dev/uinput", []byte("ttypodvirtualkbd"))

func routineHold(offset int, inputChannel <-chan bool) {
	for {
		select {
		case <-inputChannel:
			return
		default:
			switch offset {
			case pinUp:
				//println("up")
				kbd.KeyPress(uinput.KeyUp)
				time.Sleep(150 * time.Millisecond)
			case pinDown:
				//println("down")
				kbd.KeyPress(uinput.KeyDown)
				time.Sleep(150 * time.Millisecond)
			case pinLeft:
				//println("left")
				kbd.KeyPress(uinput.KeyLeft)
				time.Sleep(150 * time.Millisecond)
			case pinRight:
				//println("right")
				kbd.KeyPress(uinput.KeyEnter)
				time.Sleep(150 * time.Millisecond)
			}
		}
	}
}

func eventHandler(offset int) func(_ gpiod.LineEvent) {
	return func(evt gpiod.LineEvent) {

		if evt.Type == gpiod.LineEventRisingEdge { // button release
			exitChannel <- true

		} else if evt.Type == gpiod.LineEventFallingEdge { // button press
			go routineHold(offset, exitChannel)
		}
	}
}

func main() {
	// safely close virtual keyboard after program exits
	defer kbd.Close()

	// begin edge detection
	lineUp, _ := gpiod.RequestLine("gpiochip0", pinUp, gpiod.WithPullUp, gpiod.WithBothEdges, gpiod.WithDebounce(20*time.Millisecond), gpiod.WithEventHandler(eventHandler(pinUp)))
	defer lineUp.Close()
	lineDown, _ := gpiod.RequestLine("gpiochip0", pinDown, gpiod.WithPullUp, gpiod.WithBothEdges, gpiod.WithDebounce(20*time.Millisecond), gpiod.WithEventHandler(eventHandler(pinDown)))
	defer lineDown.Close()
	lineLeft, _ := gpiod.RequestLine("gpiochip0", pinLeft, gpiod.WithPullUp, gpiod.WithBothEdges, gpiod.WithDebounce(20*time.Millisecond), gpiod.WithEventHandler(eventHandler(pinLeft)))
	defer lineLeft.Close()
	lineRight, _ := gpiod.RequestLine("gpiochip0", pinRight, gpiod.WithPullUp, gpiod.WithBothEdges, gpiod.WithDebounce(20*time.Millisecond), gpiod.WithEventHandler(eventHandler(pinRight)))
	defer lineRight.Close()

	// run FOREVER >:)
	for {
		time.Sleep(time.Hour)
	}
}
