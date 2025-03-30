package gateway

import (
	"discord-go/api/types"
	"discord-go/utils"
	"encoding/json"
	"errors"
	"math"
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
	intentValue := IntentAll
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
		logger:        utils.NewLogger(), // Initialize logger
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

	if g.Conn != nil {
		return g.Conn.Close()
	}
	return nil
}

// Event Handling --------------------------------------------------------------

func (g *Gateway) RegisterHandler(eventType string, handlerFunc interface{}) {
	g.mu.Lock()
	defer g.mu.Unlock()

	handlerType := reflect.TypeOf(handlerFunc)
	if handlerType.Kind() != reflect.Func || handlerType.NumIn() != 1 {
		panic("handler must be a function with exactly one parameter")
	}

	eventStructType := handlerType.In(0)
	g.eventHandlers[eventType] = append(g.eventHandlers[eventType], eventHandler{
		handlerFunc: handlerFunc,
		eventType:   eventStructType,
	})
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

	// Create middleware chain
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

	// Apply middlewares in reverse order
	for i := len(middlewares) - 1; i >= 0; i-- {
		mw := middlewares[i]
		prevChain := chain
		chain = func() { mw(eventType, data, prevChain) }
	}

	chain()
}

// WebSocket Communication -----------------------------------------------------

func (g *Gateway) listen() {
	defer close(g.EventChan)
	defer g.Conn.Close()

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
				"$browser": "discord-go",
				"$device":  "discord-go",
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

	// Send initial heartbeat immediately
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
