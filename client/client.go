package client

import "github.com/Dishank-Sen/Blockchain-Scratch-layer1/types"

type Client interface{
	Get(endpoint string) (*types.Response, error)
	Post(endpoint string, body []byte) (*types.Response, error)
}