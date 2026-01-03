package cli

import (
	"fmt"
	"strings"

	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/client"
	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/utils/logger"
	"github.com/spf13/cobra"
)

func init() {
	Register("peer", Peer)
}

func Peer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "peer",
		Short: "does various peer related operations",
		Args:  cobra.NoArgs,

		RunE: peerRunE,

		// IMPORTANT: do not spam usage for runtime errors
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.Flags().BoolP("list", "l", false, "list peers")
	return cmd
}

func peerRunE(cmd *cobra.Command, args []string) error {
	list, err := cmd.Flags().GetBool("list")
	if err != nil {
		return err
	}

	switch {
	case list:
		return handleList()
	default:
		return cmd.Help()
	}
}

func handleList() error {
	c := client.NewClient()

	resp, err := c.Get("/peers")
	if err != nil {
		if isDaemonNotRunning(err) {
			logger.Error("daemon is not running")
			fmt.Println("Run `bloc connect` to start the daemon.")
			return nil // graceful exit
		}

		// unexpected error
		return err
	}

	fmt.Println(string(resp.Body))
	return nil
}

func isDaemonNotRunning(err error) bool {
	if err == nil {
		return false
	}

	msg := err.Error()

	return strings.Contains(msg, "no such file or directory") ||
		strings.Contains(msg, "connection refused")
}
