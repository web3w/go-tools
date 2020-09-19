package modules

import (
	"errors"
	"github.com/gisvr/go-tools/gdeploy/compile"
	"github.com/gisvr/golib/log"
	"os"
	"os/exec"
)

type Script struct {
	Base
	Cmds []string `yaml:"cmds"`
}

func (this *Script) Init(c *compile.CompileT, ctx *Context) error {
	this.name = getName(this)
	log.Infof("%s begin", this.name)
	return this.Base.Init(c, ctx)
}

func (this *Script) Run() error {
	defer log.Infof("%s end\n\n", this.name)
	if this.Cmds == nil {
		return errors.New("cmd is empty")
	}
	for _, cmd := range this.Cmds {
		cmd, err := this.Compile(cmd)
		if err != nil {
			return err
		}
		r := exec.Command(cmd)
		r.Stdout = os.Stdout
		r.Stderr = os.Stderr
		err = r.Run()
		if err != nil {
			return err
		}
	}
	return nil
}
