# go-natstore

**This library provides convenient abstractions for interacting with NATS and JetStream messaging systems.**

## Key Features

- **Simplified NATS messaging:**
    - Publish, request, and subscribe to messages with ease.
    - Error handling and logging for robust operations.
- **JetStream integration:**
    - Create and manage JetStream streams.
    - Publish and consume messages with reliable delivery and enhanced features.
- **Customizable configuration:**
    - Tailor messaging behavior with options for consumers and error handling.
- **Error handling and logging:**
    - Enhances debugging and monitoring for troubleshooting.

## Installation

```bash
go get github.com/your-username/natstore
```

## Usage

**Basic NATS Interactions:**

```go
package main

import (
    "context"
    "fmt"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/otyang/go-natstore"
)

func main() {
    natsConn, err := natstore.Connect("nats://localhost:4222", "my-app")
	if err != nil {
		// Handle connection error
	}

	basic := natstore.NewBasic(natsConn)

	// Publish a message
	basic.Publish(context.Background(), "hello", "Hello, world!")

	// Request a response
	var response string
	basic.Request(context.Background(), "request", "What time is it?", &response)
	fmt.Println("Response:", response)

	// Subscribe to a subject
	basic.Subscribe("updates", func(msg *nats.Msg) {
		fmt.Println("Received message:", string(msg.Data))
	})
}
```

**JetStream Interactions:**

```go
package main


import (
    "context"
    "fmt"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/otyang/go-natstore" 
	"github.com/nats-io/nats.go/jetstream" 
)

func main() {
   natsConn, err := natstore.Connect("", "jetstream_client_example") // embedded server
	if err != nil {
		// Handle connection error
	}

	jsClient, err := natstore.NewJSClient(natsConn)
	if err != nil {
		// Handle error
	}

	var (
		ctx          = context.TODO()
		consumerName = "consumername"
	)

	// You have to create a stream before publish/consume
	err = jsClient.StreamCreateOrUpdate(ctx, jetstream.StreamConfig{Name: "streamName", Subjects: []string{"Subject"}})
	// handle error

	// Publish a message to the stream
	err = jsClient.Publish(ctx, "subject", "hello")
	// handle error

	// Consume messages from the stream
	err = jsClient.Consume(ctx,
		"streamName", "subject", func(m jetstream.Msg) {
			fmt.Println("Received message:", string(m.Subject()))
		},
		natstore.WithDurableName(consumerName))
	// handle error
}
```

## Additional Information

- **Error Handling:** The library utilizes a `handleError` function for consistent error handling and logging.
- **Customization:** Options are available to customize consumer behavior and error handling strategies.

## Contributing

Contributions are welcome! Please refer to the contribution guidelines for details.

**For more detailed usage and examples, please refer to the package documentation and test cases.**
