package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

/*
================================================================================

Template renderers:
    /                   -> root, index, main page
    /browse             -> generates browsable listing
    /config             -> views configuration

Ajax 'setters', i.e. don't return data
    /ajax/play          -> plays the file
    /ajax/stop          -> stops playback
    /ajax/pause         -> pauses playback
    /ajax/volume        -> sets volume
    /ajax/seek          -> seeks an x amount of seconds?
    /ajax/mute          -> mutes volume
    /ajax/properties    -> writes properties to file

Ajax 'getters', i.e. return data for async display.
    /ajax/get_list      -> returns a JSONified list of directories and files.
                           This is planned for a possible future Android/iOS app.
                           Probably like:
        {
            "directories": ["clips", "movies", "series"],
            "files": [ "hi.mpg", "test.mp4", "song.mp3" ]
        }

    /ajax/get_status    -> returns JSON data as follows?
        { 
            "muted": true,
            "volume": 50,
            "file": "/media/share/movie.mp4",
            "stopped": false,
            "properties" {
                "MediaDirectory": "/media/share",
                "BindAddress": ":8080",
                "MplayerFifo": "/tmp/mplayer.fifo"
            }
        }
================================================================================
*/

type route struct {
	pattern *regexp.Regexp
	handler http.Handler
}

// RegexHandler is a custom Http handler.
type RegexHandler struct {
	routeList []route
}

// Adds a http.Handler as handler instead of a function.
func (h *RegexHandler) Handle(pattern string, handler http.Handler) {
	rex, err := regexp.Compile(pattern)
	if err != nil {
		log.Printf("Unable to compile regex `%s'. This URI will not be available!\n")
	} else {
		log.Printf("Registering URI pattern `%s'\n", pattern)
		h.routeList = append(h.routeList, route{rex, handler})
	}
}

// Adds a function as a handler.
func (h *RegexHandler) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	rex, err := regexp.Compile(pattern)
	if err != nil {
		log.Printf("Unable to compile regex `%s'. This URI will not be available!\n")
	} else {
		log.Printf("Registering URI pattern `%s'\n", pattern)
		h.routeList = append(h.routeList, route{rex, http.HandlerFunc(handler)})
	}
}

// Interface method expected for type http.Handler. 
func (h *RegexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range h.routeList {
		if route.pattern.MatchString(r.URL.Path) {
			route.handler.ServeHTTP(w, r)
			return
		}
	}

	http.NotFound(w, r)
}

// Registers URI handlers once.
func registerHttpHandlers(handler *RegexHandler) {
	// regular handlers:
	handler.HandleFunc("^/$", handlerIndex)
	handler.HandleFunc("^/browse.*$", handlerBrowse)
	handler.HandleFunc("^/config$", handlerConfig)

	handler.HandleFunc("^/ajax/play", handlerPlay)
	handler.HandleFunc("^/ajax/stop", handlerStop)
	handler.HandleFunc("^/ajax/pause", handlerPause)
	handler.HandleFunc("^/ajax/volume", handlerVolume)
	handler.HandleFunc("^/ajax/mute", handlerMute)
	handler.HandleFunc("^/ajax/seek", handlerSeek)
	handler.HandleFunc("^/ajax/get_status", handlerGetStatus)

	// static (JS, CSS) content handler:
	pwd, err := os.Getwd()
	pwd = filepath.Join(pwd, "/site/")

	if err != nil {
		log.Fatalf("Unable to get current working directory: %s", err)
		return
	}

	handler.Handle("^/static/.*", http.StripPrefix("/static/", http.FileServer(http.Dir(pwd))))
}
