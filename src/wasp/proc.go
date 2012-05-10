package main 

import (
    "errors"
    "os"
    "fmt"
    "io/ioutil"
)

// Accessors to Mplayer slave commands
type Mplayer struct {
    // FIFO path.
    PathFifo string
}

func (m* Mplayer) sendCommand(cmd string) (err error) {
    err = m.FifoOk()
    if err != nil {
        return err
    }

    // write the string to the fifo
    err = ioutil.WriteFile(m.PathFifo, []byte(cmd), 0644)
    return err
}

func (m* Mplayer) FifoOk() (err error) {
    fileinfo, err := os.Stat(m.PathFifo)
    if err != nil {
        return err
    }

    filemode := fileinfo.Mode()
    if filemode & os.ModeNamedPipe != os.ModeNamedPipe {
        desc := fmt.Sprintf("%s is not a named pipe (FIFO)", m.PathFifo)
        err = errors.New(desc)
        return err
    }

    return nil
}

// Loads a new file, from `file'. Returns an error when the named pipe
// could not be found, or written to.
//
// Mplayer slave command: 
// loadfile <file|url>\n
func (m* Mplayer) Loadfile(file string) (err error) {
    return m.sendCommand(fmt.Sprintf("loadfile %s\n", file))
}

// Toggles sound on/off.
//
// Mplayer slave command: 
// mute\n
func (m* Mplayer) Mute() (err error) {
    return m.sendCommand(fmt.Sprintf("mute\n"))
}

// Toggles Pausing.
func (m* Mplayer) Pause() (err error) {
    return m.sendCommand(fmt.Sprintf("pause\n"))
}
