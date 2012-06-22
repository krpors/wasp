package main

import (
    "log"
    "net/http"
    "os"
    "path/filepath"
    "regexp"
)

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
func (h *RegexHandler) ServeHTTP (w http.ResponseWriter, r *http.Request) {
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
    handler.HandleFunc("^/listing.*$", handlerListing)
    handler.HandleFunc("^/config$", handlerConfig)

    handler.HandleFunc("^/test", handlerTest)
    handler.HandleFunc("^/play", handlerPlay)
    handler.HandleFunc("^/stop", handlerStop)
    handler.HandleFunc("^/pause", handlerPause)
    handler.HandleFunc("^/volume", handlerVolume)
    handler.HandleFunc("^/mute", handlerMute)
    
    // static (JS, CSS) content handler:
    pwd, err := os.Getwd()
    pwd = filepath.Join(pwd, "/site/")

    if err != nil {
        log.Fatalf("Unable to get current working directory: %s", err)
        return
    }

    handler.Handle("/static/.*", http.StripPrefix("/static/", http.FileServer(http.Dir(pwd))))
}
