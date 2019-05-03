package app

import (
	"github.com/gyuhwankim/go-gin-starterkit/config"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingPong(t *testing.T) {
	client, teardown := setup()
	defer teardown()

	res, err := client.Get("/ping")
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Error occured: %s", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error occured: %s", err)
	}

	assert.Equal(t, body, "pong")
}

func setup() (*http.Client, func()) {
	conf := config.Configuration{
		Addr: ":8080",
	}

	appServer := New(conf)
	testServer := httptest.NewServer(appServer.core)

	return testServer.Client(), func() {
		testServer.Close()
	}
}
