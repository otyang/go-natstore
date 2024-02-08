package natstore

import (
	"errors"
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

func startUpTestServer() (*server.Server, error) {
	nserver, err := server.NewServer(&server.Options{
		ServerName: "unit_test_nat",
		DontListen: true, // Set this if you don't want a listening TCP socket
	})
	if err != nil {
		return nil, errors.New("starting nats server:" + err.Error())
	}

	go nserver.Start()

	if !nserver.ReadyForConnections(time.Second * 5) {
		return nil, errors.New("NATS server didn't start")
	}

	return nserver, nil
}

func TestConnect(t *testing.T) {
	server, err := startUpTestServer()
	assert.NoError(t, err)
	assert.NotNil(t, server)

	defer server.Shutdown()

	nconn, err := Connect("", "name:unit_test_nat", nats.InProcessServer(server))
	assert.NoError(t, err)
	assert.NotNil(t, nconn)

	err = nconn.Publish("foo", []byte("bar"))
	assert.NoError(t, err)
}

func TestStartEmbeddedServer(t *testing.T) {
	server, err := StartEmbeddedServer(&server.Options{
		ServerName: "unit_test_nat",
		DontListen: true, // Set this if you don't want a listening TCP socket
	}, time.Second*2)

	assert.NoError(t, err)
	assert.NotNil(t, server)

	defer server.Shutdown()

	nconn, err := Connect("", "name:unit_test_nat", nats.InProcessServer(server))
	assert.NoError(t, err)
	assert.NotNil(t, nconn)

	err = nconn.Publish("foo", []byte("bar"))
	assert.NoError(t, err)
}
