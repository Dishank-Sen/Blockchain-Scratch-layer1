package peer

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/client"
	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/constants"
	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/types"
	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/utils"
	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/utils/logger"
	"github.com/quic-go/quic-go"
)

type Peer struct{
	id string
	addr string
	ctx context.Context
	cancel context.CancelFunc
	session *quic.Conn
}

func NewPeer(parentCtx context.Context) *Peer{
	ctx, cancel := context.WithCancel(parentCtx)

	id, err := getID()
	if err != nil{
		panic(fmt.Sprintf("no user metadata - run bloc init, error: %v", err))
	}
	return &Peer{
		id: id,
		addr: "127.0.0.1:4242",  // "100.48.90.87:4242" -> public server address
		ctx: ctx,
		cancel: cancel,
	}
}

func getID() (string, error){
	filePath := path.Join(".bloc", "identity", "metadata.json")
	byteData, err := os.ReadFile(filePath)
	if err != nil{
		return "", err
	}

	var m types.Metadata
	if err := json.Unmarshal(byteData, &m); err != nil{
		return "", err
	}
	return m.ID, nil
}

func (p *Peer) Connect() error{
	if !isDaemonRunning(constants.SocketPath){
		cmd := exec.Command("blocd")
		if err := cmd.Start(); err != nil{
			return err
		}

		if err := waitForDaemon(constants.SocketPath, 3*time.Second); err != nil {
			return err
		}
	}

	client := client.NewClient(constants.SocketPath)

	id, err := utils.GetPeerID()
	if err != nil{
		return err
	}
	r := &types.RegisterBody{
		ID: id,
	}
	byteData, err := json.Marshal(r)
	if err != nil{
		return err
	}
	resp, err := client.Post("/register", byteData)
	if err != nil{
		logger.Debug("error in /register")
		return err
	}
	
	// logger.Debug(fmt.Sprintf("response message: %s", resp.Message))
	logger.Info(fmt.Sprintf("response body: %s", string(resp.Body)))
	return nil
}

func isDaemonRunning(sockPath string) bool {
	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		_ = os.Remove(sockPath)
		return false
	}
	conn.Close()
	return true
}

func waitForDaemon(sockPath string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		conn, err := net.Dial("unix", sockPath)
		if err == nil {
			conn.Close()
			// logger.Debug("daemon ready")
			return nil // daemon is ready
		}
		// logger.Debug("waiting for daemon")
		time.Sleep(50 * time.Millisecond)
	}

	return fmt.Errorf("daemon did not become ready in time")
}
