package natstore

import (
	"errors"
	"log/slog"
	"os"
	"strings"

	"github.com/nats-io/nats.go"
)

var defaultLogger *slog.Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

var (
	ErrNilNatsConnObject = errors.New("nats connection error: nats connection object not returned")
	ErrNilLoggerObject   = errors.New("nats: logger not initialised")
)

var (
	WithName                 = nats.Name
	WithTimeout              = nats.Timeout
	WithMaxReconnects        = nats.MaxReconnects
	WithReconnectWait        = nats.ReconnectWait
	WithRetryOnFailedConnect = nats.RetryOnFailedConnect
	WithAuthToken            = nats.Token
	WithAuthUserInfo         = nats.UserInfo
	WithAuthCredentialsFile  = nats.UserCredentials
	WithClientCert           = nats.ClientCert
	WithRootCAs              = nats.RootCAs
	WithAuthNkey             = withAuthNkey
	WithClosedHandler        = withClosedHandler
	WithDisconnectErrHandler = withDisconnectErrHandler
	WithReconnectHandler     = withReconnectHandler
)

func withAuthNkey(filePath string) nats.Option {
	natsOption, err := nats.NkeyOptionFromSeed(filePath)
	if err != nil {
		defaultLogger.Error("nats error: auth nkey option", "error", err.Error())
		os.Exit(1)
	}
	return natsOption
}

func withClosedHandler(handler func(nc *nats.Conn)) nats.Option {
	if handler == nil {
		handler = func(nc *nats.Conn) {
			defaultLogger.Error("nats client exiting", "error", nc.LastError())
		}
	}
	return nats.ClosedHandler(handler)
}

func withDisconnectErrHandler(handler func(nc *nats.Conn, err error)) nats.Option {
	if handler == nil {
		handler = func(nc *nats.Conn, err error) {
			defaultLogger.Error("nats client disconnected", "error", err.Error())
		}
	}
	return nats.DisconnectErrHandler(handler)
}

func withReconnectHandler(handler func(nc *nats.Conn)) nats.Option {
	if handler == nil {
		handler = func(nc *nats.Conn) {
			defaultLogger.Info("nats client reconnected", "url", nc.ConnectedUrl())
		}
	}
	return nats.ReconnectHandler(handler)
}

func handleError(location string, err error, logger *slog.Logger) error {
	if err != nil {
		if strings.Contains(err.Error(), "no responders available for request") {
			logger.Error(location + ": " + nats.ErrJetStreamNotEnabled.Error())
			return errors.New(location + ": " + nats.ErrJetStreamNotEnabled.Error())
		}

		logger.Error(location + ": " + err.Error())
		return errors.New(location + ": " + err.Error())
	}

	return nil
}
