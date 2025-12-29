package cli

import (
	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/pkg/peer"
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
	peer := peer.NewPeer(cmd.Context())
	return peer.Dial()
}