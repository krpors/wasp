package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
)

var proc *os.Process

func startHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hai thar!")

    if proc == nil {
        procattr := new(os.ProcAttr)
        var err error
        proc, err = os.StartProcess("/usr/bin/xev", nil, procattr)
        if err != nil {
            log.Fatal("Unable to start process")
        } else {
            log.Printf("XEV successfully started (pid %d)\n", proc.Pid)
        }
    } else {
    }
}

func stopHandler(w http.ResponseWriter, r *http.Request) {
    if proc != nil {
        log.Printf("Stopping process %d\n", proc.Pid)
        proc.Kill()
        proc.Release()
        proc = nil
    }
}

func main() {
/*
    log.Println("Wasp starting")

    http.HandleFunc("/start", startHandler)
    http.HandleFunc("/stop", stopHandler)

    log.Println("Listening on localhost:8080")
    http.ListenAndServe(":8080", nil)
*/
    m := Mplayer{"/home/krpors/test.fifo"}
    err := m.Loadfile("somefile.avi")
    if err != nil {
        log.Printf("Can't load some file: %s", err)
    }
}
