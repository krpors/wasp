package main

import (
    "log"
    "net/http"

    "wasp/conf"
    "wasp/mplayer"
)

//================================================================================

// The Mplayer we're about to use.
var mpl mplayer.Mplayer

// The 'global' configuration properties
var properties conf.Properties

//================================================================================


// Entry point. Start it up.
func main() {
    log.Println("Wasp starting")

    propFile := conf.DefaultFileName()

    properties = make(conf.Properties)
    // if default filename exists, read from that
    if conf.FileExists(propFile) {
        log.Printf("Loading properties from '%s'", propFile)
        properties.Load(propFile)

    // else set default properties, write them back
    } else {
        log.Printf("'%s' does not exist yet, setting default properties...", propFile)
        properties.SetDefaults()
        properties.Save(propFile)
    }

    mpl = mplayer.Mplayer{}
    mpl.PathFifo = properties.GetString("MplayerFifo", "/tmp/mplayer.fifo")

    log.Printf("Media directory is %s", properties.GetString("MediaDirectory", "/"))
    log.Printf("Starting to listen on '%s'", properties.GetString("BindAddress", ":8080"))

    registerHandlers()
    err := http.ListenAndServe(properties.GetString("BindAddress", ":8080"), nil)
    if err != nil {
        log.Fatalf("Failed to bind to address '%s': %s", properties.GetString("BindAddress", ":8080"), err)
    }
}
