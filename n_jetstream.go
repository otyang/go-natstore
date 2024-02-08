package natstore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// JSClient encapsulates JetStream interactions, providing a convenient abstraction layer.
type JSClient struct {
	jsCtx          jetstream.JetStream
	natsConnection *nats.Conn
	logger         *slog.Logger
}

// NewJSClient instantiates a JSClient with natsconnection, jetStream context & logger.
func NewJSClient(natsConn *nats.Conn) (*JSClient, error) {
	if natsConn == nil {
		return nil, fmt.Errorf("nats connection error: No nats connection object returned")
	}

	jsCtx, err := jetstream.New(natsConn)
	if err != nil {
		return nil, fmt.Errorf("error creating jetstream context: %w", err)
	}

	s := &JSClient{
		jsCtx:          jsCtx,
		natsConnection: natsConn,
		logger:         defaultLogger,
	}

	return s, nil
}

// WithLogger configures the logger used by this JSClient instance.
func (s *JSClient) WithLogger(logger *slog.Logger) *JSClient {
	s.logger = logger
	return s
}

// NatsJS exposes the underlying JetStream context, enabling direct interaction with JetStream APIs
// for advanced use cases or fine-grained control over messaging operations.
//
// Caution: Exercise care when using this method, as it bypasses the client's convenience
// abstractions and necessitates a deeper understanding of JetStream's protocols and APIs.
func (s *JSClient) NatsJS() jetstream.JetStream {
	return s.jsCtx
}

// checkConnections verifies the necessary connections for JSClient operations.
// It returns an error if any of the connections are missing or invalid.
func (s *JSClient) checkConnections() error {
	if s.natsConnection == nil {
		return ErrNilNatsConnObject
	}

	if s.jsCtx == nil {
		return errors.New("jetstream connection not established")
	}

	if s.logger == nil {
		return ErrNilLoggerObject
	}

	return nil
}

// StreamCreateOrUpdate creates or updates a JetStream stream with the specified configuration.
// It handles potential errors and logs them for debugging.
func (s *JSClient) StreamCreateOrUpdate(ctx context.Context, streamConfig jetstream.StreamConfig) error {
	if err := s.checkConnections(); err != nil {
		return handleError("stream create update", err, s.logger)
	}

	_, err := s.jsCtx.CreateOrUpdateStream(ctx, streamConfig)
	if err != nil {
		return handleError("stream create update", err, s.logger)
	}

	return nil
}

// Publish publishes a message to a JetStream subject, ensuring proper connection and error handling.
// It serializes the data into JSON before publishing.
func (s *JSClient) Publish(ctx context.Context, subject string, data any) error {
	if err := s.checkConnections(); err != nil {
		return handleError("publish", err, s.logger)
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return handleError("publish", err, s.logger)
	}

	_, err = s.jsCtx.Publish(ctx, subject, dataBytes)
	if err != nil {
		return handleError("publish", err, s.logger)
	}

	return nil
}

// Consume(ctx, "auth", "user.signup", notifyUser(), WithDurableName("masko"))
// Consume establishes a message subscription on a JetStream stream and processes incoming messages.
// It allows for configuration options to customize delivery behavior.
func (s *JSClient) Consume(
	ctx context.Context, streamName string, subject string, handler func(msg jetstream.Msg), opts ...ConsumerConfigOption,
) error {
	if err := s.checkConnections(); err != nil {
		return handleError("consume", err, s.logger)
	}

	consumerConfig := jetstream.ConsumerConfig{
		AckPolicy:     jetstream.AckExplicitPolicy, // Require explicit acknowledgement of messages
		FilterSubject: subject,                     // Filter messages based on the specified subject
		BackOff: []time.Duration{ // Define backoff intervals for redelivery attempts
			5 * time.Second,
			10 * time.Second,
		},
		MaxDeliver: 4, // Set a maximum number of delivery attempts
	}

	for _, opt := range opts { // Apply any provided configuration options
		if opt != nil {
			opt(&consumerConfig)
		}
	}

	consumer, err := s.jsCtx.CreateOrUpdateConsumer(ctx, streamName, consumerConfig)
	if err != nil {
		return handleError("consume", err, s.logger)
	}

	_, err = consumer.Consume(handler)
	if err != nil {
		return handleError("consume", err, s.logger)
	}

	return nil
}

// ConsumerConfigOption is a function type for customizing JetStream consumer configuration.
// It receives a pointer to a jetstream.ConsumerConfig and modifies its properties.
type ConsumerConfigOption func(*jetstream.ConsumerConfig)

// WithDurableName creates a ConsumerConfigOption that sets the durable name for a consumer.
// This allows the consumer to resume processing from its last acknowledged message.
func WithDurableName(name string) ConsumerConfigOption {
	return func(o *jetstream.ConsumerConfig) {
		o.Durable = name // Set the durable name in the configuration
	}
}
