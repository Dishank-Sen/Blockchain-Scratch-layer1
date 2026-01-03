package peer

import (
	"os/exec"
	"runtime"
	"time"
	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/utils/logger"
)

func ensureDaemonRunning(timeout time.Duration) error {
	running, err := isDaemonRunning()
	if err != nil {
		return err
	}

	if running {
		return nil
	}

	logger.Info("daemon not running, starting...")

	cmdName := "blocd"
	if runtime.GOOS == "windows" {
		cmdName = "blocd.exe"
	}

	cmd := exec.Command(cmdName)
	if err := cmd.Start(); err != nil {
		return err
	}

	return waitForDaemon(timeout)
}
