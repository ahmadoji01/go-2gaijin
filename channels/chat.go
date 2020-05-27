package channels

type Message struct {
	Data []byte
	Room string
}

type subscription struct {
	conn *connection
	Room string
}

// hub maintains the set of active connections and Broadcasts messages to the
// connections.
type hub struct {
	// Registered connections.
	rooms map[string]map[*connection]bool

	// Inbound messages from the connections.
	Broadcast chan Message

	// Register requests from the connections.
	register chan subscription

	// Unregister requests from connections.
	unregister chan subscription
}

var H = hub{
	Broadcast:  make(chan Message),
	register:   make(chan subscription),
	unregister: make(chan subscription),
	rooms:      make(map[string]map[*connection]bool),
}

func (h *hub) Run() {
	for {
		select {
		case s := <-h.register:
			connections := h.rooms[s.Room]
			if connections == nil {
				connections = make(map[*connection]bool)
				h.rooms[s.Room] = connections
			}
			h.rooms[s.Room][s.conn] = true
		case s := <-h.unregister:
			connections := h.rooms[s.Room]
			if connections != nil {
				if _, ok := connections[s.conn]; ok {
					delete(connections, s.conn)
					close(s.conn.send)
					if len(connections) == 0 {
						delete(h.rooms, s.Room)
					}
				}
			}
		case m := <-h.Broadcast:
			connections := h.rooms[m.Room]
			for c := range connections {
				select {
				case c.send <- m.Data:
				default:
					close(c.send)
					delete(connections, c)
					if len(connections) == 0 {
						delete(h.rooms, m.Room)
					}
				}
			}
		}
	}
}
