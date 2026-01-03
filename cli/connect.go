package cli

import (
	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/pkg/peer"
	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/utils/logger"
	"github.com/spf13/cobra"
)

func init(){
	Register("connect", Connect)
}

func Connect() *cobra.Command{
	return &cobra.Command{
		Use: "connect",
		Short: "connects to bootstrap server",
		RunE: connectRunE,
	}
}

func connectRunE(cmd *cobra.Command, args []string) error{
	peer, err := peer.NewPeer(cmd.Context())
	if err != nil{
		return err
	}
	logger.Info("daemon started..")
	if err := peer.Connect(); err != nil{
		return err
	}
	return nil
}