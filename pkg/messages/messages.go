package messages

import (
	// Standard
	"crypto/rsa"
	"encoding/gob"
	"fmt"

	// 3rd Party
	"github.com/satori/go.uuid"
)

// init registers message types with gob that are an interface for Base.Payload
func init() {
	gob.Register(KeyExchange{})
	gob.Register(AgentInfo{})
	gob.Register(SysInfo{})
}

const (

	// To Server

	// CHECKIN is used by the Agent to identify that it is checking in with the server
	CHECKIN = 1 // StatusCheckIn
	// OPAQUE is used to denote that embedded message contains an opaque structure
	OPAQUE = 2
	// JOBS is used to denote that the embedded message contains a list of job structures
	JOBS = 3
	// KEYEXCHANGE is used to denote that embedded message contains a key exchange structure
	KEYEXCHANGE = 4

	// To Agent

	// IDLE is used to notify the Agent that server has no tasks and that the Agent should idle
	IDLE = 10
)

// Base is the base JSON Object for HTTP POST payloads
type Base struct {
	Version float32     `json:"version"`
	ID      uuid.UUID   `json:"id"`
	Type    int         `json:"type"`
	Payload interface{} `json:"payload,omitempty"`
	Padding string      `json:"padding"`
	Token   string      `json:"token,omitempty"`
}

// KeyExchange is a JSON payload used to exchange public keys for encryption
type KeyExchange struct {
	PublicKey rsa.PublicKey `json:"publickey"`
}

// SysInfo is a JSON payload containing information about the system where the agent is running
type SysInfo struct {
	Platform     string   `json:"platform,omitempty"`
	Architecture string   `json:"architecture,omitempty"`
	UserName     string   `json:"username,omitempty"`
	UserGUID     string   `json:"userguid,omitempty"`
	HostName     string   `json:"hostname,omitempty"`
	Pid          int      `json:"pid,omitempty"`
	Ips          []string `json:"ips,omitempty"`
	Domain       string   `json:"domain,omitempty"`
}

// AgentInfo is a JSON payload containing information about the agent and its configuration
type AgentInfo struct {
	Version       string  `json:"version,omitempty"`
	Build         string  `json:"build,omitempty"`
	WaitTime      string  `json:"waittime,omitempty"`
	PaddingMax    int     `json:"paddingmax,omitempty"`
	MaxRetry      int     `json:"maxretry,omitempty"`
	FailedCheckin int     `json:"failedcheckin,omitempty"`
	Skew          int64   `json:"skew,omitempty"`
	Proto         string  `json:"proto,omitempty"`
	SysInfo       SysInfo `json:"sysinfo,omitempty"`
	KillDate      int64   `json:"killdate,omitempty"`
	JA3           string  `json:"ja3,omitempty"`
}

// String returns the text representation of a message constant
func String(messageType int) string {
	switch messageType {
	case KEYEXCHANGE:
		return "KeyExchange"
	case CHECKIN:
		return "StatusCheckIn"
	case JOBS:
		return "Jobs"
	case OPAQUE:
		return "OPAQUE"
	case IDLE:
		return "Idle"
	default:
		return fmt.Sprintf("Invalid: %d", messageType)
	}
}

