package bind

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"

	"github.com/spf13/cobra"
	"main.go/users/types"
)

type BindCmd struct{
	Cmd *cobra.Command
}

func NewBindCmd() *BindCmd{
	cmd := &cobra.Command{
		Use: "bind",
		Short: "adds the ip and port as bootstrap nodes",
		RunE: RunE,
	}

	cmd.Flags().StringP("addr", "a", "", "ip address and port")

	return &BindCmd{
		Cmd: cmd,
	}
}

func RunE(cmd *cobra.Command, args []string) error{
	// fmt.Println("bind command executed..")
	addr, err := cmd.Flags().GetString("addr")
	if err != nil{
		return err
	}

	return SaveAddr(addr)
}

func SaveAddr(addr string) error{
	path := "bootstrap.json"

	if _, err := os.Stat(path); os.IsNotExist(err) {
		empty := []byte(`{"peers": []}`)
		os.WriteFile(path, empty, 0644)
	}

	prevData, err := os.ReadFile(path)
	if err != nil{
		return err
	}

	var u types.UsersIdentity
	err = json.Unmarshal(prevData, &u)
	if err != nil{
		return err
	}
	peers := u.Peers
	if slices.Contains(peers, addr){
		return fmt.Errorf("this address already exist in bootstrap")
	}
	u.Peers = append(peers, addr)

	newData, err := json.Marshal(u)
	if err != nil{
		return err
	}

	return os.WriteFile(path, newData, 0644)
}