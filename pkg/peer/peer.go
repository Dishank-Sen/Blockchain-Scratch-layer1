package peer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/types"
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

func (p *Peer) Dial() error{
	tlsConfig := getTlsConfig()
	quicConfig := getQuicConfig()

	session, err := quic.DialAddr(p.ctx, p.addr, tlsConfig, quicConfig)
	if err != nil{
		return err
	}
	p.session = session
	return nil
}