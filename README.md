# gpbuttond 
gpbuttond is an efficient, lightweight daemon for mapping GPIO events (such as button pressed) to keyboard keystrokes.

It and the Go modules it imports are 100% native Go (no CGO) for easy cross-compilation, as the daemon was designed for use on Raspberry Pi devices.

Simply compile the single binary with `go mod tidy && go build`, place it in your $PATH, and create a way to launch it. I recommend using OpenRC and placing the included "15gpbuttond.start" file in /etc/local.d (ensure the local service is part of either the boot or default runlevel).
