package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

// Handles the /ajax/play URI. Formvalue 'file' is the relative filename to
// be playing. The base directory is the config.MediaDir.
func handlerPlay(w http.ResponseWriter, r *http.Request) {
	file := r.FormValue("file")
	if file == "" {
		// TODO: return an error or something for the page to display.
		// Preferably in JSON?
		return
	}

	mediadir := properties.GetString(PROPERTY_MEDIA_DIR, "/")
	file = filepath.Join(mediadir, file)

	err := mpl.Loadfile(file)
	if err != nil {
		log.Printf("Unable to start file '%s'. Error is: %s", file, err)
		return
	}

	log.Printf("Playing '%s'", file)
}

// Handles the /ajax/pause URI. No form, POST or GET value is used, it's just simply
// pause, or unpause.
func handlerPause(w http.ResponseWriter, r *http.Request) {
	log.Println("Toggling pause")

	err := mpl.Pause()
	if err != nil {
		log.Printf("Unable to pause Mplayer: %s", err)
	}
}

// Handles the /ajax/stop URI. Stops the stream.
func handlerStop(w http.ResponseWriter, r *http.Request) {
	log.Println("Stopping playback")

	err := mpl.Stop()
	if err != nil {
		log.Printf("Unable to stop Mplayer: %s", err)
	}
}

// Handles the /ajax/volume URI. It reads the form POST value 'volume'.
func handlerVolume(w http.ResponseWriter, r *http.Request) {
	log.Println("Changing volume.")

	vol, err := strconv.ParseFloat(r.FormValue("volume"), 32)
	if err != nil {
		// if we fail to convert the volume, set it to 50.0
		vol = 50.0
	}

	// use a percentage as volume (it will be clamped automatically)
	log.Printf("Volume is %4.1f", Percentage(vol).Clamped())
	err = mpl.SetVolume(Percentage(vol))
	if err != nil {
		log.Printf("Volume changing failed: %s", err)
	}
}

// Handles the /ajax/mute URI. Mutes or unmutes the volume.
func handlerMute(w http.ResponseWriter, r *http.Request) {
	log.Println("Muting")
	muting, err := strconv.ParseBool(r.FormValue("mute"))
	if err != nil {
		muting = false
	}

	mpl.Mute(muting)
}

// Handles the /ajax/seek URI. Seeks in the current stream if applicable.
// The form value 'seek' is used, and will allow relative seeking (so no
// absolute position). We're unable to query the current position in a 
// normal way.
func handlerSeek(w http.ResponseWriter, r *http.Request) {
	val, err := strconv.ParseInt(r.FormValue("seek"), 10, 16)
	if err != nil {
		log.Printf("Unable to parse integer for seeking: %s", err)
		return
	}

	log.Printf("Seeking relatively %d seconds", val)

	mpl.SeekRelative(int16(val))
}

// Handles the /ajax/get_status URI. A JSON object with information of the 
// currently playing file, set volume, whether it's muted, and all that jazz.
func handlerGetStatus(w http.ResponseWriter, r *http.Request) {
	// this struct will be jsonified
	type Status struct {
		Muted      bool
		Volume     float32
		File       string
		Properties struct {
			MediaDirectory string
			BindAddress    string
			MplayerFifo    string
		}
	}

	s := Status{}
	s.Muted = mpl.Muted
	s.Volume = mpl.VolumeValue.Clamped()
	s.File = mpl.File
	s.Properties.MediaDirectory = properties.GetString(PROPERTY_MEDIA_DIR, "/")
	s.Properties.BindAddress = properties.GetString(PROPERTY_BIND_ADDRESS, ":8080")
	s.Properties.MplayerFifo = properties.GetString(PROPERTY_MPLAYER_FIFO, "/tmp/mplayer.fifo")

	bytes, err := json.Marshal(s)
	if err != nil {
		log.Printf("Unable to marshal struct to JSON data: %s", err)
		return
	}

	fmt.Fprintf(w, "%s\n", bytes)
}

// Handles the /ajax/get_dirlist URI. Returns a directory list using the 
// DirListData struct, except it's JSON marshaled.
func handlerGetDirList(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	requestPath := values.Get("p")

    dld, _ := getDirectoryList(requestPath) 
    bytes, err := json.Marshal(dld)
    if err != nil {
        log.Printf("Unable to marshal struct to JSON data: %s", err)
        return
    }

    fmt.Fprintf(w, "%s\n", bytes)
}
