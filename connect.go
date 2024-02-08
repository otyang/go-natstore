package natstore

import (
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

// Connect establishes a connection to a NATS server using a specified list of server URLs.
// It accepts optional configuration options through the 'options' parameter.
//
// Parameters:
//   - serverURLs: A slice of NATS server URLs to connect to.
//   - connName: A unique name for the connection, used for identification.
//   - options: Additional NATS connection options (optional).
//
// Return values:
//   - natsConn: The established NATS connection object, or nil on error.
//   - error: Any error encountered during connection establishment, or nil on success.
//
// Notes:
//   - The provided server URLs should be valid NATS server addresses.
//   - The connection name should be unique to avoid conflicts.
//   - An error is returned if the connection fails or the returned connection object is nil.
//   - The provided options can be used to customize connection parameters.
func Connect(serverURL string, connName string, options ...nats.Option) (*nats.Conn, error) {
	options = append(options, nats.Name(connName))

	natsConn, err := nats.Connect(serverURL, options...)
	if err != nil || !natsConn.IsConnected() {
		return nil, fmt.Errorf("nats connection error: %w", err)
	}

	if natsConn == nil {
		return nil, fmt.Errorf("nats connection error: No nats connection object returned")
	}

	return natsConn, nil
}

// StartEmbeddedServer starts an embedded NATS server and waits for it to be ready for connections.
// It takes the server options as input and a timeout duration for waiting for the server to be ready.
// It returns the NATS server instance and any error encountered during the startup process.
//
// StartEmbeddedServer starts an embedded NATS Server with the provided options.
// It waits for the server to become ready for connections within a timeout period.
//
// Notes:
//   - The provided options should be pre-configured with the desired server settings.
//   - The timeout value determines the maximum time to wait for the server to be ready.
//
// Example usage:
//
//	opts := &server.Options{
//	   JetStream: true,
//	   Host: "localhost",
//	   Port: 4222,
//	}
//	timeout := 10 * time.Second
//
//	nserver, err := StartEmbeddedServer(opts, timeout)
//	if err != nil {
//	   // handle error
//	}
//
// defer nserver.Shutdown()
func StartEmbeddedServer(opts *server.Options, timeout time.Duration) (*server.Server, error) {
	nserver, err := server.NewServer(opts)
	if err != nil {
		return nil, fmt.Errorf("nats server startup error: %w", err)
	}

	if nserver == nil {
		return nil, errors.New("nats server error: no server object returned")
	}

	go nserver.Start()

	// Wait for server to be ready for connections
	if !nserver.ReadyForConnections(timeout) {
		return nil, errors.New("nats server error: not ready for connection")
	}

	return nserver, nil
}
