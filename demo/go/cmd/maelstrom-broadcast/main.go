package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type BroadcastMessage struct {
	// Message type.
	Type string `json:"type,omitempty"`

	// Error message, if an error occurred.
	Message int `json:"message,omitempty"`
}

type TopologyMessage struct {
	// Message type
	Type string `json:"type,omitempty"`

	// Topology data
	Topology map[string][]string `json:"topology,omitempty"`
}

type Store struct {
	index  map[int]bool
	cached []int
	mu     sync.RWMutex
}

func (s *Store) InsertIfNotPresent(value int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.index[value]
	if !ok {
		s.index[value] = true
	}
	return !ok
}

func (s *Store) Read() []int {

	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.cached) == len(s.index) {
		return s.cached
	}

	result := make([]int, len(s.index))
	i := 0
	for k := range s.index {
		result[i] = k
		i++
	}

	s.cached = result
	return result
}

func Broadcast(message int, neighbours []string, node *maelstrom.Node) {
	response := map[string]any{}
	response["type"] = "broadcast"
	response["message"] = message

	for _, neighbour := range neighbours {
		node.Send(neighbour, response)
	}

}

func NewStore() *Store {
	index := map[int]bool{}
	cached := []int{}
	return &Store{index: index, cached: cached}
}

func main() {
	n := maelstrom.NewNode()
	var neighbours []string

	store := NewStore()

	// Register a handler for the "echo" message that responds with an "echo_ok".
	n.Handle("broadcast", func(msg maelstrom.Message) error {
		// Unmarshal the message body as an loosely-typed map.
		var body BroadcastMessage
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		response := map[string]string{}
		message := body.Message

		if ok := store.InsertIfNotPresent(message); ok {
			Broadcast(message, neighbours, n)
		}

		response["type"] = "broadcast_ok"

		// Update the message type.
		// body["type"] = "broadcast_ok"

		// Echo the original message back with the updated message type.
		return n.Reply(msg, response)
	})

	n.Handle("broadcast_ok", func(msg maelstrom.Message) error {
		// Unmarshal the message body as an loosely-typed map.
		return nil
	})
	n.Handle("read", func(msg maelstrom.Message) error {
		// Unmarshal the message body as an loosely-typed map.
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		// Update the message type.
		body["type"] = "read_ok"
		body["messages"] = store.Read()

		// Echo the original message back with the updated message type.
		return n.Reply(msg, body)
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		// Unmarshal the message message as an loosely-typed map.
		var message TopologyMessage
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			return err
		}

		neighbours = message.Topology[n.ID()]
		response := map[string]string{}
		// Update the message type.
		response["type"] = "topology_ok"

		// Echo the original message back with the updated message type.
		return n.Reply(msg, response)
	})

	// Execute the node's message loop. This will run until STDIN is closed.
	if err := n.Run(); err != nil {
		log.Printf("ERROR: %s", err)
		os.Exit(1)
	}
}
