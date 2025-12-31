package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"slices"

	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/types"
	"github.com/Dishank-Sen/Blockchain-Scratch-layer1/utils/logger"
	"github.com/spf13/cobra"
)

func init(){
	Register("add", Add)
}

func Add() *cobra.Command{
	cmd := &cobra.Command{
		Use: "add",
		Short: "adds a bootstrap user",
		RunE: addRunE,
	}

	cmd.Flags().StringP("addr", "a", "", "ip address and port")
	return cmd
}

func addRunE(cmd *cobra.Command, args []string) error{
	addr, err := cmd.Flags().GetString("addr")
	if err != nil{
		return err
	}

	if err := SaveAddr(addr); err != nil{
		logger.Info("failed to save address")
		return err
	}else{
		logger.Info("address saved successfully")
		return nil
	}
}

func SaveAddr(addr string) error{
	path := path.Join(".bloc", "bootstrap.json")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		empty := []byte(`{"peers": []}`)
		os.WriteFile(path, empty, 0700)
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

	return os.WriteFile(path, newData, 0700)
}