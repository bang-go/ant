package cmd

import (
	"github.com/bang-go/kit/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type JobCmd struct {
	Cmd *cobra.Command
	run func(args []string) error
}

type JobOptions struct {
	CmdUse   string
	CmdShort string
	CmdArgs  cobra.PositionalArgs
}

func NewJob(opt *JobOptions) Cmder {
	return &JobCmd{
		Cmd: &cobra.Command{
			Use:   util.If(opt.CmdUse != "", opt.CmdUse, "job"),
			Short: util.If(opt.CmdShort != "", opt.CmdShort, "start job"),
			Args:  opt.CmdArgs,
		},
	}
}
func (g *JobCmd) GetFlagSet() *pflag.FlagSet {
	return g.Cmd.Flags()
}
func (g *JobCmd) AddRun(fc RunFunc) {
	g.run = fc
}
func (g *JobCmd) GetCmd() Command {
	return g.Cmd
}
func (g *JobCmd) Register() {
	g.Cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return g.run(args)
	}
}
func (g *JobCmd) Args() {
}
