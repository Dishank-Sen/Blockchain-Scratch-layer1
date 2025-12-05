package main

import (
	"github.com/spf13/cobra"
	"main.go/cmd/root"
)

func main(){
	rootCmd := root.NewRootCmd()
	rootCmd.RegisterCmd()
	cobra.CheckErr(rootCmd.Cmd.Execute())
}