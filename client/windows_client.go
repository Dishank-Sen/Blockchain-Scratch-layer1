//go:build windows
package client

import (
	"context"
	"net"
	"time"

	"github.com/Microsoft/go-winio"

	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/types"
)

type windowsClient struct {
	windowsPipeName string
}

func newWindowsClient(windowsPipeName string) *windowsClient {
	return &windowsClient{
		windowsPipeName: windowsPipeName,
	}
}

func (c *windowsClient) dial() (net.Conn, error) {
	// Dial with timeout to avoid hanging forever
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return winio.DialPipeContext(ctx, c.windowsPipeName)
}

func (c *windowsClient) Get(endpoint string) (*types.Response, error) {
	conn, err := c.dial()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if err := writeRequest(conn, "GET", endpoint, nil); err != nil {
		return nil, err
	}

	return readResponse(conn)
}

func (c *windowsClient) Post(endpoint string, body []byte) (*types.Response, error) {
	conn, err := c.dial()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if err := writeRequest(conn, "POST", endpoint, body); err != nil {
		return nil, err
	}

	return readResponse(conn)
}
