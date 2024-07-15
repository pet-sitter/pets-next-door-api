package chat

type WebSocketServer struct {
	StateManager   StateManager
	RegisterChan   chan *Client
	UnregisterChan chan *Client
	BroadcastChan  chan []byte
}

func NewWebSocketServer(stateManager StateManager) *WebSocketServer {
	return &WebSocketServer{
		StateManager:   stateManager,
		RegisterChan:   make(chan *Client),
		UnregisterChan: make(chan *Client),
		BroadcastChan:  make(chan []byte),
	}
}

func (server *WebSocketServer) Run() {
	for {
		select {
		case client := <-server.RegisterChan:
			server.StateManager.RegisterClient(client)
		case client := <-server.UnregisterChan:
			server.StateManager.UnregisterClient(client)
		case message := <-server.BroadcastChan:
			server.StateManager.BroadcastToClients(message)
		}
	}
}
