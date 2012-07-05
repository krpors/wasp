package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"syscall"
)

// Type Percentage is a float32 which can be used to set certain values in
// mplayer. They have certain boundaries, so they must be clamped to either
// 0.0 or 100.0, and not exceed these limits.
type Percentage float32

// Returns a float32 in the boundary [0, 100].
func (p Percentage) Clamped() float32 {
	if p < 0.0 {
		return 0.0
	}

	if p > 100.0 {
		return 100.0
	}

	return float32(p)
}

// Accessors to Mplayer slave commands
type Mplayer struct {
	// FIFO path.
	PathFifo string

	// Volume previously set. Used to set every new playable item to
	// the previously set volume.
	VolumeValue Percentage

	// The last/currently played file
	File string

	// Muted or not
	Muted bool
}

// Generic interface to send a command to the Mplayer FIFO.
func (m *Mplayer) sendCommand(cmd string) (err error) {
	err = m.FifoOk()
	if err != nil {
		return err
	}

	// write the string to the fifo
	err = ioutil.WriteFile(m.PathFifo, []byte(cmd), 0644)
	return err
}

func (m *Mplayer) FifoOk() (err error) {
	if m.PathFifo == "" {
		return errors.New("FIFO path is empty")
	}

	fileinfo, err := os.Stat(m.PathFifo)
	if err != nil {
		return err
	}

	filemode := fileinfo.Mode()
	if filemode&os.ModeNamedPipe != os.ModeNamedPipe {
		desc := fmt.Sprintf("%s is not a named pipe (FIFO)", m.PathFifo)
		err = errors.New(desc)
		return err
	}

	return nil
}

// Attempts to create the FIFO automatically. Can return an error if it failed,
// for instance if permission failures occur (no permission to write to path).
func (m *Mplayer) CreateFifo() (err error) {
	err = m.FifoOk()
	if err != nil {
		err = syscall.Mkfifo(m.PathFifo, 0644)
		if err != nil {
			log.Printf("Could not create FIFO: `%s'", m.PathFifo)
			return err
		} else {
			log.Printf("FIFO created: `%s'", m.PathFifo)
			return nil
		}
	}

	return err
}

// Loads a new file, from `file'. Returns an error when the named pipe
// could not be found, or written to.
//
// Mplayer slave command: loadfile <file|url>\n
func (m *Mplayer) Loadfile(file string) (err error) {
	m.File = file
	return m.sendCommand(fmt.Sprintf("loadfile \"%s\"\n", file))
}

// Toggles sound on/off.
//
// Mplayer slave command: mute\n
func (m *Mplayer) Mute(mute bool) (err error) {
	m.Muted = mute
	var muteness string
	if mute {
		muteness = "1"
	} else {
		muteness = "0"
	}

	return m.sendCommand(fmt.Sprintf("mute %s\n", muteness))
}

// Toggles Pausing.
//
// Mplayer slave command: pause\n
func (m *Mplayer) Pause() (err error) {
	return m.sendCommand(fmt.Sprintf("pause\n"))
}

// Stops playback.
//
// Mplayer slave command: stop\n
func (m *Mplayer) Stop() (err error) {
	return m.sendCommand(fmt.Sprintf("stop\n"))
}

// Seeks in the current file in a relative manner. When seconds is negative,
// seek -seconds. If seconds is positive, seek +seconds in the current stream.
//
// Mplayer slave command: seek <+/-value> 0\n
func (m *Mplayer) SeekRelative(seconds int16) (err error) {
	return m.sendCommand(fmt.Sprintf("seek %d 0\n", seconds))
}

// Seeks in the current file in an absolute manner, using percentages.
//
// Mplayer slave command: seek <value> 1\n
func (m *Mplayer) SeekPercentage(value Percentage) (err error) {
	return m.sendCommand(fmt.Sprintf("seek %5.1f 1\n", value.Clamped()))
}

// Sets volume in an absolute manner. I.e. 0 = silent, 100 = full volum.
// In mplayer (well, mplayer1, not sure about mplayer2), the volume can only 
// be succesfully changed when a file is loaded.
//
// Mplayer slave command: volume <value> [abs]
func (m *Mplayer) SetVolume(value Percentage) (err error) {
	m.VolumeValue = value
	return m.sendCommand(fmt.Sprintf("volume %f 1\n", value.Clamped()))
}

// Writes an OSD (on-screen-display) string. The text cannot contain double
// quotes. If double quotes exist, they are escaped by using a single quote.
// Duration is in milliseconds.
//
// Mplayer slave command: osd_showtext "<string>" [duration]
func (m *Mplayer) OsdShowText(text string, duration uint16) (err error) {
	newtext := strings.Replace(text, "\"", "'", -1)
	return m.sendCommand(fmt.Sprintf("osd_show_text \"%s\" %d\n", newtext, duration))
}

// Displays the current playing file on OSD. The property is fetched from internal
// mplayer source.
func (m *Mplayer) OsdShowFilename(duration uint16) (err error) {
	return m.sendCommand(fmt.Sprintf("osd_show_property_text ${filename} %d\n", duration))
}
