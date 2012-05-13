package main

import (
    "fmt"
    "log"
    "net/http"
    "html/template"

    "wasp/conf"
)

// The Mplayer we're about to use.
var mplayer Mplayer

func indexHandler(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("./templates/index.html")
    t.Execute(w, nil)
}

func startHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, ":D")
    err := mplayer.Loadfile("/media/share/video/clips/fey.mp4")
    if err != nil {
        fmt.Fprintf(w, "Fifo couldn't be stat")
    }
}

func stopHandler(w http.ResponseWriter, r *http.Request) {
    mplayer.Stop()
}

// Entry point. Start it up.
func main() {
    mplayer = Mplayer{"/tmp/mplayer.fifo"}

    log.Println("Wasp starting")

    http.HandleFunc("/", indexHandler)
    http.HandleFunc("/start", startHandler)
    http.HandleFunc("/stop", stopHandler)

    log.Println("Starting to listen on :8080")
    /*
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("Failed to open server socket: ", err)
    }
    */

    config, conferr := conf.Load()
    if conferr != nil {
        fmt.Println(conferr)
        conf.SaveDefaults()
    }

    fmt.Printf("Media directory is %s\n", config.MediaDir)
}
