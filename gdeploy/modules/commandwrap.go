package modules

import (
	"errors"
	"fmt"
	"github.com/gisvr/go-tools/gdeploy/compile"
	"gopkg.in/yaml.v2"
)

type CommandWrap struct {
	Command *Command `yaml:"command"`
	Script  *Script  `yaml:"script"`
	runner  []IBase
}

func (this *CommandWrap) Init(c *compile.CompileT, ctx *Context) error {
	this.runner = make([]IBase, 0)
	if this.Command != nil {
		this.runner = append(this.runner, this.Command)
	}
	if this.Script != nil {
		this.runner = append(this.runner, this.Script)
	}
	if len(this.runner) > 1 {
		out, _ := yaml.Marshal(this)
		return errors.New(fmt.Sprintf("[CommandWrap] synax error,\n%s", out))
	}
	for _, r := range this.runner {
		err := r.Init(c, ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *CommandWrap) Run() error {
	for _, r := range this.runner {
		err := r.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *CommandWrap) SetVar(key string, value interface{}) {
	for _, r := range this.runner {
		r.SetVar(key, value)
	}
}
