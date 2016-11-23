// Package websocket contains helpers and implementation for backend functions
// that interact with the frontend over websockets.
package websocket

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

// maxMessageSize is used to set a read limit on the websocket and prevent
// clients from flooding us with data.
const maxMessageSize int64 = 8096

// authType is the type string used for auth messages.
const authType string = "auth"

// errType is the type string used for error messages.
const errType string = "error"

// defaultTimeout is the default timeout that should be used for sending and
// receiving over the websocket. It is used unless Conn.Timeout is set
// explicitly after Upgrade is called.
const defaultTimeout time.Duration = 3 * time.Second

// JSONMessage is a wrapper struct for messages that will be sent across the wire
// as JSON.
type JSONMessage struct {
	// Type is a string indicating which message type the data contains
	Type string `json:"type"`
	// Data contains the arbitrarily schemaed JSON data. Type should
	// indicate how this should be deserialized.
	Data interface{} `json:"data"`
}

// Conn is a wrapper for a standard websocket connection with utility methods
// added for interacting with Kolide specific message types.
type Conn struct {
	*websocket.Conn
	Timeout time.Duration
}

// Upgrade is used to upgrade a normal HTTP request to a websocket connection.
func Upgrade(w http.ResponseWriter, r *http.Request) (*Conn, error) {
	var upgrader = websocket.Upgrader{
		HandshakeTimeout: defaultTimeout,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, errors.Wrap(err, "upgrading connection")
	}

	conn.SetReadLimit(maxMessageSize)

	return &Conn{conn, defaultTimeout}, nil
}

// WriteJSONMessage writes the provided data as JSON (using the Message struct),
// returning any error condition from the connection.
func (c *Conn) WriteJSONMessage(typ string, data interface{}) error {
	c.SetWriteDeadline(time.Now().Add(c.Timeout))
	defer c.SetWriteDeadline(time.Time{})

	return c.WriteJSON(JSONMessage{Type: typ, Data: data})
}

// WriteJSONError writes an error (Message struct with Type="error"), returning any
// error condition from the connection.
func (c *Conn) WriteJSONError(data interface{}) error {
	return c.WriteJSONMessage(errType, data)
}

// ReadJSONMessage reads an incoming Message from JSON. Note that the
// Message.Data field is guaranteed to be *json.RawMessage, and so unchecked
// type assertions may be performed as in:
//  msg, err := c.ReadJSONMessage()
//  if err == nil && msg.Type == "foo" {
//  	var foo fooData
//  	json.Unmarshal(*(msg.Data.(*json.RawMessage)), &foo)
//  }
func (c *Conn) ReadJSONMessage() (*JSONMessage, error) {
	c.SetReadDeadline(time.Now().Add(c.Timeout))
	defer c.SetReadDeadline(time.Time{})

	mType, data, err := c.ReadMessage()
	if err != nil {
		return nil, errors.Wrap(err, "reading from websocket")
	}
	if mType != websocket.TextMessage {
		return nil, errors.Errorf("unsupported websocket message type: %d", mType)
	}

	msg := &JSONMessage{Data: &json.RawMessage{}}

	if err := json.Unmarshal(data, msg); err != nil {
		return nil, errors.Wrap(err, "parsing msg json")
	}

	if msg.Type == "" {
		return nil, errors.New("missing message type")
	}

	return msg, nil
}

// authData defines the data used to authenticate a Kolide frontend client over
// a websocket connection.
type authData struct {
	Token string `json:"token"`
}

// ReadAuthToken reads from the websocket, returning an auth token embedded in
// a JSONMessage with type "auth" and data that can be unmarshalled to
// authData.
func (c *Conn) ReadAuthToken() (string, error) {
	msg, err := c.ReadJSONMessage()
	if err != nil {
		return "", errors.Wrap(err, "read auth token")
	}
	if msg.Type != authType {
		return "", errors.Errorf(`message type not "%s": "%s"`, authType, msg.Type)
	}

	var auth authData
	if err := json.Unmarshal(*(msg.Data.(*json.RawMessage)), &auth); err != nil {
		return "", errors.Wrap(err, "unmarshal auth data")
	}

	return auth.Token, nil
}
