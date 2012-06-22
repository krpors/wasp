package main

import (
    "log"
    "net/http"
    "os"
    "path/filepath"
)

// Registers URI handlers once.
func registerHttpHandlers() {
    // regular handlers:
    http.HandleFunc("/", handlerIndex)
    http.HandleFunc("/listing", handlerListing)
    http.HandleFunc("/index", handlerIndex)
    http.HandleFunc("/config", handlerConfig)

    http.HandleFunc("/test", handlerTest)
    http.HandleFunc("/play", handlerPlay)
    http.HandleFunc("/stop", handlerStop)
    http.HandleFunc("/pause", handlerPause)
    http.HandleFunc("/volume", handlerVolume)
    http.HandleFunc("/mute", handlerMute)

    // static (JS, CSS) content handler:
    pwd, err := os.Getwd()
    pwd = filepath.Join(pwd, "/site/")

    if err != nil {
        log.Fatalf("Unable to get current working directory: %s", err)
        return
    }

    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(pwd))))
}
