package peer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/client"
	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/types"
	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/utils"
	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/utils/logger"
)

type Peer struct {
	id     string
	ctx    context.Context
	cancel context.CancelFunc
}

func NewPeer(parentCtx context.Context) (*Peer, error) {
	ctx, cancel := context.WithCancel(parentCtx)

	id, err := getID()
	if err != nil {
		cancel()
		return nil, fmt.Errorf(
			"peer identity not found; run `bloc init`: %w",
			err,
		)
	}

	return &Peer{
		id:     id,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

func getID() (string, error) {
	filePath := path.Join(".bloc", "identity", "metadata.json")

	byteData, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	var m types.Metadata
	if err := json.Unmarshal(byteData, &m); err != nil {
		return "", err
	}

	if m.ID == "" {
		return "", errors.New("empty peer ID in metadata")
	}

	return m.ID, nil
}

func (p *Peer) Connect() error {
	if err := ensureDaemonRunning(3 * time.Second); err != nil {
		return err
	}

	c := client.NewClient()

	id, err := utils.GetPeerID()
	if err != nil {
		return err
	}

	req := &types.RegisterBody{ID: id}
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := c.Post("/register", body)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("daemon response: %s", string(resp.Body)))
	return nil
}

func (p *Peer) Disconnect() error {
	running, err := isDaemonRunning()
	if err != nil {
		return err
	}
	if !running {
		logger.Info("daemon already stopped")
		return nil
	}
	// Force kill if running
	logger.Warn("force killing daemon")
	return forceKillDaemon()
}
