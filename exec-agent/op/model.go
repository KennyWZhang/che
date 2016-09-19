// Provides a lightweight implementation of jsonrpc2.0 protocol.
// The original jsonrpc2.0 specification - http://www.jsonrpc.org/specification
//
// The implementations does not fully implement the protocol
// and introduces a few modifications to its terminology in
// term of exec-agent transport needs.
//
// From the specification:
// The Client is defined as the origin of Request objects and the handler of Response objects.
// The Server is defined as the origin of Response objects and the handler of Request objects.
//
// Exec-agent serves as both, server and client as it receives
// Responses and sends Notifications in the same time.
//
// Request.
// It's a message from the physical websocket connection client to the exec-agent server.
// Request(in it's origin form) is considered is considered to be unidirectional.
// WS Client =---> WS Server.
//
// Response.
// It's a message from the the exec-agent server to a websocket client,
// indicates the result of the operation execution requested by certain request.
// Response doesn't exist without request. The response is considered to be unidirectional.
// WS Client <---= WS Server
//
// Event.
// Is a message from the exec-agent server to a websocket client, the analogue
// from the specification is Notification, which is defined as a request
// which doesn't need any response, that's also true for events.
// Events may happen periodically and don't need to be indicated by request.
// WS Client <---X WS Server
package op

import (
	"encoding/json"
	"time"
)

// Describes named operation which is called
// on the websocket client's side and performed
// on the servers's side, if appropriate Route exists.
type Request struct {

	// Version of this request e.g. '2.0'.
	Version string `json:"jsonrpc"`

	// The method name which should be proceeded by this call
	// usually dot separated resource and action e.g. 'process.start'.
	Method string `json:"method"`

	// The unique identifier of this operation request.
	// If a client needs to identify the result of the operation execution,
	// the id should be passed by the client, then it is guaranteed
	// that the client will receive the result frame with the same id.
	// The uniqueness of the identifier must be controlled by the client,
	// if client doesn't specify the identifier in the operation call,
	// the response won't contain the identifier as well.
	//
	// It is preferable to specify identifier for those calls which may
	// either validate data, or produce such information which can't be
	// identified by itself.
	Id interface{} `json:"id"`

	// Request data, parameters which are needed for operation execution.
	RawBody json.RawMessage `json:"params"`
}

// A message from the server to the client,
// which represents the result of the certain operation execution.
// The result is sent to the client only once per operation.
type Response struct {

	// Version of this response e.g. '2.0'.
	Version string `json:"jsonrpc"`

	// The operation call identifier, will be set only
	// if the operation contains it. See 'op.Call.Id'
	Id interface{} `json:"id"`

	// The actual result data, the operation execution result.
	Body interface{} `json:"result,omitempty"`

	// Body and Error are mutual exclusive.
	// Present only if the operation execution fails due to an error.
	Error *Error `json:"error,omitempty"`
}

// A message from the server to the client,
// which may notify client about any activity that the client is interested in.
// The difference from the 'op.Response' is that the event may happen periodically,
// before or even after some operation calls, while the 'op.Response' is more like
// result of the operation call execution, which is sent to the client immediately
// after the operation execution is done.
type Event struct {

	// A type of this operation event, must be always set.
	// The type must be generally unique.
	EventType string

	// Event related data.
	Body Periodical
}

// Event has to be periodical
type Periodical interface {
	SetTime(time time.Time)
}

// Holds a value of time.
type EventBody struct {
	time time.Time
}

// Implements Periodical interface.
func (th *EventBody) SetTime(t time.Time) {
	th.time = t
}

func NewEventNow(eType string, body Periodical) *Event {
	body.SetTime(time.Now())
	return &Event{
		EventType: eType,
		Body:      body,
	}
}

func NewEvent(eType string, body Periodical, time time.Time) *Event {
	body.SetTime(time)
	return &Event{
		EventType: eType,
		Body:      body,
	}
}
