package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

// HTTP request timeout
const Timeout = 5 * time.Second

type Portainer struct {
	Endpoint string
	Token    string
}

// AuthRequest is authorization form struct for Portainer API
type AuthRequest struct {
	Username string
	Password string
}

// AuthResponse is response format for authorization in Portainer API
type AuthResponse struct {
	JWT string
	Err string
}

// SwarmEndpoint is response format for GET /endpoints call in Portainer API
type SwarmEndpoint struct {
	ID         int
	Name       string
	PublicURL  string
	Containers []DockerContainer
}

// DockerContainer is response format for GET /endpoints/{endpointId}/docker/containers call in Portainer API
type DockerContainer struct {
	ID    string
	Image string
	State string
}

// NewPortainer initializes connection to Portainer API and obtain JWT access token
func NewPortainer(username string, password string, endpoint string) *Portainer {

	url := endpoint + "/api/auth"

	login := AuthRequest{Username: username, Password: password}
	data, err := json.Marshal(login)

	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: Timeout}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	token := &AuthResponse{}

	err = json.Unmarshal(body, token)
	if err != nil {
		log.Error(err)
	}

	if token.Err != "" {
		log.Fatal(token.Err)
	}
	if token.JWT == "" {
		log.Fatal("Portainer server didn't return JWT token")
	}

	log.Info("Successfully logged in as " + username)

	return &Portainer{Token: token.JWT, Endpoint: endpoint}
}

// GetSwarmEndpoints makes request to Portainer API and returns []SwarmEndpoint
func (p *Portainer) GetSwarmEndpoints() ([]SwarmEndpoint, error) {

	var endpoints []SwarmEndpoint

	resp, err := p.makeRequest("GET", "/api/endpoints", nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &endpoints)
	if err != nil {
		return nil, err
	}

	return endpoints, nil

}

// GetDockerContainers makes request to Portainer API and returns []DockerContainer with label "name=factomd" for requested endpointId
func (p *Portainer) GetDockerContainers(endpointID int) ([]DockerContainer, error) {

	var containers []DockerContainer

	resp, err := p.makeRequest("GET", "/api/endpoints/"+strconv.Itoa(endpointID)+"/docker/containers/json?all=1&filters={\"label\":[\"name=factomd\"]}", nil)
	if err != nil {
		return nil, err
	}
	if len(resp) == 0 {
		return nil, fmt.Errorf("Empty response received from Portainer API")
	}

	err = json.Unmarshal(resp, &containers)
	if err != nil {
		return nil, err
	}

	if len(containers) == 0 {
		return nil, fmt.Errorf("No containers with label 'name=factomd'")
	}

	return containers, nil

}

// RestartDockerContainer makes restart request to Portainer API for requested endpointId and containerId
func (p *Portainer) RestartDockerContainer(endpointID int, containerID string) error {

	// nothing returned = success
	_, err := p.makeRequest("POST", "/api/endpoints/"+strconv.Itoa(endpointID)+"/docker/containers/"+containerID+"/restart?t=5", nil)
	if err != nil {
		return err
	}

	return nil

}

// Low level function making requests to Portainer API
func (p *Portainer) makeRequest(method string, path string, data []byte) ([]byte, error) {

	url := p.Endpoint + path

	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.Token)

	client := &http.Client{Timeout: Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)

}
