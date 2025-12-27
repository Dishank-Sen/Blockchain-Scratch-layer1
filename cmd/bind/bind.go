package bind

import (
	"fmt"

	"github.com/spf13/cobra"
)

type BindCmd struct{
	Cmd *cobra.Command
}

func NewBindCmd() *BindCmd{
	cmd := &cobra.Command{
		Use: "bind",
		Short: "binds two peers",
		Long: "connects two peers provided and doing the handshake",
		RunE: RunE,
	}

	cmd.Flags().StringP("addr1", "a1", "", "ip address and port for peer 1")
	cmd.Flags().StringP("addr2", "a2", "", "ip address and port for peer 2")

	return &BindCmd{
		Cmd: cmd,
	}
}

func RunE(cmd *cobra.Command, args []string) error{
	fmt.Println("bind cmd executed")

	return nil
}