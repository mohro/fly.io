package services

import (
	"encoding/json"
	"maelstrom-broadcast/db"
	"maelstrom-broadcast/messages"
	"math/rand"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Handlers struct {
	node       *maelstrom.Node
	store      *db.Store
	neighbours []string
}

// Unmarshal the message body as an Broadcast Message.
func (h *Handlers) HandleBroadcast(msg maelstrom.Message) error {
	var body messages.BroadcastMessage
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	h.store.Insert(body.Message)
	response := map[string]string{}
	response["type"] = "broadcast_ok"

	// Respond with broadcast_ok message
	return h.node.Reply(msg, response)
}

func (h *Handlers) HandleRead(msg maelstrom.Message) error {
	// Unmarshal the message body as an loosely-typed map.
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	// Update the message type.
	body["type"] = "read_ok"
	body["messages"] = h.store.Read()

	// Echo the original message back with the updated message type.
	return h.node.Reply(msg, body)
}

func (h *Handlers) HandleTopology(msg maelstrom.Message) error {
	// Unmarshal the message message as an loosely-typed map.
	var message messages.TopologyMessage
	if err := json.Unmarshal(msg.Body, &message); err != nil {
		return err
	}

	h.neighbours = message.Topology[h.node.ID()]
	response := map[string]string{}
	// Update the message type.
	response["type"] = "topology_ok"

	// Echo the original message back with the updated message type.
	return h.node.Reply(msg, response)
}

func (h *Handlers) HandleGossip(msg maelstrom.Message) error {
	// Unmarshal the message message as an loosely-typed map.
	var message messages.GossipMessage
	if err := json.Unmarshal(msg.Body, &message); err != nil {
		return err
	}

	data := message.Data
	h.store.Merge(data)
	response := map[string]string{}
	// Update the message type.
	response["type"] = "gossip_ok"

	// Echo the original message back with the updated message type.
	return h.node.Reply(msg, response)
}

func (h *Handlers) HandleInit(msg maelstrom.Message) error {
	repeatDuration := 100 + rand.Intn(100)
	ticker := time.NewTicker(time.Duration(repeatDuration) * time.Millisecond)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				response := map[string]any{}
				response["type"] = "gossip"
				response["data"] = h.store.Read()
				for _, neighbour := range h.neighbours {
					h.node.Send(neighbour, response)
				}
			}
		}
	}()

	return nil
}

func (h *Handlers) EmptyHandler(msg maelstrom.Message) error {
	// Unmarshal the message message as an loosely-typed map.
	return nil
}

func NewHandlers(node *maelstrom.Node, store *db.Store) *Handlers {
	return &Handlers{
		node:  node,
		store: store,
	}
}
