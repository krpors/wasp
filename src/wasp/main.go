package main

import (
    "log"
    "net/http"
)

//================================================================================

// The Mplayer struct we're about to use.
var mpl Mplayer

// The 'global' configuration properties
var properties Properties

//================================================================================


// Entry point. Start it up.
func main() {
    log.Println("Wasp starting")

    properties = make(Properties)
    propFile := properties.DefFileName()

    // if default filename exists, read from that
    if properties.FileExists() {
        log.Printf("Loading properties from '%s'", propFile)
        properties.Load(propFile)

    // else set default properties, write them back
    } else {
        log.Printf("'%s' does not exist yet, setting default properties...", propFile)
        properties.SetDefaults()
        properties.Save(propFile)
    }

    mpl = Mplayer{}
    mpl.PathFifo = properties.GetString(PROPERTY_MPLAYER_FIFO, "/tmp/mplayer.fifo")
    ferr := mpl.CreateFifo()
    if ferr != nil {
        log.Fatalf("Cannot ", ferr)
    }

    log.Printf("Media directory is %s", properties.GetString(PROPERTY_MEDIA_DIR, "/"))
    log.Printf("Starting to listen on '%s'", properties.GetString(PROPERTY_BIND_ADDRESS, ":8080"))
    log.Printf("Input FIFO filename is '%s'", properties.GetString(PROPERTY_MPLAYER_FIFO, "/tmp/mplayer.fifo"))
    log.Printf("Make sure MPlayer is configured to read its input from this FIFO!\n")

    registerHttpHandlers()
    err := http.ListenAndServe(properties.GetString(PROPERTY_BIND_ADDRESS, ":8080"), nil)
    if err != nil {
        log.Fatalf("Failed to bind to address '%s': %s", properties.GetString(PROPERTY_BIND_ADDRESS, ":8080"), err)
    }
}
