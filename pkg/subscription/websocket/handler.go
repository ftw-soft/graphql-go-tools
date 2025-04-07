package websocket

import (
	"net/http"
	"time"

	"github.com/jensneuse/abstractlogger"

	"github.com/jensneuse/graphql-go-tools/pkg/subscription"
)

const (
	DefaultConnectionInitTimeOut = "15s"

	HeaderSecWebSocketProtocol = "Sec-WebSocket-Protocol"
)

// Protocol defines the protocol names as type.
type Protocol string

const (
	ProtocolUndefined          Protocol = ""
	ProtocolGraphQLWS          Protocol = "graphql-ws"
	ProtocolGraphQLTransportWS Protocol = "graphql-transport-ws"
)

var DefaultProtocol = ProtocolGraphQLTransportWS

// HandleOptions can be used to pass options to the websocket handler.
type HandleOptions struct {
	Logger                           abstractlogger.Logger
	Protocol                         Protocol
	WebSocketInitFunc                InitFunc
	CustomClient                     subscription.TransportClient
	CustomKeepAliveInterval          time.Duration
	CustomSubscriptionUpdateInterval time.Duration
	CustomConnectionInitTimeOut      time.Duration
	CustomReadErrorTimeOut           time.Duration
	CustomSubscriptionEngine         subscription.Engine
}

// HandleOptionFunc can be used to define option functions.
type HandleOptionFunc func(opts *HandleOptions)

// WithLogger is a function that sets a logger for the websocket handler.
func WithLogger(logger abstractlogger.Logger) HandleOptionFunc {
	return func(opts *HandleOptions) {
		opts.Logger = logger
	}
}

// WithInitFunc is a function that sets the init function for the websocket handler.
func WithInitFunc(initFunc InitFunc) HandleOptionFunc {
	return func(opts *HandleOptions) {
		opts.WebSocketInitFunc = initFunc
	}
}

// WithCustomClient is a function that set a custom transport client for the websocket handler.
func WithCustomClient(client subscription.TransportClient) HandleOptionFunc {
	return func(opts *HandleOptions) {
		opts.CustomClient = client
	}
}

// WithCustomKeepAliveInterval is a function that sets a custom keep-alive interval for the websocket handler.
func WithCustomKeepAliveInterval(keepAliveInterval time.Duration) HandleOptionFunc {
	return func(opts *HandleOptions) {
		opts.CustomKeepAliveInterval = keepAliveInterval
	}
}

// WithCustomSubscriptionUpdateInterval is a function that sets a custom subscription update interval for the
// websocket handler.
func WithCustomSubscriptionUpdateInterval(subscriptionUpdateInterval time.Duration) HandleOptionFunc {
	return func(opts *HandleOptions) {
		opts.CustomSubscriptionUpdateInterval = subscriptionUpdateInterval
	}
}

// WithCustomConnectionInitTimeOut is a function that sets a custom connection init time out.
func WithCustomConnectionInitTimeOut(connectionInitTimeOut time.Duration) HandleOptionFunc {
	return func(opts *HandleOptions) {
		opts.CustomConnectionInitTimeOut = connectionInitTimeOut
	}
}

// WithCustomReadErrorTimeOut is a function that sets a custom read error time out for the
// websocket handler.
func WithCustomReadErrorTimeOut(readErrorTimeOut time.Duration) HandleOptionFunc {
	return func(opts *HandleOptions) {
		opts.CustomReadErrorTimeOut = readErrorTimeOut
	}
}

// WithCustomSubscriptionEngine is a function that sets a custom subscription engine for the websocket handler.
func WithCustomSubscriptionEngine(subscriptionEngine subscription.Engine) HandleOptionFunc {
	return func(opts *HandleOptions) {
		opts.CustomSubscriptionEngine = subscriptionEngine
	}
}

// WithProtocol is a function that sets the protocol.
func WithProtocol(protocol Protocol) HandleOptionFunc {
	return func(opts *HandleOptions) {
		opts.Protocol = protocol
	}
}

// WithProtocolFromRequestHeaders is a function that sets the protocol based on the request headers.
// It fallbacks to the DefaultProtocol if the header can't be found, the value is invalid or no request
// was provided.
func WithProtocolFromRequestHeaders(req *http.Request) HandleOptionFunc {
	return func(opts *HandleOptions) {
		if req == nil {
			opts.Protocol = DefaultProtocol
			return
		}

		protocolHeaderValue := req.Header.Get(HeaderSecWebSocketProtocol)
		switch Protocol(protocolHeaderValue) {
		case ProtocolGraphQLWS:
			opts.Protocol = ProtocolGraphQLWS
		case ProtocolGraphQLTransportWS:
			opts.Protocol = ProtocolGraphQLTransportWS
		default:
			opts.Protocol = DefaultProtocol
		}
	}
}

func CreateProtocolHandler(handleOptions HandleOptions, client subscription.TransportClient) (protocolHandler subscription.Protocol, err error) {
	protocol := handleOptions.Protocol
	if protocol == ProtocolUndefined {
		protocol = DefaultProtocol
	}

	switch protocol {
	case ProtocolGraphQLWS:
		protocolHandler, err = NewProtocolGraphQLWSHandlerWithOptions(client, ProtocolGraphQLWSHandlerOptions{
			Logger:                  handleOptions.Logger,
			WebSocketInitFunc:       handleOptions.WebSocketInitFunc,
			CustomKeepAliveInterval: handleOptions.CustomKeepAliveInterval,
		})
	default:
		protocolHandler, err = NewProtocolGraphQLTransportWSHandlerWithOptions(client, ProtocolGraphQLTransportWSHandlerOptions{
			Logger:                    handleOptions.Logger,
			WebSocketInitFunc:         handleOptions.WebSocketInitFunc,
			CustomKeepAliveInterval:   handleOptions.CustomKeepAliveInterval,
			CustomInitTimeOutDuration: handleOptions.CustomConnectionInitTimeOut,
		})
	}

	return protocolHandler, err
}
