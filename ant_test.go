package ant

import (
	"fmt"
	"github.com/bang-go/ant/cmd"
	"github.com/bang-go/kit/base/bint"
	"github.com/bang-go/kit/berror"
	"github.com/bang-go/kit/blog"
	"github.com/bang-go/kit/bviper"
	"github.com/bang-go/kit/env"
	"github.com/spf13/cobra"
	"testing"
	"time"
)

func TestGin(t *testing.T) {
	var err error
	defer berror.PanicRecover()
	artisan := NewWithOption(&Options{AllowLogLevel: blog.InfoLevel})
	artisan.AddBlock(Block{
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
			_, err := blog.New(&blog.Options{Default: blog.ConfigProd})
			return err
		},
		Close: func() error {
			return blog.Sync()
		},
	})
	//添加CMD
	cmdGin := cmd.NewGin(&cmd.GinOptions{})
	cmdGin.GetCmd()
	artisan.AddCmd(cmd.NewGin(&cmd.GinOptions{}))
	if err = artisan.Start(); err != nil {
		fmt.Println(err)
	}
	time.Sleep(5 * time.Second)
	_ = artisan.Stop()
}

func TestJob(t *testing.T) {
	var err error
	defer berror.PanicRecover()
	artisan := NewWithOption(&Options{AllowLogLevel: blog.ErrorLevel})
	artisan.AddBlock(Block{
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
			_, err := blog.New(&blog.Options{Default: blog.ConfigProd})
			return err
		},
		Close: func() error {
			return blog.Sync()
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
	artisan.AddCmd(cmdJob)
	if err = artisan.Start(); err != nil {
		fmt.Println(err)
	}
	_ = artisan.Stop()
}
