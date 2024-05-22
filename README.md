# gpbuttond 
gpbuttond is an efficient, lightweight daemon for mapping GPIO events (such as button presses) to keyboard keystrokes.

It is a simplistic alternative to [pikeyd](https://github.com/mmoller2k/pikeyd) that functions on more devices and on more modern kernel versions (Linux 5.10+ _**required**_).

As the daemon was designed for use on Raspberry Pi devices, it and the Go modules it imports are 100% native Go (no CGO) for easy cross-compilation.

The daemon can be configured by setting the `GPBD_MAP` environment variable. The format for this configuration will be provided upon running the program. Alongside this information, gpbuttond will provide other information pertaining to proper use of the program upon running without a configuration.

For launching gpbuttond, I recommend using OpenRC and placing the included `15gpbuttond.start` file in `/etc/local.d` (ensure the local service is part of either the boot or default runlevel). This method is not required and is purely a recommendation.

# Important!
gpbuttond relies on the `uinput` kernel module to create its virtual keyboard. Ensure this kernel module is loaded, otherwise gpbuttond _**will error out**_ on startup.

gpbuttond will need to be run as root on most systems.