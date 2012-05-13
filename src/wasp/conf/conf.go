package conf

import(
    "encoding/json"
    "errors"
    "fmt"
    "io/ioutil"
    "os"
    "os/user"
)

const CONFIG_DIR string = "/.wasp"
const CONFIG_FILE string = "/.wasp/config.json"

type ConfigError struct {
    Desc string
}

func (c *ConfigError) Error() string {
    return fmt.Sprintf("Configuration error: %s", c.Desc)
}

// The Config struct denotes Wasp configuration. Decoded as JSON.
type Config struct {
    MediaDir string
}

// Util function to get the configuration 
func getConfigFile() (file* os.File, err error) {
    usr, uerr := user.Current()
    if uerr != nil {
        return nil, errors.New("1")
    }

    fullfile := usr.HomeDir + CONFIG_FILE
    f, derr := os.Open(fullfile)
    if derr != nil {
        return nil, errors.New("2")
    }

    return f, nil
}

// This function loads the configuration from the current user's 
// ${HOME}/.wasp/config.json file. If it does not exist, return an
// error.
func Load() (conf Config, err error) {
    file, cerr := getConfigFile()
    if cerr != nil {
        return Config{}, errors.New("3")
    }


    // TODO: unmarshal the json configuration
    fmt.Printf("Filename to unmarshal: %s\n", file.Name())
    var c Config
    bytes, err := ioutil.ReadFile(file.Name())
    if err != nil {

    }
    jerr := json.Unmarshal(bytes, &c)
    if jerr != nil {
        fmt.Printf("Unable to unmarshall file %s", file.Name())
    }

    conf = c
    err = nil
    return
}

// Saves the given configuration to the current user's configuration file.
// The config file is saved in a subdirectory "~/.wasp/". If this directory
// does not exist yet, create it. 
func Save(conf Config) error {
    usr, uerr := user.Current()
    if uerr != nil {
        return uerr
    }

    // forcefully create directory. Does nothing if it already exists.
    derr := os.MkdirAll(usr.HomeDir + CONFIG_DIR, 0755)
    if derr != nil {
        return derr
    }

    // TODO: marshal the config.
    b, err := json.Marshal(conf)
    if err != nil {
        return err
    }

    return ioutil.WriteFile(usr.HomeDir + CONFIG_FILE, b, 0644)
}

func SaveDefaults() error {
    // TODO: Config struct with default values, store them through Save()
    config := Config{}
    config.MediaDir = "/home/media/"
    return Save(config)
}
