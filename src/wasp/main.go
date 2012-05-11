package main

import (
    "fmt"
    "log"
    "net/http"
)

// The Mplayer we're about to use.
var mplayer Mplayer

func startHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Epic lulz0r!")
    err := mplayer.Loadfile("/home/krpors/fey.mp4")
    if err != nil {
        log.Fatal(err)
    }
}

func stopHandler(w http.ResponseWriter, r *http.Request) {
    mplayer.Stop()
}

func main() {
    mplayer = Mplayer{"/tmp/mplayer.fifo"}

    log.Println("Wasp starting")

    http.HandleFunc("/start", startHandler)
    http.HandleFunc("/stop", stopHandler)

    log.Println("Listening on localhost:8080")
    http.ListenAndServe(":8080", nil)
}
