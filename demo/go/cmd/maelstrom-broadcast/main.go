package main

import (
	"encoding/json"
	"log"
	"os"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type BroadcastMessage struct {
	// Message type.
	Type string `json:"type,omitempty"`

	// Error message, if an error occurred.
	Message int `json:"message,omitempty"`
}

func main() {
	n := maelstrom.NewNode()
	store := map[int]bool{}

	// Register a handler for the "echo" message that responds with an "echo_ok".
	n.Handle("broadcast", func(msg maelstrom.Message) error {
		// Unmarshal the message body as an loosely-typed map.
		var body BroadcastMessage
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		response := map[string]string{}
		message := body.Message

		store[message] = true
		response["type"] = "broadcast_ok"

		// Update the message type.
		// body["type"] = "broadcast_ok"

		// Echo the original message back with the updated message type.
		return n.Reply(msg, response)
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		// Unmarshal the message body as an loosely-typed map.
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		keys := make([]int, len(store))

		i := 0
		for k := range store {
			keys[i] = k
			i++
		}

		// Update the message type.
		body["type"] = "read_ok"
		body["messages"] = keys

		// Echo the original message back with the updated message type.
		return n.Reply(msg, body)
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		// Unmarshal the message body as an loosely-typed map.
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

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
