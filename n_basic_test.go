package natstore

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

func setupNatsServerConnectAndEstablishClient(t *testing.T) (*server.Server, *nats.Conn) {
	server, err := startUpTestServer()
	if err != nil {
		log.Fatal(err)
	}

	nconn, err := Connect("", "name:unit_test_nat", nats.InProcessServer(server))
	if err != nil {
		log.Fatal(err)
	}

	return server, nconn
}

func TestBasic_Request(t *testing.T) {
	server, nc := setupNatsServerConnectAndEstablishClient(t)
	defer server.Shutdown()

	basicStore := NewBasic(nc)

	requestPayload := "==4444=="
	responsePayload := "hello, hahaha"

	basicStore.Subscribe("greet.*", func(msg *nats.Msg) {
		// request
		var req string
		err := json.Unmarshal(msg.Data, &req)
		assert.NoError(t, err)
		assert.Equal(t, requestPayload, req)

		// response
		responsePayloadBytes, err := json.Marshal(responsePayload)
		assert.NoError(t, err)
		err = msg.Respond(responsePayloadBytes)
		assert.NoError(t, err)
	})

	//== request
	var resp string
	err := basicStore.Request(context.TODO(), "greet.joe", requestPayload, &resp)
	assert.NoError(t, err)
	assert.Equal(t, responsePayload, resp)
}

func TestBasic_Publish(t *testing.T) {
	server, nc := setupNatsServerConnectAndEstablishClient(t)
	defer server.Shutdown()

	basic := NewBasic(nc)

	basic.Subscribe("basic.*", func(msg *nats.Msg) {
		// subject name
		name := msg.Subject[6:]

		// Request Data
		var b string
		err := json.Unmarshal(msg.Data, &b)
		assert.NoError(t, err)

		fmt.Println("from subscribe handler: ", name, b)
	})

	// == publish
	err := basic.Publish(context.TODO(), "basic.morning", "trying basic publish")
	assert.NoError(t, err)
}

func TestBasic_Subscribe(t *testing.T) {
	TestBasic_Publish(t)
	TestBasic_Request(t)
}
