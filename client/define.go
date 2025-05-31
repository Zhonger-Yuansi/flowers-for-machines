package client

import (
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/bunker/auth"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft"
	"github.com/Happy2018new/the-last-problem-of-the-humankind/core/minecraft/protocol/packet"
)

// ------------------------- Config -------------------------

// Config ..
type Config struct {
	AuthServerAddress    string
	AuthServerToken      string
	RentalServerCode     string
	RentalServerPasscode string
}

// ------------------------- Client -------------------------

// Client ..
type Client struct {
	connection            *minecraft.Conn
	authClient            *auth.Client
	getCheckNumEverPassed bool
	cachedPacket          chan packet.Packet
}

// Conn ..
func (c Client) Conn() *minecraft.Conn {
	return c.connection
}

// CachedPacket ..
func (c Client) CachedPacket() chan packet.Packet {
	return c.cachedPacket
}

// ------------------------- MCPCheckChallengesSolver -------------------------

// MCPCheckChallengesSolver ..
type MCPCheckChallengesSolver struct {
	client *Client
}

// NewChallengeSolver ..
func NewChallengeSolver(client *Client) *MCPCheckChallengesSolver {
	return &MCPCheckChallengesSolver{client: client}
}
