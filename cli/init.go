package cli

import (
	"context"
	"errors"
	"os"

	initdir "github.com/Dishank-Sen/Blockchain-Scratch-layer1/cli/initDIr"
	initfiles "github.com/Dishank-Sen/Blockchain-Scratch-layer1/cli/initFiles"
	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/utils"
	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/utils/logger"
	"github.com/spf13/cobra"
)

var ErrSkipRun = errors.New("cli: skip runE")

func init(){
	Register("init", Init)
}

func Init() *cobra.Command{
	return &cobra.Command{
		Use: "init",
		Short: "initialize a new bloc repository",
		RunE: initRunE,
		PersistentPreRunE: initPersistentPreRunE,
		SilenceUsage: true,     // prevents usage on error
		SilenceErrors: true,    // prevents printing sentinel error
	}
}

func initPersistentPreRunE(cmd *cobra.Command, args []string)error{
	rootDir := ".bloc"
	parentCtx := cmd.Context()
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	if utils.CheckDirExist(rootDir){
		logger.Info("Reinitializing bloc repository")
		if err := reinitialize(ctx, cancel); err != nil {
			return err // real error
		}
		return ErrSkipRun // signals to skip RunE
	}

	return nil
}

func initRunE(cmd *cobra.Command, args []string) error{
	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()

	// Create hidden .bloc directory
	if err := os.Mkdir(".bloc", 0700); err != nil{
		return err
	}

    // Create directories
    if err := createDir(ctx, cancel, false); err != nil{
		return err
	}

    // Create files
    if err := createFiles(ctx, cancel, false); err != nil{
		return err
	}

	logger.Info("Initialized empty bloc repository")
	return nil
}

func createFiles(ctx context.Context, cancel context.CancelFunc, reinit bool) error{
	for _, f := range initfiles.InitFiles{
		err := f(ctx, cancel, reinit)
		if err != nil{
			return err
		}
	}
	return nil
}

func createDir(ctx context.Context, cancel context.CancelFunc, reinit bool) error{
	for _, f := range initdir.InitDirectories{
		err := f(ctx, cancel, reinit)
		if err != nil{
			return err
		}
	}
	return nil
}

func reinitialize(ctx context.Context, cancel context.CancelFunc) error{
    // Create directories
    if err := createDir(ctx, cancel, true); err != nil{
		return err
	}

    // Create files
    if err := createFiles(ctx, cancel, true); err != nil{
		return err
	}

	logger.Info("Reinitialized bloc repository")
	return nil
}