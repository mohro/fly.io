package messages

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

type GossipMessage struct {
	// Message type
	Type string `json:"type,omitempty"`

	// Topology data
	Data []int `json:"data,omitempty"`
}
