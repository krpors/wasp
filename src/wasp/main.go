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
    err := mplayer.Loadfile("/home/krpors/fey.mp4")
    if err != nil {
        fmt.Fprintf(w, "Fifo couldn't be stat")
    }
}

func stopHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Stopping...")
    mplayer.Stop()
}

// Entry point. Start it up.
func main() {
    mplayer = Mplayer{"/tmp/mplayer.fifo"}

    log.Println("Wasp starting")

    filename, err := conf.FileName()
    if err != nil {
        // If we can't get the current user's directory, just stop for now.
        // We need a place to put our configuration in, and that should be
        // in the executing user's home dir. For now :) This should probably
        // a command line option in the future? I.e. `./wasp -c /opt/config.json'
        log.Fatal(err)
    }

    if !conf.Exists() {
        log.Printf("Creating default configuration file '%s'", filename)
        conf.SaveDefaults()
    }

    log.Printf("Loading configuration from '%s'", filename)

    config, conferr := conf.Load()
    if conferr != nil {
        log.Println("Unable to load configuration: ", conferr)
    }

    http.HandleFunc("/", indexHandler)
    http.HandleFunc("/start", startHandler)
    http.HandleFunc("/stop", stopHandler)

    log.Printf("Media directory is %s", config.MediaDir)

    log.Printf("Starting to listen on '%s'", config.BindAddress)

    err = http.ListenAndServe(config.BindAddress, nil)
    if err != nil {
        log.Fatalf("Failed to bind to address '%s': %s", config.BindAddress, err)
    }

}
