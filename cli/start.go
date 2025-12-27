package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/types"
	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/users"
	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/utils/logger"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func init(){
	Register("start", Start)
}

func Start() *cobra.Command{
	return &cobra.Command{
		Use:   "start",
		Short: "start the bootstrap nodes",
		Long:  "it runs independent processes which bring up the bootstrap nodes",
		RunE:  startRunE,
	}
}

func startRunE(cmd *cobra.Command, args []string) error{
	peers, err := getPeers()
	if err != nil {
		return err
	}

	// root context for all peers
	rootCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// single signal handler for the whole process
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		logger.Info("Daemon stopping...")
		cancel()
	}()

	// use errgroup with the root context so all goroutines see cancellation
	g, gctx := errgroup.WithContext(rootCtx)

	for _, p := range peers {
		p := p // capture loop variable
		g.Go(func() error {
			// users.MakeLive should return when ctx is done
			return users.MakeLive(gctx, p)
		})
	}

	// wait for all peers to finish or return error
	if err := g.Wait(); err != nil {
		return fmt.Errorf("one or more errors: %w", err)
	}

	logger.Info("All bootstrap nodes exited cleanly")
	return nil
}

func getPeers() ([]string, error) {
	bootstrapPath := path.Join(".bloc", "bootstrap.json")
	prevData, err := os.ReadFile(bootstrapPath)
	if err != nil {
		return nil, err
	}

	var u types.UsersIdentity
	if err := json.Unmarshal(prevData, &u); err != nil {
		return nil, err
	}
	return u.Peers, nil
}