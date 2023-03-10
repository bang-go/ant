package ant

import (
	"github.com/bang-go/ant/global"
	"github.com/bang-go/kit/blog"
)

type BlockExec func() error
type Block struct {
	Name  string    //阶段名
	Init  BlockExec //初始化
	Close BlockExec //结束
}

// exec blocks
func execBlocks(blocks ...Block) error {
	if len(blocks) > 0 {
		for _, v := range blocks {
			if err := v.Init(); err != nil {
				global.ALog.Error("init failed", blog.String("name", v.Name), blog.String("err", err.Error()))
				return err
			}
			global.ALog.Info("init successful", blog.String("name", v.Name))
		}
	}
	return nil
}
