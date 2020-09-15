package main

import (
	"flag"
	"fmt"
	"os/user"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

func main() {

	var err error
	var live bool
	var conf *Config

	log.SetLevel(5)

	usr, err := user.Current()
	if err != nil {
		log.Error(err)
	}

	// default config location
	configFile := usr.HomeDir + "/.factom-restart/config.yaml"

	// check if custom config location passed as flag
	flag.StringVar(&configFile, "c", configFile, "config.yaml path")

	// parse dry run flag
	flag.BoolVar(&live, "live", false, "live bool")

	flag.Parse()

	log.Info("Using config: ", configFile)

	// load config
	if conf, err = NewConfig(configFile); err != nil {
		log.Fatal(err)
	}

	log.Info("Starting restart system")
	log.Info("Portainer endpoint: ", conf.Endpoint)

	if live {
		log.Warn("LIVE mode: network will be restarted!")
	} else {
		log.Info("DRY RUN mode: simulating restart")
	}

	// initialize portainer
	p := NewPortainer(conf.Username, conf.Password, conf.Endpoint)

	// get all endpoints from portainer
	endpoints, err := p.GetSwarmEndpoints()
	if err != nil {
		log.Fatal(err)
	}

	// get factomd containers for each endpoint
	for i := range endpoints {
		if endpoints[i].PublicURL != "" {
			fmt.Printf("------------------------------\n")
			log.Debug("Trying connecting to ", endpoints[i].Name, " (", endpoints[i].PublicURL, "), ID=", endpoints[i].ID)
			endpoints[i].Containers, err = p.GetDockerContainers(endpoints[i].ID)
			if err != nil {
				log.Error(err)
			}
			for j := range endpoints[i].Containers {
				version := strings.Split(endpoints[i].Containers[j].Image, ":")
				if endpoints[i].Containers[j].State == "running" {
					log.Debug("factomd container is ", endpoints[i].Containers[j].State, "\n", endpoints[i].Containers[j].ID, "\n", version[1])
				} else {
					log.Warn("factomd container is ", endpoints[i].Containers[j].State, "\n", endpoints[i].Containers[j].ID, "\n", version[1])
				}
			}
		}
	}

	// RESTART START
	fmt.Printf("------------------------------\n")
	if live {
		log.Info("RESTARTING NOW")
	} else {
		log.Info("SIMULATING RESTART (DRY-RUN)")
	}

	var wg sync.WaitGroup

	// Restart each factomd container
	for _, e := range endpoints {
		if e.PublicURL != "" {

			if len(e.Containers) > 0 {
				// for all online hosts with containers, restart factomd container(s)
				for _, c := range e.Containers {
					if live {
						// if live mode, async restarting
						wg.Add(1)
						go func(name string, endpointID int, containerID string) {
							restart(p, name, endpointID, containerID)
							wg.Done()
						}(e.Name, e.ID, c.ID)
					} else {
						// if dry-run mode, just print endpoint and containerId
						log.Info("OK ", e.Name, " (", c.ID, ")")
					}
				}
			} else {
				// log skipped endpoints with no containers
				log.Warn("SKIP ", e.Name)
			}

		}
	}

	wg.Wait()

}

func restart(p Portainer, name string, endpointID int, containerID string) {

	err := p.RestartDockerContainer(endpointID, containerID)
	if err != nil {
		// if restart request failed, print error
		log.Error("ERROR ", name, " (", containerID, ")")
		log.Error(err)
	} else {
		// no error = success restart, print OK
		log.Info("OK ", name, " (", containerID, ")")
	}

}
