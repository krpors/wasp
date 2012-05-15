package main

import (
    "fmt"
    "container/list"
    "log"
    "net/http"
    "os"
    "html/template"
    "path"
    "strconv"

    "wasp/conf"
)

// The Mplayer we're about to use.
var mplayer Mplayer
// The 'global' configuration.
var config conf.Config

func indexHandler(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("./templates/index.html")
    t.Execute(w, nil)

    log.Println("Index handlr")
}

func startHandler(w http.ResponseWriter, r *http.Request) {
    err := mplayer.Loadfile("/home/krpors/fey.mp4")
    if err != nil {
        fmt.Fprintf(w, "Fifo couldn't be stat")
    }
}

func pauseHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Toggling pause")

    err := mplayer.Pause()
    if err != nil {
        log.Println("Couldn't pause")
    }
}

func stopHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Stopping playback")
    mplayer.Stop()
}

func volumeHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Should handle volume.")
    //log.Println("Content: ", r.FormValue("volume"))
    vol, err := strconv.ParseFloat(r.FormValue("volume"), 32)
    if err != nil {
        // if we fail to convert the volume, set it to 50.0
        vol = 50.0
    }

    // use a percentage as volume (it will be clamped automatically)
    log.Printf("Volume is %4.1f", Percentage(vol).Clamped())
    mplayer.Volume(Percentage(vol), true)
}

func listingHandler(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("./templates/listing.html")

    // 1. request path from URI.
    // 2. concat it with config.MediaDir + p=...
    // 3. disallow using ".." in the path
    // 4. fetch dirlist from path
    // 5. build list with paths
    // 6. add to template execution

    values := r.URL.Query()
    requestPath := values.Get("p")

    log.Printf("Requesting path '%s'", requestPath)

    dir, err := os.Open(path.Join(config.MediaDir, requestPath))
    if err != nil {
        log.Println("Can't open directory")
        return
    }

    fileinfos, err := dir.Readdir(0)
    if err != nil {
        log.Println("Can't list directory")
        return
    }

    // list holding files
    dirList := list.New()
    fileList := list.New()
    for _, fi := range(fileinfos) {
        // ignore 'hidden' directories/files, starting with a dot.
        if fi.Name()[0] == '.' {
            continue
        }

        // TODO: only add media files. So probably a set of allowed extensions.
        // Match them case insensitive. If dir, or media file, add them. Needs
        // sorting too on directories first, then files, and on alphabetical order.

        if fi.IsDir() {
            dirList.PushBack(fi.Name())
        } else {
            fileList.PushBack(fi.Name())
        }

    }

    dirs := make([]string, dirList.Len())
    files := make([]string, fileList.Len())
    i := 0
    for e := dirList.Front(); e != nil; e = e.Next() {
        dirs[i] = e.Value.(string)
        i++
    }

    i = 0
    for e := fileList.Front(); e != nil; e = e.Next() {
        files[i] = e.Value.(string)
        i++
    }

    type ListingData struct {
        ParentDir string
        Directories []string
        Files []string
    }

    data := ListingData{requestPath, dirs, files}

    t.Execute(w, data)
}

func registerHandlers() {
    http.HandleFunc("/listing", listingHandler)
    http.HandleFunc("/start", startHandler)
    http.HandleFunc("/stop", stopHandler)
    http.HandleFunc("/pause", pauseHandler)
    http.HandleFunc("/volume", volumeHandler)
    http.HandleFunc("/index", indexHandler)
}

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
    mplayer = Mplayer{config.MplayerFifo}

    log.Printf("Media directory is %s", config.MediaDir)

    log.Printf("Starting to listen on '%s'", config.BindAddress)

    registerHandlers()
    err = http.ListenAndServe(config.BindAddress, nil)
    if err != nil {
        log.Fatalf("Failed to bind to address '%s': %s", config.BindAddress, err)
    }

}
