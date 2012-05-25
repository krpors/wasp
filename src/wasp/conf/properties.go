// Properties, inspired by goproperties, Java Properties. Highly, highly 
// simplified, but meets the needs.

package conf

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
    P_MEDIA_DIR string = "MediaDirectory"
    P_BIND_ADDRESS string = "BindAddress"
    P_MPLAYER_FIFO string = "MplayerFifo"
)

//================================================================================

// Gets the filename as a string, using the home directory.
func DefaultFileName() (filename string) {
    usr, uerr := user.Current()
    if uerr != nil {
        log.Fatalf("Unable to get the current user's directory: %s", uerr)
    }

    return usr.HomeDir + PROPERTIES_FILE
}

func FileExists(filename string) bool {
    _, err := os.Stat(filename)
    if err != nil {
        return false
    }

    return true
}

//================================================================================

// Sets default properties
func (p *Properties) SetDefaults() {
    q := *p
    q[P_MEDIA_DIR] = "/"
    q[P_BIND_ADDRESS] = ":8080"
    q[P_MPLAYER_FIFO] = "/tmp/mplayer.fifo"
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
    for _, line := range(lines) {
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
    for key, value := range(*p) {
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
