//go:build windows
package peer

import (
	"context"
	"time"

	"github.com/Microsoft/go-winio"

	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/constants"
	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/errors"
)

func isDaemonRunning() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	conn, err := winio.DialPipeContext(ctx, constants.WindowsPipeName)
	if err != nil {
		return false, nil
	}
	conn.Close()
	return true, nil
}

func waitForDaemon(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		conn, err := winio.DialPipeContext(ctx, constants.WindowsPipeName)
		cancel()

		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(50 * time.Millisecond)
	}

	return errors.ErrDaemonTimeout
}
