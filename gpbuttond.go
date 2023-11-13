package main

import (
	"fmt"
	"github.com/bendahl/uinput"
	"github.com/warthog618/gpiod"
	"os"
	"strconv"
	"strings"
	"time"
)

// pin to key mapping, global
var pinMap [10][2]int

// create channel for signaling end of routineHold
var exitChannel = make(chan bool)

// create virtual uinput keyboard for simulating keystrokes
var kbd, _ = uinput.CreateKeyboard("/dev/uinput", []byte("gpbuttondvirtualkbd"))

// meant to run as a Go routine (launched from eventHandler) - repeats keystroke until GPIO button is no longer held down
func routineHold(keycode int, inputChannel <-chan bool) {
	for {
		select {
		case <-inputChannel:
			return
		default:
			kbd.KeyPress(keycode)
			time.Sleep(150 * time.Millisecond)
		}
	}
}

// called whenever a GPIO event is detected on a watched line - uses routineHold as a Go routine to repeat keystrokes for as long as buttons are held down
func eventHandler(keycode int) func(_ gpiod.LineEvent) {
	return func(evt gpiod.LineEvent) {
		if evt.Type == gpiod.LineEventRisingEdge { // button release
			// send signal to exit Go routine (stops repeating keystrokes)
			exitChannel <- true

		} else if evt.Type == gpiod.LineEventFallingEdge { // button press
			// launch Go routine to repeat keystrokes until a LineEventRisingEdge event is received
			go routineHold(keycode, exitChannel)
		}
	}
}

func main() {
	// ensure closure of keyboard after program exits
	defer kbd.Close()

	// check environment variables (settings)
	// pin to key mapping - local

	var mapEnv, mapPresent = os.LookupEnv("GPBD_MAP")
	var mapEnvSplit []string
	var pairCount int
	if mapPresent {
		mapEnvSplit = strings.Split(mapEnv, ",")
		pairCount = len(mapEnvSplit)
	} else {
		fmt.Println("Error: no pairings provided - see \"gpbuttond help\"\n\nNOTE THAT \"gpbuttond help\" FUNCTIONALITY HAS NOT YET BEEN ADDED AND YOU ARE USING AN UNTAGGED RELEASE")
		os.Exit(1)
	}

	for i, pair := range mapEnvSplit {
		var intConvert [2]int
		for i, num := range strings.Split(pair, ":") {
			intConvert[i], _ = strconv.Atoi(num)
		}
		pinMap[i] = [2]int{intConvert[0], intConvert[1]}
	}

	// begin edge detection - will call eventHandler whenever a watched GPIO line changes state
	// TODO EDIT HERE TO ADD MORE BUTTONS (1/2)
	const buttonCount = 10
	// TODO END EDIT ZONE (1/2)
	for i := 0; i < min(pairCount, buttonCount); i++ {
		switch i + 1 {
		case 1:
			line1, _ := gpiod.RequestLine("gpiochip0", pinMap[i][0], gpiod.WithPullUp, gpiod.WithBothEdges, gpiod.WithDebounce(20*time.Millisecond), gpiod.WithEventHandler(eventHandler(pinMap[i][1])))
			defer line1.Close()
		case 2:
			line2, _ := gpiod.RequestLine("gpiochip0", pinMap[i][0], gpiod.WithPullUp, gpiod.WithBothEdges, gpiod.WithDebounce(20*time.Millisecond), gpiod.WithEventHandler(eventHandler(pinMap[i][1])))
			defer line2.Close()
		case 3:
			line3, _ := gpiod.RequestLine("gpiochip0", pinMap[i][0], gpiod.WithPullUp, gpiod.WithBothEdges, gpiod.WithDebounce(20*time.Millisecond), gpiod.WithEventHandler(eventHandler(pinMap[i][1])))
			defer line3.Close()
		case 4:
			line4, _ := gpiod.RequestLine("gpiochip0", pinMap[i][0], gpiod.WithPullUp, gpiod.WithBothEdges, gpiod.WithDebounce(20*time.Millisecond), gpiod.WithEventHandler(eventHandler(pinMap[i][1])))
			defer line4.Close()
		case 5:
			line5, _ := gpiod.RequestLine("gpiochip0", pinMap[i][0], gpiod.WithPullUp, gpiod.WithBothEdges, gpiod.WithDebounce(20*time.Millisecond), gpiod.WithEventHandler(eventHandler(pinMap[i][1])))
			defer line5.Close()
		case 6:
			line6, _ := gpiod.RequestLine("gpiochip0", pinMap[i][0], gpiod.WithPullUp, gpiod.WithBothEdges, gpiod.WithDebounce(20*time.Millisecond), gpiod.WithEventHandler(eventHandler(pinMap[i][1])))
			defer line6.Close()
		case 7:
			line7, _ := gpiod.RequestLine("gpiochip0", pinMap[i][0], gpiod.WithPullUp, gpiod.WithBothEdges, gpiod.WithDebounce(20*time.Millisecond), gpiod.WithEventHandler(eventHandler(pinMap[i][1])))
			defer line7.Close()
		case 8:
			line8, _ := gpiod.RequestLine("gpiochip0", pinMap[i][0], gpiod.WithPullUp, gpiod.WithBothEdges, gpiod.WithDebounce(20*time.Millisecond), gpiod.WithEventHandler(eventHandler(pinMap[i][1])))
			defer line8.Close()
		case 9:
			line9, _ := gpiod.RequestLine("gpiochip0", pinMap[i][0], gpiod.WithPullUp, gpiod.WithBothEdges, gpiod.WithDebounce(20*time.Millisecond), gpiod.WithEventHandler(eventHandler(pinMap[i][1])))
			defer line9.Close()
		case 10:
			line10, _ := gpiod.RequestLine("gpiochip0", pinMap[i][0], gpiod.WithPullUp, gpiod.WithBothEdges, gpiod.WithDebounce(20*time.Millisecond), gpiod.WithEventHandler(eventHandler(pinMap[i][1])))
			defer line10.Close()
			// TODO EDIT HERE TO ADD MORE BUTTONS (2/2)
			// TODO END EDIT ZONE (2/2)
		}
	}

	// run FOREVER >:)
	for {
		time.Sleep(time.Hour)
	}
}
