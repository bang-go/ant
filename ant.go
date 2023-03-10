package ant

import (
	"github.com/bang-go/ant/cmd"
	"github.com/bang-go/ant/global"
	"github.com/bang-go/kit/log"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var ant *Artisan
var (
	AddBlock = ant.AddBlock
	DoBlock  = ant.DoBlock
	Start    = ant.Start
	Stop     = ant.Stop
)

// IAnt 类型定义
type IAnt interface {
	AddBlock(...Block)
	DoBlock(...Block) error
	AddCmd(...cmd.Cmder)
	Start() error
	Stop() error
}

type Artisan struct {
	opt      *Options
	Blocks   []Block
	DoBlocks []Block
	Cmds     []cmd.Cmder
	command  *cobra.Command
}

type Options struct {
	AllowLogLevel log.Level //允许的log level -1:Debug info:0 1:warn 2:error 3:dpanic 4 panic 5 fatal
}

// New creates a new ant instance.
func New() IAnt {
	opt := &Options{}
	return NewWithOption(opt)
}

func NewWithOption(opt *Options) IAnt {
	var err error
	ant = &Artisan{opt: opt, Blocks: []Block{}, DoBlocks: []Block{}, command: cmd.RootCmd}
	if err = ant.initAnt(); err != nil { //框架预加载
		panic(err)
	}
	return ant
}

func (a *Artisan) Start() error {
	var err error
	if err := execBlocks(a.Blocks...); err != nil {
		return err
	}
	if len(a.Cmds) > 0 {
		for _, v := range a.Cmds {
			v.Register()
			a.command.AddCommand(v.GetCmd())
		}
		return a.command.Execute()
	}
	return err
}

func (a *Artisan) initAnt() error {
	var err error
	//初始化日志客户端
	if global.ALog, err = initAntLog(a.opt.AllowLogLevel); err != nil {
		return err
	}
	return nil
}

// InitAntLog 初始化框架log
func initAntLog(logLevel log.Level) (*log.Logger, error) {
	//框架本身加载
	return log.New(&log.Options{
		Config: log.Config{
			Level:       zap.NewAtomicLevelAt(logLevel),
			Development: false,
			Sampling: &zap.SamplingConfig{
				Initial:    100,
				Thereafter: 100,
			},
			Encoding: "json",
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:        "time",
				LevelKey:       "level",
				NameKey:        "ant-logger",
				CallerKey:      "caller",
				FunctionKey:    zapcore.OmitKey,
				MessageKey:     "msg",
				StacktraceKey:  "stacktrace",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.LowercaseLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.SecondsDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		},
	})
}

// Stop 停止
func (a *Artisan) Stop() error {
	//框架相关
	_ = global.ALog.Sync()
	blocks := append(a.Blocks, a.DoBlocks...)
	//应用相关
	if len(blocks) > 0 {
		for _, v := range a.Blocks {
			if v.Close != nil { //todo:打印日志
				if err := v.Init(); err != nil {
					return err
				}
				global.ALog.Info("closed success", log.String("name", v.Name))
			}
		}
	}
	return nil
}

// AddBlock 添加Block
func (a *Artisan) AddBlock(blocks ...Block) {
	a.Blocks = append(a.Blocks, blocks...)
}

// DoBlock 添加并运行block
func (a *Artisan) DoBlock(blocks ...Block) error {
	a.DoBlocks = append(a.DoBlocks, blocks...)
	if err := execBlocks(a.DoBlocks...); err != nil {
		return err
	}
	return nil
}

func (a *Artisan) AddCmd(cmds ...cmd.Cmder) {
	a.Cmds = append(a.Cmds, cmds...)
}
