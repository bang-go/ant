package job

import "errors"

type Func func()

var jobs map[string]Func

// Register 工作注册
func Register(wk map[string]Func) {
	jobs = wk
	return
}

// Do 执行工作函数
func Do(id string) error {
	if _, ok := jobs[id]; ok {
		jobs[id]()
	}
	return errors.New("undefined job")
}
