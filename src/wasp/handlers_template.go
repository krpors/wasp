package main

import (
	"html/template"
	"log"
	"net/http"
)

var templateIndex = template.Must(template.ParseFiles("./site/templates/index.html"))
var templateBrowse = template.Must(template.ParseFiles("./site/templates/browse.html"))
var templateConfig = template.Must(template.ParseFiles("./site/templates/config.html"))

//==============================================================================

func handlerIndex(w http.ResponseWriter, r *http.Request) {
	log.Printf("Index requested\n")
	templateIndex.Execute(w, nil)
}

//================================================================================

func handlerBrowse(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	requestPath := values.Get("p")

    dld, _ := getDirectoryList(requestPath)
    templateBrowse.Execute(w, dld)
}

//================================================================================

func handlerConfig(w http.ResponseWriter, r *http.Request) {
	type ConfigData struct {
		MediaDir    string
		MplayerFifo string
		BindAddress string
	}

	cd := ConfigData{}
	cd.MediaDir = properties.GetString(PROPERTY_MEDIA_DIR, "/")
	cd.BindAddress = properties.GetString(PROPERTY_BIND_ADDRESS, ":8080")
	cd.MplayerFifo = properties.GetString(PROPERTY_MPLAYER_FIFO, "/tmp/mplayer.fifo")

	templateConfig.Execute(w, cd)
}
