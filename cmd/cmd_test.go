package cmd

import (
	"github.com/bang-go/ant"
	"github.com/bang-go/kit/blog"
	"github.com/spf13/cobra"
	"testing"
)

func TestCmd(t *testing.T) {
	artisan := ant.NewWithOption(&ant.Options{AllowLogLevel: blog.ErrorLevel})

	//添加CMD1
	cmdJob1 := NewJob(&JobOptions{CmdUse: "refresh_cache", CmdShort: "刷新缓存", CmdArgs: cobra.ExactArgs(2)})
	cmdJob1.AddRun(func(args []string) error {
		return nil
	})
	artisan.AddCmd(cmdJob1)
	artisan.AddCmd(Pack(NewJob(&JobOptions{CmdUse: "refresh_cache", CmdShort: "刷新缓存", CmdArgs: cobra.ExactArgs(2)}), func(args []string) error {
		return nil
	}))
}
