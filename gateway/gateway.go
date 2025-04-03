package gateway

import (
	"bytes"
	"discord-go/api"
	"discord-go/api/types"
	"discord-go/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"reflect"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type ConnectionState int

const (
	StateDisconnected ConnectionState = iota
	StateConnecting
	StateConnected
	StateResuming
)

const (
	baseReconnectDelay = 1 * time.Second
	maxReconnectDelay  = 60 * time.Second
)

type Gateway struct {
	Token         string
	Conn          *websocket.Conn
	EventChan     chan json.RawMessage
	stopHeartbeat chan struct{}
	cache         *SessionCache
	eventHandlers map[string][]eventHandler
	middlewares   []MiddlewareFunc
	logger        utils.Logger

	// protected by mutex
	mu                sync.RWMutex
	sequence          *int64
	sessionID         string
	intents           int
	state             ConnectionState
	reconnectAttempts int
	heartbeatInterval time.Duration
}

type eventHandler struct {
	handlerFunc interface{}
	eventType   reflect.Type
}

type MiddlewareFunc func(eventType string, data json.RawMessage, next func())

type SessionCache struct {
	Guilds   map[string]types.Guild
	Channels map[string]types.Channel
	Users    map[string]types.User
	Messages map[string]types.Message
	mu       sync.RWMutex
}

func NewGateway(token string, intents ...int) *Gateway {
	intentValue := api.IntentAll
	if len(intents) > 0 {
		intentValue = intents[0]
	}

	return &Gateway{
		Token:         token,
		EventChan:     make(chan json.RawMessage, 100),
		stopHeartbeat: make(chan struct{}),
		intents:       intentValue,
		state:         StateDisconnected,
		eventHandlers: make(map[string][]eventHandler),
		logger:        utils.NewLogger(),
		cache: &SessionCache{
			Guilds:   make(map[string]types.Guild),
			Channels: make(map[string]types.Channel),
			Users:    make(map[string]types.User),
			Messages: make(map[string]types.Message),
		},
	}
}

// Connection Management --------------------------------------------------------

func (g *Gateway) Connect(url string) error {
	g.mu.Lock()
	if g.EventChan == nil {
		g.EventChan = make(chan json.RawMessage, 100)
	}
	g.mu.Unlock()

	g.setState(StateConnecting)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		g.setState(StateDisconnected)
		return err
	}
	g.Conn = conn

	if err := g.sendIdentify(); err != nil {
		g.setState(StateDisconnected)
		return err
	}

	g.setState(StateConnected)
	go g.listen()
	return nil
}

func (g *Gateway) Close() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.setState(StateDisconnected)
	close(g.stopHeartbeat)

	if g.EventChan != nil {
		close(g.EventChan)
		g.EventChan = nil
	}

	if g.Conn != nil {
		return g.Conn.Close()
	}
	return nil
}

// Event Handling --------------------------------------------------------------

func (g *Gateway) RegisterHandler(eventType string, handlerFunc interface{}, eventStruct ...interface{}) {
	g.mu.Lock()
	defer g.mu.Unlock()

	handlerType := reflect.TypeOf(handlerFunc)
	if handlerType.Kind() != reflect.Func || handlerType.NumIn() != 1 {
		panic("handler must be a function with exactly one parameter")
	}

	var eventTypeReflect reflect.Type
	if len(eventStruct) > 0 {
		eventTypeReflect = reflect.TypeOf(eventStruct[0])
	} else {
		eventTypeReflect = handlerType.In(0)
	}

	g.eventHandlers[eventType] = append(g.eventHandlers[eventType], eventHandler{
		handlerFunc: handlerFunc,
		eventType:   eventTypeReflect,
	})
}

func (g *Gateway) RemoveHandler(eventType string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.eventHandlers, eventType)
}

func (g *Gateway) Use(middleware MiddlewareFunc) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.middlewares = append(g.middlewares, middleware)
}

func (g *Gateway) handleEvent(eventType string, data json.RawMessage) {
	g.mu.RLock()
	handlers := g.eventHandlers[eventType]
	middlewares := g.middlewares
	g.mu.RUnlock()

	if len(handlers) == 0 {
		return
	}

<<<<<<< HEAD
	//   middleware chain
=======
	// create 	 middleware chain
>>>>>>> b641ce9 (should be ok)
	chain := func() {
		for _, handler := range handlers {
			eventPtr := reflect.New(handler.eventType).Interface()
			if err := json.Unmarshal(data, eventPtr); err != nil {
				g.logger.Errorf("Failed to unmarshal %s event: %v", eventType, err)
				continue
			}

			reflect.ValueOf(handler.handlerFunc).Call([]reflect.Value{
				reflect.ValueOf(eventPtr).Elem(),
			})
		}
	}

<<<<<<< HEAD
	// apply middlewares in reverse order
=======
	// reverse order
>>>>>>> b641ce9 (should be ok)
	for i := len(middlewares) - 1; i >= 0; i-- {
		mw := middlewares[i]
		prevChain := chain
		chain = func() { mw(eventType, data, prevChain) }
	}

	chain()
}

// WebSocket Communication -----------------------------------------------------

func (g *Gateway) listen() {
	defer func() {
		g.mu.Lock()
		defer g.mu.Unlock()
		if g.EventChan != nil {
			close(g.EventChan)
			g.EventChan = nil
		}
	}()

	for {
		_, message, err := g.Conn.ReadMessage()
		if err != nil {
			g.logger.Errorf("Gateway read err: %v", err)
			g.reconnect()
			return
		}

		var baseEvent struct {
			OP int             `json:"op"`
			T  string          `json:"t"`
			D  json.RawMessage `json:"d"`
			S  int64           `json:"s"`
		}

		if err := json.Unmarshal(message, &baseEvent); err != nil {
			g.logger.Errorf("Failed to unmarshal base event: %v", err)
			continue
		}

		if baseEvent.S != 0 {
			g.mu.Lock()
			g.sequence = &baseEvent.S
			g.mu.Unlock()
		}

		switch baseEvent.OP {
		case 0: // Dispatch
			g.handleEvent(baseEvent.T, baseEvent.D)
		case 10: // Hello
			g.handleHello(baseEvent.D)
		case 11: // Heartbeat ACK
			g.logger.Debug("Heartbeat acknowledged")
		case 7: // Reconnect
			g.logger.Info("Server requested reconnect")
			g.reconnect()
		case 9: // Invalid Session
			g.logger.Warn("Invalid session, reconnecting...")
			time.Sleep(5 * time.Second)
			g.reconnect()
		default:
			g.logger.Debugf("Unhandled OP code: %d", baseEvent.OP)
		}
	}
}

func (g *Gateway) handleHello(data json.RawMessage) {
	var hello struct {
		HeartbeatInterval int `json:"heartbeat_interval"`
	}
	if err := json.Unmarshal(data, &hello); err != nil {
		g.logger.Errorf("Failed to parse hello: %v", err)
		return
	}

	g.heartbeatInterval = time.Duration(hello.HeartbeatInterval) * time.Millisecond
	go g.startHeartbeat(g.heartbeatInterval)
	g.logger.Infof("Heartbeat started with interval: %v", g.heartbeatInterval)
}

func (g *Gateway) sendIdentify() error {
	payload := map[string]interface{}{
		"op": 2,
		"d": map[string]interface{}{
			"token": g.Token,
			"properties": map[string]string{
				"$os":      "linux",
				"$browser": "nyrilol/discord-go",
				"$device":  "nyrilol/discord-go",
			},
			"intents": g.intents,
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return g.Conn.WriteMessage(websocket.TextMessage, data)
}

func (g *Gateway) sendHeartbeat() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.Conn == nil {
		return errors.New("connection not established")
	}

	seq := g.sequence
	if seq == nil {
		seq = new(int64)
	}

	payload := map[string]interface{}{
		"op": 1,
		"d":  *seq,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return g.Conn.WriteMessage(websocket.TextMessage, data)
}

func (g *Gateway) startHeartbeat(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// send heartbeat
	if err := g.sendHeartbeat(); err != nil {
		g.logger.Errorf("Failed to send heartbeat: %v", err)
		return
	}

	for {
		select {
		case <-ticker.C:
			if err := g.sendHeartbeat(); err != nil {
				g.logger.Errorf("Failed to send heartbeat:", err)
				return
			}
		case <-g.stopHeartbeat:
			return
		}
	}
}

// State Management ------------------------------------------------------------

func (g *Gateway) GetState() ConnectionState {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.state
}

func (g *Gateway) setState(state ConnectionState) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.state = state
}

// Reconnection Handling --------------------------------------------------------

func (g *Gateway) reconnect() {
	g.setState(StateConnecting)
	g.reconnectAttempts++

	delay := calculateBackoff(g.reconnectAttempts)
	g.logger.Infof("Reconnecting attempt %d in %v...", g.reconnectAttempts, delay)

	time.Sleep(delay)

	if err := g.Connect("wss://gateway.discord.gg/?v=10&encoding=json"); err != nil {
		g.logger.Errorf("Reconnection failed: %v", err)
		g.reconnect()
		return
	}

	g.reconnectAttempts = 0
}

func calculateBackoff(attempt int) time.Duration {
	minDelay := math.Pow(2, float64(attempt)) * float64(baseReconnectDelay)
	maxDelay := math.Min(minDelay*1.5, float64(maxReconnectDelay))
	jitter := (maxDelay - minDelay) * 0.2

	return time.Duration(minDelay + jitter)
}

// Interaction Handling --------------------------------------------------------

func (g *Gateway) SendInteractionResponse(interactionID types.Snowflake, interactionToken string, response types.InteractionResponse) error {
	payload := map[string]interface{}{
		"type": response.Type,
		"data": response.Data,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://discord.com/api/v10/interactions/%s/%s/callback", interactionID, interactionToken)
	_, err = g.makeHTTPRequest("POST", url, data)
	return err
}

func (g *Gateway) SendFollowupMessage(interactionToken string, message types.WebhookMessage) error {
	appID, err := g.getApplicationID()
	if err != nil {
		g.logger.Errorf("failed to get application ID: %v", err)

		return err
	}

	g.mu.RLock()
	defer g.mu.RUnlock()

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://discord.com/api/v10/webhooks/%s/%s/messages", appID, interactionToken)
	_, err = g.makeHTTPRequest("POST", url, data)
	return err
}

func (g *Gateway) EditOriginalInteractionResponse(interactionToken, content string) error {
	appID, err := g.getApplicationID()
	if err != nil {
		g.logger.Errorf("failed to get application ID: %v", err)

		return err
	}

	g.mu.RLock()
	defer g.mu.RUnlock()

	payload := map[string]interface{}{
		"content": content,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://discord.com/api/v10/webhooks/%s/%s/messages/@original", appID, interactionToken)
	_, err = g.makeHTTPRequest("PATCH", url, data)
	return err
}

func (g *Gateway) CreateGlobalApplicationCommand(command types.ApplicationCommand) error {
	return g.createApplicationCommand("", command)
}

func (g *Gateway) CreateGuildApplicationCommand(guildID types.Snowflake, command types.ApplicationCommand) error {
	return g.createApplicationCommand(guildID, command)
}

func (g *Gateway) createApplicationCommand(guildID types.Snowflake, command types.ApplicationCommand) error {
	appID, err := g.getApplicationID()
	if err != nil {
		g.logger.Errorf("failed to get application ID: %v", err)

		return err
	}

	data, err := json.Marshal(command)
	if err != nil {
		return err
	}

	var url string
	if guildID != "" {
		url = fmt.Sprintf("https://discord.com/api/v10/applications/%s/guilds/%s/commands", appID, guildID)
	} else {
		url = fmt.Sprintf("https://discord.com/api/v10/applications/%s/commands", appID)
	}

	_, err = g.makeHTTPRequest("POST", url, data)
	return err
}

func (g *Gateway) getApplicationID() (string, error) {
	url := "https://discord.com/api/v10/users/@me"
	resp, err := g.makeHTTPRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	var botUser struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(resp, &botUser); err != nil {
		return "", err
	}

	return botUser.ID, nil
}

func (g *Gateway) makeHTTPRequest(method, url string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bot "+g.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "DiscordBot (https://github.com/nyrilol/discord-go, 1.0)")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, resp.Status)
	}

	return io.ReadAll(resp.Body)
}
