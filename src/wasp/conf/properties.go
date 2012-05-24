package conf

import (
    "fmt"
    "io/ioutil"
    "strconv"
    "strings"
)

// The base directory (/ = user's home directory)
const PROPERTIES_DIR string = "/.wasp"
// The actual configuration/properties file used to r/w from/to.
const PROPERTIES_FILE string = "/.wasp/wasp.properties"

// Our custom type Properties, declared as a map of k/v string/string.
type Properties map[string]string

// Loads properties from a file.
func (p *Properties) Load(file string) (err error) {
    bytes, err := ioutil.ReadFile(file)
    if err != nil {
        fmt.Printf("Could not load file: %s", err)
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
func (p Properties) Save(file string) (err error) {
    return nil
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
