package cmd

import (
	"fmt"
	"github.com/bang-go/ant/global"
	"github.com/bang-go/kit/bgin"
	"github.com/bang-go/kit/blog"
	"github.com/bang-go/kit/butil"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type GinCmd struct {
	opt    *GinOptions
	Cmd    *cobra.Command
	Client bgin.IGin
	run    RunFunc
}

type GinOptions struct {
	bgin.Options
	CmdUse   string
	CmdShort string
	CmdArgs  cobra.PositionalArgs
}

var DefineServerAddr = ":8080"
var flagServerPort string

func NewGin(opt *GinOptions) *GinCmd {
	client := bgin.New(&bgin.Options{Mode: opt.Mode, Addr: opt.Addr})
	return &GinCmd{
		opt:    opt,
		Client: client,
		Cmd: &cobra.Command{
			Use:   butil.If(opt.CmdUse != "", opt.CmdUse, "gin"),
			Short: butil.If(opt.CmdShort != "", opt.CmdShort, "Listen and serve http server"),
			Args:  opt.CmdArgs,
		},
	}
}

func (g *GinCmd) GetFlagSet() *pflag.FlagSet {
	return g.Cmd.Flags()
}
func (g *GinCmd) AddRun(fc RunFunc) {
	g.run = fc
}
func (g *GinCmd) GetCmd() Command {
	return g.Cmd
}
func (g *GinCmd) Register() {
	//添加默认flag
	g.addDefaultFlags()
	//添加默认Run func
	g.addDefaultRun()
	g.Cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return g.run(args)
	}
}

// 添加默认flag
func (g *GinCmd) addDefaultFlags() {
	g.Cmd.Flags().StringVarP(&flagServerPort, "port", "p", "", "http server port")
}

// 添加默认Run func
func (g *GinCmd) addDefaultRun() {
	if g.run != nil {
		return
	}
	g.run = func(args []string) error {
		var addr string
		if flagServerPort != "" { //优先级:命令行flag > 参数传递 > 默认
			addr = fmt.Sprintf(":%s", flagServerPort)
		} else {
			if g.opt.Addr != "" {
				addr = g.opt.Addr
			} else {
				addr = DefineServerAddr
			}
		}
		global.ALog.Info("http server start", blog.String("addr", addr))
		return g.Client.Run()
	}
}
