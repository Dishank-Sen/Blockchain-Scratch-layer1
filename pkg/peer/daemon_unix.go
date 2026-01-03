//go:build !windows

package peer

import (
	"net"
	"time"

	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/constants"
	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/errors"
)

func isDaemonRunning() (bool, error) {
	conn, err := net.Dial("unix", constants.SocketPath)
	if err != nil {
		return false, nil
	}
	conn.Close()
	return true, nil
}

func waitForDaemon(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		conn, err := net.Dial("unix", constants.SocketPath)
		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(50 * time.Millisecond)
	}

	return errors.ErrDaemonTimeout
}
