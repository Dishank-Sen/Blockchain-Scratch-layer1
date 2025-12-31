package peer

import (
	"time"

	"github.com/quic-go/quic-go"
)

func getQuicConfig() *quic.Config{
	return &quic.Config{
		MaxIdleTimeout: 30 * time.Second,
	}
}