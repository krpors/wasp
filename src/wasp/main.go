package main

import (
    "log"
    "net/http"

    "wasp/conf"
    "wasp/mplayer"
)

// The Mplayer we're about to use.
var mpl mplayer.Mplayer

// The 'global' configuration.
var config conf.Config

// Entry point. Start it up.
func main() {
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

    var conferr error
    config, conferr = conf.Load()
    if conferr != nil {
        log.Println("Unable to load configuration: ", conferr)
    }

    log.Printf("Using fifo path %s. Make sure Mplayer uses this same named pipe.", config.MplayerFifo)
    mpl = mplayer.Mplayer{}
    mpl.PathFifo = config.MplayerFifo

    log.Printf("Media directory is %s", config.MediaDir)

    log.Printf("Starting to listen on '%s'", config.BindAddress)

    registerHandlers()
    err = http.ListenAndServe(config.BindAddress, nil)
    if err != nil {
        log.Fatalf("Failed to bind to address '%s': %s", config.BindAddress, err)
    }

}
