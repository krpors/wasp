package main

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

// Struct to generate a directory listing
type DirListData struct {
	ParentDir   string   // parent directory
	RequestPath string   // requested path
	Directories []string // slice of directories
	Files       []string // slice of files
	Error       string   // possible error. May be empty.
}

// list of allowed extensions for the listing handler to display.
var allowedExtensions = map[string]bool{
	".mp4":  true,
	".avi":  true,
	".mp3":  true,
	".ogg":  true,
	".oga":  true,
	".mpg":  true,
	".wmv":  true,
	".flv":  true,
	".swf":  true,
	".vob":  true,
	".flac": true,
	".mpeg": true,
}

// Simple type with allowed extensions.
type AllowedExtensions map[string]bool

// Parses a string in the form of ".ext;.otherext;.bla". Splits based on
// semicolon, and each entry is added to the map value.
func (a *AllowedExtensions) Parse(s string) {
	log.Printf("Parsing extensions: %s", s)
	for _, ext := range strings.Split(s, ";") {
		(*a)[ext] = true
	}
}

// This function gets a directory listing of the requested path. The requested
// path is relative with regard to the media directory location (from the 
// properties instance). The function will return a struct instance of 
// DirListdata. For example, if the media directory property is set to "/home/user"
// and the requested path is "/opt/media", this function will attempt to list the
// directory contents of "/home/user/opt/media".
//
// The files and directories in the struct will be alphabetically sorted.
//
// Upon error (non-existant directory, unreadable, whatevs), an error will be
// returned, along with an empty DirListData{}.

func getDirectoryList(requestPath string) (DirListData, error) {
	// our dir list data instance. We're going to fill and return this one in
	// this function.
	dld := DirListData{}
	dld.ParentDir = path.Clean(path.Dir(requestPath))
	dld.RequestPath = path.Clean(requestPath)

	mediadir := properties.GetString(PROPERTY_MEDIA_DIR, "/")
	// Get a directory listing of the selected directory. First, concat
	// the media directory with the request path so we have an absolute path.
	fullpath := path.Join(mediadir, requestPath)
	dir, err := os.Open(fullpath)
	if err != nil {
		// This might happen if we aren't allowed to open a directory
		// due to permission issues.
		log.Printf("Unable to open directory: %s", err)
		dld.Error = "Contents could not be listed."
		return dld, err
	}

	log.Printf("Listing directory `%s'", fullpath)

	// Fetch the actual file information slice.
	fileinfos, err := dir.Readdir(0)
	if err != nil {
		log.Printf("Can't list directory: %s", err)
		return dld, err
	}

	for _, fi := range fileinfos {
		// ignore 'hidden' directories/files, starting with a dot.
		if fi.Name()[0] == '.' {
			continue
		}

		if fi.IsDir() {
			//dirList.PushBack(fi.Name())
			dld.Directories = append(dld.Directories, fi.Name())
		} else {
			// only allow certain kind of extensions.
			extension := filepath.Ext(fi.Name())
			if extensionsVideo[extension] || extensionsAudio[extension] {
				dld.Files = append(dld.Files, fi.Name())
			}
		}
	}

	sort.Strings(dld.Directories)
	sort.Strings(dld.Files)

	return dld, nil
}
