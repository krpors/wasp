// Properties, inspired by goproperties, Java Properties. Highly, highly 
// simplified, but meets the needs.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"strconv"
	"strings"
)

const (
	// The default configuration/properties file used to r/w from/to.
	PROPERTIES_FILE string = "/.wasp/wasp.properties"
)

//================================================================================

// Our custom type Properties, declared as a map of k/v string/string.
type Properties map[string]string

// Constants (keys) we're using as property names.
const (
	PROPERTY_MEDIA_DIR    string = "MediaDirectory"
	PROPERTY_BIND_ADDRESS string = "BindAddress"
	PROPERTY_MPLAYER_FIFO string = "MplayerFifo"
	PROPERTY_VIDEO_EXTS string = "VideoExtensions"
	PROPERTY_AUDIO_EXTS string = "AudioExtensions"
)

//================================================================================

// This function returns the default configuration filename, based on 
// the constant PROPERTIES_FILE. It attempts to get the current user's
// directory, and then concatenate it with the PROPERTIES_FILE value.
// The returned string is an absolute path to the file.
func (p *Properties) DefFileName() string {
	usr, uerr := user.Current()
	if uerr != nil {
		log.Fatalf("Unable to get the current user's directory: %s", uerr)
	}

	return usr.HomeDir + PROPERTIES_FILE
}

// Check whether the default properties file exists, by stat-ing it. The
// filename being tested is returned by DefFileName(). On error, return false
// otherwise return true.
func (p *Properties) FileExists() bool {
	_, err := os.Stat(p.DefFileName())
	if err != nil {
		return false
	}

	return true
}

// Sets default properties. Mainly used at first initialization of Wasp.
func (p *Properties) SetDefaults() {
	q := *p
	q[PROPERTY_MEDIA_DIR] = "/"
	q[PROPERTY_BIND_ADDRESS] = ":8080"
	q[PROPERTY_MPLAYER_FIFO] = "/tmp/mplayer.fifo"
	q[PROPERTY_VIDEO_EXTS] = ".mp4;.avi;.mpg;.mpeg;.wmv;.flv;.swf"
	q[PROPERTY_AUDIO_EXTS] = ".mp3;.ogg;.oga;.flac;.wav"

}

// Loads properties from a file.
func (p *Properties) Load(file string) (err error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	// split the file's contents. Each line by a \n
	lines := strings.Split(string(bytes), "\n")
	// Iterate over the lines, extract properties.
	for _, line := range lines {
		eqindex := strings.Index(line, "=")
		if eqindex > 0 {
			propname := line[0:eqindex]
			propval := line[eqindex+1:]
			(*p)[propname] = propval
		}
	}

	return
}

// Saves properties to file.
func (p *Properties) Save(file string) (err error) {
	str := ""
	for key, value := range *p {
		s := fmt.Sprintf("%s=%s\n", key, value)
		str += s
	}

	dir := path.Clean(path.Dir(file))
	// forcefully create directory. Does nothing if it already exists.
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, []byte(str), 0644)
}

// Gets a property as an unsigned integer, size 8 bits
func (p *Properties) GetUint8(propname string, def uint8) (val uint8) {
	propval := (*p)[propname]
	tempval, berr := strconv.ParseUint(propval, 10, 8)
	if berr != nil {
		return def
	}

	return uint8(tempval)
}

// Gets a value simply as a string. Just here to simply adhere to the rest of
// the "getters". One may also just do (propertiesInstance["propertyname"]) of
// course.
func (p *Properties) GetString(propname string, def string) (val string) {
	propval := (*p)[propname]
	if propval == "" {
		return def
	}

	return propval
}

// Gets a property as a boolean, or returns def when the property can't
// be found, or loaded. The parsing errors are discarded.
func (p *Properties) GetBool(propname string, def bool) (val bool) {
	propval := (*p)[propname]
	val, berr := strconv.ParseBool(propval)
	if berr != nil {
		return def
	}

	return val
}
