package start

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"main.go/users"
	"main.go/users/types"
)

type StartCmd struct{
	Cmd *cobra.Command
}

func NewStartCmd() *StartCmd{
	cmd := &cobra.Command{
		Use: "start",
		Short: "start the bootstrap nodes",
		Long: "it runs independent processes which bring up the bootstrap nodes",
		RunE: RunE,
	}
	return &StartCmd{
		Cmd: cmd,
	}
}

func getPeers() ([]string, error){
	path := "bootstrap.json"
	prevData, err := os.ReadFile(path)
	if err != nil{
		return []string{}, err
	}

	var u types.UsersIdentity
	err = json.Unmarshal(prevData, &u)
	if err != nil{
		return []string{}, err
	}
	return u.Peers, nil
}

func RunE(cmd *cobra.Command, args []string) error{
	fmt.Println("start cmd executed")
	peers, err := getPeers()
	fmt.Println(peers)
	if err != nil{
		panic(err)
	}
	var g errgroup.Group

	for _, peer := range peers {
		peer := peer
		g.Go(func() error {
			return users.MakeLive(peer)
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Println("one or more errors:", err)
	}
	return nil
}

