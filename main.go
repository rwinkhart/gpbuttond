package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bendahl/uinput"
	gpio "github.com/warthog618/go-gpiocdev"
)

// line to key mapping - global
var lineMap [][3]int

// repeat timer configuration - global
var repeatDuration time.Duration

// long press timer configuration - global
var longPressDuration time.Duration

// create channel for signaling end of routineHold
var exitChannel = make(chan bool)

// create sync.Mutex object for preventing routineHold conflicts
var mutex = &sync.Mutex{}

// create virtual uinput keyboard for simulating keystrokes
var kbd, _ = uinput.CreateKeyboard("/dev/uinput", []byte("gpbuttondvirtualkbd"))

// meant to run as a Go routine (launched from routineHoldLong) - acts as a timer to determine whether a GPIO button has been held for long enough to actuate the specified long keycode
func routineLongpressTimer(timerChannel chan<- bool) {
	time.Sleep(longPressDuration)
	timerChannel <- true
}

// meant to run as a Go routine (launched from eventHandler) - repeats keystroke until GPIO button is no longer held down
func routineHoldShort(keycode int, exitChannel <-chan bool) {
	for {
		select {
		case <-exitChannel:
			return
		default:
			kbd.KeyPress(keycode)
			time.Sleep(repeatDuration)
		}
	}
}

// meant to run as a Go routine (launched from eventHandler) - executes a keystroke based on a timer
func routineHoldLong(keycode int, longKeycode int, exitChannel <-chan bool) {
	var timerChannel = make(chan bool)
	var timerTriggered bool

	// launch an additional Go routine to track if the button has been held long enough to register the long keycode
	go routineLongpressTimer(timerChannel)

	for {
		select {
		case <-exitChannel:
			if !timerTriggered { // only run if the long keycode has not been executed
				kbd.KeyPress(keycode)
			}
			return
		case <-timerChannel:
			kbd.KeyPress(longKeycode)
			timerTriggered = true
		}
	}
}

// called whenever a GPIO event is detected on a watched line - uses routineHold as a Go routine to repeat keystrokes for as long as buttons are held down
func eventHandler(keycode int, longKeycode int) func(_ gpio.LineEvent) {
	return func(evt gpio.LineEvent) {
		// set mutex.Lock() to prevent routineHold conflicts
		mutex.Lock()
		defer mutex.Unlock()

		if evt.Type == gpio.LineEventRisingEdge { // button release
			// send signal to exit Go routine (stops repeating keystrokes)
			exitChannel <- true

		} else if evt.Type == gpio.LineEventFallingEdge { // button press
			if longKeycode == 0 { // if no long keycode was specified
				// launch Go routine to repeat keystrokes until a LineEventRisingEdge event is received
				go routineHoldShort(keycode, exitChannel)
			} else { // if a long keycode was specified
				// launch a Go routine to execute a keystroke based on a timer
				go routineHoldLong(keycode, longKeycode, exitChannel)
			}
		}
	}
}

func main() {
	// ensure closure of keyboard after program exits
	defer kbd.Close()

	// display version and licensing information
	fmt.Print("\n==================================================================================================\n\n" +
		"gpbuttond v0.3.1 - Copyright 2023-2024 (Randall Winkhart) - https://github.com/rwinkhart/gpbuttond\n\n" +
		"This program is free software: you can redistribute it and/or modify it under the terms of\n" +
		"version 3 (only) of the GNU General Public License as published by the Free Software Foundation.\n\n" +
		"This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY;\n" +
		"without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.\n\n" +
		"See the GNU General Public License for more details:\n" +
		"https://opensource.org/licenses/GPL-3.0\n\n" +
		"==================================================================================================\n")

	// check environment variables (settings)

	// line to key mapping - local
	var mapEnv1, mapPresent1 = os.LookupEnv("GPBD_MAP")
	var mapEnv1Split []string
	var mapCount int
	if mapPresent1 {
		mapEnv1Split = strings.Split(mapEnv1, ",")
		mapCount = len(mapEnv1Split)
	} else {
		fmt.Print("\nERROR: No mappings provided!\n\n" +
			"GPIO lines must be mapped to keycodes through the setting of the GPBD_MAP environment variable.\n" +
			"The keycode for any given key can be found by using the widely available \"showkey\" command in a raw TTY.\n\n" +
			"The format for setting GPBD_MAP is as follows:\n" +
			" export GPBD_MAP=<GPIO line #>:<decimal keycode>:[optional long press decimal keycode],<GPIO line #>:<decimal keycode>:[optional long press decimal keycode], etc.\n\n" +
			"Example:\n" +
			" export GPBD_MAP=19:103:1,6:108,26:105,5:28\n\n" +
			"Additional things to consider:\n\n" +
			"LINE NUMBERING\n" +
			" Note that gpbuttond uses the GPIO line numbering reported by \"/dev/gpiochip0\", which typically refers to\n" +
			" the internal CPU/SoC numbering of the GPIO lines rather than the numbering as it relates to\n" +
			" the physical layout of the pins on the board. Be sure you are using the correct numbering scheme!\n\n" +
			"LINE PULL DIRECTION\n" +
			" By default, it is likely your GPIO lines are not all pulled in the same direction.\n" +
			" gpbuttond makes the assumption that all lines are pulled up by default.\n" +
			" The process of matching this behavior varies between devices.\n\n" +
			" On a Raspberry Pi, this can be done by modifying your config.txt file.\n" +
			" For example, lines 5 and 6 can be pulled up by default with the following:\n" +
			"  gpio=5,6=pu\n\n" +
			"OTHER SUPPORTED CONFIGURATIONS (ENVIRONMENT VARIABLES)\n" +
			" GPBD_LONG can optionally be set to alter how long a button must be held for a long press to be registered. Set equal to any integer (measured in milliseconds).\n" +
			" GPBD_DEBOUNCE can optionally be set to apply a custom debounce time. Set equal to any integer (measured in milliseconds).\n" +
			" GPBD_REPEAT can optionally be set to alter the time before registering multiple keystrokes when a button is held down. Set equal to any integer (measured in milliseconds).\n\n")
		os.Exit(1)
	}

	for _, mapTrio := range mapEnv1Split {
		var intConvert [3]int
		for i, num := range strings.Split(mapTrio, ":") {
			intConvert[i], _ = strconv.Atoi(num)
		}
		lineMap = append(lineMap, [3]int{intConvert[0], intConvert[1], intConvert[2]})
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

	// long press timer configuration - local
	var mapEnv4, mapPresent4 = os.LookupEnv("GPBD_LONG")
	if mapPresent4 {
		longPressMultiplier, _ := strconv.Atoi(mapEnv4)
		longPressDuration = time.Duration(longPressMultiplier) * time.Millisecond
	} else {
		longPressDuration = 500 * time.Millisecond // default to 500-millisecond long press timer
	}

	// begin edge detection - will call eventHandler whenever a watched GPIO line changes state
	var gpioLines []*gpio.Line
	for i := 0; i < mapCount; i++ {
		gpioLine, _ := gpio.RequestLine("gpiochip0", lineMap[i][0], gpio.WithPullUp, gpio.WithBothEdges, gpio.WithDebounce(debounceDuration), gpio.WithEventHandler(eventHandler(lineMap[i][1], lineMap[i][2])))
		gpioLines = append(gpioLines, gpioLine)
	}

	// run FOREVER >:)
	for {
		fmt.Print("\nThis program is a background daemon - run it with \"gpbuttond &\" to properly background it.\n")
		time.Sleep(time.Hour)
	}
}
