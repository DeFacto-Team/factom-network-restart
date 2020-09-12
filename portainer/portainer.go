package portainer

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

type Context struct {
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

type Portainer interface {
	GetSwarmEndpoints() ([]SwarmEndpoint, error)
	GetDockerContainers(endpointID int) ([]DockerContainer, error)
	RestartDockerContainer(endpointID int, containerID string) error
	GetToken() string
}

// NewPortainer initializes connection to Portainer API and obtain JWT access token
func NewPortainer(username string, password string, endpoint string) Portainer {

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
		log.Fatal("Portainer server doesn't return JWT token")
	}

	log.Info("Successfully logged in as " + username)

	return &Context{Token: token.JWT, Endpoint: endpoint}
}

// GetSwarmEndpoints makes request to Portainer API and returns []SwarmEndpoint
func (c *Context) GetSwarmEndpoints() ([]SwarmEndpoint, error) {

	var endpoints []SwarmEndpoint

	resp, err := c.makeRequest("GET", "/api/endpoints", nil)
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
func (c *Context) GetDockerContainers(endpointID int) ([]DockerContainer, error) {

	var containers []DockerContainer

	resp, err := c.makeRequest("GET", "/api/endpoints/"+strconv.Itoa(endpointID)+"/docker/containers/json?all=1&filters={\"label\":[\"name=factomd\"]}", nil)
	if err != nil || len(resp) == 0 {
		return nil, fmt.Errorf("Can not connect to remote host")
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
func (c *Context) RestartDockerContainer(endpointID int, containerID string) error {

	// nothing returned = success
	_, err := c.makeRequest("POST", "/api/endpoints/"+strconv.Itoa(endpointID)+"/docker/containers/"+containerID+"/restart?t=5", nil)
	if err != nil {
		return err
	}

	return nil

}

// Low level function making requests to Portainer API
func (c *Context) makeRequest(method string, path string, data []byte) ([]byte, error) {

	url := c.Endpoint + path

	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Token)

	client := &http.Client{Timeout: Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)

}

// Token() returns context token
func (c *Context) GetToken() string {
	return c.Token
}
