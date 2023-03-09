package job

import (
	"fmt"
	"testing"
)

func TestJob(t *testing.T) {
	Register(map[string]Func{
		"h1": func() {
			abc()
		},
	})
	err := Do("h1")
	fmt.Println(err)
}

func abc() {

}
