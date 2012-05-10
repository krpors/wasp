package main 

import (
    "errors"
)

// Accessors to Mplayer slave commands
type Mplayer struct {
    // FIFO path.
    Path string
}

// Loads a new file, from `file'. Returns an error when the file could
// not be found, or whatevs.
//
// Mplayer slave command: 
// loadfile <file|url>\n
func (m* Mplayer) Loadfile(file string) (err error) {
    err = errors.New("File not found")
    return
}

// Toggles sound on/off.
//
// Mplayer slave command: 
// mute\n
func (m* Mplayer) Mute() (err error) {
    err = errors.New("FIFO unavailable")
    return
}

// Toggles Pausing.
func (m* Mplayer) Pause() (err error) {
    err = errors.New("FIFO unavailable")
    return
}
