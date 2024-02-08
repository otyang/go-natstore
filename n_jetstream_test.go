package natstore

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
)

func testStartEmbeddedServer(t *testing.T) *server.Server {
	server, err := StartEmbeddedServer(&server.Options{
		JetStream:  true,
		ServerName: "streamCreateUpdate",
		DontListen: true,
	}, time.Second*3)

	assert.NoError(t, err)
	assert.NotNil(t, server)
	return server
}

func TestJSClient_StreamCreateOrUpdate(t *testing.T) {
	server := testStartEmbeddedServer(t)
	defer server.Shutdown()

	nconn, err := Connect("", "", nats.InProcessServer(server))
	assert.NoError(t, err)
	assert.NotNil(t, nconn)

	var (
		ctx        = context.Background()
		streamName = "stream_name"
	)

	s, err := NewJSClient(nconn)
	assert.NoError(t, err)

	err = s.StreamCreateOrUpdate(ctx, jetstream.StreamConfig{Name: streamName, Subjects: []string{streamName + ".>"}})
	assert.NoError(t, err)

	streamObject, err := s.NatsJS().Stream(ctx, streamName)
	assert.NoError(t, err)
	assert.NotEmpty(t, streamObject)

	err = s.NatsJS().DeleteStream(ctx, streamName)
	assert.NoError(t, err)
}

func TestJSClient_Publish(t *testing.T) {
	server := testStartEmbeddedServer(t)
	defer server.Shutdown()

	nc, err := Connect("", "", nats.InProcessServer(server))
	assert.NoError(t, err)
	assert.NotNil(t, nc)

	var (
		ctx        = context.TODO()
		streamName = "stream_name"
		subject    = "MAST.page_loaded"
	)

	s, err := NewJSClient(nc)
	assert.NoError(t, err)

	err = s.StreamCreateOrUpdate(ctx, jetstream.StreamConfig{Name: streamName, Subjects: []string{subject}})
	assert.NoError(t, err)

	err = s.Publish(ctx, "MAST.page_loaded", "hello")
	assert.NoError(t, err)
}

func TestJSClient_Consume(t *testing.T) {
	server := testStartEmbeddedServer(t)
	defer server.Shutdown()

	nc, err := Connect("", "", nats.InProcessServer(server))
	assert.NoError(t, err)
	assert.NotNil(t, nc)

	var (
		ctx          = context.TODO()
		streamName   = "stream_name"
		subject      = "stream_name.subject"
		consumerName = "consumername"
		handler      = func(m jetstream.Msg) {
			fmt.Println(m.Subject())
		}
	)

	s, err := NewJSClient(nc)
	assert.NoError(t, err)

	err = s.StreamCreateOrUpdate(ctx, jetstream.StreamConfig{Name: streamName, Subjects: []string{subject}})
	assert.NoError(t, err)

	err = s.Publish(context.Background(), subject, "hello")
	assert.NoError(t, err)

	err = s.Consume(ctx, streamName, subject, handler, WithDurableName(consumerName))
	assert.NoError(t, err)
}
