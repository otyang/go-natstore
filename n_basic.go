package natstore

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
)

// DefaultRequestTimeout defines the default duration to wait for a response in request operations.
var DefaultRequestTimeout = 10 * time.Millisecond

// Basic provides fundamental NATS messaging capabilities for publishing, requesting, and subscribing.
type Basic struct {
	nconn  *nats.Conn
	logger *slog.Logger
}

// NewBasic creates a new Basic instance with the given NATS connection.
func NewBasic(natsConn *nats.Conn) *Basic {
	return &Basic{
		nconn:  natsConn,
		logger: defaultLogger,
	}
}

// WithLogger sets a custom logger for the Basic instance.
func (b *Basic) WithLogger(logger *slog.Logger) *Basic {
	b.logger = logger
	return b
}

// NConn exposes the underlying NATS connection for direct access when needed.
func (b *Basic) NConn() *nats.Conn {
	return b.nconn
}

// Publish sends a message to the specified subject with the given data.
func (b *Basic) Publish(ctx context.Context, subject string, data any) error {
	// Validate connections before proceeding:
	if err := b.checkConnections(); err != nil {
		return handleError("publish", err, b.logger)
	}

	// Marshal data into JSON format:
	dataByte, err := json.Marshal(data)
	if err != nil {
		return handleError("publish", err, b.logger)
	}

	// Publish the message using the NATS connection:
	err = b.nconn.Publish(subject, dataByte)
	if err != nil {
		return handleError("publish", err, b.logger)
	}

	return nil
}

// Request sends a request message to the specified subject and waits for a response.
func (b *Basic) Request(ctx context.Context, subj string, requestData any, responseDataPtr any) error {
	if err := b.checkConnections(); err != nil {
		return err
	}

	// Marshal request data into JSON:
	dataByte, err := json.Marshal(requestData)
	if err != nil {
		return handleError("request", err, b.logger)
	}

	// Send the request and wait for a response:
	msg, err := b.nconn.Request(subj, dataByte, DefaultRequestTimeout)
	if err != nil {
		return handleError("request", err, b.logger)
	}

	// Unmarshal the response data:
	if err := json.Unmarshal(msg.Data, responseDataPtr); err != nil {
		return handleError("request", err, b.logger)
	}

	return nil
}

// Subscribe creates a subscription to the specified subject, invoking the provided handler for received messages.
func (b *Basic) Subscribe(subject string, handler func(msg *nats.Msg)) {
	if err := b.checkConnections(); err != nil {
		handleError("subscribe", err, b.logger)
		return
	}

	// Establish the subscription:
	_, err := b.nconn.Subscribe(subject, handler)
	handleError("subscribe", err, b.logger)
}

// checkConnections verifies the existence of necessary connections and logger.
func (b *Basic) checkConnections() error {
	if b.nconn == nil {
		return ErrNilNatsConnObject
	}

	if b.logger == nil {
		return ErrNilLoggerObject
	}

	return nil
}
