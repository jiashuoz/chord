package chord

import (
	"crypto/sha1"
	"google.golang.org/grpc"
	"hash"
	"time"
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

	return config
}

// []grpc.DialOption{
// 	grpc.WithInsecure(), // TODO(ricky): find a better way to use this for testing
// },
