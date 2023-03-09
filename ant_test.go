package ant

import (
	"fmt"
	"github.com/bang-go/ant/cmd"
	"github.com/bang-go/kit/base/bint"
	"github.com/bang-go/kit/berror"
	"github.com/bang-go/kit/bviper"
	"github.com/bang-go/kit/env"
	"github.com/bang-go/kit/log"
	"github.com/spf13/cobra"
	"testing"
	"time"
)

func TestGin(t *testing.T) {
	var err error
	defer berror.PanicRecover()
	ant := NewWithOption(&Options{AllowLogLevel: log.InfoLevel})
	ant.AddBlock(Block{
		Name: "env",
		Init: func() error {
			return env.Configure()
		},
	}, Block{
		Name: "viper",
		Init: func() error {
			return bviper.Configure(&bviper.Options{ConfigType: "yaml", ConfigPaths: []string{"./config"}, ConfigNames: []string{"default", env.GetAppEnv()}})
		},
	}, Block{
		Name: "log",
		Init: func() error {
			_, err := log.New(&log.Options{Default: log.ConfigProd})
			return err
		},
		Close: func() error {
			return log.Sync()
		},
	})
	//添加CMD
	ant.AddCmd(cmd.NewGin(&cmd.GinOptions{}))
	if err = ant.Start(); err != nil {
		fmt.Println(err)
	}
	time.Sleep(5 * time.Second)
	_ = ant.Stop()
}

func TestJob(t *testing.T) {
	var err error
	defer berror.PanicRecover()
	ant := NewWithOption(&Options{AllowLogLevel: log.ErrorLevel})
	ant.AddBlock(Block{
		Name: "env",
		Init: func() error {
			return env.Configure()
		},
	}, Block{
		Name: "viper",
		Init: func() error {
			return bviper.Configure(&bviper.Options{ConfigType: "yaml", ConfigPaths: []string{"./config"}, ConfigNames: []string{"default", env.GetAppEnv()}})
		},
	}, Block{
		Name: "log",
		Init: func() error {
			_, err := log.New(&log.Options{Default: log.ConfigProd})
			return err
		},
		Close: func() error {
			return log.Sync()
		},
	})
	//添加CMD
	cmdJob := cmd.NewJob(&cmd.JobOptions{
		CmdArgs: cobra.MaximumNArgs(2),
	})
	cmdJob.AddRun(func(args []string) error {
		fmt.Println("I am job " + bint.String(bint.RandRange(3, 100)))
		fmt.Println(args)
		return nil
	})
	ant.AddCmd(cmdJob)
	if err = ant.Start(); err != nil {
		fmt.Println(err)
	}
	_ = ant.Stop()
}
