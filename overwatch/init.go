package overwatch


import (
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

// Config is general init config file to be parsed from disk upon boot
type Config struct {
	Version	string
	Title 	string
	Services []Service
}

func readConfig(filepath string) *Config{
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Panicln("Could not read crucial config file for init", err)
	}
	
	config := &Config{}
	
    err = yaml.Unmarshal(bytes, config)
	if err != nil {
		log.Panicln("Could not parse config file!", err)
	}
	
	return config
}


// Bootstrap starts the environment and reads the config file for services
func Bootstrap(){
    log.Println("The gonit is booting. Go is now in control.")
		
	// PrepareDisk()
	config := readConfig("config.yml")
	ConfigureServices(config)
		
	for {
		
	}
	// Not expected to reach here, then you get a kernel panic.	
}