package overwatch

import (
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/go-ps"
	"github.com/shirou/gopsutil/host"
)

// runningServices holds the pids of the services
var runningServices map[string]int

// Service describes how should a service be run, what its name is and what it does TODO: add respawn and something like runlevels? maybe dependencies
type Service struct {
	Name        string
	Description string
}

// setupSignalHandler catches SIGCHILD for now only.
func setupSignalHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt,
		syscall.SIGCHLD,
		syscall.SIGTERM,
		syscall.SIGKILL)

	go func() {
		for sig := range c {
			log.Printf("captured %v", sig)
		}
	}()
}

func printStatus() {
	h, err := host.HostInfo()
	if err != nil {
		log.Println("Could not get host info", err)
	}
	log.Println(h.String())

	procs, err := ps.Processes()
	if err != nil {
		log.Println("Could not get processes")
	}

	log.Println(len(procs))

	for _, e := range procs {
		log.Println(e.Executable(), e.Pid(), e.PPid())
	}
}

// ConfigureServices stars the services, registers the SIGCHILD handler
func ConfigureServices(c *Config) {
	log.WithFields(log.Fields{
		"title":   c.Title,
		"version": c.Version,
	}).Info("Starting services")

	setupSignalHandler()

	for _, serv := range c.Services {
		log.WithFields(log.Fields{
			"service": serv.Name,
		}).Info("Starting service")
		spawn(serv.Name)
	}

	// printStatus()
}

// spawn spawns the given filename in background using syscall
func spawn(name string) {
	filepath := path.Join("/init/services", name)

	stdinpath := path.Join("/logs/", name+".stdin")
	stdoutpath := path.Join("/logs/", name+".stdout")
	stderrpath := path.Join("/logs/", name+".stderr")

	os.MkdirAll(path.Dir(stdinpath), 0777)
	os.MkdirAll(path.Dir(stdoutpath), 0777)
	os.MkdirAll(path.Dir(stderrpath), 0777)

	fstdin, err := os.Create(stdinpath)
	if err != nil {
		log.Println("waat", err)
	}
	fstdout, err := os.Create(stdoutpath)
	if err != nil {
		log.Println("waat", err)
	}
	fstderr, err := os.Create(stderrpath)
	if err != nil {
		log.Println("waat", err)
	}

	// Open Files for stdout, stderr
	procAttr := &syscall.ProcAttr{
		Dir:   "/",
		Env:   []string{"MYVAR=345"},
		Files: []uintptr{fstdin.Fd(), fstdout.Fd(), fstderr.Fd()},
		Sys:   nil,
	}

	pid, err := syscall.ForkExec(filepath, nil, procAttr)
	if err != nil {
		log.WithFields(log.Fields{
			"service": filepath,
			"error":   err,
		}).Error("Could not start service.")
	} else {
		log.WithFields(log.Fields{
			"service": filepath,
			"pid":     pid,
		}).Info("Started service succesfully")
	}

	log.Info("Waiting for 3 seconds")
	time.Sleep(3 * time.Second)

	a, err1 := ioutil.ReadFile(stdoutpath)
	b, err2 := ioutil.ReadFile(stderrpath)
	if err1 != nil || err2 != nil {
		log.Error("Could not read", err1, err2)
	} else {
		log.WithFields(log.Fields{
			"service": name,
			"stdout":  string(a),
			"stderr":  string(b),
		}).Info("Service ended.")
	}
}
