//go:build !windows

package client

import "github.com/Dishank-Sen/Blockchain-Scratch-layer1/constants"

func NewClient() Client{
	return newUnixClient(constants.SocketPath)
}