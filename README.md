# gpbuttond 
gpbuttond is an efficient, lightweight daemon for mapping GPIO events (such as button presses) to keyboard keystrokes.

As the daemon was designed for use on Raspberry Pi devices, it and the Go modules it imports are 100% native Go (no CGO) for easy cross-compilation.

Currently, the daemon is hard-coded for the GPIO setup of a specific Pi hat I am using, but I plan on changing it so that it accepts pin mappings from command-line arguments or environment variables. The mapping is currently done by changing the values of the variables at the very top of the single .go file.

Once you have adjusted the mappings for your use-case, simply compile the single binary with `go mod tidy && go build gpbutton.go`, place it in your $PATH, and create a way to launch it. I recommend using OpenRC and placing the included "15gpbuttond.start" file in /etc/local.d (ensure the local service is part of either the boot or default runlevel).
