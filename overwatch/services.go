package overwatch

import (
  	log "github.com/Sirupsen/logrus"
	// "syscall"
	"syscall"
	"path"
)

// runningServices holds the pids of the services
var runningServices map[string]int

// Service describes how should a service be run, what its name is and what it does TODO: add respawn and something like runlevels? maybe dependencies
type Service struct {
	Name string
	Description string		
}

// ConfigureServices stars the services, registers the SIGCHILD handler
func ConfigureServices(c *Config){
	log.WithFields(log.Fields{
		"title": c.Title,
		"version": c.Version,
	}).Info("Starting services")
		
	for _, serv := range c.Services {
		log.WithFields(log.Fields{
			"service": serv.Name,
		}).Info("Starting service")
		spawn(serv.Name)
	}
}

// spawn spawns the given filename in background using syscall
func spawn(name string){
	filepath := path.Join("/init/services", name)
	pid, err := syscall.ForkExec(filepath, nil, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"service": filepath,
			"error": err,
		}).Error("Could not start service.")	
	} else {
		log.WithFields(log.Fields{
			"service": filepath,
			"pid": pid,
		}).Info("Started service succesfully")			
	}
}