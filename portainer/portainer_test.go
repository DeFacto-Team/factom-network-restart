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

	// expected success response
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

	assert := assert.New(t)

	// expected success response
	success := []SwarmEndpoint{{ID: 1, Name: "Test1", PublicURL: "1.1.1.1"}, {ID: 2, Name: "Test2", PublicURL: "2.2.2.2"}}
	resp, _ := json.Marshal(success)

	// start test auth server
	authResponse := &AuthResponse{JWT: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImp0aSI6ImNlN2EyY2U0LWMxYWYtNGYyMy1iNjM4LTczNjUyMDQ3ZDQyNSIsImlhdCI6MTU5OTg5MDc4MiwiZXhwIjoxNTk5ODk0MzgyfQ.kdFpPMkUYlTJa_K4h2bWTmCZYfLW9LTqXhlCZulcPxc"}
	authResponseByte, _ := json.Marshal(authResponse)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write(authResponseByte)
	}))

	defer srv.Close()

	// start portainer client
	p := NewPortainer("test", "test", srv.URL)

	// update test server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write(resp)
	}))

	defer srv.Close()

	endpoints, err := p.GetSwarmEndpoints()

	assert.Equal(err, assert.Nil)
	assert.Equal(endpoints, success)

}
