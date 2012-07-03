package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

// TODO: /ajax/properties    -> writes properties to file

type route struct {
	pattern *regexp.Regexp
	handler http.Handler
}

// RegexHandler is a custom Http handler. Allows us to use regexes as URI handlers,
// pretty much like Ruby on Rails, or Django.
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

// Registers URI handlers once. This makes the custom RegexHandler, which enables
// you to register regular expressions as URIs.
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
	handler.HandleFunc("^/ajax/get_dirlist", handlerGetDirList)

	// static (JS, CSS) content handler:
	pwd, err := os.Getwd()
	pwd = filepath.Join(pwd, "/site/")

	if err != nil {
		log.Fatalf("Unable to get current working directory: %s", err)
		return
	}

	handler.Handle("^/static/.*", http.StripPrefix("/static/", http.FileServer(http.Dir(pwd))))
}
