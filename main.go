package main

import (
	"flag"
	"log"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
)

const (
	ColorRed    = "\033[31m"
	ColorReset  = "\033[0m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
)

func main() {

	var live bool
	var conf *Config

	usr, err := user.Current()
	if err != nil {
		log.Println(err)
	}

	// default config location
	configFile := filepath.Join(usr.HomeDir, ".factom-restart/config.yaml")

	// check if custom config location passed as flag
	flag.StringVar(&configFile, "c", configFile, "config.yaml path")

	// parse dry run flag
	flag.BoolVar(&live, "live", false, "live bool")

	flag.Parse()

	log.Printf("Using config: %s\n", configFile)

	// load config
	if conf, err = NewConfig(configFile); err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting restart system\n")
	log.Printf("Portainer endpoint: %s\n", conf.Endpoint)

	if live {
		log.Println(string(ColorYellow), "LIVE mode: network will be restarted!", string(ColorReset))
	} else {
		log.Println(string(ColorBlue), "DRY RUN mode: simulating restart", string(ColorReset))
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
			log.Println("------------------------------")
			log.Printf("Trying connecting to %s (%s), ID=%d\n", endpoints[i].Name, endpoints[i].PublicURL, endpoints[i].ID)
			endpoints[i].Containers, err = p.GetDockerContainers(endpoints[i].ID)
			if err != nil {
				log.Println(string(ColorRed), err, string(ColorReset))
			}
			for j := range endpoints[i].Containers {
				version := strings.Split(endpoints[i].Containers[j].Image, ":")
				if endpoints[i].Containers[j].State == "running" {
					log.Printf("factomd container is %s\n%s\n%s\n", endpoints[i].Containers[j].State, endpoints[i].Containers[j].ID, version[1])
				} else {
					log.Printf("factomd container is %s\n%s\n%s\n", endpoints[i].Containers[j].State, endpoints[i].Containers[j].ID, version[1])
				}
			}
		}
	}

	// RESTART START
	log.Println("------------------------------")
	if live {
		log.Println("RESTARTING NOW")
	} else {
		log.Println("SIMULATING RESTART (DRY-RUN)")
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
						log.Printf("OK %s (%s)\n", e.Name, c.ID)
					}
				}
			} else {
				// log skipped endpoints with no containers
				log.Printf("SKIP %s", e.Name)
			}

		}
	}

	wg.Wait()

}

func restart(p *Portainer, name string, endpointID int, containerID string) {

	err := p.RestartDockerContainer(endpointID, containerID)
	if err != nil {
		// if restart request failed, print error
		log.Printf("ERROR %s (%s)\n", name, containerID)
		log.Println(string(ColorRed), err, string(ColorReset))
	} else {
		// no error = success restart, print OK
		log.Printf("OK %s (%s)\n", name, containerID)
	}

}
