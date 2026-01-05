package cli

import (
	"fmt"

	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/pkg/peer"
	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/utils/logger"
	"github.com/spf13/cobra"
)

func init(){
	Register("disconnect", Disconnect)
}

func Disconnect() *cobra.Command{
	return &cobra.Command{
		Use: "disconnect",
		Short: "disconnects to bootstrap server",
		RunE: disconnectRunE,
	}
}

func disconnectRunE(cmd *cobra.Command, args []string) error{
	peer, err := peer.NewPeer(cmd.Context())
	if err != nil{
		return err
	}
	if err := peer.Disconnect(); err != nil{
		logger.Error(fmt.Sprintf("error while stopping daemon: %v", err))
		return err
	}
	return nil
}