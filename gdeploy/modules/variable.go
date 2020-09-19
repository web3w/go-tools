package modules

import (
	"github.com/gisvr/go-tools/gdeploy/compile"
	"github.com/gisvr/golib/log"
	"reflect"
)

type Variable struct {
	Vars map[string]interface{} `yaml:"vars"`
}

func (this *Variable) Init(c *compile.CompileT, ctx *Context) error {
	cnt := len(this.Vars)
	if cnt > 0 {
		log.Infof("define var begin")
	}
	err := compile.LoadVarsFromMap(this.Vars, c)
	if err != nil {
		return err
	}
	base := &Base{}
	base.Init(c, ctx)
	for k, v := range this.Vars {
		val, err := base.Compile(reflect.ValueOf(v).String())
		if err != nil {
			return err
		}
		log.Infof("var{%s} = %s", k, val)
	}
	if cnt > 0 {
		log.Infof("define var end")
	}
	return nil
}
