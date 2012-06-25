package main

import (
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

// Handles the /play URI. Formvalue 'file' is the relative filename to
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

func handlerPause(w http.ResponseWriter, r *http.Request) {
	log.Println("Toggling pause")

	err := mpl.Pause()
	if err != nil {
		log.Printf("Unable to pause Mplayer: %s", err)
	}
}

func handlerStop(w http.ResponseWriter, r *http.Request) {
	log.Println("Stopping playback")

	err := mpl.Stop()
	if err != nil {
		log.Printf("Unable to stop Mplayer: %s", err)
	}
}

func handlerVolume(w http.ResponseWriter, r *http.Request) {
	log.Println("Changing volume.")

	//log.Println("Content: ", r.FormValue("volume"))
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

func handlerMute(w http.ResponseWriter, r *http.Request) {
	log.Println("Muting")
	muting, err := strconv.ParseBool(r.FormValue("mute"))
	if err != nil {
		muting = false
	}

	mpl.Mute(muting)
}

func handlerSeek(w http.ResponseWriter, r *http.Request) {
	val, err := strconv.ParseInt(r.FormValue("seek"), 10, 16)
	if err != nil {
		log.Printf("Unable to parse integer for seeking: %s", err)
		return
	}

	log.Printf("Seeking relatively %d seconds", val)

	mpl.SeekRelative(int16(val))
}
