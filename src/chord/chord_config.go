package chord

import (
	"crypto/sha1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/benchmark/latency"
	"hash"
	"time"
)

var (
	//Local simulates local network.
	Local = latency.Network{0, 0, 0}
	//LAN simulates local area network network.
	LAN = latency.Network{100 * 1024, 2 * time.Millisecond, 1500}
	//WAN simulates wide area network.
	WAN = latency.Network{20 * 1024, 30 * time.Millisecond, 1500}
	//Longhaul simulates bad network.
	Longhaul = latency.Network{1000 * 1024, 200 * time.Millisecond, 9000}
)

// Config contains configs for chord instance
type Config struct {
	ip string

	stabilizeTime time.Duration
	fixFingerTime time.Duration

	Hash     func() hash.Hash
	ringSize int

	timeout time.Duration

	DialOptions []grpc.DialOption

	networkSimulator *latency.Network
}

var defaultConfig = DefaultConfig()

// DefaultConfig returns a default configuration for chord
func DefaultConfig() *Config {
	config := &Config{
		stabilizeTime: 50 * time.Millisecond,
		fixFingerTime: 50 * time.Millisecond,
	}

	config.Hash = sha1.New
	config.ringSize = 4 // sha1: 160 bits
	config.DialOptions = []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTimeout(10 * time.Second),
		grpc.FailOnNonTempDialError(true),
		grpc.WithInsecure(),
	}

	config.networkSimulator = &WAN

	return config
}
