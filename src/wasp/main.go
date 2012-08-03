package main

import (
	"log"
	"net"
	"net/http"
	"os"
)

//================================================================================

// The Mplayer struct we're about to use.
var mpl Mplayer

// The 'global' configuration properties.
var properties Properties

// The allowed video extensions.
var extensionsVideo AllowedExtensions

// The allowed audio extensions
var extensionsAudio AllowedExtensions

//================================================================================

// This function just simply outputs the IP addresses for debugging purposes, or
// to just quickly view on what IP WASP is being listening on. 
func logMachineIPAddresses() {
	_, _ = os.Hostname()
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Println("Unable to get a list of network interfaces")
	}

	log.Println("Available interfaces:")
	for _, ifs := range interfaces {
		addresses, err := ifs.Addrs()
		if err != nil {
			log.Printf("Unable to get IP address for interface %s", ifs.Name)
			continue
		}

		var strAddresses string = ""
		for i, addr := range addresses {
			strAddresses += addr.String()
			if i < len(addresses)-1 {
				strAddresses += ", "
			}
		}
		log.Printf("%d: %s (%s)", ifs.Index, ifs.Name, strAddresses)
	}
}

// This function initializes the application with all its properties. Will
// re-read configuration files if necessary.
func initialize() {
	log.Println("Initializing")

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
		log.Fatalf("Cannot forcefully create necessary FIFO: ", ferr)
	}

	extensionsVideo = make(AllowedExtensions)
	extensionsVideo.Parse(properties.GetString(PROPERTY_VIDEO_EXTS, ""))
	extensionsAudio = make(AllowedExtensions)
	extensionsAudio.Parse(properties.GetString(PROPERTY_AUDIO_EXTS, ""))

	log.Println("Initializiation done!")
}

// Entry point. Start it up.
func main() {
	log.Println("Wasp starting")

	logMachineIPAddresses()

	initialize()

	log.Printf("Media directory is %s", properties.GetString(PROPERTY_MEDIA_DIR, "/"))
	log.Printf("Starting to listen on '%s'", properties.GetString(PROPERTY_BIND_ADDRESS, ":8080"))
	log.Printf("Input FIFO filename is '%s'", properties.GetString(PROPERTY_MPLAYER_FIFO, "/tmp/mplayer.fifo"))
	log.Printf("Make sure MPlayer is configured to read its input from this FIFO!\n")

	handler := RegexHandler{}
	registerHttpHandlers(&handler)

	err := http.ListenAndServe(properties.GetString(PROPERTY_BIND_ADDRESS, ":8080"), &handler)
	if err != nil {
		log.Fatalf("Failed to bind to address '%s': %s", properties.GetString(PROPERTY_BIND_ADDRESS, ":8080"), err)
	}
}
