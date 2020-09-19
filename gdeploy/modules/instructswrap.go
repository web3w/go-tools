package modules

import (
	"errors"
	"fmt"
	"github.com/gisvr/go-tools/gdeploy/compile"
	"github.com/gisvr/golib/log"
	"gopkg.in/yaml.v2"
)

type InstructionsWrap struct {
	Copy     *Copy     `yaml:"copy"`
	Replace  *Replace  `yaml:"replace"`
	Template *Template `yaml:"template"`
	runner   []IBase
}

func (this *InstructionsWrap) Init(c *compile.CompileT, ctx *Context) error {
	this.runner = make([]IBase, 0)
	if this.Copy != nil {
		this.runner = append(this.runner, this.Copy)
	}
	if this.Replace != nil {
		this.runner = append(this.runner, this.Replace)
	}
	if this.Template != nil {
		log.Infof("template")
		this.runner = append(this.runner, this.Template)
	}
	if len(this.runner) > 1 {
		out, _ := yaml.Marshal(this)
		return errors.New(fmt.Sprintf("[InstructionsWrap] synax error,\n%s", out))
	}
	log.Infof("%d", len(this.runner))
	for _, m := range this.runner {

		err := m.Init(c, ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *InstructionsWrap) Run() error {
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

func (this *InstructionsWrap) SetVar(key string, value interface{}) {
	for _, r := range this.runner {
		r.SetVar(key, value)
	}
}
