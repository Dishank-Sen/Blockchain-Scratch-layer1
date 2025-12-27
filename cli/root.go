package cli

import (
	"context"
	"errors"
	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/utils"
	"github.com/spf13/cobra"
)

func Root(ctx context.Context) *cobra.Command{
	var rootCmd = &cobra.Command{
		Use: "bloc",
		Short: "this is simple decentralized application",
		PersistentPreRunE: persistentPreRunE,
	}

	rootCmd.SetContext(ctx)

	// loop which register all the commands
	for _, cmd := range Registered{
		c := cmd()
		rootCmd.AddCommand(c)
	}

	return rootCmd
}

func persistentPreRunE(cmd *cobra.Command, args []string) error{
	if cmd.Name() == "init"{
		return nil
	}
	
	// if .rec is not created prompt user to run init command
	if !utils.CheckDirExist(".bloc"){
		// log.Info(cmd.Context(), "debug-2")
		return errors.New("not a bloc repository, run 'bloc init' to initialize a empty bloc repository")
	}
	
	return nil
}