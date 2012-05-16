package main

import (
    "fmt"
    "container/list"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "sort"
    "html/template"
    "path"
    "strconv"

    "wasp/conf"
    "wasp/mplayer"
)

// The Mplayer we're about to use.
var mpl mplayer.Mplayer

// The 'global' configuration.
var config conf.Config

//==============================================================================

func indexHandler(w http.ResponseWriter, r *http.Request) {
    t, err := template.ParseFiles("./site/templates/index.html")
    if err != nil {
        log.Fatalf("Can not parse template: %s", err)
        return
    }

    t.Execute(w, nil)

    log.Println("Index handlr")
}

func startHandler(w http.ResponseWriter, r *http.Request) {
    err := mpl.Loadfile("/home/krpors/fey.mp4")
    if err != nil {
        fmt.Fprintf(w, "Fifo couldn't be stat")
    }
}

func pauseHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Toggling pause")

    err := mpl.Pause()
    if err != nil {
        log.Printf("Unable to pause Mplayer: %s", err)
    }
}

func stopHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Stopping playback")

    err := mpl.Stop()
    if err != nil {
        log.Printf("Unable to stop Mplayer: %s", err)
    }
}

func volumeHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Changing volume.")

    //log.Println("Content: ", r.FormValue("volume"))
    vol, err := strconv.ParseFloat(r.FormValue("volume"), 32)
    if err != nil {
        // if we fail to convert the volume, set it to 50.0
        vol = 50.0
    }

    // use a percentage as volume (it will be clamped automatically)
    log.Printf("Volume is %4.1f", mplayer.Percentage(vol).Clamped())
    err = mpl.Volume(mplayer.Percentage(vol))
    if err != nil {
        log.Printf("Volume changing failed: %s", err)
    }
}

func muteHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Muting!?")
    muting, err := strconv.ParseBool(r.FormValue("mute"))
    if err != nil {
        muting = false
    }

    mpl.Mute(muting)
}

// The listing handler generates a list of directories and files
// which can be clicked on to browse with. The request path is 
// given in the http.Request using the parameter name `p'.
func listingHandler(w http.ResponseWriter, r *http.Request) {
    t, err := template.ParseFiles("./site/templates/listing.html")
    if err != nil {
        log.Fatalf("Can not parse template: %s", err)
    }

    values := r.URL.Query()
    requestPath := values.Get("p")

    log.Printf("Requesting path '%s'", requestPath)

    // Get a directory listing of the selected directory. First, concat
    // the media directory with the request path so we have an absolute path.
    dir, err := os.Open(path.Join(config.MediaDir, requestPath))
    if err != nil {
        log.Printf("Can't open directory: %s", err)
        return
    }

    // Fetch the actual file information slice.
    fileinfos, err := dir.Readdir(0)
    if err != nil {
        log.Printf("Can't list directory: %s", err)
        return
    }

    // Generate a (doubly linked) list with files.
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

    // Copy them back to a slice.
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

    sort.Strings(dirs)
    sort.Strings(files)

    // 'Temporary' struct to use for the template
    type ListingData struct {
        ParentDir string
        RequestPath string
        Directories []string
        Files []string
    }

    // Create a struct with content. path.Dir() gets the parent directory, and
    // is used to navigate to back up one directory. The requestPath is used
    // to browse to a new directory. The dirs and files slices contains the 
    // directories and the files respectively. TODO: sort these, alphabetically.
    data := ListingData{
        // The parent directory, so we can go back.
        path.Clean(path.Dir(requestPath)),
        // The requested, current path
        path.Clean(requestPath),
        // The directories in requestPath
        dirs,
        // The files in requestPath
        files,
    }

    // Execute the template, write outcome to `w'.
    t.Execute(w, data)
}

type Lol string

func registerHandlers() {
    // regular handlers:
    http.HandleFunc("/listing", listingHandler)
    http.HandleFunc("/start", startHandler)
    http.HandleFunc("/stop", stopHandler)
    http.HandleFunc("/pause", pauseHandler)
    http.HandleFunc("/volume", volumeHandler)
    http.HandleFunc("/mute", muteHandler)
    http.HandleFunc("/index", indexHandler)

    // static (JS, CSS) content handler:
    pwd, err := os.Getwd()
    pwd = filepath.Join(pwd, "/site/")
    log.Printf("It's %s", pwd)

    if err != nil {
        log.Fatalf("Unable to get current working directory: %s", err)
        return
    }

    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(pwd))))
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
