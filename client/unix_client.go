//go:build !windows
package client

import (
	"net"

	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/types"
)

type unixClient struct {
	socketPath string
}

func newUnixClient(socketPath string) *unixClient {
	return &unixClient{
		socketPath: socketPath,
	}
}

func (c *unixClient) Get(endpoint string) (*types.Response, error) {
	conn, err := net.Dial("unix", c.socketPath)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if err := writeRequest(conn, "GET", endpoint, nil); err != nil {
		return nil, err
	}

	return readResponse(conn)
}

func (c *unixClient) Post(endpoint string, body []byte) (*types.Response, error) {
	conn, err := net.Dial("unix", c.socketPath)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if err := writeRequest(conn, "POST", endpoint, body); err != nil {
		return nil, err
	}

	return readResponse(conn)
}
