package client

import (
	"net"

	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/types"
)

type Client struct {
	socketPath string
}

func NewClient(socketPath string) *Client {
	return &Client{
		socketPath: socketPath,
	}
}

func (c *Client) Get(endpoint string) (*types.Response, error) {
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

func (c *Client) Post(endpoint string, body []byte) (*types.Response, error) {
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
