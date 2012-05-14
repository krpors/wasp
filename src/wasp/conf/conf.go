package conf

import(
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "os/user"
)

const CONFIG_DIR string = "/.wasp"
const CONFIG_FILE string = "/.wasp/config.json"

type ConfigError struct {
    // Our description
    Desc string
}

func (c *ConfigError) Error() string {
    return c.Desc
}

// The Config struct denotes Wasp configuration. Decoded as JSON.
type Config struct {
    // Base media directory
    MediaDir string
    // Webserver's bind address. Default is ":8080"
    BindAddress string
    // Mplayer FIFO location. Default is /tmp/mplayer.fifo
    MplayerFifo string
}

func FileName() (filename string, err error) {
    usr, uerr := user.Current()
    if uerr != nil {
        return "", &ConfigError{ "Unable to get current user information" }
    }

    return usr.HomeDir + CONFIG_FILE, nil
}

// Checks if the file exists and all that jazz.
func Exists() bool {
    usr, uerr := user.Current()
    if uerr != nil {
        return false
    }

    fullfile := usr.HomeDir + CONFIG_FILE
    _, derr := os.Open(fullfile)
    if derr != nil {
        return false
    }

    return true
}

// This function loads the configuration from the current user's 
// ${HOME}/.wasp/config.json file. If it does not exist, return an
// error.
func Load() (conf Config, err error) {
    file, cerr := FileName()
    if cerr != nil {
        return Config{}, cerr
    }

    bytes, err := ioutil.ReadFile(file)
    if err != nil {
        errorDesc := fmt.Sprintf("Unable to read file contents from %s", file)
        return Config{}, &ConfigError { errorDesc }
    }

    err = json.Unmarshal(bytes, &conf)
    if err != nil {
        errorDesc := fmt.Sprintf("Unable to read JSON configuration: %s", err)
        return Config{}, &ConfigError { errorDesc }
    }

    return conf, nil
}

// Saves the given configuration to the current user's configuration file.
// The config file is saved in a subdirectory "~/.wasp/". If this directory
// does not exist yet, create it. 
func Save(conf *Config) error {
    usr, err := user.Current()
    if err != nil {
        return err
    }

    dir := usr.HomeDir + CONFIG_DIR
    file := usr.HomeDir + CONFIG_FILE

    // forcefully create directory. Does nothing if it already exists.
    err = os.MkdirAll(dir, 0755)
    if err != nil {
        errorDesc := fmt.Sprintf("Unable to create directory '%s': %s", dir, err.Error())
        return &ConfigError{ errorDesc }
    }

    var b []byte
    b, err = json.Marshal(conf)
    if err != nil {
        errorDesc := fmt.Sprintf("Unable to marshal configuration to JSON: %s", err.Error())
        return &ConfigError{ errorDesc }
    }

    err = ioutil.WriteFile(file, b, 0644)
    if err != nil {
        errorDesc := fmt.Sprintf("Unable to save configuration to '%s': %s", file, err.Error())
        return &ConfigError{ errorDesc }
    }

    // A-okay!
    return nil
}

func SaveDefaults() error {
    config := Config{}
    config.MediaDir = "/home/media/"
    config.BindAddress = ":8080"
    config.MplayerFifo = "/tmp/mplayer.fifo"
    return Save(&config)
}
