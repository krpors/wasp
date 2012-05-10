package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
)

func startHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Epic lulz0r!")
}

func stopHandler(w http.ResponseWriter, r *http.Request) {
}

func main() {
    /*
    log.Println("Wasp starting")

    http.HandleFunc("/start", startHandler)
    http.HandleFunc("/stop", stopHandler)

    log.Println("Listening on localhost:8080")
    http.ListenAndServe("192.168.1.100:8080", nil)
    */

    m := Mplayer{"/tmp/mplayer.fifo"}
    err := m.FifoOk()
    if err != nil {
        log.Printf("Fifo is not ok: %s\n", err)
        os.Exit(1)
    }

    //m.Loadfile("/home/krpors/downloads/850.avi")
}
