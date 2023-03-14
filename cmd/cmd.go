package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var RootCmd = &cobra.Command{
	Use: "root",
	//Short: "Git is a distributed version control system.",
	Args: cobra.ExactArgs(0),
}

type Command *cobra.Command
type RunFunc func(args []string) error
type Cmder interface {
	GetFlagSet() *pflag.FlagSet
	GetCmd() Command
	AddRun(fc RunFunc)
	Register()
}

func Pack(c Cmder, fc RunFunc) Cmder {
	c.AddRun(fc)
	return c
}
