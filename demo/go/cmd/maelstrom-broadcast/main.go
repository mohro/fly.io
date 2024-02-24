package main

import (
	"log"
	"maelstrom-broadcast/db"
	"maelstrom-broadcast/services"
	"os"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

// NewNode returns a new instance of Node connected to STDIN/STDOUT.

func main() {
	n := maelstrom.NewNode()

	store := db.NewStore()
	handlers := services.NewHandlers(n, store)

	n.Handle("init", handlers.HandleInit)
	n.Handle("broadcast", handlers.HandleBroadcast)
	n.Handle("read", handlers.HandleRead)
	n.Handle("topology", handlers.HandleTopology)
	n.Handle("broadcast_ok", handlers.EmptyHandler)
	n.Handle("gossip", handlers.HandleGossip)
	n.Handle("gossip_ok", handlers.EmptyHandler)

	// Execute the node's message loop. This will run until STDIN is closed.
	if err := n.Run(); err != nil {
		log.Printf("ERROR: %s", err)
		os.Exit(1)
	}
}
