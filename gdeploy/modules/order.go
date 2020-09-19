package modules

import (
	"errors"
	"fmt"
	"github.com/gisvr/go-tools/gdeploy/compile"
	"gopkg.in/yaml.v2"
)

//支持所有命令按脚本定义的顺序执行
type OrderWrap struct {
	BuildWrap *BuildWrap `yaml:"build"`
	Command   *Command   `yaml:"command"`
	Script    *Script    `yaml:"script"`
	Copy      *Copy      `yaml:"copy"`
	Replace   *Replace   `yaml:"replace"`
	Template  *Template  `yaml:"template"`

	runner []IBase
}

func (this *OrderWrap) Init(c *compile.CompileT, ctx *Context) error {
	this.runner = make([]IBase, 0)
	if this.BuildWrap != nil {
		this.runner = append(this.runner, this.BuildWrap)
	}
	if this.Command != nil {
		this.runner = append(this.runner, this.Command)
	}
	if this.Script != nil {
		this.runner = append(this.runner, this.Script)
	}

	if this.Copy != nil {
		this.runner = append(this.runner, this.Copy)
	}
	if this.Replace != nil {
		this.runner = append(this.runner, this.Replace)
	}
	if this.Template != nil {
		this.runner = append(this.runner, this.Template)
	}

	if len(this.runner) > 1 {
		out, _ := yaml.Marshal(this)
		return errors.New(fmt.Sprintf("[OrderWrap] synax error,\n%s", out))
	}
	for _, m := range this.runner {
		err := m.Init(c, ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *OrderWrap) Run() error {
	if len(this.runner) > 1 {
		return errors.New("synax error")
	}
	for _, m := range this.runner {
		err := m.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *OrderWrap) SetVar(key string, value interface{}) {
	for _, r := range this.runner {
		r.SetVar(key, value)
	}
}
