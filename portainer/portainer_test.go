package portainer

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPortainer(t *testing.T) {

	assert := assert.New(t)

	// expected response
	success := &AuthResponse{JWT: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImp0aSI6ImNlN2EyY2U0LWMxYWYtNGYyMy1iNjM4LTczNjUyMDQ3ZDQyNSIsImlhdCI6MTU5OTg5MDc4MiwiZXhwIjoxNTk5ODk0MzgyfQ.kdFpPMkUYlTJa_K4h2bWTmCZYfLW9LTqXhlCZulcPxc"}
	resp, _ := json.Marshal(success)

	// start test server
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write(resp)
	}))

	defer srv.Close()

	// start portainer client
	p := NewPortainer("test", "test", srv.URL)

	assert.Equal(p.GetToken(), success.JWT)

}

func TestGetSwarmEndpoints(t *testing.T) {

	var resp []byte

	assert := assert.New(t)

	// expected responses
	endpointsResp := []SwarmEndpoint{{ID: 1, Name: "Test1", PublicURL: "1.1.1.1"}, {ID: 2, Name: "Test2", PublicURL: "2.2.2.2"}}
	authResp := &AuthResponse{JWT: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImp0aSI6ImNlN2EyY2U0LWMxYWYtNGYyMy1iNjM4LTczNjUyMDQ3ZDQyNSIsImlhdCI6MTU5OTg5MDc4MiwiZXhwIjoxNTk5ODk0MzgyfQ.kdFpPMkUYlTJa_K4h2bWTmCZYfLW9LTqXhlCZulcPxc"}

	// setup test server
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case "/api/endpoints":
			resp, _ = json.Marshal(endpointsResp)
		case "/api/auth":
			resp, _ = json.Marshal(authResp)
		}
		w.WriteHeader(200)
		w.Write(resp)
	}))

	defer srv.Close()

	// start portainer client
	p := NewPortainer("test", "test", srv.URL)

	endpoints, err := p.GetSwarmEndpoints()

	assert.Equal(err, nil)
	assert.Equal(endpoints, endpointsResp)

}

func TestGetDockerContainers(t *testing.T) {

	var resp []byte

	assert := assert.New(t)

	// expected responses
	containersResp := []DockerContainer{{ID: "c55e61d8a89b363c450bef7163cc46ea44479ed2245c2c3fde5e3c818e7dd0ef", Image: "v6.6.0", State: "running"}}
	authResp := &AuthResponse{JWT: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImp0aSI6ImNlN2EyY2U0LWMxYWYtNGYyMy1iNjM4LTczNjUyMDQ3ZDQyNSIsImlhdCI6MTU5OTg5MDc4MiwiZXhwIjoxNTk5ODk0MzgyfQ.kdFpPMkUYlTJa_K4h2bWTmCZYfLW9LTqXhlCZulcPxc"}

	// setup test server
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case "/api/endpoints/1/docker/containers/json?all=1&filters={\"label\":[\"name=factomd\"]}":
			resp, _ = json.Marshal(containersResp)
		case "/api/auth":
			resp, _ = json.Marshal(authResp)
		}
		w.WriteHeader(200)
		w.Write(resp)
	}))

	defer srv.Close()

	// start portainer client
	p := NewPortainer("test", "test", srv.URL)

	containers, err := p.GetDockerContainers(1)

	assert.Equal(err, nil)
	assert.Equal(containers, containersResp)

}

func TestRestartDockerContainer(t *testing.T) {

	var resp []byte

	assert := assert.New(t)

	// expected responses
	authResp := &AuthResponse{JWT: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImp0aSI6ImNlN2EyY2U0LWMxYWYtNGYyMy1iNjM4LTczNjUyMDQ3ZDQyNSIsImlhdCI6MTU5OTg5MDc4MiwiZXhwIjoxNTk5ODk0MzgyfQ.kdFpPMkUYlTJa_K4h2bWTmCZYfLW9LTqXhlCZulcPxc"}

	// setup test server
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case "/api/endpoints/1/docker/containers/c55e61d8a89b363c450bef7163cc46ea44479ed2245c2c3fde5e3c818e7dd0ef/restart?t=5":
			resp = make([]byte, 0)
		case "/api/auth":
			resp, _ = json.Marshal(authResp)
		}
		w.WriteHeader(200)
		w.Write(resp)
	}))

	defer srv.Close()

	// start portainer client
	p := NewPortainer("test", "test", srv.URL)

	err := p.RestartDockerContainer(1, "c55e61d8a89b363c450bef7163cc46ea44479ed2245c2c3fde5e3c818e7dd0ef")

	assert.Equal(err, nil)

}
