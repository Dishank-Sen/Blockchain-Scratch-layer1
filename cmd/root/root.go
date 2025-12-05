package root

import (
	"fmt"

	"github.com/spf13/cobra"
	"main.go/cmd/bind"
	"main.go/cmd/start"
)

type RootCmd struct{
	Cmd *cobra.Command
}

func NewRootCmd() *RootCmd{
	cmd := &cobra.Command{
		Use: "bloc",
		Short: "this is simple decentralized application",
		Long: "it is a CLI which can add bootstrap user and help them to discover other peers",
		RunE: RunE,
	}
	return &RootCmd{
		Cmd: cmd,
	}
}

func RunE(cmd *cobra.Command, args []string) error{
	fmt.Println("CLI started...")
	return nil
}

func (r *RootCmd) RegisterCmd(){
	r.Cmd.AddCommand(start.NewStartCmd().Cmd)
	r.Cmd.AddCommand(bind.NewBindCmd().Cmd)
}