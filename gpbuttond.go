package main

import (
	"fmt"
	"github.com/bendahl/uinput"
	"github.com/warthog618/gpiod"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// TODO EDIT HERE TO ADD MORE BUTTONS (1/2, scroll down for second edit zone)
const buttonCount = 17

// TODO END EDIT ZONE (1/2, scroll down for second edit zone)

// line to key mapping - global
var lineMap [10][2]int

// repeat timer configuration - global
var repeatDuration time.Duration

// create channel for signaling end of routineHold
var exitChannel = make(chan bool)

var mutex = &sync.Mutex{}

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
			time.Sleep(repeatDuration)
		}
	}
}

// called whenever a GPIO event is detected on a watched line - uses routineHold as a Go routine to repeat keystrokes for as long as buttons are held down
func eventHandler(keycode int) func(_ gpiod.LineEvent) {
	return func(evt gpiod.LineEvent) {
		mutex.Lock()
		defer mutex.Unlock()
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

	// display version and licensing information
	fmt.Print("\n================================================================================================\n\n" +
		"gpbuttond v0.2.0 - Copyright 2023 (Randall Winkhart) - https://github.com/rwinkhart/gpbuttond\n\n" +
		"This program is free software: you can redistribute it and/or modify it under the terms of\n" +
		"version 3 (only) of the GNU General Public License as published by the Free Software Foundation.\n\n" +
		"This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY;\n" +
		"without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.\n\n" +
		"See the GNU General Public License for more details:\n" +
		"https://opensource.org/licenses/GPL-3.0\n\n" +
		"================================================================================================\n")

	// check environment variables (settings)

	// line to key mapping - local
	var mapEnv1, mapPresent1 = os.LookupEnv("GPBD_MAP")
	var mapEnv1Split []string
	var pairCount int
	if mapPresent1 {
		mapEnv1Split = strings.Split(mapEnv1, ",")
		pairCount = len(mapEnv1Split)
	} else {
		fmt.Printf("\nERROR: No pairings provided!\n\n"+
			"GPIO lines must be mapped to keycodes through the setting of the GPBD_MAP environment variable.\n"+
			"The keycode for any given key can be found by using the widely available \"showkey\" command in a raw TTY.\n\n"+
			"The format for setting GPBD_MAP is as follows:\n"+
			" export GPBD_MAP=<GPIO line #>:<decimal keycode>,<GPIO line #>:<decimal keycode>, etc.\n\n"+
			"Example:\n"+
			" export GPBD_MAP=19:103,6:108,26:105,5:28\n\n"+
			"Additional things to consider:\n\n"+
			"LINE NUMBERING\n"+
			" Note that gpbuttond uses the GPIO line numbering reported by \"/dev/gpiochip0\", which typically refers to\n"+
			" the internal CPU/SoC numbering of the GPIO lines rather than the numbering as it relates to\n"+
			" the physical layout of the pins on the board. Be sure you are using the correct numbering scheme!\n\n"+
			"LINE PULL DIRECTION\n"+
			" By default, it is likely your GPIO lines are not all pulled in the same direction.\n"+
			" gpbuttond makes the assumption that all lines are pulled up by default.\n"+
			" The process of matching this behavior varies between devices.\n\n"+
			" On a Raspberry Pi, this can be done by modifying your config.txt file.\n"+
			" For example, lines 5 and 6 can be pulled up by default with the following:\n"+
			"  gpio=5,6=pu\n\n"+
			"MAXIMUM SUPPORTED BUTTONS\n"+
			" Note that this compiled version of gpbuttond supports a maximum of %d line-to-button pairings.\n"+
			" More can be easily added through simple modification of the source code.\n"+
			" The lines to edit are clearly marked with \"// TODO\" comments.\n\n"+
			"OTHER SUPPORTED CONFIGURATIONS (ENVIRONMENT VARIABLES)\n"+
			" GPBD_DEBOUNCE can optionally be set to apply a custom debounce time. Set equal to any integer (measured in milliseconds).\n"+
			" GPBD_REPEAT can optionally be set to alter the time before registering multiple keystrokes when a button is held down. Set equal to any integer (measured in milliseconds).\n\n", buttonCount)
		os.Exit(1)
	}

	for i, pair := range mapEnv1Split {
		var intConvert [2]int
		for i, num := range strings.Split(pair, ":") {
			intConvert[i], _ = strconv.Atoi(num)
		}
		lineMap[i] = [2]int{intConvert[0], intConvert[1]}
	}

	// debounce timer configuration (time button must be held before it is registered as a keystroke, in milliseconds)
	var mapEnv2, mapPresent2 = os.LookupEnv("GPBD_DEBOUNCE")
	var debounceDuration time.Duration
	if mapPresent2 {
		debounceMultiplier, _ := strconv.Atoi(mapEnv2)
		debounceDuration = time.Duration(debounceMultiplier) * time.Millisecond
	} else {
		debounceDuration = 20 * time.Millisecond // default to 20-millisecond debounce timer
	}

	// repeat timer configuration - local
	var mapEnv3, mapPresent3 = os.LookupEnv("GPBD_REPEAT")
	if mapPresent3 {
		repeatMultiplier, _ := strconv.Atoi(mapEnv3)
		repeatDuration = time.Duration(repeatMultiplier) * time.Millisecond
	} else {
		repeatDuration = 150 * time.Millisecond // default to 150-millisecond repeat timer
	}

	// begin edge detection - will call eventHandler whenever a watched GPIO line changes state
	for i := 0; i < min(pairCount, buttonCount); i++ {
		switch i + 1 {
		case 1:
			line1, _ := gpiod.RequestLine("gpiochip0", lineMap[i][0], gpiod.WithPullUp, gpiod.WithBothEdges, gpiod.WithDebounce(debounceDuration), gpiod.WithEventHandler(eventHandler(lineMap[i][1])))
			defer line1.Close()
		case 2:
			line2, _ := gpiod.RequestLine("gpiochip0", lineMap[i][0], gpiod.WithPullUp, gpiod.WithBothEdges, gpiod.WithDebounce(debounceDuration), gpiod.WithEventHandler(eventHandler(lineMap[i][1])))
			defer line2.Close()
		case 3:
			line3, _ := gpiod.RequestLine("gpiochip0", lineMap[i][0], gpiod.WithPullUp, gpiod.WithBothEdges, gpiod.WithDebounce(debounceDuration), gpiod.WithEventHandler(eventHandler(lineMap[i][1])))
			defer line3.Close()
		case 4:
			line4, _ := gpiod.RequestLine("gpiochip0", lineMap[i][0], gpiod.WithPullUp, gpiod.WithBothEdges, gpiod.WithDebounce(debounceDuration), gpiod.WithEventHandler(eventHandler(lineMap[i][1])))
			defer line4.Close()
		case 5:
			line5, _ := gpiod.RequestLine("gpiochip0", lineMap[i][0], gpiod.WithPullUp, gpiod.WithBothEdges, gpiod.WithDebounce(debounceDuration), gpiod.WithEventHandler(eventHandler(lineMap[i][1])))
			defer line5.Close()
		case 6:
			line6, _ := gpiod.RequestLine("gpiochip0", lineMap[i][0], gpiod.WithPullUp, gpiod.WithBothEdges, gpiod.WithDebounce(debounceDuration), gpiod.WithEventHandler(eventHandler(lineMap[i][1])))
			defer line6.Close()
		case 7:
			line7, _ := gpiod.RequestLine("gpiochip0", lineMap[i][0], gpiod.WithPullUp, gpiod.WithBothEdges, gpiod.WithDebounce(debounceDuration), gpiod.WithEventHandler(eventHandler(lineMap[i][1])))
			defer line7.Close()
		case 8:
			line8, _ := gpiod.RequestLine("gpiochip0", lineMap[i][0], gpiod.WithPullUp, gpiod.WithBothEdges, gpiod.WithDebounce(debounceDuration), gpiod.WithEventHandler(eventHandler(lineMap[i][1])))
			defer line8.Close()
		case 9:
			line9, _ := gpiod.RequestLine("gpiochip0", lineMap[i][0], gpiod.WithPullUp, gpiod.WithBothEdges, gpiod.WithDebounce(debounceDuration), gpiod.WithEventHandler(eventHandler(lineMap[i][1])))
			defer line9.Close()
		case 10:
			line10, _ := gpiod.RequestLine("gpiochip0", lineMap[i][0], gpiod.WithPullUp, gpiod.WithBothEdges, gpiod.WithDebounce(debounceDuration), gpiod.WithEventHandler(eventHandler(lineMap[i][1])))
			defer line10.Close()
		case 11:
			line11, _ := gpiod.RequestLine("gpiochip0", lineMap[i][0], gpiod.WithPullUp, gpiod.WithBothEdges, gpiod.WithDebounce(debounceDuration), gpiod.WithEventHandler(eventHandler(lineMap[i][1])))
			defer line11.Close()
		case 12:
			line12, _ := gpiod.RequestLine("gpiochip0", lineMap[i][0], gpiod.WithPullUp, gpiod.WithBothEdges, gpiod.WithDebounce(debounceDuration), gpiod.WithEventHandler(eventHandler(lineMap[i][1])))
			defer line12.Close()
		case 13:
			line13, _ := gpiod.RequestLine("gpiochip0", lineMap[i][0], gpiod.WithPullUp, gpiod.WithBothEdges, gpiod.WithDebounce(debounceDuration), gpiod.WithEventHandler(eventHandler(lineMap[i][1])))
			defer line13.Close()
		case 14:
			line14, _ := gpiod.RequestLine("gpiochip0", lineMap[i][0], gpiod.WithPullUp, gpiod.WithBothEdges, gpiod.WithDebounce(debounceDuration), gpiod.WithEventHandler(eventHandler(lineMap[i][1])))
			defer line14.Close()
		case 15:
			line15, _ := gpiod.RequestLine("gpiochip0", lineMap[i][0], gpiod.WithPullUp, gpiod.WithBothEdges, gpiod.WithDebounce(debounceDuration), gpiod.WithEventHandler(eventHandler(lineMap[i][1])))
			defer line15.Close()
		case 16:
			line16, _ := gpiod.RequestLine("gpiochip0", lineMap[i][0], gpiod.WithPullUp, gpiod.WithBothEdges, gpiod.WithDebounce(debounceDuration), gpiod.WithEventHandler(eventHandler(lineMap[i][1])))
			defer line16.Close()
		case 17:
			line17, _ := gpiod.RequestLine("gpiochip0", lineMap[i][0], gpiod.WithPullUp, gpiod.WithBothEdges, gpiod.WithDebounce(debounceDuration), gpiod.WithEventHandler(eventHandler(lineMap[i][1])))
			defer line17.Close()
			// TODO EDIT HERE TO ADD MORE BUTTONS (2/2)
			// TODO END EDIT ZONE (2/2)
		}
	}

	// run FOREVER >:)
	for {
		fmt.Print("\nThis program is a background daemon - run it with \"gpbuttond &\" to properly background it.\n")
		time.Sleep(time.Hour)
	}
}
