package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/client"
	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/utils/logger"
	"github.com/spf13/cobra"
)

func init() {
	Register("peer", Peer)
}

type PeerInfo struct {
	ID      string `json:"id"`
	Addr string `json:"addr"`
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
			return nil
		}
		return err
	}

	var peers []PeerInfo
	if err := json.Unmarshal(resp.Body, &peers); err != nil {
		return fmt.Errorf("invalid peer list response: %w", err)
	}

	printPeersTable(peers)
	return nil
}


func printPeersTable(peers []PeerInfo) {
	if len(peers) == 0 {
		fmt.Println("No peers connected.")
		return
	}

	// Column headers
	idHeader := "PEER ID"
	addrHeader := "ADDRESS"

	// Compute max widths
	idWidth := len(idHeader)
	addrWidth := len(addrHeader)

	for _, p := range peers {
		if len(p.ID) > idWidth {
			idWidth = len(p.ID)
		}
		if len(p.Addr) > addrWidth {
			addrWidth = len(p.Addr)
		}
	}

	// Row number column width
	indexWidth := len(fmt.Sprintf("%d", len(peers)))

	// Header
	fmt.Printf(
		"%-*s  %-*s  %-*s\n",
		indexWidth, "#",
		idWidth, idHeader,
		addrWidth, addrHeader,
	)

	// Separator
	fmt.Println(strings.Repeat("-", indexWidth+idWidth+addrWidth+6))

	// Rows
	for i, p := range peers {
		fmt.Printf(
			"%-*d  %-*s  %-*s\n",
			indexWidth, i+1,
			idWidth, p.ID,
			addrWidth, p.Addr,
		)
	}
}



func isDaemonNotRunning(err error) bool {
	if err == nil {
		return false
	}

	msg := err.Error()

	return strings.Contains(msg, "no such file or directory") ||
		strings.Contains(msg, "connection refused")
}
