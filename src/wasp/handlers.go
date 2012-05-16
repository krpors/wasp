package main

import (
    "container/list"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "sort"
    "html/template"
    "path"
    "strconv"

    "wasp/mplayer"
)

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

func testHandler(w http.ResponseWriter, r *http.Request) {
    err := mpl.Loadfile("/home/krpors/fey.mp4")
    if err != nil {
        log.Printf("Unable to start media: %s", err)
    }
}

func playHandler(w http.ResponseWriter, r *http.Request) {
    file := r.FormValue("file")
    if file == "" {
        // TODO: return an error or something for the page to display.
        // Preferably in JSON?
        return
    }

    file = filepath.Join(config.MediaDir, file)

    err := mpl.Loadfile(file)
    if err != nil {
        log.Printf("Unable to start file '%s'. Error is: %s", file, err)
    }
    log.Printf("Playing '%s'", file)
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
    // 'Temporary' struct to use for the template
    type ListingData struct {
        ParentDir string        // parent directory
        RequestPath string      // requested path
        Directories []string    // slice of directories
        Files []string          // slice of files
        Error string            // possible error. May be empty.
    }

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
        // This might happen if we aren't allowed to open a directory
        // due to permission issues.
        log.Printf("Unable to open directory: %s", err)
        data := ListingData{}
        data.ParentDir = path.Clean(path.Dir(requestPath))
        data.Error = "Contents could not be listed."
        t.Execute(w, data)
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
        // The error. Empty, cus OK!
        "",
    }

    // Execute the template, write outcome to `w'.
    t.Execute(w, data)
}

func registerHandlers() {
    // regular handlers:
    http.HandleFunc("/listing", listingHandler)
    http.HandleFunc("/test", testHandler)
    http.HandleFunc("/play", playHandler)
    http.HandleFunc("/stop", stopHandler)
    http.HandleFunc("/pause", pauseHandler)
    http.HandleFunc("/volume", volumeHandler)
    http.HandleFunc("/mute", muteHandler)
    http.HandleFunc("/index", indexHandler)

    // static (JS, CSS) content handler:
    pwd, err := os.Getwd()
    pwd = filepath.Join(pwd, "/site/")

    if err != nil {
        log.Fatalf("Unable to get current working directory: %s", err)
        return
    }

    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(pwd))))
}
